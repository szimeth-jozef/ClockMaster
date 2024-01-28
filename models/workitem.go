package models

import (
	"time"
)

type WorkItemStatus uint

const (
	Todo WorkItemStatus = iota
	InProgress
	Done
)

type WorkItem struct {
	ID          uint `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	Status      WorkItemStatus
	PeriodYear  int
	PeriodMonth int
	IsInvoiced  bool
	WorkDays    []WorkDay
}

func (wi WorkItem) GetTotalTime() time.Duration {
	var totalTime time.Duration
	for _, workDay := range wi.WorkDays {
		totalTime += workDay.TotalDuration
	}
	return totalTime
}

func (wi WorkItem) GetWorkDayForDate(date time.Time) *WorkDay {
	for _, workDay := range wi.WorkDays {
		if workDay.IsSameDateAs(date) {
			return &workDay
		}
	}
	return nil
}

func (wi WorkItem) IsRunning() bool {
	for _, workDay := range wi.WorkDays {
		if workDay.LastStartedAt != nil {
			return true
		}
	}
	return false
}

func (wi WorkItem) GetWorkDayByID(id uint) *WorkDay {
	for _, workDay := range wi.WorkDays {
		if workDay.ID == id {
			return &workDay
		}
	}

	return nil
}
