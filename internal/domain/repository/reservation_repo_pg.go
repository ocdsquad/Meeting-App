package repository

import (
	"E-Meeting/internal/domain/entity"
	"E-Meeting/pkg/helper"
	"E-Meeting/pkg/reason"
	"E-Meeting/presenter/model"
	"context"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type ReservationRepo struct {
	DB *sqlx.DB
}

func NewReservationRepository(db *sqlx.DB) ReservationRepository {
	return &ReservationRepo{
		DB: db,
	}
}

func (r *ReservationRepo) Save(ctx context.Context, reservation *entity.ReservationRooms) error {
	query := `INSERT INTO reservation_rooms (room_id, user_id,date, start_time, end_time, status, category_snack_id, Name, phone, total_participant, organization, notes, total_duration, grand_total, created_at, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16 ) RETURNING id`
	log.Println(ctx.Value("email"))
	err := r.DB.QueryRowContext(ctx, query, reservation.RoomID, reservation.UserID, reservation.Date, reservation.StartTime, reservation.EndTime, reservation.Status, reservation.SnackID, reservation.Name, reservation.Phone, reservation.TotalParticipant, reservation.Organization, reservation.Note, reservation.TotalDuration, reservation.GrandTotal, reservation.CreatedAt, ctx.Value("email")).Scan(&reservation.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReservationRepo) GetDetailReservation(ctx context.Context, reservationID int) (*model.ReservationDetailService, error) {
	query := `
	SELECT r.id AS reservation_id,
    r.status,
    r.date,
    r.start_time,
    r.end_time,
    r.total_participant,
    r.notes,
	r.phone,
	r.organization,
	r.total_duration,
    r.grand_total,
    r.name,
    ro.id AS room_id,
    ro.name AS room_name,
    ro.price_hour AS room_price,
	ro.capacity AS room_capacity,
	rt.name AS room_type_name,
    s.id AS snack_id,
    s.name AS snack_name,
    s.price AS snack_price
	FROM reservation_rooms r
	LEFT JOIN rooms ro ON r.room_id = ro.id
	LEFT JOIN users u ON r.user_id = u.id
	LEFT JOIN category_snacks s ON r.category_snack_id = s.id
	LEFT JOIN room_types rt ON ro.room_type_id = rt.id
	WHERE r.id = $1;
	`
	var reservation model.ReservationDetailService
	row := r.DB.QueryRowContext(ctx, query, reservationID)
	err := row.Scan(&reservation.ID, &reservation.Status, &reservation.Date,
		&reservation.StartTime, &reservation.EndTime,
		&reservation.TotalParticipant, &reservation.Note, &reservation.Phone, &reservation.Organization, &reservation.TotalDuration, &reservation.GrandTotal,
		&reservation.Name,
		&reservation.Room.ID, &reservation.Room.Name, &reservation.Room.PriceHour, &reservation.Room.Capacity, &reservation.Room.RoomTypeName,
		&reservation.Snack.ID, &reservation.Snack.Name, &reservation.Snack.Price)
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *ReservationRepo) GetAllReservation(ctx context.Context, startDate time.Time, endDate time.Time) ([]*model.ReservationGetAllResponse, error) {
	query := `
	SELECT r.id AS reservation_id,
	r.date,
	r.start_time,
	r.end_time,
	r.organization,
	ro.id AS room_id,
	ro.name AS room_name
	FROM reservation_rooms r
	LEFT JOIN rooms ro ON r.room_id = ro.id
	WHERE r.status = 'paid'
	`
	var args []interface{}
	query, args = helper.VerifyDateFilter(startDate, endDate, query, args...)
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []*model.ReservationGetAllResponse
	for rows.Next() {
		var reservation model.ReservationGetAllResponse
		err = rows.Scan(&reservation.ID, &reservation.Date, &reservation.StartTime,
			&reservation.EndTime, &reservation.Organization, &reservation.RoomID,
			&reservation.RoomName)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, &reservation)
	}
	return reservations, nil
}

func (r *ReservationRepo) GetHistoryReservationByUserID(ctx context.Context, userID int) ([]*model.ReservationHistoryResponse, error) {
	query := `
	SELECT r.id AS reservation_id,
	r.status,
	r.date,
	ro.name AS room_name,
	rt.name AS room_type_name
	FROM reservation_rooms r
	LEFT JOIN rooms ro ON r.room_id = ro.id
	LEFT JOIN room_types rt ON ro.room_type_id = rt.id
	WHERE r.user_id = $1;
	`
	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []*model.ReservationHistoryResponse
	for rows.Next() {
		var reservation model.ReservationHistoryResponse
		err = rows.Scan(&reservation.ID, &reservation.Status, &reservation.Date, &reservation.RoomName, &reservation.RoomType)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, &reservation)
	}
	return reservations, nil
}

