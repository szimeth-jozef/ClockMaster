package models

import (
	"time"
)

type WorkDay struct {
	ID            uint `gorm:"primarykey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	WorkItemID    uint
	LastStartedAt *time.Time
	TotalDuration time.Duration
}

func (wd WorkDay) IsSameDateAs(date time.Time) bool {
	dateFormat := "2006-01-02"
	return wd.CreatedAt.Format(dateFormat) == date.Format(dateFormat)
}

func (wd WorkDay) GetNextDay(hour, min, sec, nsec int) WorkDay {
	nextDate := wd.CreatedAt.AddDate(0, 0, 1)
	next := time.Date(
		nextDate.Year(),
		nextDate.Month(),
		nextDate.Day(),
		hour,
		min,
		sec,
		nsec,
		wd.CreatedAt.Location(),
	)

	return WorkDay{
		CreatedAt:     next,
		WorkItemID:    wd.WorkItemID,
		TotalDuration: time.Duration(0),
	}
}
