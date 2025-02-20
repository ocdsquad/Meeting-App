package entity

import (
	"E-Meeting/pkg/utils"
	"encoding/json"

	"github.com/guregu/null"
)

type Room struct {
	ID               null.Int        `db:"id" json:"id"`
	Name             null.String     `db:"name" json:"name"`
	PriceHour        null.Float      `db:"price_hour" json:"price_hour"`
	IsActive         null.Bool       `db:"is_active" json:"is_active"`
	Description      null.String     `db:"description" json:"description"`
	RoomTypeID       null.Int        `db:"room_type_id" json:"room_type_id"`
	CapacityID       null.Int        `db:"capacity_id" json:"capacity_id"`
	Capacity         null.Int        `db:"capacity" json:"capacity"`
	AttachmentURL    null.String     `db:"attachment_url" json:"attachment_url"`
	CreatedAt        null.Time       `db:"created_at" json:"created_at"`
	CreatedBy        null.String     `db:"created_by,omitempty" json:"created_by"`
	UpdatedAt        null.Time       `db:"updated_at" json:"updated_at"`
	UpdatedBy        null.String     `db:"updated_by,omitempty" json:"updated_by"`
	RoomTypeName     null.String     `db:"room_type_name" json:"room_type_name"`
	RoomTypeIsActive null.Bool       `db:"room_type_is_active" json:"room_type_is_active"`
	ValueMinimum     null.Int        `db:"capacity_value_minimum" json:"value_minimum"`
	ValueMaximum     null.Int        `db:"capacity_value_maximum" json:"value_maximum"`
	Uom              null.String     `db:"capacity_uom" json:"uom"`
	CapacityIsActive null.Bool       `db:"capacity_is_active" json:"capacity_is_active"`
	ReservationsJson json.RawMessage `db:"reservations" json:"-"`
	Reservations     []Reservation
}

type Reservation struct {
	ReservationRoomID int64  `json:"reservation_room_id"`
	Date              string `json:"date"`
	StartTime         string `json:"start_time"`
	EndTime           string `json:"end_time"`
}

const TableRoomName = "rooms"
const TableRoomAliasName = "r"

type RoomsDataAccessObject struct {
	Rooms      []Room
	QueryCount utils.QueryCount
}
