package usecase

import (
	"E-Meeting/internal/domain/entity"
	"E-Meeting/internal/domain/repository"
	"E-Meeting/pkg/cache"
	"E-Meeting/pkg/helper"
	"E-Meeting/pkg/reason"
	"E-Meeting/presenter/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/guregu/null"

	"github.com/go-playground/validator/v10"
)

type reservationUseCase struct {
	repo      repository.ReservationRepository
	repoRoom  repository.RoomRepository
	repoSnack repository.SnackRepository
	validator *validator.Validate
}

func NewReservationUseCase(repo repository.ReservationRepository, repoRoom repository.RoomRepository, repoSnack repository.SnackRepository) ReservationUseCase {
	return &reservationUseCase{
		repo:      repo,
		repoRoom:  repoRoom,
		repoSnack: repoSnack,
		validator: validator.New(),
	}
}

func (u *reservationUseCase) Save(ctx context.Context, request *model.ReservationCodeRequest) error {
	err := u.validator.Struct(request)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("reservation_%s", strings.TrimSpace(request.Code))
	result, ok := cache.MyCache.Get(cacheKey)
	if (!ok) || (result == nil) {
		log.Println(fmt.Sprintf("message : failed insertt data because data order not found | service : reservation_usecase_impl.save | method : cache.MyCache.Get"))
		return reason.ErrFailedInsertData
	}

	var entityReservation entity.ReservationRooms

	dataByte, err := json.Marshal(result)
	if err != nil {
		log.Println(fmt.Sprintf("message : failed marshal json | service : reservation_usecase_impl.save | method : json.Marshal"))
		return reason.ErrFailedInsertData
	}

	err = json.Unmarshal(dataByte, &entityReservation)
	if err != nil {
		log.Println(fmt.Sprintf("message : failed unmarshal json | service : reservation_usecase_impl.save | method : json.Unmarshal(dataByte, &entityReservation)"))
		return reason.ErrFailedInsertData
	}

	err = u.repo.Save(ctx, &entityReservation)
	if err != nil {
		return err
	}

	return nil
}

