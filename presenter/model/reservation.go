package model

import (
	"E-Meeting/internal/domain/entity"
	"encoding/json"
	"strings"
	"time"

	"github.com/guregu/null"
)

type CustomDate time.Time

const dateFormat = "2006-01-02" // Format tanggal (YYYY-MM-DD)

func (d *CustomDate) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	parsedDate, err := time.Parse(dateFormat, str)
	if err != nil {
		return err
	}
	*d = CustomDate(parsedDate)
	return nil
}

func (d CustomDate) MarshalJSON() ([]byte, error) {
	formatted := time.Time(d).Format(dateFormat)
	return json.Marshal(formatted)
}

func (d CustomDate) String() string {
	return time.Time(d).Format(dateFormat)
}

type CustomTime time.Time

const TimeFormat = "15:04:05" // Only time format (HH:mm:ss)

func (t *CustomTime) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	parsedTime, err := time.Parse(TimeFormat, str)
	if err != nil {
		return err
	}
	*t = CustomTime(parsedTime)
	return nil
}

func (t CustomTime) MarshalJSON() ([]byte, error) {
	formatted := time.Time(t).Format(TimeFormat)
	return json.Marshal(formatted)
}

func (t CustomTime) String() string {
	return time.Time(t).Format(TimeFormat)
}

type ReservationCreateRequest struct {
	RoomID           int         `json:"room_id" validate:"required"`
	SnackID          null.Int    `json:"category_snack_id"`
	Name             string      `json:"name" validate:"required"`
	Date             CustomDate  `json:"date" validate:"required"`
	StartTime        CustomTime  `json:"start_time" validate:"required"`
	EndTime          CustomTime  `json:"end_time" validate:"required"`
	Phone            string      `json:"phone" validate:"required"`
	TotalParticipant int         `json:"total_participant" validate:"required"`
	Organization     string      `json:"organization" validate:"required"`
	Note             null.String `json:"notes"`
}

type ReservationCreateServiceResponse struct {
	Room                          *entity.Room
	Snack                         *entity.Snack
	RequestInput                  ReservationCreateRequest
	Duration                      int64
	GrandTotalPrice               float64
	TotalPriceReservationDuration float64
	TotalPriceSnack               float64
	ReservationCode               string
}

type ReservationCreateResponse struct {
	RoomName                      string      `json:"room_name"`
	RoomType                      string      `json:"room_type"`
	RoomCapacity                  int64       `json:"room_capacity"`
	RoomPrice                     float64     `json:"room_price"`
	ReservationName               string      `json:"reservation_name"`
	ReservationPhone              string      `json:"reservation_phone"`
	ReservationOrganization       string      `json:"reservation_organization"`
	ReservationDate               string      `json:"date"`
	Duration                      int64       `json:"duration"`
	TotalParticipant              int64       `json:"total_participant"`
	Snack                         *Snack      `json:"snack"`
	GrandTotalPrice               float64     `json:"grand_total_price"`
	TotalPriceReservationDuration float64     `json:"total_price_reservation_duration"`
	TotalPriceSnack               float64     `json:"total_price_snack"`
	Note                          null.String `json:"notes"`
	ReservationCode               string      `json:"reservation_code"`
}

type ReservationCodeRequest struct {
	Code string `json:"code" validate:"required"`
}

type ReservationCreateRequestSwagger struct {
	RoomID           int        `json:"room_id" validate:"required"`
	SnackID          *int       `json:"category_snack_id"`
	Name             string     `json:"name" validate:"required"`
	Date             CustomDate `json:"date" validate:"required"`
	StartTime        CustomTime `json:"start_time" validate:"required"`
	EndTime          CustomTime `json:"end_time" validate:"required"`
	Phone            string     `json:"phone" validate:"required"`
	TotalParticipant int        `json:"total_participant" validate:"required"`
	Organization     string     `json:"organization" validate:"required"`
	Note             *string    `json:"notes"`
}

type ReservationGetResponse struct {
	ID               int       `json:"id"`
	User             User      `json:"user"`
	Room             Room      `json:"room"`
	Status           string    `json:"status"`
	Snack            Snack     `json:"snack"`
	Name             string    `json:"name"`
	Date             time.Time `json:"date"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	Phone            string    `json:"phone"`
	TotalParticipant int       `json:"total_participant"`
	Organization     string    `json:"organization"`
	Note             string    `json:"notes"`
	TotalDuration    int       `json:"total_duration"`
	GrandTotal       int       `json:"grand_total"`
	CreatedAt        time.Time `json:"created_at"`
}
type ReservationDetailService struct {
	ID               int         `json:"id"`
	Room             RoomDetail  `json:"room"`
	Status           string      `json:"status"`
	Snack            SnackDetail `json:"snack"`
	Name             string      `json:"name"`
	Date             CustomDate  `json:"date"`
	StartTime        CustomTime  `json:"start_time"`
	EndTime          CustomTime  `json:"end_time"`
	Phone            string      `json:"phone"`
	TotalParticipant int         `json:"total_participant"`
	Organization     string      `json:"organization"`
	Note             string      `json:"notes"`
	TotalDuration    int         `json:"total_duration"`
	GrandTotal       float64     `json:"grand_total"`
}

type ReservationDetailResponse struct {
	ID               int         `json:"id"`
	Room             RoomDetail  `json:"room"`
	Status           string      `json:"status"`
	Snack            SnackDetail `json:"snack"`
	Name             string      `json:"name"`
	Date             time.Time   `json:"date"`
	StartTime        time.Time   `json:"start_time"`
	EndTime          time.Time   `json:"end_time"`
	Phone            string      `json:"phone"`
	TotalParticipant int         `json:"total_participant"`
	Organization     string      `json:"organization"`
	Note             string      `json:"notes"`
	TotalDuration    int         `json:"total_duration"`
	GrandTotal       float64     `json:"grand_total"`
}

type ReservationHistoryResponse struct {
	ID       int       `json:"id"`
	Date     time.Time `json:"date"`
	RoomType string    `json:"room_type_name"`
	RoomName string    `json:"room_name"`
	Status   string    `json:"status"`
}

type ReservationListByRoomIdResponse struct {
	ID            int       `json:"id"`
	Status        string    `json:"status"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	TotalDuration int       `json:"total_duration"`
}

type ReservationUpdateRequest struct {
	Status string `json:"status" validate:"required"`
}

type ReservationGetAllResponse struct {
	ID           int       `json:"id"`
	Organization string    `json:"organization"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Date         time.Time `json:"date"`
	RoomID       int       `json:"room_id"`
	RoomName     string    `json:"room_name"`
}
