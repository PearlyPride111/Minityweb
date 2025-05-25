package models

import "time"

// UserRole определяет возможные роли пользователя
type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	PhoneNumber  string    `json:"phone_number"`
	Email        *string   `json:"email,omitempty"`
	PasswordHash string    `json:"-"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserRegisterInput struct {
	Name        string  `json:"name" validate:"required,min=2,max=100"`
	PhoneNumber string  `json:"phone_number" validate:"required,e164"`
	Email       *string `json:"email,omitempty" validate:"omitempty,email"`
	Password    string  `json:"password" validate:"required,min=8,max=72"`
}

type UserLoginInput struct {
	PhoneNumber string `json:"phone_number" validate:"required,e164"`
	Password    string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	User         *User  `json:"user,omitempty"`
}

type ErrorResponse struct {
	Error      string            `json:"error"`
	Details    map[string]string `json:"details,omitempty"`
	StatusCode int               `json:"status_code,omitempty"`
}

type EstablishmentType string

const (
	RestaurantEstablishment EstablishmentType = "restaurant"
	CoworkingEstablishment  EstablishmentType = "coworking"
)

type Establishment struct {
	ID              int64             `json:"id"`
	OwnerUserID     int64             `json:"owner_user_id"`
	Name            string            `json:"name" validate:"required,min=3,max=100"`
	Type            EstablishmentType `json:"type" validate:"required"`
	Address         string            `json:"address" validate:"omitempty,max=255"`
	WorkingHours    string            `json:"working_hours" validate:"omitempty,max=100"`
	Description     string            `json:"description" validate:"omitempty,max=1000"`
	Photos          []string          `json:"photos,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

type EstablishmentUpdateInput struct {
	Name         *string `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Address      *string `json:"address,omitempty" validate:"omitempty,max=255"`
	WorkingHours *string `json:"working_hours,omitempty" validate:"omitempty,max=100"`
	Description  *string `json:"description,omitempty" validate:"omitempty,max=1000"`
}

type EstablishmentWithDetails struct {
	Establishment
	Halls []HallWithPlaces `json:"halls"`
	Menu  []MenuCategory   `json:"menu,omitempty"`
}

type Hall struct {
	ID                int64     `json:"id"`
	EstablishmentID   int64     `json:"establishment_id"`
	Name              string    `json:"name" validate:"required,min=2,max=100"`
	Description       *string   `json:"description,omitempty" validate:"omitempty,max=500"`
	Capacity          int       `json:"capacity" validate:"omitempty,min=0"`
	HasAirConditioner bool      `json:"has_air_conditioner"`
	Photos            []string  `json:"photos,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type HallCreateInput struct {
	Name              string  `json:"name" validate:"required,min=2,max=100"`
	Description       *string `json:"description,omitempty" validate:"omitempty,max=500"`
	Capacity          int     `json:"capacity" validate:"omitempty,min=0"`
	HasAirConditioner bool    `json:"has_air_conditioner"`
}

type HallUpdateInput struct {
	Name              *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description       *string `json:"description,omitempty" validate:"omitempty,max=500"`
	Capacity          *int    `json:"capacity,omitempty" validate:"omitempty,min=0"`
	HasAirConditioner *bool   `json:"has_air_conditioner,omitempty"`
}

type HallWithPlaces struct {
	Hall
	Places []Place `json:"places"`
}

type PlaceStatus string
const (
	PlaceStatusFree        PlaceStatus = "free"
	PlaceStatusBooked      PlaceStatus = "booked"
	PlaceStatusOccupied    PlaceStatus = "occupied"
	PlaceStatusUnavailable PlaceStatus = "unavailable"
)

type Place struct {
	ID           int64       `json:"id"`
	HallID       int64       `json:"hall_id"`
	Name         string      `json:"name" validate:"required,min=1,max=50"`
	Type         string      `json:"type,omitempty" validate:"omitempty,max=50"`
	Status       PlaceStatus `json:"status" validate:"required,oneof=free booked occupied unavailable"`
	VisualInfo   string      `json:"visual_info,omitempty"` 
	IconFree     *string     `json:"icon_free_url,omitempty"`
	IconBooked   *string     `json:"icon_booked_url,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

// PlaceCreateInput - структура для создания нового места
type PlaceCreateInput struct {
    Name       string  `json:"name" validate:"required,min=1,max=50"`
    Type       *string `json:"type,omitempty" validate:"omitempty,max=50"`
    VisualInfo string  `json:"visual_info" validate:"required,json"` // Ожидаем JSON строку { "x": 10, "y": 20 }
    // Статус по умолчанию 'free' при создании
}

// PlaceUpdateInput - структура для обновления места (например, только visual_info или статус)
type PlaceUpdateInput struct {
    Name       *string      `json:"name,omitempty" validate:"omitempty,min=1,max=50"`
    Type       *string      `json:"type,omitempty" validate:"omitempty,max=50"`
    Status     *PlaceStatus `json:"status,omitempty" validate:"omitempty,oneof=free booked occupied unavailable"`
    VisualInfo *string      `json:"visual_info,omitempty" validate:"omitempty,json"`
}


type MenuCategory struct {
	Category string     `json:"category"`
	Items    []MenuItem `json:"items"`
}

type MenuItem struct {
	ID              int64     `json:"id"`
	EstablishmentID int64     `json:"establishment_id"`
	CategoryName    string    `json:"category_name" validate:"required,max=100"`
	Name            string    `json:"name" validate:"required,min=2,max=100"`
	Price           float64   `json:"price" validate:"required,min=0"`
	Description     *string   `json:"description,omitempty" validate:"omitempty,max=500"`
	PhotoURL        *string   `json:"photo_url,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type BookingStatus string
const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusCompleted BookingStatus = "completed"
)

type Booking struct {
	ID              int64         `json:"id"`
	UserID          int64         `json:"user_id"`
	EstablishmentID int64         `json:"establishment_id"`
	HallID          int64         `json:"hall_id"`
	PlaceIDs        []int64       `json:"place_ids"`
	BookingTime     time.Time     `json:"booking_time"`
	DurationMinutes int           `json:"duration_minutes,omitempty"`
	PeopleCount     int           `json:"people_count"`
	Status          BookingStatus `json:"status"`
	Notes           *string       `json:"notes,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

type Review struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	EstablishmentID int64     `json:"establishment_id"`
	Rating          int       `json:"rating"`
	Text            string    `json:"text"`
	PhotoURL        *string   `json:"photo_url,omitempty"`
	IsModerated     bool      `json:"is_moderated"`
	IsApproved      bool      `json:"is_approved"`
	CreatedAt       time.Time `json:"created_at"`
}