func (u *reservationUseCase) Inquiry(ctx context.Context, request *model.ReservationCreateRequest, userID int64) (*model.ReservationCreateServiceResponse, error) {
	err := u.validator.Struct(request)
	if err != nil {
		return nil, err
	}

	var snack *entity.Snack
	if request.SnackID.Valid {
		// 	check snack ID
		snackResult, err := u.repoSnack.GetOneByID(ctx, request.SnackID.Int64)
		if err != nil {
			log.Println(fmt.Sprintf("message : category snack is not found | service : reservation_usecase_impl.inqueiry | method : u.repoSnack.GetOneByID"))
			return nil, reason.ErrFailedInsertData
		}

		// Validate roomType
		if (!snackResult.IsActive.Valid) || (!snackResult.IsActive.Bool) {
			log.Println(fmt.Sprintf("message : snack is not active | service : reservation_usecase_impl.inqueiry | validate : category_snack_is_active"))
			return nil, reason.ErrFailedInsertData
		}

		snack = snackResult

	}

	var wg sync.WaitGroup
	var room *entity.Room
	var roomErr error

	var reservations []entity.ReservationRooms
	var reservationErr error

	// START GOROUTINE TO CHECK ROOM_TYPE
	wg.Add(1)
	go func() {
		defer wg.Done()
		room, roomErr = u.repoRoom.FindOneByID(ctx, int64(request.RoomID))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		reservations, roomErr = u.repo.GetAllReservationByRoomIDAndDate(ctx, int64(request.RoomID), request.Date.String(), request.StartTime.String(), request.EndTime.String())
	}()

	// wait goroutine
	wg.Wait()

	// Handle error in data room_type AND capacity
	if roomErr != nil {
		log.Println(fmt.Sprintf("message : failed get data room | service : reservation_usecase_impl.inqueiry | error : %s", roomErr))
		return nil, reason.ErrFailedInsertData
	}

	if (!room.IsActive.Valid) || (!room.IsActive.Bool) {
		log.Println(fmt.Sprintf("message : room is not active | service : reservation_usecase_impl.inqueiry | validate : room_is_active"))
		return nil, errors.New("room is not active")
	}

	if reservationErr != nil {

		if reservationErr.Error() == "sql: no rows in result set" {
			log.Println(fmt.Sprintf("message : failed get data reservation room by room_id | service : reservation_usecase_impl.inqueiry | error : %s", reservationErr))
			return nil, reason.ErrFailedInsertData
		}

	}

	if reservations != nil && len(reservations) >= 0 {
		log.Println(fmt.Sprintf("message : date time is schduled | service : reservation_usecase_impl.inqueiry | validate : date time is scheduled"))
		return nil, reason.ErrDateIsScheduled

	}

	var grandTotalPrice float64
	var totalPriceReservationDuration float64
	totalDuration := 0
	var totalPriceSnack float64

	// parsing time
	start, err := time.Parse(model.TimeFormat, request.StartTime.String())
	if err != nil {
		log.Println(fmt.Sprintf("message : Error parsing start time | service : reservation_usecase_impl.inqueiry | error : %s", err))
		return nil, reason.ErrInvalidFormatTime
	}

	end, err := time.Parse(model.TimeFormat, request.EndTime.String())
	if err != nil {
		log.Println(fmt.Sprintf("message : Error parsing end time | service : reservation_usecase_impl.inqueiry | error : %s", err))
		return nil, reason.ErrInvalidFormatTime
	}

	// calculate duration
	totalDuration = int(end.Sub(start).Hours())

	if totalDuration < 0 {
		log.Println(fmt.Sprintf("message : Error duration must be > 0 | service : reservation_usecase_impl.inqueiry | validate : total_duration "))
		return nil, reason.ErrDurationMinus
	}

	// calculate price for duration
	totalPriceReservationDuration = room.PriceHour.Float64 * float64(totalDuration)

	// calculate total price snack
	if (snack != nil) && (snack.ID > 0) {
		totalPriceSnack = snack.Price.Float64 * float64(request.TotalParticipant)
	}

	grandTotalPrice = totalPriceReservationDuration + totalPriceSnack

	// generate random data
	randomString, err := helper.GenerateRandomID(12)
	if err != nil {
		log.Println(fmt.Sprintf("message : failed generate random code | service : reservation_usecase_impl.inqueiry | validate : total_duration "))
		return nil, reason.ErrFailedInsertData
	}

	// Get current Unix timestamp in milliseconds
	currentTime := time.Now()
	unixTimestampMilliseconds := currentTime.UnixNano() / 1e6 // Convert to milliseconds

	// Combine the random string and Unix timestamp
	finalString := fmt.Sprintf("%s-%d", randomString, unixTimestampMilliseconds)

	data := model.ReservationCreateServiceResponse{
		Room:                          room,
		Snack:                         snack,
		RequestInput:                  *request,
		GrandTotalPrice:               grandTotalPrice,
		TotalPriceReservationDuration: totalPriceReservationDuration,
		TotalPriceSnack:               totalPriceSnack,
		Duration:                      int64(totalDuration),
		ReservationCode:               finalString,
	}

	date, err := time.Parse("2006-01-02", request.Date.String())
	if err != nil {

		log.Printf("message : Failed to parse date | service : reservation_usecase_impl.inqueiry | err : %s ", err)
		return nil, reason.ErrFailedInsertData
	}
	startTime, err := time.Parse("15:04:05", request.StartTime.String())
	if err != nil {

		log.Println(fmt.Sprintf("message : Failed to parse start time | service : reservation_usecase_impl.inqueiry | err : %s ", err))
		return nil, reason.ErrFailedInsertData
	}

	endTime, err := time.Parse("15:04:05", request.EndTime.String())
	if err != nil {

		log.Println(fmt.Sprintf("message : Failed to parse end time | service : reservation_usecase_impl.inqueiry | err : %s ", err))
		return nil, reason.ErrFailedInsertData

	}

	dataEntity := entity.ReservationRooms{
		UserID:           userID,
		RoomID:           int64(request.RoomID),
		Status:           "booked",
		SnackID:          request.SnackID,
		Name:             request.Name,
		Date:             date,
		StartTime:        startTime,
		EndTime:          endTime,
		Phone:            request.Phone,
		TotalParticipant: int64(request.TotalParticipant),
		Organization:     request.Organization,
		Note:             request.Note,
		TotalDuration:    int64(totalDuration),
		GrandTotal:       grandTotalPrice,
		CreatedAt:        null.TimeFrom(time.Now()),
	}

	// add to cache
	cacheKey := fmt.Sprintf("reservation_%s", finalString)
	log.Printf("cache key: %s\n", cacheKey)
	cache.MyCache.Set(cacheKey, dataEntity, 5*time.Minute)

	return &data, nil
}

