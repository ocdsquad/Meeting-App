package entity

import (
	"github.com/guregu/null"
	"time"
)

type ReservationRooms struct {
	ID               int64       `json:"id" db:"id"`
	UserID           int64       `json:"user_id" db:"user_id"`
	RoomID           int64       `json:"room_id" db:"room_id"`
	Status           string      `json:"status" db:"status"`
	SnackID          null.Int    `json:"category_snack_id" db:"category_snack_id"`
	Name             string      `json:"name" db:"name"`
	Date             time.Time   `json:"date" db:"date"`
	StartTime        time.Time   `json:"start_time" db:"start_time"`
	EndTime          time.Time   `json:"end_time" db:"end_time"`
	Phone            string      `json:"phone" db:"phone"`
	TotalParticipant int64       `json:"total_participant" db:"total_participant"`
	Organization     string      `json:"organization" db:"organization"`
	Note             null.String `json:"notes" db:"notes"`
	TotalDuration    int64       `json:"total_duration" db:"total_duration"`
	GrandTotal       float64     `json:"grand_total" db:"grand_total"`
	CreatedBy        null.String `json:"created_aby" db:"created_by"`
	CreatedAt        null.Time   `json:"created_at" db:"created_at"`
	UpdatedAt        null.Time   `json:"updated_at" db:"updated_at"`
}

// type ReservationResultDataAccessObject struct {
// 	Reservation []ReservationRooms
// 	QueryCount  utils.QueryCount
// }
