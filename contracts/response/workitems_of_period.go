package response

import (
	"szimeth-jozef/clockmaster/models"
	"szimeth-jozef/clockmaster/services/period"
)

type WorkItemsOfPeriod struct {
	Period    period.InvoicePeriod `json:"period"`
	WorkItems []models.WorkItem    `json:"work_items"`
}
