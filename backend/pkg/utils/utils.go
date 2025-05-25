package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"minityweb/backend/pkg/models" // ИЗМЕНЕНО
)

// HashPassword хэширует пароль с использованием bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash проверяет, соответствует ли предоставленный пароль хэшу
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// RespondWithError отправляет JSON-ответ с ошибкой
func RespondWithError(w http.ResponseWriter, code int, message string, details map[string]string) {
	errResp := models.ErrorResponse{
		Error:   message,
		Details: details,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(errResp); err != nil {
		log.Printf("Error encoding error response: %v", err)
		http.Error(w, `{"error":"Failed to encode error response"}`, http.StatusInternalServerError)
	}
}

// RespondWithJSON отправляет JSON-ответ
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to marshal JSON response", nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		log.Printf("Error writing JSON response: %v", err)
	}
}