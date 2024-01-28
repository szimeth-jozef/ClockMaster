package request

type RoundedWorkDayData struct {
	WorkDayID              uint `json:"workday_id"`
	RoundedDurationInHours uint `json:"rounded_duration_in_hours"`
}

type MarkAsDoneRequest struct {
	WorkDays []RoundedWorkDayData `json:"workdays"`
}
