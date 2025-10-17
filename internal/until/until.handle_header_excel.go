package until

import "time"

// GetDateRange returns all dates between startDay and endDay (inclusive)
func GetDateRange(start, end *time.Time) []string {
	const layout = "2006-01-02"
	var dates []string
	for d := *start; !d.After(*end); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format(layout))
	}
	return dates
}

func GetHeaderSummary() []string {
	headers := []string{
		"Month", "Date", "Username", "Full name", "Client", "Job", "Item", "Role",
		"Start", "End", "Hour", "Minute", "Break Time", "Total Time(m)", "OT(m)",
		"Volume(件)", "Note",
	}
	return headers
}
