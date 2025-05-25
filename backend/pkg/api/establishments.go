package api

// Этот файл предназначен для обработчиков API, связанных с заведениями (establishments).
// Например, получение списка заведений, информации о конкретном заведении,
// а также CRUD-операции для управления заведениями (для админ-панели).

// import (
// 	"net/http"
// 	"minityweb/backend/pkg/models"
// 	"minityweb/backend/pkg/utils"
// 	"strconv" // Для преобразования ID из URL или query params
// 	// "encoding/json" // Для декодирования тела запроса при создании/обновлении
// 	// "github.com/gorilla/mux" // Если используете mux для извлечения параметров из URL
// 	// "errors" // Для errors.Is
// 	// "database/sql" // Для sql.ErrNoRows
// 	// "log"
// )

/*
Примеры заглушек для обработчиков (некоторые из них уже могут быть реализованы в handlers.go):

// GetPublicEstablishmentsHandler (для публичного API)
// Логика, похожая на h.GetEstablishments из handlers.go, может быть здесь.
func (h *APIHandler) GetPublicEstablishmentsHandler(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// queryParams := r.URL.Query()

	// limitStr := queryParams.Get("limit")
	// offsetStr := queryParams.Get("offset")
	// typeStr := queryParams.Get("type")

	// limit, err := strconv.Atoi(limitStr)
	// if err != nil || limit <= 0 {
	// 	limit = 10
	// }
	// if limit > 100 {
	// 	limit = 100
	// }

	// offset, err := strconv.Atoi(offsetStr)
	// if err != nil || offset < 0 {
	// 	offset = 0
	// }

	// filters := make(map[string]interface{})
	// if typeStr != "" {
	// 	validType := models.EstablishmentType(typeStr)
	// 	if validType == models.RestaurantEstablishment || validType == models.CoworkingEstablishment {
	// 		filters["type"] = validType
	// 	} else {
	// 		utils.RespondWithError(w, http.StatusBadRequest, "Invalid establishment type filter", nil)
	// 		return
	// 	}
	// }

	// establishments, err := h.Store.GetEstablishments(ctx, limit, offset, filters)
	// if err != nil {
	// 	log.Printf("Error fetching establishments: %v", err)
	// 	utils.RespondWithError(w, http.StatusInternalServerError, "Could not fetch establishments", nil)
	// 	return
	// }

	// if establishments == nil {
	// 	establishments = []*models.Establishment{}
	// }
	// utils.RespondWithJSON(w, http.StatusOK, establishments)
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "GetPublicEstablishmentsHandler (Not Implemented Yet, see handlers.go)"})
}

// GetPublicEstablishmentByIDHandler (для публичного API)
// Логика, похожая на h.GetEstablishmentByID из handlers.go, может быть здесь.
func (h *APIHandler) GetPublicEstablishmentByIDHandler(w http.ResponseWriter, r *http.Request) {
	// // Пример для gorilla/mux:
	// // vars := mux.Vars(r)
	// // idStr, ok := vars["id"]
	// // if !ok { ... }

	// // Пример для query param:
	// idStr := r.URL.Query().Get("id")
	// if idStr == "" {
	// 	utils.RespondWithError(w, http.StatusBadRequest, "Establishment ID is missing", nil)
	// 	return
	// }
	// id, err := strconv.ParseInt(idStr, 10, 64)
	// if err != nil {
	// 	utils.RespondWithError(w, http.StatusBadRequest, "Invalid establishment ID format", nil)
	// 	return
	// }

	// establishment, err := h.Store.GetEstablishmentByID(r.Context(), id)
	// if err != nil {
	// 	if errors.Is(err, sql.ErrNoRows) {
	// 		utils.RespondWithError(w, http.StatusNotFound, "Establishment not found", nil)
	// 		return
	// 	}
	// 	log.Printf("Error fetching establishment by ID %d: %v", id, err)
	// 	utils.RespondWithError(w, http.StatusInternalServerError, "Could not fetch establishment", nil)
	// 	return
	// }
	// if establishment == nil {
	// 	utils.RespondWithError(w, http.StatusNotFound, "Establishment not found", nil)
	// 	return
	// }
	// utils.RespondWithJSON(w, http.StatusOK, establishment)
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "GetPublicEstablishmentByIDHandler (Not Implemented Yet, see handlers.go)"})
}

// AdminCreateEstablishmentHandler (для админ-панели)
func (h *APIHandler) AdminCreateEstablishmentHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать создание заведения администратором
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "AdminCreateEstablishmentHandler (Not Implemented Yet)"})
}
*/