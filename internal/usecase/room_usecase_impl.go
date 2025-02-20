package usecase

import (
	"E-Meeting/internal/domain/entity"
	"E-Meeting/internal/domain/repository"
	"E-Meeting/pkg/reason"
	"E-Meeting/pkg/utils"
	"E-Meeting/presenter/model"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"sort"
	"sync"
	"time"

	"github.com/guregu/null"
)

type roomUseCase struct {
	repo         repository.RoomRepository
	repoRoomType repository.RoomTypeRepository
	repoCapacity repository.CapacityRepository
}

func NewRoomUseCase(repo repository.RoomRepository, repoRoomType repository.RoomTypeRepository, repoCapacity repository.CapacityRepository) RoomUseCase {
	return &roomUseCase{repo: repo, repoRoomType: repoRoomType, repoCapacity: repoCapacity}
}

func (u *roomUseCase) FindAllRoom(ctx context.Context, queryPageLimit utils.QueryPageLimit, filter *model.FilterDataRoomRequest) (*entity.RoomsDataAccessObject, error) {

	rooms, err := u.repo.FindALl(ctx, queryPageLimit, filter)
	if err != nil {
		log.Println(fmt.Sprintf("message : error in service | service : room_usecase_impl.FindAllRoom | error : %s", err))
		return nil, err
	}

	return rooms, nil

}

func (u *roomUseCase) GetDashboard(ctx context.Context) ([]*model.RoomDashboards, error) {

	rooms, err := u.repo.GetDashboard(ctx)
	if err != nil {
		log.Println(fmt.Sprintf("message : error in service | service : room_usecase_impl.GetDashboard | error : %s", err))
		return nil, err
	}

	return rooms, nil
}

func (u *roomUseCase) Insert(ctx context.Context, input model.RoomRequest, file *multipart.FileHeader) error {
	/*
	* CARA 1
	 */
	// // CHECK active room type
	// roomType, err := u.repoRoomType.GetOneByID(ctx, int64(input.RoomTypeID))
	// if err != nil {
	// 	log.Println(fmt.Sprintf("message : failed get data room_type | service : room_usecase_impl | error : %s", err))
	// 	return reason.ErrFailedInsertData
	// }
	//
	// if (!roomType.IsActive.Valid) || (!roomType.IsActive.Bool) {
	// 	log.Println(fmt.Sprintf("message : room type is not active | service : room_usecase_impl | validate : room_type_is_active"))
	// 	return reason.ErrFailedInsertData
	// }
	//
	// // CHECK active capacity
	// capacity, err := u.repoCapacity.GetOneByID(ctx, int64(input.CapacityID))
	// if err != nil {
	// 	log.Println(fmt.Sprintf("message : failed get data capacity | service : room_usecase_impl | error : %s", err))
	// 	return reason.ErrFailedInsertData
	// }
	//
	// if (!capacity.IsActive.Valid) || (!capacity.IsActive.Bool) {
	// 	log.Println(fmt.Sprintf("message : capacity is not active | service : room_usecase_impl | validate : capacity_is_active"))
	// 	return reason.ErrFailedInsertData
	// }

	/*
	* CARA CONCURRENT
	 */
	var wg sync.WaitGroup
	var roomType *entity.RoomType
	var capacities *entity.CapacityResultDataAccessObject
	var roomTypeErr error
	var capacityErr error

	// START GOROUTINE TO CHECK ROOM_TYPE
	wg.Add(1)
	go func() {
		defer wg.Done()
		roomType, roomTypeErr = u.repoRoomType.GetOneByID(ctx, int64(input.RoomTypeID))
	}()

	// START GOROUTINE TO CHECK CAPACITY
	wg.Add(1)
	go func() {
		defer wg.Done()
		queryPagination := utils.QueryPageLimit{
			Page:    0,
			Limit:   0,
			OrderBy: "id",
			SortBy:  "desc",
		}
		capacities, capacityErr = u.repoCapacity.FindAll(ctx, queryPagination)

	}()

	// wait goroutine
	wg.Wait()

	// Handle error in data room_type AND capacity
	if roomTypeErr != nil {
		log.Println(fmt.Sprintf("message : failed get data room_type | service : room_usecase_impl.Insert | error : %s", roomTypeErr))
		return reason.ErrFailedInsertData
	}

	if capacityErr != nil {
		log.Println(fmt.Sprintf("message : failed get data capacity | service : room_usecase_impl.Insert | error : %s", capacityErr))
		return reason.ErrFailedInsertData
	}

	// Validate roomType
	if (!roomType.IsActive.Valid) || (!roomType.IsActive.Bool) {
		log.Println(fmt.Sprintf("message : room type is not active | service : room_usecase_impl | validate : room_type_is_active"))
		return reason.ErrFailedInsertData
	}

	// // Validate capacity
	// if (!capacity.IsActive.Valid) || (!capacity.IsActive.Bool) {
	// 	log.Println(fmt.Sprintf("message : capacity is not active | service : room_usecase_impl | validate : capacity_is_active"))
	// 	return reason.ErrFailedInsertData
	// }

	var capacityID int64
	if len(capacities.Capacity) > 0 {

		capacityResults := capacities.Capacity

		sort.Slice(capacities.Capacity, func(i, j int) bool {
			// Ambil nilai value_minimum, default -1 jika null
			min1 := -1
			if capacities.Capacity[i].ValueMinimum.Valid {
				min1 = int(capacities.Capacity[i].ValueMinimum.Int64)
			}
			min2 := -1
			if capacities.Capacity[j].ValueMinimum.Valid {
				min2 = int(capacities.Capacity[j].ValueMinimum.Int64)
			}

			// Jika value_minimum berbeda, urutkan berdasarkan value_minimum
			if min1 != min2 {
				return min1 < min2
			}

			// Jika value_minimum sama atau keduanya null, urutkan berdasarkan value_maximum
			max1 := int(capacities.Capacity[i].ValueMaximum.Int64)
			max2 := int(capacities.Capacity[j].ValueMaximum.Int64)
			return max1 < max2
		})

		for _, capacity := range capacityResults {

			if !capacity.ValueMinimum.Valid {

				if capacity.ValueMaximum.Valid {

					if int64(input.Capacity) <= capacity.ValueMaximum.Int64 {
						capacityID = int64(capacity.ID)
						break
					}
				}
			}

			if (capacity.ValueMinimum.Valid) && (capacity.ValueMaximum.Valid) {

				if int64(input.Capacity) >= capacity.ValueMinimum.Int64 && (int64(input.Capacity) <= capacity.ValueMaximum.Int64) {
					capacityID = int64(capacity.ID)
					break
				}

			}

			if capacity.ValueMinimum.Valid {
				if !capacity.ValueMaximum.Valid {

					if int64(input.Capacity) >= capacity.ValueMinimum.Int64 {
						capacityID = int64(capacity.ID)
						break
					}
				}
			}

		}

	}

	isFile := false

	if file != nil {
		isFile = true
	}

	if isFile {

		fileURL, err := utils.SaveFile(file, "room")
		if err != nil {
			log.Println(fmt.Sprintf("message : capacity is not active | service : room_usecase_impl | validate : capacity_is_active"))
			return reason.ErrFailedInsertData
		}

		input.AttachmentURL = null.StringFrom(fileURL)

	}

	inputRepo := &entity.Room{
		Name:          null.StringFrom(input.Name),
		PriceHour:     null.FloatFrom(input.Price),
		IsActive:      null.BoolFrom(true),
		RoomTypeID:    null.IntFrom(int64(input.RoomTypeID)),
		CapacityID:    null.IntFrom(capacityID),
		Capacity:      null.IntFrom(int64(input.Capacity)),
		AttachmentURL: input.AttachmentURL,
		Description:   input.Description,
		CreatedAt:     null.TimeFrom(time.Now()),
	}

	lastInsertID, err := u.repo.Insert(ctx, inputRepo)
	if err != nil {
		log.Println(fmt.Sprintf("message : error in service | service : room_usecase_impl.Insert | error : %s", err))
		return reason.ErrFailedInsertData
	}

	if lastInsertID == 0 {
		log.Println(fmt.Sprintf("message : error in service | service : room_usecase_impl.Insert | method : Insert"))
		return reason.ErrFailedInsertData
	}

	return nil

}

