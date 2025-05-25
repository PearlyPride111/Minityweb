package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"minityweb/backend/pkg/models"
	"minityweb/backend/pkg/store"
	"minityweb/backend/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

type Claims struct {
	UserID int64          `json:"user_id"`
	Role   models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

type APIHandler struct {
	Store     store.DataStore
	Validator *validator.Validate
	JWTSecret []byte
}

func NewAPIHandler(s store.DataStore, jwtSecret string) *APIHandler {
	return &APIHandler{
		Store:     s,
		Validator: validator.New(),
		JWTSecret: []byte(jwtSecret),
	}
}

func (h *APIHandler) generateJWT(user *models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour); claims := &Claims{UserID: user.ID, Role: user.Role, RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(expirationTime), IssuedAt:  jwt.NewNumericDate(time.Now()), Issuer: "minity-app"}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims); tokenString, err := token.SignedString(h.JWTSecret); if err != nil { return "", err }; return tokenString, nil
}

func (h *APIHandler) HealthCheck(w http.ResponseWriter, r *http.Request) { /* ... без изменений ... */ utils.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "ok", "timestamp": time.Now().String()}) }
func (h *APIHandler) RegisterUser(w http.ResponseWriter, r *http.Request) { /* ... без изменений ... */ 
	var input models.UserRegisterInput; if err := json.NewDecoder(r.Body).Decode(&input); err != nil { utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", nil); return }; defer r.Body.Close()
	if err := h.Validator.Struct(input); err != nil { validationErrors := make(map[string]string); for _, errVal := range err.(validator.ValidationErrors) { validationErrors[errVal.Field()] = errVal.Tag() }; utils.RespondWithError(w, http.StatusBadRequest, "Validation failed", validationErrors); return }
	existingUser, err := h.Store.GetUserByPhoneNumber(r.Context(), input.PhoneNumber); if err != nil && !errors.Is(err, sql.ErrNoRows) { log.Printf("Error checking existing user: %v", err); utils.RespondWithError(w, http.StatusInternalServerError, "Could not process registration", nil); return }; if existingUser != nil { utils.RespondWithError(w, http.StatusConflict, "User with this phone number already exists", nil); return }
	hashedPassword, err := utils.HashPassword(input.Password); if err != nil { log.Printf("Error hashing password: %v", err); utils.RespondWithError(w, http.StatusInternalServerError, "Could not process registration", nil); return }
	newUser := models.User{Name: input.Name, PhoneNumber: input.PhoneNumber, Email: input.Email, PasswordHash: hashedPassword, Role: models.RoleUser }
	userID, err := h.Store.CreateUser(r.Context(), &newUser)
	if err != nil { if pgErr, ok := err.(*pq.Error); ok && pgErr.Code.Name() == "unique_violation" { utils.RespondWithError(w, http.StatusConflict, "User with this phone number or email already exists", nil) } else { log.Printf("Failed to create user (handler): %v", err); utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create user", nil) }; return }
	utils.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{"message": "User registered successfully. Please login.", "user_id": userID})
}
func (h *APIHandler) LoginUser(w http.ResponseWriter, r *http.Request) { /* ... без изменений ... */ 
	var input models.UserLoginInput; if err := json.NewDecoder(r.Body).Decode(&input); err != nil { utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", nil); return }; defer r.Body.Close()
	if err := h.Validator.Struct(input); err != nil { validationErrors := make(map[string]string); for _, errVal := range err.(validator.ValidationErrors) { validationErrors[errVal.Field()] = errVal.Tag() }; utils.RespondWithError(w, http.StatusBadRequest, "Validation failed", validationErrors); return }
	user, err := h.Store.GetUserByPhoneNumber(r.Context(), input.PhoneNumber); if err != nil { if errors.Is(err, sql.ErrNoRows) { utils.RespondWithError(w, http.StatusUnauthorized, "Invalid phone number or password", nil); return }; log.Printf("Error fetching user for login: %v", err); utils.RespondWithError(w, http.StatusInternalServerError, "Login failed", nil); return }; if user == nil { utils.RespondWithError(w, http.StatusUnauthorized, "Invalid phone number or password", nil); return }
	if !utils.CheckPasswordHash(input.Password, user.PasswordHash) { utils.RespondWithError(w, http.StatusUnauthorized, "Invalid phone number or password", nil); return }
	tokenString, err := h.generateJWT(user); if err != nil { log.Printf("Error generating JWT token: %v", err); utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token", nil); return }
	user.PasswordHash = ""; utils.RespondWithJSON(w, http.StatusOK, models.AuthResponse{ Token: tokenString, User: user, })
}
func (h *APIHandler) GetEstablishments(w http.ResponseWriter, r *http.Request) { /* ... без изменений ... */ 
	ctx := r.Context(); queryParams := r.URL.Query(); limitStr := queryParams.Get("limit"); offsetStr := queryParams.Get("offset"); typeStr := queryParams.Get("type")
	limit, err := strconv.Atoi(limitStr); if err != nil || limit <= 0 { limit = 10 }; if limit > 100 { limit = 100 }
	offset, err := strconv.Atoi(offsetStr); if err != nil || offset < 0 { offset = 0 }
	filters := make(map[string]interface{}); if typeStr != "" { validType := models.EstablishmentType(typeStr); if validType == models.RestaurantEstablishment || validType == models.CoworkingEstablishment { filters["type"] = validType } else { utils.RespondWithError(w, http.StatusBadRequest, "Invalid establishment type filter", nil); return } }
	establishments, err := h.Store.GetEstablishments(ctx, limit, offset, filters); if err != nil { log.Printf("Error fetching establishments: %v", err); utils.RespondWithError(w, http.StatusInternalServerError, "Could not fetch establishments", nil); return }
	if establishments == nil { establishments = []*models.Establishment{} }; utils.RespondWithJSON(w, http.StatusOK, establishments)
}
func (h *APIHandler) GetEstablishmentByIDPublic(w http.ResponseWriter, r *http.Request) { /* ... без изменений ... */ 
    vars := mux.Vars(r); idStr, ok := vars["id"]; if !ok { utils.RespondWithError(w, http.StatusBadRequest, "Establishment ID is missing in URL path", nil); return }
    id, err := strconv.ParseInt(idStr, 10, 64); if err != nil { utils.RespondWithError(w, http.StatusBadRequest, "Invalid establishment ID format in URL path", nil); return }
    establishmentDetails, err := h.Store.GetEstablishmentDetailsByID(r.Context(), id)
    if err != nil { if errors.Is(err, sql.ErrNoRows) { utils.RespondWithError(w, http.StatusNotFound, "Establishment not found", nil); return }; log.Printf("Error fetching est details by ID %d: %v", id, err); utils.RespondWithError(w, http.StatusInternalServerError, "Could not fetch est details", nil); return }; if establishmentDetails == nil { utils.RespondWithError(w, http.StatusNotFound, "Establishment not found", nil); return }; utils.RespondWithJSON(w, http.StatusOK, establishmentDetails)
}
func (h *APIHandler) GetMyEstablishmentAdminMVP(w http.ResponseWriter, r *http.Request) { /* ... без изменений ... */ 
	ownerIDStr := r.URL.Query().Get("owner_id"); var ownerIDToFetch int64 = 1
	if ownerIDStr != "" { parsedID, err := strconv.ParseInt(ownerIDStr, 10, 64); if err == nil && parsedID > 0 { ownerIDToFetch = parsedID } else { log.Printf("Admin MVP GET: Invalid owner_id: %s. Defaulting to %d.", ownerIDStr, ownerIDToFetch) } }
	log.Printf("Admin MVP GET: Fetching est for owner_id: %d", ownerIDToFetch)
	establishmentDetails, err := h.Store.GetEstablishmentDetailsByOwnerID(r.Context(), ownerIDToFetch)
	if err != nil { if errors.Is(err, sql.ErrNoRows) { log.Printf("Admin MVP GET: No est for owner_id %d", ownerIDToFetch); utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Заведение для админа ID %d не найдено.", ownerIDToFetch), nil); return }; log.Printf("Admin MVP GET: Error fetching est for owner_id %d: %v", ownerIDToFetch, err); utils.RespondWithError(w, http.StatusInternalServerError, "Не удалось загрузить данные.", nil); return }; if establishmentDetails == nil { log.Printf("Admin MVP GET: Est details nil for owner_id %d", ownerIDToFetch); utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Заведение для админа ID %d не найдено (nil).", ownerIDToFetch), nil); return }
	log.Printf("Admin MVP GET: Fetched '%s' for owner_id %d", establishmentDetails.Name, ownerIDToFetch); utils.RespondWithJSON(w, http.StatusOK, establishmentDetails)
}
func (h *APIHandler) UpdateMyEstablishmentAdminMVP(w http.ResponseWriter, r *http.Request) { /* ... без изменений ... */ 
	establishmentIDStr := r.URL.Query().Get("establishment_id"); if establishmentIDStr == "" { utils.RespondWithError(w, http.StatusBadRequest, "establishment_id required", nil); return }; establishmentIDToUpdate, err := strconv.ParseInt(establishmentIDStr, 10, 64); if err != nil || establishmentIDToUpdate <= 0 { utils.RespondWithError(w, http.StatusBadRequest, "Invalid establishment_id", nil); return }
	ownerIDStr := r.URL.Query().Get("owner_id"); var adminOwnerIDForDemo int64
	if ownerIDStr != "" { parsedID, errP := strconv.ParseInt(ownerIDStr, 10, 64); if errP != nil || parsedID <= 0 { utils.RespondWithError(w, http.StatusBadRequest, "Invalid owner_id", nil); return }; adminOwnerIDForDemo = parsedID } else { log.Println("Admin MVP PUT Est: owner_id required"); utils.RespondWithError(w, http.StatusBadRequest, "owner_id required", nil); return }
	var input models.EstablishmentUpdateInput; if err := json.NewDecoder(r.Body).Decode(&input); err != nil { utils.RespondWithError(w, http.StatusBadRequest, "Invalid payload: "+err.Error(), nil); return }; defer r.Body.Close()
	if err := h.Validator.Struct(input); err != nil { validationErrors := make(map[string]string); for _, errVal := range err.(validator.ValidationErrors) { validationErrors[errVal.Field()] = errVal.Tag() }; utils.RespondWithError(w, http.StatusBadRequest, "Validation failed", validationErrors); return }
	log.Printf("Admin MVP PUT Est: Updating est ID: %d by owner ID: %d", establishmentIDToUpdate, adminOwnerIDForDemo)
	updatedEstablishment, err := h.Store.UpdateEstablishment(r.Context(), establishmentIDToUpdate, adminOwnerIDForDemo, &input)
	if err != nil { if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "permission denied") { log.Printf("Admin MVP PUT Est: Est ID %d not found or permission denied for owner %d. Err: %v", establishmentIDToUpdate, adminOwnerIDForDemo, err); utils.RespondWithError(w, http.StatusNotFound, "Заведение не найдено или отказано в доступе.", nil); return }; log.Printf("Admin MVP PUT Est: Error updating est ID %d: %v", establishmentIDToUpdate, err); utils.RespondWithError(w, http.StatusInternalServerError, "Не удалось обновить.", nil); return }
	log.Printf("Admin MVP PUT Est: Updated est ID: %d, Name: %s", updatedEstablishment.ID, updatedEstablishment.Name); utils.RespondWithJSON(w, http.StatusOK, updatedEstablishment)
}
func (h *APIHandler) AdminCreateHallHandler(w http.ResponseWriter, r *http.Request) { /* ... без изменений ... */ 
	vars := mux.Vars(r); establishmentIDStr, ok := vars["establishment_id"]; if !ok { utils.RespondWithError(w, http.StatusBadRequest, "establishment_id missing", nil); return }; establishmentID, err := strconv.ParseInt(establishmentIDStr, 10, 64); if err != nil || establishmentID <= 0 { utils.RespondWithError(w, http.StatusBadRequest, "Invalid establishment_id", nil); return }
	ownerIDStr := r.URL.Query().Get("owner_id"); var simulatedOwnerID int64
	if ownerIDStr != "" { parsedID, errP := strconv.ParseInt(ownerIDStr, 10, 64); if errP != nil || parsedID <= 0 { utils.RespondWithError(w, http.StatusBadRequest, "Invalid owner_id for permission", nil); return }; simulatedOwnerID = parsedID } else { utils.RespondWithError(w, http.StatusBadRequest, "owner_id required for permission", nil); return }
	currentEst, err := h.Store.GetEstablishmentByID(r.Context(), establishmentID); if err != nil { log.Printf("AdminCreateHall: Err fetching est %d: %v", establishmentID, err); utils.RespondWithError(w, http.StatusInternalServerError, "Could not verify est.", nil); return }; if currentEst == nil { utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Est ID %d not found.", establishmentID), nil); return }
	if currentEst.OwnerUserID != simulatedOwnerID { log.Printf("AdminCreateHall: Permission denied. Est owner %d, simulated admin %d", currentEst.OwnerUserID, simulatedOwnerID); utils.RespondWithError(w, http.StatusForbidden, "No permission.", nil); return }
	var input models.HallCreateInput; if err := json.NewDecoder(r.Body).Decode(&input); err != nil { utils.RespondWithError(w, http.StatusBadRequest, "Invalid payload for hall: "+err.Error(), nil); return }; defer r.Body.Close()
	if err := h.Validator.Struct(input); err != nil { validationErrors := make(map[string]string); for _, errVal := range err.(validator.ValidationErrors) { validationErrors[errVal.Field()] = errVal.Tag() }; utils.RespondWithError(w, http.StatusBadRequest, "Validation failed for hall", validationErrors); return }
	hallToCreate := &models.Hall{EstablishmentID: establishmentID, Name: input.Name, Description: input.Description, Capacity: input.Capacity, HasAirConditioner: input.HasAirConditioner }
	createdHall, err := h.Store.CreateHall(r.Context(), hallToCreate); if err != nil { log.Printf("AdminCreateHall: Err creating hall for est %d: %v", establishmentID, err); utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create hall.", nil); return }
	log.Printf("AdminCreateHall: Created hall '%s' (ID: %d) for est ID: %d", createdHall.Name, createdHall.ID, establishmentID); utils.RespondWithJSON(w, http.StatusCreated, createdHall)
}
func (h *APIHandler) AdminUpdateHallHandler(w http.ResponseWriter, r *http.Request) { /* ... без изменений ... */ 
	vars := mux.Vars(r); establishmentIDStr, okEst := vars["establishment_id"]; hallIDStr, okHall := vars["hall_id"]; if !okEst || !okHall { utils.RespondWithError(w, http.StatusBadRequest, "establishment_id and hall_id missing", nil); return }
	establishmentID, err := strconv.ParseInt(establishmentIDStr, 10, 64); if err != nil || establishmentID <= 0 { utils.RespondWithError(w, http.StatusBadRequest, "Invalid establishment_id", nil); return }; hallID, err := strconv.ParseInt(hallIDStr, 10, 64); if err != nil || hallID <= 0 { utils.RespondWithError(w, http.StatusBadRequest, "Invalid hall_id", nil); return }
	ownerIDStr := r.URL.Query().Get("owner_id"); var simulatedOwnerID int64
	if ownerIDStr != "" { parsedID, errP := strconv.ParseInt(ownerIDStr, 10, 64); if errP != nil || parsedID <= 0 { utils.RespondWithError(w, http.StatusBadRequest, "Invalid owner_id for permission", nil); return }; simulatedOwnerID = parsedID } else { utils.RespondWithError(w, http.StatusBadRequest, "owner_id required for permission", nil); return }
	var input models.HallUpdateInput; if err := json.NewDecoder(r.Body).Decode(&input); err != nil { utils.RespondWithError(w, http.StatusBadRequest, "Invalid payload for hall update: "+err.Error(), nil); return }; defer r.Body.Close()
	if err := h.Validator.Struct(input); err != nil { validationErrors := make(map[string]string); for _, errVal := range err.(validator.ValidationErrors) { validationErrors[errVal.Field()] = errVal.Tag() }; utils.RespondWithError(w, http.StatusBadRequest, "Validation failed for hall update", validationErrors); return }
	log.Printf("Admin MVP PUT Hall: Updating hall ID: %d for est ID: %d by owner ID: %d", hallID, establishmentID, simulatedOwnerID)
	updatedHall, err := h.Store.UpdateHall(r.Context(), hallID, establishmentID, simulatedOwnerID, &input)
	if err != nil { if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "permission denied") { log.Printf("Admin MVP PUT Hall: Hall ID %d not found or permission denied. Err: %v", hallID, err); utils.RespondWithError(w, http.StatusNotFound, "Зал не найден или отказано в доступе.", nil); return }; log.Printf("Admin MVP PUT Hall: Error updating hall ID %d: %v", hallID, err); utils.RespondWithError(w, http.StatusInternalServerError, "Не удалось обновить зал.", nil); return }
	log.Printf("Admin MVP PUT Hall: Updated hall ID: %d, Name: %s", updatedHall.ID, updatedHall.Name); utils.RespondWithJSON(w, http.StatusOK, updatedHall)
}
func (h *APIHandler) AdminDeleteHallHandler(w http.ResponseWriter, r *http.Request) { /* ... без изменений ... */ 
	vars := mux.Vars(r); establishmentIDStr, okEst := vars["establishment_id"]; hallIDStr, okHall := vars["hall_id"]; if !okEst || !okHall { utils.RespondWithError(w, http.StatusBadRequest, "establishment_id and hall_id missing", nil); return }
	establishmentID, err := strconv.ParseInt(establishmentIDStr, 10, 64); if err != nil || establishmentID <= 0 { utils.RespondWithError(w, http.StatusBadRequest, "Invalid establishment_id", nil); return }; hallID, err := strconv.ParseInt(hallIDStr, 10, 64); if err != nil || hallID <= 0 { utils.RespondWithError(w, http.StatusBadRequest, "Invalid hall_id", nil); return }
	ownerIDStr := r.URL.Query().Get("owner_id"); var simulatedOwnerID int64
	if ownerIDStr != "" { parsedID, errP := strconv.ParseInt(ownerIDStr, 10, 64); if errP != nil || parsedID <= 0 { utils.RespondWithError(w, http.StatusBadRequest, "Invalid owner_id for permission", nil); return }; simulatedOwnerID = parsedID } else { utils.RespondWithError(w, http.StatusBadRequest, "owner_id required for permission", nil); return }
	log.Printf("Admin MVP DELETE Hall: Deleting hall ID: %d from est ID: %d by owner ID: %d", hallID, establishmentID, simulatedOwnerID)
	err = h.Store.DeleteHall(r.Context(), hallID, establishmentID, simulatedOwnerID)
	if err != nil { if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "permission denied") { log.Printf("Admin MVP DELETE Hall: Hall ID %d not found or permission denied. Err: %v", hallID, err); utils.RespondWithError(w, http.StatusNotFound, "Зал не найден или отказано в доступе.", nil); return }; log.Printf("Admin MVP DELETE Hall: Error deleting hall ID %d: %v", hallID, err); utils.RespondWithError(w, http.StatusInternalServerError, "Не удалось удалить зал.", nil); return }
	log.Printf("Admin MVP DELETE Hall: Deleted hall ID: %d from est ID: %d", hallID, establishmentID); utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Hall deleted successfully"})
}


// --- НОВЫЕ ОБРАБОТЧИКИ ДЛЯ МЕСТ (PLACES) ---

// AdminCreatePlaceHandler создает новое место в зале
func (h *APIHandler) AdminCreatePlaceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	establishmentIDStr, okEst := vars["establishment_id"]
	hallIDStr, okHall := vars["hall_id"]
	if !okEst || !okHall {
		utils.RespondWithError(w, http.StatusBadRequest, "establishment_id and hall_id are required in URL path", nil)
		return
	}
	establishmentID, err := strconv.ParseInt(establishmentIDStr, 10, 64)
	if err != nil { utils.RespondWithError(w, http.StatusBadRequest, "Invalid establishment_id", nil); return }
	hallID, err := strconv.ParseInt(hallIDStr, 10, 64)
	if err != nil { utils.RespondWithError(w, http.StatusBadRequest, "Invalid hall_id", nil); return }

	// MVP: Проверка прав (упрощенная)
	ownerIDStr := r.URL.Query().Get("owner_id")
	var simulatedOwnerID int64
	if ownerIDStr != "" {
		parsedID, _ := strconv.ParseInt(ownerIDStr, 10, 64)
		simulatedOwnerID = parsedID // Ошибку парсинга обработает проверка ниже, если ID = 0
	}
	if simulatedOwnerID == 0 { // Проверка, что owner_id передан и валиден
		utils.RespondWithError(w, http.StatusBadRequest, "Valid owner_id is required for permission check", nil)
		return
	}

	// Проверяем, что заведение и зал существуют и принадлежат этому админу
	hall, err := h.Store.GetHallByID(r.Context(), hallID)
	if err != nil || hall == nil || hall.EstablishmentID != establishmentID {
		utils.RespondWithError(w, http.StatusNotFound, "Hall not found or does not belong to the establishment.", nil)
		return
	}
	est, err := h.Store.GetEstablishmentByID(r.Context(), establishmentID)
	if err != nil || est == nil || est.OwnerUserID != simulatedOwnerID {
		utils.RespondWithError(w, http.StatusForbidden, "Permission denied to manage this hall's places.", nil)
		return
	}

	var input models.PlaceCreateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload for place: "+err.Error(), nil)
		return
	}
	defer r.Body.Close()

	if err := h.Validator.Struct(input); err != nil {
		// ... (обработка ошибок валидации)
		utils.RespondWithError(w, http.StatusBadRequest, "Validation failed for place data", nil) // Упрощено
		return
	}

	placeToCreate := &models.Place{
		HallID:     hallID,
		Name:       input.Name,
		Type:       "", // Инициализируем пустой строкой, если nil
		VisualInfo: input.VisualInfo,
		Status:     models.PlaceStatusFree, // Новые места по умолчанию свободны
	}
    if input.Type != nil {
        placeToCreate.Type = *input.Type
    }


	createdPlace, err := h.Store.CreatePlace(r.Context(), placeToCreate)
	if err != nil {
		log.Printf("AdminCreatePlace: Error creating place for hall %d: %v", hallID, err)
		if strings.Contains(err.Error(), "already exists") { // Проверка на уникальность имени места в зале
			utils.RespondWithError(w, http.StatusConflict, err.Error(), nil)
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create place.", nil)
		}
		return
	}
	log.Printf("AdminCreatePlace: Successfully created place '%s' (ID: %d) in hall ID: %d", createdPlace.Name, createdPlace.ID, hallID)
	utils.RespondWithJSON(w, http.StatusCreated, createdPlace)
}

