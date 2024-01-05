package workitem

import (
	"fmt"
	"math"
	"szimeth-jozef/clockmaster/contracts/request"
	"szimeth-jozef/clockmaster/models"
	"szimeth-jozef/clockmaster/services/period"
	"time"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type WorkItemService struct {
	DB *gorm.DB
}

func (s WorkItemService) Create(data request.CreateWorkItemData) (*models.WorkItem, error) {
	p, err := period.FromPart(data.PeriodYear, data.PeriodMonth)
	if err != nil {
		return nil, err
	}

	workItem := models.WorkItem{
		Name:       data.Name,
		Status:     models.Todo,
		Period:     p.String(),
		IsInvoiced: false,
	}

	s.DB.Create(&workItem)

	if data.InitTotalDurationInSeconds != 0 {
		initTotalDurationSeconds := int(math.Abs(float64(data.InitTotalDurationInSeconds)))
		initTotalDuration := time.Duration(initTotalDurationSeconds) * time.Second

		workDay := models.WorkDay{
			WorkItemID:    workItem.ID,
			TotalDuration: initTotalDuration,
		}

		s.DB.Create(&workDay)
	}

	return &workItem, nil
}

func (s WorkItemService) Stop() ([]models.WorkItem, error) {
	now := time.Now()

	var workDays []models.WorkDay
	s.DB.Where("last_started_at IS NOT NULL").Find(&workDays)

	if len(workDays) == 0 {
		return nil, fmt.Errorf("no running work item")
	}

	var stoppedWorkItems []models.WorkItem

	if len(workDays) == 1 {
		// Expected scenario, stop the work item
		workDay := workDays[0]
		var workItem models.WorkItem
		s.DB.Preload("WorkDays").First(&workItem, workDay.WorkItemID)

		s.stopSingleWorkItem(&workItem, &workDay, now)

		stoppedWorkItems := append(stoppedWorkItems, workItem)
		return stoppedWorkItems, nil
	}

	// More than one started work item (to mitigate, run stop more times)
	log.Info("unexpected scenario: more than one started work item (to mitigate, run stop more times)")
	workItemCache := make(map[uint]models.WorkItem)

	for _, workDay := range workDays {
		if workItem, ok := workItemCache[workDay.WorkItemID]; ok {
			// Work item already cached
			s.stopSingleWorkItem(&workItem, &workDay, now)
		} else {
			// Work item not cached, cache it
			var workItem models.WorkItem
			s.DB.Preload("WorkDays").First(&workItem, workDay.WorkItemID)
			workItemCache[workDay.WorkItemID] = workItem
			s.stopSingleWorkItem(&workItem, &workDay, now)
		}
	}

	for _, workItem := range workItemCache {
		var wi models.WorkItem
		s.DB.Preload("WorkDays").First(&wi, workItem.ID)
		stoppedWorkItems = append(stoppedWorkItems, wi)
	}

	return stoppedWorkItems, nil
}

func (s WorkItemService) stopSingleWorkItem(
	runningWorkItem *models.WorkItem,
	timedWorkDay *models.WorkDay,
	stopTime time.Time) error {

	if timedWorkDay.IsSameDateAs(stopTime) {
		// If the work day is the same as today, handle stopping normally
		timedWorkDay.TotalDuration += stopTime.Sub(*timedWorkDay.LastStartedAt)
		timedWorkDay.LastStartedAt = nil
		s.DB.Save(&timedWorkDay)
		return nil
	}

	// If the work day is not the same as today, split the work day
	// into as many work days as needed
	log.Info("Work item stopped on a different day than started")

	wd := timedWorkDay
	for !wd.IsSameDateAs(stopTime) {
		nextWorkDay := wd.GetNextDay(0, 0, 0, 0)
		wd.TotalDuration += nextWorkDay.CreatedAt.Sub(*wd.LastStartedAt)
		wd.LastStartedAt = nil
		s.DB.Save(&wd)

		nextWorkDay.LastStartedAt = &nextWorkDay.CreatedAt
		s.DB.Create(&nextWorkDay)
		wd = &nextWorkDay
	}

	wd.TotalDuration += stopTime.Sub(*wd.LastStartedAt)
	wd.LastStartedAt = nil
	s.DB.Save(&wd)

	return nil
}
