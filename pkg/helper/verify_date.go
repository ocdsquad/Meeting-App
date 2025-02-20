package helper

import "time"

func VerifyDateFilter(startDateTime, endDateTime time.Time, query string, args ...any) (string, []any) {
	if !startDateTime.IsZero() && !endDateTime.IsZero() {
		query += "AND r.date BETWEEN $1 AND $2"
		args = append(args, startDateTime, endDateTime)
	} else if !startDateTime.IsZero() {
		query += "AND r.date >= $1"
		args = append(args, startDateTime)
	} else if !endDateTime.IsZero() {
		query += "AND r.date <= $1"
		args = append(args, endDateTime)
	}
	return query, args
}
