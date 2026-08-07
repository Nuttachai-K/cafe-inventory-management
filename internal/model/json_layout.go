package model

import "time"

const jsonTimestampDateHourLayout = "2006-01-02 15:04:05"

const jsonTimestampHourLayout = "15:04"

var jst = time.FixedZone("JST", 9*60*60)

func formatDateHourJST(t time.Time) string {
	return t.In(jst).Format(jsonTimestampDateHourLayout)
}

func formatHour(t time.Time) string {
	return t.Format(jsonTimestampHourLayout)
}

func parseClockTime(s string) (time.Time, error) {
	if t, err := time.Parse("15:04:05", s); err == nil {
		return t, nil
	}
	return time.Parse("15:04", s)
}
