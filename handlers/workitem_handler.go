package handlers

import (
	"net/http"
	"szimeth-jozef/clockmaster/contracts/request"
	"szimeth-jozef/clockmaster/contracts/response"
	"szimeth-jozef/clockmaster/models"
	"szimeth-jozef/clockmaster/services/period"
	"szimeth-jozef/clockmaster/services/workitem"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type WorkItemHandler struct {
	DB              *gorm.DB
	WorkItemService workitem.WorkItemService
}

func (h WorkItemHandler) GetWorkItems(e echo.Context) error {
	year := e.QueryParam("year")
	month := e.QueryParam("month")

	p, err := period.FromString(year + "-" + month)
	if err != nil {
		log.Warn("Parsing string to period: " + err.Error())

		p = period.New()
	}

	var dbWorkItems []models.WorkItem
	h.DB.
		Preload("WorkDays").
		Where("period_year = ?", p.Year).
		Where("period_month = ?", p.Month).
		Find(&dbWorkItems)

	var workItems []response.WorkItemResponse
	for _, workItem := range dbWorkItems {
		workItems = append(workItems, response.WorkItemResponse{
			ID:                   workItem.ID,
			Created:              workItem.CreatedAt,
			Name:                 workItem.Name,
			Status:               workItem.Status,
			Period:               period.InvoicePeriod{Year: workItem.PeriodYear, Month: workItem.PeriodMonth},
			IsInvoiced:           workItem.IsInvoiced,
			TotalTimeNanoseconds: workItem.GetTotalTime(),
			IsRunning:            workItem.IsRunning(),
		})
	}

	response := response.WorkItemsOfPeriod{
		Period:    p,
		WorkItems: workItems,
	}

	return e.JSON(http.StatusOK, response)
}

func (h WorkItemHandler) GetWorkItem(e echo.Context) error {
	var workItem models.WorkItem
	if err := h.DB.Preload("WorkDays").First(&workItem, e.Param("id")).Error; err != nil {
		log.Error(err)
		return e.JSON(http.StatusNotFound, nil)
	}

	response := response.WorkItemResponse{
		ID:                   workItem.ID,
		Created:              workItem.CreatedAt,
		Name:                 workItem.Name,
		Status:               workItem.Status,
		Period:               period.InvoicePeriod{Year: workItem.PeriodYear, Month: workItem.PeriodMonth},
		IsInvoiced:           workItem.IsInvoiced,
		TotalTimeNanoseconds: workItem.GetTotalTime(),
		IsRunning:            workItem.IsRunning(),
	}

	return e.JSON(http.StatusOK, response)
}

func (h WorkItemHandler) CreateWorkItem(e echo.Context) error {
	var workItemData request.CreateWorkItemData
	if err := e.Bind(&workItemData); err != nil {
		log.Error(err)
		return e.JSON(http.StatusBadRequest, nil)
	}

	workItem, err := h.WorkItemService.Create(workItemData)
	if err != nil {
		log.Error(err)
		return e.JSON(http.StatusBadRequest, nil)
	}

	response := response.WorkItemResponse{
		ID:                   workItem.ID,
		Created:              workItem.CreatedAt,
		Name:                 workItem.Name,
		Status:               workItem.Status,
		Period:               period.InvoicePeriod{Year: workItem.PeriodYear, Month: workItem.PeriodMonth},
		IsInvoiced:           workItem.IsInvoiced,
		TotalTimeNanoseconds: workItem.GetTotalTime(),
		IsRunning:            workItem.IsRunning(),
	}

	return e.JSON(http.StatusOK, response)
}

func (h WorkItemHandler) StartWorkItem(e echo.Context) error {
	now := time.Now()

	// Try to get the specified work item
	var workItem models.WorkItem
	if err := h.DB.Preload("WorkDays").First(&workItem, e.Param("id")).Error; err != nil {
		log.Error(err)
		return e.JSON(http.StatusNotFound, nil)
	}

	// Cannot start a work item when
	// - a work item is already started (doesnt matter even if it is the same)
	// - it is already invoiced

	if workItem.IsInvoiced {
		log.Error("cannot start an invoiced work item")
		return e.JSON(http.StatusBadRequest, nil)
	}

	if h.ExistsRunningWorkItem() {
		log.Error("cannot start a work item when another one is already started")
		return e.JSON(http.StatusBadRequest, nil)
	}

	// Try find a work day for today
	workDay := workItem.GetWorkDayForDate(now)
	if workDay == nil {
		// If there is no work day for today, create one
		workDay = &models.WorkDay{
			WorkItemID:    workItem.ID,
			TotalDuration: time.Duration(0),
		}
		h.DB.Create(&workDay)
		log.Info("created work day for today")
	} else {
		log.Info("found work day for today")
	}

	// Set the last started
	workDay.LastStartedAt = &now
	h.DB.Save(&workDay)

	// Set the work item status to in progress
	if workItem.Status != models.InProgress {
		workItem.Status = models.InProgress
		h.DB.Save(&workItem)
	}

	// Refresh the work item
	h.DB.Preload("WorkDays").First(&workItem, e.Param("id"))

	return e.JSON(http.StatusOK, workItem)
}

func (h WorkItemHandler) StopWorkItem(e echo.Context) error {
	stoppedWorkItems, err := h.WorkItemService.Stop()
	if err != nil {
		log.Error(err)
		return e.JSON(http.StatusBadRequest, nil)
	}

	return e.JSON(http.StatusOK, stoppedWorkItems)
}

func (h WorkItemHandler) DeleteWorkItem(e echo.Context) error {
	var workItem models.WorkItem
	if err := h.DB.Preload("WorkDays").First(&workItem, e.Param("id")).Error; err != nil {
		log.Error(err)
		return e.JSON(http.StatusNotFound, nil)
	}

	// Cannot delete running work item
	// TODO: questionable, maybe implement later

	// Cannot delete an invoiced work item
	if workItem.IsInvoiced {
		log.Error("cannot delete an invoiced work item")
		return e.JSON(http.StatusBadRequest, nil)
	}

	if len(workItem.WorkDays) > 0 {
		// Cascade delete, delete first the work days
		h.DB.Delete(&workItem.WorkDays)
	}

	// Delete the work item
	h.DB.Delete(&workItem)

	return e.JSON(http.StatusOK, nil)
}

func (h WorkItemHandler) GetStatus(e echo.Context) error {
	workItem, workDay := h.GetRunningWorkItemWithWorkDay()
	if workItem == nil || workDay == nil {
		return e.JSON(http.StatusOK, response.StatusResponse{
			IsRunning:                false,
			DeltaDurationNanoseconds: nil,
			WorkItem:                 nil,
		})
	}

	deltaDuration := time.Since(*workDay.LastStartedAt)

	return e.JSON(http.StatusOK, response.StatusResponse{
		IsRunning:                true,
		DeltaDurationNanoseconds: &deltaDuration,
		WorkItem: &response.WorkItemResponse{
			ID:                   workItem.ID,
			Created:              workItem.CreatedAt,
			Name:                 workItem.Name,
			Status:               workItem.Status,
			Period:               period.InvoicePeriod{Year: workItem.PeriodYear, Month: workItem.PeriodMonth},
			IsInvoiced:           workItem.IsInvoiced,
			TotalTimeNanoseconds: workItem.GetTotalTime(),
			IsRunning:            workItem.IsRunning(),
		},
	})
}

func (h WorkItemHandler) GetWorkDays(e echo.Context) error {
	var workItem models.WorkItem
	if err := h.DB.Preload("WorkDays").First(&workItem, e.Param("id")).Error; err != nil {
		log.Error(err)
		return e.JSON(http.StatusNotFound, nil)
	}

	return e.JSON(http.StatusOK, workItem.WorkDays)
}

func (h WorkItemHandler) MarkAsDone(e echo.Context) error {
	var workItem models.WorkItem
	if err := h.DB.Preload("WorkDays").First(&workItem, e.Param("id")).Error; err != nil {
		log.Error(err)
		return e.JSON(http.StatusNotFound, nil)
	}

	var markAsDoneData request.MarkAsDoneRequest
	if err := e.Bind(&markAsDoneData); err != nil {
		log.Error(err)
		return e.JSON(http.StatusBadRequest, nil)
	}

	// Work item cannot be marked as done when
	// - it is is running
	if workItem.IsRunning() {
		log.Error("cannot mark as done a work item that is running")
		return e.JSON(http.StatusBadRequest, "cannot mark as done a work item that is running")
	}

	// Validation
	if len(workItem.WorkDays) != len(markAsDoneData.WorkDays) {
		log.Error("work day count mismatch")
		return e.JSON(http.StatusBadRequest, "work day count mismatch")
	}

	for _, workDayData := range markAsDoneData.WorkDays {
		dbWorkDay := workItem.GetWorkDayByID(workDayData.WorkDayID)
		if dbWorkDay == nil {
			log.Error("specified work day is not part of the work item")
			return e.JSON(http.StatusBadRequest, "specified work day is not part of the work item")
		}
	}

	// Update the work days
	for _, workDayData := range markAsDoneData.WorkDays {
		dbWorkDay := workItem.GetWorkDayByID(workDayData.WorkDayID)
		dbWorkDay.TotalDuration = time.Duration(workDayData.RoundedDurationInHours) * time.Hour
		h.DB.Save(&dbWorkDay)
	}

	// Update the work item
	workItem.Status = models.Done
	h.DB.Save(&workItem)

	return e.JSON(http.StatusOK, nil)
}

func (h WorkItemHandler) ExistsRunningWorkItem() bool {
	var workDays []models.WorkDay

	h.DB.Where("last_started_at IS NOT NULL").Find(&workDays)

	return len(workDays) > 0
}

func (h WorkItemHandler) GetRunningWorkItemWithWorkDay() (*models.WorkItem, *models.WorkDay) {
	var workDays []models.WorkDay

	h.DB.Where("last_started_at IS NOT NULL").Find(&workDays)

	if len(workDays) == 0 {
		return nil, nil
	}

	if len(workDays) != 1 {
		log.Error("more than one started work item (to mitigate, run stop more times)")
		return nil, nil
	}

	workDay := workDays[0]
	var workItem models.WorkItem

	h.DB.Preload("WorkDays").First(&workItem, workDay.WorkItemID)

	return &workItem, &workDay
}