// AdminUpdatePlaceHandler обновляет существующее место
func (h *APIHandler) AdminUpdatePlaceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	establishmentIDStr, _ := vars["establishment_id"]
	hallIDStr, _ := vars["hall_id"]
	placeIDStr, okPlace := vars["place_id"]
	if !okPlace { utils.RespondWithError(w, http.StatusBadRequest, "place_id required", nil); return }

	establishmentID, _ := strconv.ParseInt(establishmentIDStr, 10, 64)
	hallID, _ := strconv.ParseInt(hallIDStr, 10, 64)
	placeID, err := strconv.ParseInt(placeIDStr, 10, 64)
	if err != nil { utils.RespondWithError(w, http.StatusBadRequest, "Invalid place_id", nil); return }

	ownerIDStr := r.URL.Query().Get("owner_id")
	var simulatedOwnerID int64
	if ownerIDStr != "" { parsedID, _ := strconv.ParseInt(ownerIDStr, 10, 64); simulatedOwnerID = parsedID }
	if simulatedOwnerID == 0 { utils.RespondWithError(w, http.StatusBadRequest, "Valid owner_id required", nil); return }
	
	// Проверки прав и существования родительских сущностей выполняются в store.UpdatePlace

	var input models.PlaceUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid payload for place update: "+err.Error(), nil); return
	}
	defer r.Body.Close()

	if err := h.Validator.Struct(input); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Validation failed for place update", nil); return // Упрощено
	}

	log.Printf("AdminUpdatePlace: Updating place ID: %d in hall ID: %d", placeID, hallID)
	updatedPlace, err := h.Store.UpdatePlace(r.Context(), placeID, hallID, establishmentID, simulatedOwnerID, &input)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "permission denied") || strings.Contains(err.Error(), "not found") {
			utils.RespondWithError(w, http.StatusNotFound, "Место не найдено или отказано в доступе.", nil); return
		}
        if strings.Contains(err.Error(), "already exists") {
			utils.RespondWithError(w, http.StatusConflict, err.Error(), nil); return
        }
		log.Printf("AdminUpdatePlace: Error updating place ID %d: %v", placeID, err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Не удалось обновить место.", nil); return
	}
	log.Printf("AdminUpdatePlace: Updated place ID: %d, Name: %s", updatedPlace.ID, updatedPlace.Name)
	utils.RespondWithJSON(w, http.StatusOK, updatedPlace)
}

