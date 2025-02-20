package reason

import "errors"

var (
	ErrDataNotFound      = errors.New("data not found")
	ErrFailedInsertData  = errors.New("failed insert data")
	ErrFailedDeleteData  = errors.New("failed delete data")
	ErrInvalidFormatTime = errors.New("invalid format time")
	ErrDurationMinus     = errors.New("duration must be > 0")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrDateIsScheduled   = errors.New("there is already a schedule")
)
