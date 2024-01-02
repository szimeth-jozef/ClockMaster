package request

type CreateWorkItemData struct {
	Name                       string `json:"name"`
	PeriodMonth                int    `json:"period_month"`
	PeriodYear                 int    `json:"period_year"`
	InitTotalDurationInSeconds int    `json:"init_total_duration_in_seconds"`
}
