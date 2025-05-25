package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"minityweb/backend/pkg/models"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type PGStore struct {
	DB *sql.DB
}

func NewPGStore(db *sql.DB) *PGStore {
	return &PGStore{DB: db}
}

// --- User Methods (без изменений) ---
func (s *PGStore) CreateUser(ctx context.Context, user *models.User) (int64, error) {
	query := `INSERT INTO users (name, phone_number, email, password_hash, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	var userID int64; now := time.Now()
	err := s.DB.QueryRowContext(ctx, query, user.Name, user.PhoneNumber, user.Email, user.PasswordHash, user.Role, now, now).Scan(&userID)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code.Name() == "unique_violation" { return 0, errors.New("user with this phone number or email already exists") }
		log.Printf("Error creating user in DB: %v", err); return 0, fmt.Errorf("could not create user: %w", err)
	}
	return userID, nil
}
func (s *PGStore) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error) {
	query := `SELECT id, name, phone_number, email, password_hash, role, created_at, updated_at FROM users WHERE phone_number = $1`
	user := &models.User{}; err := s.DB.QueryRowContext(ctx, query, phoneNumber).Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil { if errors.Is(err, sql.ErrNoRows) { return nil, nil }; log.Printf("Error getting user by phone: %v", err); return nil, fmt.Errorf("could not get user by phone: %w", err) }
	return user, nil
}
func (s *PGStore) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	query := `SELECT id, name, phone_number, email, password_hash, role, created_at, updated_at FROM users WHERE id = $1`
	user := &models.User{}; err := s.DB.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil { if errors.Is(err, sql.ErrNoRows) { return nil, nil }; log.Printf("Error getting user by ID: %v", err); return nil, fmt.Errorf("could not get user by ID: %w", err) }
	return user, nil
}

// --- Establishment Methods (без изменений) ---
func (s *PGStore) GetEstablishments(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*models.Establishment, error) {
	baseQuery := `SELECT id, owner_user_id, name, type, address, working_hours, description, photos, created_at, updated_at FROM establishments`
	var conditions []string; var args []interface{}; argID := 1
	if estType, ok := filters["type"].(models.EstablishmentType); ok && estType != "" { conditions = append(conditions, fmt.Sprintf("type = $%d", argID)); args = append(args, estType); argID++ }
	query := baseQuery
	if len(conditions) > 0 { query += " WHERE " + strings.Join(conditions, " AND ") }
	query += fmt.Sprintf(" ORDER BY name ASC LIMIT $%d OFFSET $%d", argID, argID+1); args = append(args, limit, offset)
	rows, err := s.DB.QueryContext(ctx, query, args...); if err != nil { log.Printf("Error querying establishments: %v", err); return nil, fmt.Errorf("could not query establishments: %w", err) }; defer rows.Close()
	var establishments []*models.Establishment
	for rows.Next() {
		est := &models.Establishment{}; err := rows.Scan(&est.ID, &est.OwnerUserID, &est.Name, &est.Type, &est.Address, &est.WorkingHours, &est.Description, pq.Array(&est.Photos), &est.CreatedAt, &est.UpdatedAt)
		if err != nil { log.Printf("Error scanning est row: %v", err); return nil, fmt.Errorf("could not scan est: %w", err) }
		establishments = append(establishments, est)
	}
	return establishments, rows.Err()
}
func (s *PGStore) GetEstablishmentByID(ctx context.Context, id int64) (*models.Establishment, error) {
    query := `SELECT id, owner_user_id, name, type, address, working_hours, description, photos, created_at, updated_at FROM establishments WHERE id = $1`
    est := &models.Establishment{}; err := s.DB.QueryRowContext(ctx, query, id).Scan(&est.ID, &est.OwnerUserID, &est.Name, &est.Type, &est.Address, &est.WorkingHours, &est.Description, pq.Array(&est.Photos), &est.CreatedAt, &est.UpdatedAt)
    if err != nil { if errors.Is(err, sql.ErrNoRows) { return nil, nil }; return nil, fmt.Errorf("could not get est by ID: %w", err) }
    return est, nil
}
func (s *PGStore) getDetailsForEstablishment(ctx context.Context, est *models.Establishment) (*models.EstablishmentWithDetails, error) {
	details := &models.EstablishmentWithDetails{Establishment: *est, Halls: []models.HallWithPlaces{}, Menu: []models.MenuCategory{}}
	queryHalls := `SELECT id, establishment_id, name, description, capacity, has_air_conditioner, photos, created_at, updated_at FROM halls WHERE establishment_id = $1 ORDER BY name`
	hallRows, err := s.DB.QueryContext(ctx, queryHalls, est.ID); if err != nil { return nil, fmt.Errorf("err fetching halls for est ID %d: %w", est.ID, err) }; defer hallRows.Close()
	for hallRows.Next() {
		var h models.Hall; var hallPhotos pq.StringArray
		if err := hallRows.Scan(&h.ID, &h.EstablishmentID, &h.Name, &h.Description, &h.Capacity, &h.HasAirConditioner, &hallPhotos, &h.CreatedAt, &h.UpdatedAt); err != nil { return nil, fmt.Errorf("err scanning hall for est ID %d: %w", est.ID, err) } 
		h.Photos = []string(hallPhotos)
		hallWithPlaces := models.HallWithPlaces{Hall: h, Places: []models.Place{}}
		queryPlaces := `SELECT id, hall_id, name, type, status, visual_info, created_at, updated_at FROM places WHERE hall_id = $1 ORDER BY name` 
		placeRows, errP := s.DB.QueryContext(ctx, queryPlaces, h.ID); if errP != nil { return nil, fmt.Errorf("err fetching places for hall ID %d: %w", h.ID, errP) }
		for placeRows.Next() {
			var p models.Place
			if errPScan := placeRows.Scan(&p.ID, &p.HallID, &p.Name, &p.Type, &p.Status, &p.VisualInfo, &p.CreatedAt, &p.UpdatedAt); errPScan != nil { placeRows.Close(); return nil, fmt.Errorf("err scanning place for hall ID %d: %w", h.ID, errPScan) }
			hallWithPlaces.Places = append(hallWithPlaces.Places, p)
		}
		if errPRows := placeRows.Err(); errPRows != nil { placeRows.Close(); return nil, fmt.Errorf("err after iterating places for hall ID %d: %w", h.ID, errPRows) }
		placeRows.Close(); details.Halls = append(details.Halls, hallWithPlaces)
	}
	if err = hallRows.Err(); err != nil { return nil, fmt.Errorf("err after iterating halls for est ID %d: %w", est.ID, err) }
	if est.Type == models.RestaurantEstablishment {
		queryMenu := `SELECT id, establishment_id, category_name, name, price, description, photo_url, created_at, updated_at FROM menu_items WHERE establishment_id = $1 ORDER BY category_name, name`
		menuRows, errM := s.DB.QueryContext(ctx, queryMenu, est.ID); if errM != nil { return nil, fmt.Errorf("err fetching menu for est ID %d: %w", est.ID, errM) }; defer menuRows.Close()
		menuMap := make(map[string][]models.MenuItem)
		for menuRows.Next() {
			var mi models.MenuItem
			if errMScan := menuRows.Scan(&mi.ID, &mi.EstablishmentID, &mi.CategoryName, &mi.Name, &mi.Price, &mi.Description, &mi.PhotoURL, &mi.CreatedAt, &mi.UpdatedAt); errMScan != nil { return nil, fmt.Errorf("err scanning menu item for est ID %d: %w", est.ID, errMScan) }
			menuMap[mi.CategoryName] = append(menuMap[mi.CategoryName], mi)
		}
		if errMRows := menuRows.Err(); errMRows != nil { return nil, fmt.Errorf("err after iterating menu items for est ID %d: %w", est.ID, errMRows) }
		for category, items := range menuMap { details.Menu = append(details.Menu, models.MenuCategory{Category: category, Items: items}) }
	}
	return details, nil
}
func (s *PGStore) GetEstablishmentDetailsByOwnerID(ctx context.Context, ownerUserID int64) (*models.EstablishmentWithDetails, error) {
	est := &models.Establishment{}; queryEst := `SELECT id, owner_user_id, name, type, address, working_hours, description, photos, created_at, updated_at FROM establishments WHERE owner_user_id = $1 LIMIT 1`
	err := s.DB.QueryRowContext(ctx, queryEst, ownerUserID).Scan(&est.ID, &est.OwnerUserID, &est.Name, &est.Type, &est.Address, &est.WorkingHours, &est.Description, pq.Array(&est.Photos), &est.CreatedAt, &est.UpdatedAt)
	if err != nil { if errors.Is(err, sql.ErrNoRows) { return nil, nil }; return nil, fmt.Errorf("err fetching est by owner ID %d: %w", ownerUserID, err) }
	return s.getDetailsForEstablishment(ctx, est)
}
func (s *PGStore) GetEstablishmentDetailsByID(ctx context.Context, establishmentID int64) (*models.EstablishmentWithDetails, error) {
	est, err := s.GetEstablishmentByID(ctx, establishmentID); if err != nil { return nil, err }; if est == nil { return nil, nil }
	return s.getDetailsForEstablishment(ctx, est)
}
func (s *PGStore) UpdateEstablishment(ctx context.Context, establishmentID int64, ownerUserID int64, input *models.EstablishmentUpdateInput) (*models.Establishment, error) {
	currentEst, err := s.GetEstablishmentByID(ctx, establishmentID); if err != nil { return nil, fmt.Errorf("failed to fetch current est for update: %w", err) }; if currentEst == nil { return nil, sql.ErrNoRows }
	if currentEst.OwnerUserID != ownerUserID { return nil, errors.New("permission denied: you can only update your own establishment") }
	newName := currentEst.Name; if input.Name != nil { newName = *input.Name }; newAddress := currentEst.Address; if input.Address != nil { newAddress = *input.Address }; newWorkingHours := currentEst.WorkingHours; if input.WorkingHours != nil { newWorkingHours = *input.WorkingHours }; newDescription := currentEst.Description; if input.Description != nil { newDescription = *input.Description }
	query := `UPDATE establishments SET name = $1, address = $2, working_hours = $3, description = $4, updated_at = NOW() WHERE id = $5 AND owner_user_id = $6 RETURNING id, owner_user_id, name, type, address, working_hours, description, photos, created_at, updated_at`
	updatedEst := &models.Establishment{}
	err = s.DB.QueryRowContext(ctx, query, newName, newAddress, newWorkingHours, newDescription, establishmentID, ownerUserID).Scan(&updatedEst.ID, &updatedEst.OwnerUserID, &updatedEst.Name, &updatedEst.Type, &updatedEst.Address, &updatedEst.WorkingHours, &updatedEst.Description, pq.Array(&updatedEst.Photos), &updatedEst.CreatedAt, &updatedEst.UpdatedAt)
	if err != nil { if errors.Is(err, sql.ErrNoRows) { return nil, errors.New("establishment not found or permission denied for update") }; log.Printf("Error updating est in DB (ID: %d): %v", establishmentID, err); return nil, fmt.Errorf("could not update est: %w", err) }
	return updatedEst, nil
}

// --- Hall Methods (CreateHall, DeleteHall, GetHallByID, UpdateHall - без изменений) ---
// ... (код методов Hall остается прежним) ...
func (s *PGStore) CreateHall(ctx context.Context, hall *models.Hall) (*models.Hall, error) {
	query := `INSERT INTO halls (establishment_id, name, description, capacity, has_air_conditioner, photos, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, establishment_id, name, description, capacity, has_air_conditioner, photos, created_at, updated_at`
	now := time.Now(); createdHall := &models.Hall{}; var photosData pq.StringArray
	if hall.Photos != nil { photosData = pq.StringArray(hall.Photos) }
	err := s.DB.QueryRowContext(ctx, query, hall.EstablishmentID, hall.Name, hall.Description, hall.Capacity, hall.HasAirConditioner, photosData, now, now).Scan(&createdHall.ID, &createdHall.EstablishmentID, &createdHall.Name, &createdHall.Description, &createdHall.Capacity, &createdHall.HasAirConditioner, pq.Array(&createdHall.Photos), &createdHall.CreatedAt, &createdHall.UpdatedAt)
	if err != nil { log.Printf("Error creating hall for est ID %d: %v", hall.EstablishmentID, err); return nil, fmt.Errorf("could not create hall: %w", err) }
	return createdHall, nil
}
func (s *PGStore) DeleteHall(ctx context.Context, hallID int64, establishmentID int64, ownerUserID int64) error {
	est, err := s.GetEstablishmentByID(ctx, establishmentID); if err != nil { return fmt.Errorf("error verifying est for hall deletion: %w", err) }; if est == nil { return sql.ErrNoRows }
	if est.OwnerUserID != ownerUserID { return errors.New("permission denied: you do not own this establishment") }
	query := `DELETE FROM halls WHERE id = $1 AND establishment_id = $2`
	result, err := s.DB.ExecContext(ctx, query, hallID, establishmentID); if err != nil { log.Printf("Error deleting hall ID %d: %v", hallID, err); return fmt.Errorf("could not delete hall: %w", err) }
	rowsAffected, err := result.RowsAffected(); if err != nil { log.Printf("Error getting affected rows for hall ID %d: %v", hallID, err); return fmt.Errorf("could not confirm hall deletion: %w", err) }
	if rowsAffected == 0 { return sql.ErrNoRows }
	return nil
}
func (s *PGStore) GetHallByID(ctx context.Context, hallID int64) (*models.Hall, error) {
	query := `SELECT id, establishment_id, name, description, capacity, has_air_conditioner, photos, created_at, updated_at FROM halls WHERE id = $1`
	hall := &models.Hall{}; var photosData pq.StringArray
	err := s.DB.QueryRowContext(ctx, query, hallID).Scan(&hall.ID, &hall.EstablishmentID, &hall.Name, &hall.Description, &hall.Capacity, &hall.HasAirConditioner, &photosData, &hall.CreatedAt, &hall.UpdatedAt)
	if err != nil { if errors.Is(err, sql.ErrNoRows) { return nil, nil }; return nil, fmt.Errorf("could not get hall by ID %d: %w", hallID, err) }
	hall.Photos = []string(photosData); return hall, nil
}
func (s *PGStore) UpdateHall(ctx context.Context, hallID int64, establishmentID int64, ownerUserID int64, input *models.HallUpdateInput) (*models.Hall, error) {
	est, err := s.GetEstablishmentByID(ctx, establishmentID); if err != nil { return nil, fmt.Errorf("error verifying est for hall update: %w", err) }; if est == nil { return nil, errors.New("parent establishment not found") }
	if est.OwnerUserID != ownerUserID { return nil, errors.New("permission denied: you do not own this establishment") }
	currentHall, err := s.GetHallByID(ctx, hallID); if err != nil { return nil, fmt.Errorf("failed to fetch current hall for update: %w", err) }; if currentHall == nil { return nil, sql.ErrNoRows }
	if currentHall.EstablishmentID != establishmentID { return nil, errors.New("hall does not belong to the specified establishment") }
	newName := currentHall.Name; if input.Name != nil { newName = *input.Name }; newDescription := currentHall.Description; if input.Description != nil { newDescription = input.Description }; newCapacity := currentHall.Capacity; if input.Capacity != nil { newCapacity = *input.Capacity }; newHasAirConditioner := currentHall.HasAirConditioner; if input.HasAirConditioner != nil { newHasAirConditioner = *input.HasAirConditioner}
	query := `UPDATE halls SET name = $1, description = $2, capacity = $3, has_air_conditioner = $4, updated_at = NOW() WHERE id = $5 AND establishment_id = $6 RETURNING id, establishment_id, name, description, capacity, has_air_conditioner, photos, created_at, updated_at`
	updatedHall := &models.Hall{}; var photosData pq.StringArray
	err = s.DB.QueryRowContext(ctx, query, newName, newDescription, newCapacity, newHasAirConditioner, hallID, establishmentID).Scan(&updatedHall.ID, &updatedHall.EstablishmentID, &updatedHall.Name, &updatedHall.Description, &updatedHall.Capacity, &updatedHall.HasAirConditioner, &photosData, &updatedHall.CreatedAt, &updatedHall.UpdatedAt)
	if err != nil { if errors.Is(err, sql.ErrNoRows) { return nil, errors.New("hall not found or permission denied for update") }; log.Printf("Error updating hall in DB (ID: %d): %v", hallID, err); return nil, fmt.Errorf("could not update hall: %w", err) }
	updatedHall.Photos = []string(photosData); return updatedHall, nil
}


// --- Place Methods ---
func (s *PGStore) CreatePlace(ctx context.Context, place *models.Place) (*models.Place, error) {
	query := `
		INSERT INTO places (hall_id, name, type, status, visual_info, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, hall_id, name, type, status, visual_info, created_at, updated_at`
	
	now := time.Now()
	createdPlace := &models.Place{}

	// Устанавливаем статус по умолчанию, если он не передан (хотя модель PlaceCreateInput может это требовать)
	statusToInsert := place.Status
	if statusToInsert == "" {
		statusToInsert = models.PlaceStatusFree
	}
	// visual_info должен быть валидной JSON строкой или NULL
	var visualInfo sql.NullString
	if place.VisualInfo != "" {
		visualInfo.String = place.VisualInfo
		visualInfo.Valid = true
	}


	err := s.DB.QueryRowContext(ctx, query,
		place.HallID,
		place.Name,
		place.Type, // Может быть пустой строкой, если тип не указан
		statusToInsert,
		visualInfo, // Передаем sql.NullString
		now,
		now,
	).Scan(
		&createdPlace.ID,
		&createdPlace.HallID,
		&createdPlace.Name,
		&createdPlace.Type,
		&createdPlace.Status,
		&createdPlace.VisualInfo, // Будет сканироваться как строка
		&createdPlace.CreatedAt,
		&createdPlace.UpdatedAt,
	)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code.Name() == "unique_violation" { // Например, (hall_id, name)
			return nil, fmt.Errorf("place with name '%s' already exists in this hall: %w", place.Name, err)
		}
		log.Printf("Error creating place in DB for hall ID %d: %v", place.HallID, err)
		return nil, fmt.Errorf("could not create place: %w", err)
	}
	return createdPlace, nil
}

func (s *PGStore) UpdatePlace(ctx context.Context, placeID int64, hallID int64, establishmentID int64, ownerUserID int64, input *models.PlaceUpdateInput) (*models.Place, error) {
	// 1. Проверка прав на изменение зала (и, следовательно, мест в нем)
	// Это можно сделать, проверив, что заведение, которому принадлежит зал, принадлежит ownerUserID
	hall, err := s.GetHallByID(ctx, hallID)
	if err != nil { return nil, fmt.Errorf("error verifying hall for place update: %w", err)}
	if hall == nil { return nil, errors.New("parent hall not found")}
	if hall.EstablishmentID != establishmentID { return nil, errors.New("hall does not belong to the specified establishment for place update")}
	
	est, err := s.GetEstablishmentByID(ctx, establishmentID)
	if err != nil { return nil, fmt.Errorf("error verifying establishment for place update: %w", err)}
	if est == nil { return nil, errors.New("parent establishment not found for place update")}
	if est.OwnerUserID != ownerUserID { return nil, errors.New("permission denied: you do not own this establishment to update places")}

	// 2. Получаем текущее место
	// (Можно создать GetPlaceByID, но для MVP обновим через UPDATE ... WHERE id=... AND hall_id=...)
	// Для простоты, будем строить запрос на обновление только переданных полей.
	
	var setClauses []string
	var args []interface{}
	argID := 1

	if input.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argID))
		args = append(args, *input.Name)
		argID++
	}
	if input.Type != nil {
		setClauses = append(setClauses, fmt.Sprintf("type = $%d", argID))
		args = append(args, *input.Type)
		argID++
	}
	if input.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", argID))
		args = append(args, *input.Status)
		argID++
	}
	if input.VisualInfo != nil {
		setClauses = append(setClauses, fmt.Sprintf("visual_info = $%d", argID))
		args = append(args, *input.VisualInfo) // VisualInfo это строка JSON
		argID++
	}

	if len(setClauses) == 0 {
		return nil, errors.New("no fields to update for place")
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argID))
	args = append(args, time.Now())
	argID++

	args = append(args, placeID, hallID) // Для WHERE

	query := fmt.Sprintf(`
		UPDATE places SET %s 
		WHERE id = $%d AND hall_id = $%d
		RETURNING id, hall_id, name, type, status, visual_info, created_at, updated_at`,
		strings.Join(setClauses, ", "), argID, argID+1)

	updatedPlace := &models.Place{}
	err = s.DB.QueryRowContext(ctx, query, args...).Scan(
		&updatedPlace.ID, &updatedPlace.HallID, &updatedPlace.Name, &updatedPlace.Type,
		&updatedPlace.Status, &updatedPlace.VisualInfo,
		&updatedPlace.CreatedAt, &updatedPlace.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("place not found or does not belong to the specified hall")
		}
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code.Name() == "unique_violation" {
             return nil, fmt.Errorf("place with this name already exists in the hall: %w", err)
        }
		log.Printf("Error updating place in DB (ID: %d): %v", placeID, err)
		return nil, fmt.Errorf("could not update place: %w", err)
	}
	return updatedPlace, nil
}

func (s *PGStore) DeletePlace(ctx context.Context, placeID int64, hallID int64, establishmentID int64, ownerUserID int64) error {
	// 1. Проверка прав (аналогично UpdatePlace)
	hall, err := s.GetHallByID(ctx, hallID)
	if err != nil { return fmt.Errorf("error verifying hall for place deletion: %w", err)}
	if hall == nil { return errors.New("parent hall not found")}
	if hall.EstablishmentID != establishmentID { return errors.New("hall does not belong to the specified establishment for place deletion")}
	
	est, err := s.GetEstablishmentByID(ctx, establishmentID)
	if err != nil { return fmt.Errorf("error verifying establishment for place deletion: %w", err)}
	if est == nil { return errors.New("parent establishment not found for place deletion")}
	if est.OwnerUserID != ownerUserID { return errors.New("permission denied: you do not own this establishment to delete places")}

	// 2. Удаляем место
	// TODO: Проверить, нет ли активных бронирований для этого места, прежде чем удалять
	query := `DELETE FROM places WHERE id = $1 AND hall_id = $2`
	result, err := s.DB.ExecContext(ctx, query, placeID, hallID)
	if err != nil {
		log.Printf("Error deleting place ID %d from hall ID %d: %v", placeID, hallID, err)
		return fmt.Errorf("could not delete place: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting affected rows after deleting place ID %d: %v", placeID, err)
		return fmt.Errorf("could not confirm place deletion: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows // Место не найдено или не принадлежит этому залу
	}
	return nil
}

func (s *PGStore) GetPlacesByHallID(ctx context.Context, hallID int64) ([]models.Place, error) {
	var places []models.Place
	queryPlaces := `SELECT id, hall_id, name, type, status, visual_info, created_at, updated_at FROM places WHERE hall_id = $1 ORDER BY name`
	placeRows, err := s.DB.QueryContext(ctx, queryPlaces, hallID)
	if err != nil { return nil, fmt.Errorf("error fetching places for hall ID %d: %w", hallID, err) }
	defer placeRows.Close()
	for placeRows.Next() {
		var p models.Place
		if err := placeRows.Scan(&p.ID, &p.HallID, &p.Name, &p.Type, &p.Status, &p.VisualInfo, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning place for hall ID %d: %w", hallID, err)
		}
		places = append(places, p)
	}
	if err = placeRows.Err(); err != nil { return nil, fmt.Errorf("error after iterating places for hall ID %d: %w", hallID, err) }
	return places, nil
}