// AdminDeletePlaceHandler удаляет место
func (h *APIHandler) AdminDeletePlaceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	establishmentIDStr, _ := vars["establishment_id"]
	hallIDStr, _ := vars["hall_id"]
	placeIDStr, okPlace := vars["place_id"]
	if !okPlace { utils.RespondWithError(w, http.StatusBadRequest, "place_id required", nil); return }

	establishmentID, _ := strconv.ParseInt(establishmentIDStr, 10, 64)
	hallID, _ := strconv.ParseInt(hallIDStr, 10, 64)
	placeID, err := strconv.ParseInt(placeIDStr, 10, 64)
	if err != nil { utils.RespondWithError(w, http.StatusBadRequest, "Invalid place_id", nil); return }

	ownerIDStr := r.URL.Query().Get("owner_id")
	var simulatedOwnerID int64
	if ownerIDStr != "" { parsedID, _ := strconv.ParseInt(ownerIDStr, 10, 64); simulatedOwnerID = parsedID }
	if simulatedOwnerID == 0 { utils.RespondWithError(w, http.StatusBadRequest, "Valid owner_id required", nil); return }

	// Проверки прав и существования родительских сущностей выполняются в store.DeletePlace
	log.Printf("AdminDeletePlace: Deleting place ID: %d from hall ID: %d", placeID, hallID)
	err = h.Store.DeletePlace(r.Context(), placeID, hallID, establishmentID, simulatedOwnerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "permission denied") {
			utils.RespondWithError(w, http.StatusNotFound, "Место не найдено или отказано в доступе.", nil); return
		}
		log.Printf("AdminDeletePlace: Error deleting place ID %d: %v", placeID, err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Не удалось удалить место.", nil); return
	}
	log.Printf("AdminDeletePlace: Deleted place ID: %d", placeID)
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Place deleted successfully"})
}