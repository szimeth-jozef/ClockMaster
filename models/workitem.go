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
	ID         uint `gorm:"primarykey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Name       string
	Status     WorkItemStatus
	Period     string
	IsInvoiced bool
	WorkDays   []WorkDay
}

func (wi WorkItem) GetWorkDayForDate(date time.Time) *WorkDay {
	for _, workDay := range wi.WorkDays {
		if workDay.IsSameDateAs(date) {
			return &workDay
		}
	}
	return nil
}