func (r *ReservationRepo) UpdateStatusReservation(ctx context.Context, reservationID int, status string) error {
	query := `UPDATE reservation_rooms SET status = $1 WHERE id = $2`
	_, err := r.DB.ExecContext(ctx, query, status, reservationID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReservationRepo) GetListReservationByRoomID(ctx context.Context, roomID int, startDateTime time.Time, endDateTime time.Time) ([]*model.ReservationListByRoomIdResponse, error) {
	query := `
	SELECT 
	id,
	status,
	start_time,
	end_time,
	total_duration
	FROM reservation_rooms
	WHERE room_id = $1`

	log.Println(startDateTime, endDateTime)

	var args []interface{}
	args = append(args, roomID)
	query, args = helper.VerifyDateFilter(startDateTime, endDateTime, query, args...)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var response []*model.ReservationListByRoomIdResponse
	for rows.Next() {
		var reservation model.ReservationListByRoomIdResponse
		err = rows.Scan(&reservation.ID, &reservation.Status, &reservation.StartTime, &reservation.EndTime, &reservation.TotalDuration)
		if err != nil {
			return nil, err
		}
		response = append(response, &reservation)
	}

	return response, nil
}

func (r *ReservationRepo) GetHistoryReservationByIsAdmin(ctx context.Context) ([]*model.ReservationHistoryResponse, error) {
	query := `
	SELECT r.id AS reservation_id,
	r.status,
	r.date,
	ro.name AS room_name,
	rt.name AS room_type_name
	FROM reservation_rooms r
	LEFT JOIN rooms ro ON r.room_id = ro.id
	LEFT JOIN room_types rt ON ro.room_type_id = rt.id;
	`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []*model.ReservationHistoryResponse
	for rows.Next() {
		var reservation model.ReservationHistoryResponse
		err = rows.Scan(&reservation.ID, &reservation.Status, &reservation.Date, &reservation.RoomName, &reservation.RoomType)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, &reservation)
	}
	return reservations, nil
}

func (r *ReservationRepo) GetDashboard(_ context.Context, startDateTime time.Time, endDateTime time.Time) (*model.DashboardResponse, error) {
	var dashboard model.DashboardResponse

	query := `
	SELECT COUNT(r.id) AS total_reservation,  COALESCE(SUM(CAST(grand_total AS NUMERIC)), 0) AS total_omset, COUNT(ro.id) as total_room, COALESCE(SUM(CAST(total_participant AS NUMERIC)), 0) AS total_visitor FROM reservation_rooms r LEFT JOIN rooms ro ON r.room_id = ro.id WHERE r.status = 'paid'
	`
	var args []interface{}
	query, args = helper.VerifyDateFilter(startDateTime, endDateTime, query, args...)

	err := r.DB.QueryRow(query, args...).Scan(&dashboard.TotalReservation, &dashboard.TotalOmset, &dashboard.TotalRoom, &dashboard.TotalVisitor)

	if err != nil {
		return nil, err
	}
	if dashboard.TotalReservation == 0 {
		return nil, reason.ErrDataNotFound
	}

	return &dashboard, nil
}

func (r *ReservationRepo) DeleteReservation(ctx context.Context, reservationID int) error {
	query := `DELETE FROM reservation_rooms WHERE id = $1`
	_, err := r.DB.ExecContext(ctx, query, reservationID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReservationRepo) GetAllReservationByRoomIDAndDate(ctx context.Context, roomID int64, date string, startTime string, endTime string) ([]entity.ReservationRooms, error) {

	var reservations []entity.ReservationRooms

	query := `
	SELECT 
	    r.id,
		r.date,
		r.start_time,
		r.end_time,
		r.organization
	FROM 
	    reservation_rooms r
	WHERE 
	    r.room_id = $1 
	AND 
	    r.date = $2
    AND 
	    r.deleted_at IS NULL
    AND (
        -- Cek jika waktu baru overlap dengan reservasi yang sudah ada
        (r.start_time <= $4 AND r.end_time > $3) -- $3 = start_time baru, $4 = end_time baru
        OR 
        (r.start_time < $4 AND r.end_time >= $4)
        OR
        (r.start_time >= $3 AND r.end_time <= $4)
    );;
	`

	err := r.DB.SelectContext(ctx, &reservations, query, roomID, date, startTime, endTime)
	if err != nil {
		return nil, err
	}

	return reservations, nil
}
