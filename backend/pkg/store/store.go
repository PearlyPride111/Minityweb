package store

import (
	"context"
	"minityweb/backend/pkg/models" // Убедитесь, что путь к моделям корректен
)

// DataStore определяет методы, которые должно реализовывать хранилище данных.
type DataStore interface {
	// User methods
	CreateUser(ctx context.Context, user *models.User) (int64, error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error)
	GetUserByID(ctx context.Context, id int64) (*models.User, error)

	// Establishment methods
	GetEstablishments(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*models.Establishment, error)
	GetEstablishmentByID(ctx context.Context, id int64) (*models.Establishment, error)
	GetEstablishmentDetailsByOwnerID(ctx context.Context, ownerUserID int64) (*models.EstablishmentWithDetails, error)
	GetEstablishmentDetailsByID(ctx context.Context, establishmentID int64) (*models.EstablishmentWithDetails, error)
	UpdateEstablishment(ctx context.Context, establishmentID int64, ownerUserID int64, input *models.EstablishmentUpdateInput) (*models.Establishment, error)

	// Hall methods
	CreateHall(ctx context.Context, hall *models.Hall) (*models.Hall, error)
	DeleteHall(ctx context.Context, hallID int64, establishmentID int64, ownerUserID int64) error
	UpdateHall(ctx context.Context, hallID int64, establishmentID int64, ownerUserID int64, input *models.HallUpdateInput) (*models.Hall, error)
	GetHallByID(ctx context.Context, hallID int64) (*models.Hall, error) 

	// Place methods (ВОТ ЭТИ МЕТОДЫ ДОЛЖНЫ БЫТЬ ЗДЕСЬ)
	CreatePlace(ctx context.Context, place *models.Place) (*models.Place, error)
	UpdatePlace(ctx context.Context, placeID int64, hallID int64, establishmentID int64, ownerUserID int64, input *models.PlaceUpdateInput) (*models.Place, error)
	DeletePlace(ctx context.Context, placeID int64, hallID int64, establishmentID int64, ownerUserID int64) error
	GetPlacesByHallID(ctx context.Context, hallID int64) ([]models.Place, error)
}