// func (u *reservationUseCase) GetHistoryReservationByUserID(context context.Context, userID int) ([]model.ReservationHistoryResponse, error) {
// 	reservations, err := u.repo.GetHistoryReservationByUserID(context, userID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var response []model.ReservationHistoryResponse
// 	for _, reservation := range reservations {
// 		response = append(response, model.ReservationHistoryResponse{
// 			ID:       reservation.ID,
// 			Status:   reservation.Status,
// 			Date:     reservation.Date,
// 			RoomType: reservation.RoomType,
// 			RoomName: reservation.RoomName,
// 		})
// 	}
// 	return response, nil
// }

func (u *reservationUseCase) GetHistoryReservation(context context.Context, userID int, isAdmin bool) ([]model.ReservationHistoryResponse, error) {
	var reservations []*model.ReservationHistoryResponse
	var err error
	if !isAdmin {
		reservations, err = u.repo.GetHistoryReservationByUserID(context, userID)
		if err != nil {
			return nil, err
		}
	} else {
		reservations, err = u.repo.GetHistoryReservationByIsAdmin(context)
		if err != nil {
			return nil, err
		}

	}

	var response []model.ReservationHistoryResponse
	for _, reservation := range reservations {
		response = append(response, model.ReservationHistoryResponse{
			ID:       reservation.ID,
			Status:   reservation.Status,
			Date:     reservation.Date,
			RoomType: reservation.RoomType,
			RoomName: reservation.RoomName,
		})
	}
	return response, nil
}

func (u *reservationUseCase) GetDetailReservation(ctx context.Context, reservationID int) (*model.ReservationDetailResponse, error) {
	reservation, err := u.repo.GetDetailReservation(ctx, reservationID)
	if err != nil {
		return nil, err
	}
	log.Println(reservationID)

	var response *model.ReservationDetailResponse
	res := reservation
	date, err := time.Parse("2006-01-02", res.Date.String())
	if err != nil {

		log.Printf("message : Failed to parse date | service : reservation_usecase_impl.inqueiry | err : %s ", err)
		return nil, err
	}
	startTime, err := time.Parse("15:04:05", res.StartTime.String())
	if err != nil {
		fmt.Println("Error parsing start_time:", err)
		return nil, err
	}

	endTime, err := time.Parse("15:04:05", res.EndTime.String())
	if err != nil {
		fmt.Println("Error parsing end_time:", err)
		return nil, err
	}
	log.Println(date)
	log.Println(startTime)
	log.Println(endTime)
	response = &model.ReservationDetailResponse{
		ID:               res.ID,
		Status:           res.Status,
		Date:             date,
		StartTime:        startTime,
		EndTime:          endTime,
		TotalParticipant: res.TotalParticipant,
		Note:             res.Note,
		GrandTotal:       res.GrandTotal,
		Name:             res.Name,
		Phone:            res.Phone,
		Organization:     res.Organization,
		TotalDuration:    res.TotalDuration,
		Room: model.RoomDetail{
			ID:           res.Room.ID,
			Name:         res.Room.Name,
			PriceHour:    res.Room.PriceHour,
			Capacity:     res.Room.Capacity,
			RoomTypeName: res.Room.RoomTypeName,
		},
		Snack: model.SnackDetail{
			ID:    res.Snack.ID,
			Name:  null.StringFrom(res.Snack.Name.String),
			Price: null.FloatFrom(res.Snack.Price.Float64),
		},
	}
	return response, nil
}

