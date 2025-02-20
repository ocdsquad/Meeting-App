package model

import (
	"E-Meeting/pkg/utils"

	"github.com/guregu/null"
)

type Room struct {
	ID            null.Int       `json:"id"`
	Name          null.String    `json:"name"`
	PriceHour     null.Float     `json:"price_hour" `
	IsActive      null.Bool      `json:"is_active"`
	Description   null.String    `json:"description"`
	RoomTypeID    null.Int       `json:"room_type_id"`
	CapacityID    null.Int       `json:"capacity_id"`
	Capacity      int64          `json:"capacity"`
	RoomTypeName  string         `json:"room_type_name"`
	AttachmentURL null.String    `json:"attachment_url"`
	Reservations  *[]Reservation `json:"reservations"`
	// CreatedAt     null.Time   `json:"created_at" db:"created_at"`
}
type RoomDetail struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	PriceHour    float64 `json:"price_hour"`
	RoomTypeName string  `json:"room_type_name"`
	Capacity     int64   `json:"capacity"`
}

type RoomDashboards struct {
	ID         int     `json:"id" db:"room_id"`
	Name       string  `json:"name" db:"room_name"`
	PriceHour  float64 `json:"price_hour" db:"price_hour"`
	Percentage float64 `json:"percentage " db:"percentage"`
}

type RoomsDataAccessObject struct {
	Rooms      []Room
	QueryCount utils.QueryCount
}

type RoomRequest struct {
	Name          string      `form:"name" validate:"required"`
	Price         float64     `form:"price" validate:"required"`
	RoomTypeID    int         `form:"room_type_id" validate:"required"`
	Capacity      int         `form:"capacity" validate:"required"`
	Description   null.String `form:"description"`
	AttachmentURL null.String `form:"attachment_url"`
}

type RoomRequestSwagger struct {
	Name          string  `form:"name" validate:"required"`
	Price         float64 `form:"price" validate:"required"`
	RoomTypeID    int     `form:"room_type_id" validate:"required"`
	Capacity      int     `form:"capacity" validate:"required"`
	AttachmentURL *string `form:"attachment_url"`
}

type FilterDataRoomRequest struct {
	Capacity *int64
	RoomType *int64
}

type Reservation struct {
	ReservationRoomID int64  `json:"reservation_room_id"`
	Date              string `json:"date"`
	StartTime         string `json:"start_time"`
	EndTime           string `json:"end_time"`
}
