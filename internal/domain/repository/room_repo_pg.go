package repository

import (
	"E-Meeting/internal/domain/entity"
	reasonError "E-Meeting/pkg/reason"
	"E-Meeting/pkg/utils"
	"E-Meeting/presenter/model"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type roomRepo struct {
	DB *sqlx.DB
}

func NewRoomRepository(db *sqlx.DB) RoomRepository {
	return &roomRepo{
		DB: db,
	}
}

func (repo *roomRepo) FindALl(_ context.Context, queryPageLimit utils.QueryPageLimit, filter *model.FilterDataRoomRequest) (*entity.RoomsDataAccessObject, error) {

	queryCount := utils.QueryCount{
		TotalData:    0,
		TotalContent: 0,
	}

	roomsDataAccessObject := entity.RoomsDataAccessObject{
		Rooms:      nil,
		QueryCount: queryCount,
	}

	queryLimit := ""
	if (queryPageLimit.Page > 0) && (queryPageLimit.Limit > 0) {

		queryPageLimit.Page = (queryPageLimit.Page - 1) * queryPageLimit.Limit

		queryLimit = fmt.Sprintf("LIMIT %d OFFSET %d", queryPageLimit.Limit, queryPageLimit.Page)
	}

	queryOrderSort := fmt.Sprintf("ORDER BY %s %s", queryPageLimit.OrderBy, queryPageLimit.SortBy)

	whereByRoomType := "1 = 1"
	whereByCapacity := "1 = 1"
	if filter.RoomType != nil {
		whereByRoomType = fmt.Sprintf("room_type_id = %d", *filter.RoomType)
	}

	if filter.Capacity != nil {
		whereByCapacity = fmt.Sprintf("capacity_id = %d", *filter.Capacity)
	}

	var rooms []entity.Room
	query := fmt.Sprintf("SELECT "+
		"r.id, r.name, r.is_active, r.description, r.price_hour, r.room_type_id, r.capacity_id,r.capacity, r.attachment_url,"+
		"c.value_minimum as capacity_value_minimum, c.value_maximum as capacity_value_maximum,c.uom as capacity_uom, c.is_active as capacity_is_active,"+
		"rt.name as room_type_name, rt.is_active as room_type_is_active, "+
		"r.created_at, r.updated_at, r.created_by, r.updated_by "+
		"FROM public.%s %s "+
		"LEFT JOIN public.capacities c on c.id = r.capacity_id "+
		"LEFT JOIN public.room_types rt on rt.id = r.room_type_id WHERE %s AND %s  %s %s", entity.TableRoomName, entity.TableRoomAliasName, whereByRoomType, whereByCapacity, queryOrderSort, queryLimit)

	err := repo.DB.Select(&rooms, query) // Mengisi slice `users` dengan hasil query
	if err != nil {
		return nil, err
	}

	if len(rooms) < 1 {
		return nil, reasonError.ErrDataNotFound
	}

	// Query for total count
	sqlCount := fmt.Sprintf("SELECT COUNT(*) FROM public.%s WHERE %s AND %s", entity.TableRoomName, whereByRoomType, whereByCapacity)

	var totalCount int
	err = repo.DB.Get(&totalCount, sqlCount)

	if err != nil {
		log.Println(fmt.Sprintf("message : error in query count | repo : room_repo_pg.FindALl | error : %s", err))
		return nil, err
	}
	roomsDataAccessObject.Rooms = rooms
	roomsDataAccessObject.QueryCount.TotalData = totalCount
	roomsDataAccessObject.QueryCount.TotalContent = len(rooms)

	return &roomsDataAccessObject, nil
}

func (repo *roomRepo) Insert(ctx context.Context, input *entity.Room) (lastInsertID int, err error) {

	query := `INSERT INTO rooms (name, room_type_id, price_hour, capacity_id,capacity, attachment_url,description, created_by)
          VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	err = repo.DB.QueryRow(query,
		input.Name,
		input.RoomTypeID,
		input.PriceHour,
		input.CapacityID,
		input.Capacity,
		input.AttachmentURL,
		input.Description,
		ctx.Value("email"),
	).Scan(&lastInsertID)
	if err != nil {
		log.Println(fmt.Sprintf("message : error in query insert data | repo : room_repo_pg | error : %s", err))
		return 0, err
	}
	return lastInsertID, nil
}

func (repo *roomRepo) GetDashboard(_ context.Context) (rooms []*model.RoomDashboards, err error) {

	query := `SELECT r.id as room_id, r.name as room_name, r.price_hour as price_hour, 
			COALESCE(
			(CAST(COUNT(CASE WHEN res.status = 'paid' THEN 1 END) AS FLOAT) /
			NULLIF(COUNT(CASE WHEN res.status IN ('paid', 'booked') THEN 1 END), 0)) * 100, 0) as percentage
			FROM rooms r LEFT JOIN reservation_rooms res ON r.id = res.room_id GROUP BY r.id, r.name, r.price_hour;`

	err = repo.DB.Select(&rooms, query)
	if err != nil {
		log.Println(fmt.Sprintf("message : error in query get dashboard | repo : room_repo_pg | error : %s", err))
		return nil, err
	}

	return rooms, nil
}

func (repo *roomRepo) Update(ctx context.Context, input *entity.Room, roomID int64) (rowAffected int, err error) {

	query := `UPDATE rooms SET name=$1, room_type_id=$2, price_hour=$3, capacity_id=$4, attachment_url=$5, updated_by=$6, updated_at=NOW() WHERE id=$7`

	_, err = repo.DB.Exec(query,
		input.Name,
		input.RoomTypeID,
		input.PriceHour,
		input.CapacityID,
		input.AttachmentURL,
		ctx.Value("email"),
		roomID,
	)
	if err != nil {
		log.Printf("message: error in query update data | repo: room_repo_pg | error: %s", err)
		return 0, err
	}

	return 1, nil
}

func (repo *roomRepo) DeleteByID(_ context.Context, roomID int64) (err error) {

	query := `DELETE FROM rooms WHERE id = $1`

	err = repo.DB.QueryRow(query, roomID).Err()
	if err != nil {
		log.Println(fmt.Sprintf("message : error in query insert data | repo : room_repo_pg | error : %s", err))
		return err
	}
	return nil
}

func (repo *roomRepo) FindOneByID(ctx context.Context, roomID int64) (*entity.Room, error) {

	var room entity.Room

	query := fmt.Sprintf(`
			SELECT 
			    r.id, r.name, r.is_active, r.description, r.price_hour, r.room_type_id, r.capacity_id, r.attachment_url,
			    rt.name as room_type_name,
			    COALESCE(
					JSON_AGG(
						JSON_BUILD_OBJECT(
							'reservation_room_id', rr.id,
							'date', rr.date,
							'start_time', rr.start_time,
							'end_time', rr.end_time
						)
					) FILTER (WHERE rr.id IS NOT NULL), '[]'
				) AS reservations
			FROM  public.rooms r 
			    LEFT JOIN public.room_types rt ON rt.id = r.room_type_id
				LEFT JOIN public.reservation_rooms rr ON rr.room_id = r.id
			WHERE 
			    r.id = $1
			GROUP BY 
			r.id, 
			--    r.name, 
			--    r.is_active, 
			--    r.description, 
			--    r.price_hour, 
			--    r.room_type_id, 
			--    r.capacity_id, 
			--    r.attachment_url, 
			rt.id
			`)

	err := repo.DB.GetContext(ctx, &room, query, roomID)
	if err != nil {
		return nil, err
	}

	// Dekode Reservations
	var reservations *[]entity.Reservation
	if err := json.Unmarshal(room.ReservationsJson, &reservations); err != nil {
		return nil, fmt.Errorf("failed to decode reservations JSON: %w", err)
	}

	if len(*reservations) > 0 {

		room.Reservations = *reservations

	}

	if room.ID.Int64 == 0 {
		return nil, reasonError.ErrDataNotFound
	}

	return &room, nil
}