func (u *reservationUseCase) UpdateStatusReservation(ctx context.Context, reservationID int, status string) error {
	err := u.repo.UpdateStatusReservation(ctx, reservationID, status)
	if err != nil {
		return err
	}
	return nil
}

func (u *reservationUseCase) GetListReservationByRoomID(ctx context.Context, roomID int, startDateTime time.Time, endDateTime time.Time) ([]model.ReservationListByRoomIdResponse, error) {

	reservations, err := u.repo.GetListReservationByRoomID(ctx, roomID, startDateTime, endDateTime)
	if err != nil {
		return nil, err
	}
	log.Println(reservations)
	var response []model.ReservationListByRoomIdResponse
	for _, reservation := range reservations {
		response = append(response, model.ReservationListByRoomIdResponse{
			ID:            reservation.ID,
			Status:        reservation.Status,
			StartTime:     reservation.StartTime,
			EndTime:       reservation.EndTime,
			TotalDuration: reservation.TotalDuration,
		})
	}
	return response, nil
}

func (u *reservationUseCase) GetDashboard(ctx context.Context, startDateTime time.Time, endDateTime time.Time) (*model.DashboardResponse, error) {
	var wg sync.WaitGroup
	var rooms []*model.RoomDashboards
	var roomErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		rooms, roomErr = u.repoRoom.GetDashboard(ctx)
	}()
	if roomErr != nil {
		log.Println(fmt.Sprintf("message : failed get data room | service : reservation_usecase_impl.inqueiry | error : %s", roomErr))
		return nil, roomErr
	}
	// wait goroutine
	wg.Wait()

	dashboard, err := u.repo.GetDashboard(ctx, startDateTime, endDateTime)
	if err != nil {
		return nil, err
	}
	var convertedRooms []struct {
		ID         int     `json:"id"`
		Name       string  `json:"name"`
		Percentage float64 `json:"percentage"`
		PriceHour  float64 `json:"price_hour"`
	}
	for _, room := range rooms {
		convertedRooms = append(convertedRooms, struct {
			ID         int     `json:"id"`
			Name       string  `json:"name"`
			Percentage float64 `json:"percentage"`
			PriceHour  float64 `json:"price_hour"`
		}{
			ID:         room.ID,
			Name:       room.Name,
			Percentage: room.Percentage,
			PriceHour:  room.PriceHour,
		})
	}

	responses := model.DashboardResponse{
		TotalRoom:        dashboard.TotalRoom,
		TotalReservation: dashboard.TotalReservation,
		TotalVisitor:     dashboard.TotalVisitor,
		TotalOmset:       dashboard.TotalOmset,
		Rooms:            convertedRooms,
	}

	return &responses, nil
}

func (u *reservationUseCase) GetAllReservation(ctx context.Context, startTime time.Time, endTime time.Time) ([]model.ReservationGetAllResponse, error) {
	reservations, err := u.repo.GetAllReservation(ctx, startTime, endTime)
	if err != nil {
		return nil, err
	}

	var response []model.ReservationGetAllResponse
	for _, reservation := range reservations {
		response = append(response, model.ReservationGetAllResponse{
			ID:           reservation.ID,
			Organization: reservation.Organization,
			Date:         reservation.Date,
			StartTime:    reservation.StartTime,
			EndTime:      reservation.EndTime,
			RoomID:       reservation.RoomID,
			RoomName:     reservation.RoomName,
		})
	}
	return response, nil

}

func (u *reservationUseCase) DeleteReservation(ctx context.Context, reservationID int) error {
	err := u.repo.DeleteReservation(ctx, reservationID)
	if err != nil {
		return err
	}
	return nil
}