func (u *roomUseCase) UpdateOneByID(ctx context.Context, input model.RoomRequest, file *multipart.FileHeader, roomID int64) error {
	_, err := u.repo.FindOneByID(ctx, roomID)
	if err != nil {
		log.Println(fmt.Sprintf("message : failed get data room by id | service : room_usecase_impl.UpdateOneByID | error : %s", err))
		return reason.ErrDataNotFound
	}

	/*
	* CARA CONCURRENT
	 */
	var wg sync.WaitGroup
	var roomType *entity.RoomType
	var capacities *entity.CapacityResultDataAccessObject
	var roomTypeErr error
	var capacityErr error

	// START GOROUTINE TO CHECK ROOM_TYPE
	wg.Add(1)
	go func() {
		defer wg.Done()
		roomType, roomTypeErr = u.repoRoomType.GetOneByID(ctx, int64(input.RoomTypeID))
	}()

	// START GOROUTINE TO CHECK CAPACITY
	wg.Add(1)
	go func() {
		defer wg.Done()
		queryPagination := utils.QueryPageLimit{
			Page:    0,
			Limit:   0,
			OrderBy: "id",
			SortBy:  "desc",
		}
		capacities, capacityErr = u.repoCapacity.FindAll(ctx, queryPagination)
	}()

	// wait goroutine
	wg.Wait()

	// Handle error in data room_type AND capacity
	if roomTypeErr != nil {
		log.Println(fmt.Sprintf("message : failed get data room_type | service : room_usecase_impl.UpdateOneByID | error : %s", roomTypeErr))
		return reason.ErrFailedInsertData
	}

	if capacityErr != nil {
		log.Println(fmt.Sprintf("message : failed get data capacity | service : room_usecase_impl.UpdateOneByID | error : %s", capacityErr))
		return reason.ErrFailedInsertData
	}

	// Validate roomType
	if (!roomType.IsActive.Valid) || (!roomType.IsActive.Bool) {
		log.Println(fmt.Sprintf("message : room type is not active | service : room_usecase_impl.UpdateOneByID | validate : room_type_is_active"))
		return reason.ErrFailedInsertData
	}

	// Validate capacity
	// if (!capacity.IsActive.Valid) || (!capacity.IsActive.Bool) {
	// 	log.Println(fmt.Sprintf("message : capacity is not active | service : room_usecase_impl | validate : capacity_is_active"))
	// 	return reason.ErrFailedInsertData
	// }
	var capacityID int64
	if len(capacities.Capacity) > 0 {

		capacityResults := capacities.Capacity

		sort.Slice(capacities.Capacity, func(i, j int) bool {
			// Ambil nilai value_minimum, default -1 jika null
			min1 := -1
			if capacities.Capacity[i].ValueMinimum.Valid {
				min1 = int(capacities.Capacity[i].ValueMinimum.Int64)
			}
			min2 := -1
			if capacities.Capacity[j].ValueMinimum.Valid {
				min2 = int(capacities.Capacity[j].ValueMinimum.Int64)
			}

			// Jika value_minimum berbeda, urutkan berdasarkan value_minimum
			if min1 != min2 {
				return min1 < min2
			}

			// Jika value_minimum sama atau keduanya null, urutkan berdasarkan value_maximum
			max1 := int(capacities.Capacity[i].ValueMaximum.Int64)
			max2 := int(capacities.Capacity[j].ValueMaximum.Int64)
			return max1 < max2
		})

		for _, capacity := range capacityResults {

			if !capacity.ValueMinimum.Valid {

				if capacity.ValueMaximum.Valid {

					if int64(input.Capacity) <= capacity.ValueMaximum.Int64 {
						capacityID = int64(capacity.ID)
						break
					}
				}
			}

			if (capacity.ValueMinimum.Valid) && (capacity.ValueMaximum.Valid) {

				if int64(input.Capacity) >= capacity.ValueMinimum.Int64 && (int64(input.Capacity) <= capacity.ValueMaximum.Int64) {
					capacityID = int64(capacity.ID)
					break
				}

			}

			if capacity.ValueMinimum.Valid {
				if !capacity.ValueMaximum.Valid {

					if int64(input.Capacity) >= capacity.ValueMinimum.Int64 {
						capacityID = int64(capacity.ID)
						break
					}
				}
			}

		}

	}

	isFile := false

	if file != nil {
		isFile = true
	}

	if isFile {

		fileURL, err := utils.SaveFile(file, "room")
		if err != nil {
			log.Println(fmt.Sprintf("message : capacity is not active | service : room_usecase_impl.UpdateOneByID | validate : capacity_is_active"))
			return reason.ErrFailedInsertData
		}

		input.AttachmentURL.String = fileURL

	}

	inputRepo := &entity.Room{
		Name:          null.StringFrom(input.Name),
		PriceHour:     null.FloatFrom(input.Price),
		IsActive:      null.BoolFrom(true),
		RoomTypeID:    null.IntFrom(int64(input.RoomTypeID)),
		CapacityID:    null.IntFrom(capacityID),
		Capacity:      null.IntFrom(int64(input.Capacity)),
		AttachmentURL: input.AttachmentURL,
		CreatedAt:     null.TimeFrom(time.Now()),
	}

	affectedRow, err := u.repo.Update(ctx, inputRepo, roomID)
	if err != nil {
		log.Println(fmt.Sprintf("message : error in service | service : room_usecase_impl | error : %s", err))
		return reason.ErrFailedInsertData
	}

	if affectedRow == 0 {
		log.Println(fmt.Sprintf("message : error in service | service : room_usecase_impl.UpdateOneByID | method : Insert"))
		return reason.ErrFailedInsertData
	}

	return nil
}

func (u *roomUseCase) DeleteOneByID(ctx context.Context, roomID int64) error {

	_, err := u.repo.FindOneByID(ctx, roomID)
	if err != nil {
		log.Println(fmt.Sprintf("message : failed get data room by id | service : room_usecase_impl.DeleteOneByID | error : %s", err))
		return reason.ErrDataNotFound
	}

	err = u.repo.DeleteByID(ctx, roomID)
	if err != nil {
		log.Println(fmt.Sprintf("message : failed delete data room by id | service : room_usecase_impl.DeleteOneByID | error : %s", err))
		return reason.ErrFailedDeleteData
	}
	return nil
}

func (u *roomUseCase) FindOneByID(ctx context.Context, roomID int64) (*entity.Room, error) {

	resultRoom, err := u.repo.FindOneByID(ctx, roomID)
	if err != nil {
		log.Println(fmt.Sprintf("message : failed get data room by id | service : room_usecase_impl.FindOneByID | error : %s", err))
		return nil, reason.ErrDataNotFound
	}

	return resultRoom, nil
}
