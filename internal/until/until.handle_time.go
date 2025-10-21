package until

import "time"

func NowUTC() time.Time {
	return time.Now().UTC()
}

func SafeTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("15:04:05")
}

func SafeTimeUTC7(t *time.Time) string {
	if t == nil {
		return ""
	}
	localTime := t.Add(7 * time.Hour)
	return localTime.Format("15:04:05")
}
