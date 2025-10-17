package until

import "time"

func SafeTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("15:04:05")
}
