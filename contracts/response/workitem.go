package response

import (
	"szimeth-jozef/clockmaster/models"
	"szimeth-jozef/clockmaster/services/period"
	"time"
)

type WorkItemResponse struct {
	ID                   uint                  `json:"id"`
	Created              time.Time             `json:"created"`
	Name                 string                `json:"name"`
	Status               models.WorkItemStatus `json:"status"`
	Period               period.InvoicePeriod  `json:"period"`
	IsInvoiced           bool                  `json:"isInvoiced"`
	TotalTimeNanoseconds time.Duration         `json:"totalTimeNanoseconds"`
	IsRunning            bool                  `json:"isRunning"`
}

type WorkItemsOfPeriod struct {
	Period    period.InvoicePeriod `json:"period"`
	WorkItems []WorkItemResponse   `json:"workItems"`
}

type StatusResponse struct {
	IsRunning                bool              `json:"isRunning"`
	DeltaDurationNanoseconds *time.Duration    `json:"deltaDurationNanoseconds"`
	WorkItem                 *WorkItemResponse `json:"workItem"`
}
