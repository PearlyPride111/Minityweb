package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"minityweb/backend/pkg/models" // Убедитесь, что путь к моделям корректен
	"minityweb/backend/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware проверяет JWT токен в заголовке Authorization.
// Если токен валиден, он извлекает UserID и Role и добавляет их в контекст запроса.
func (h *APIHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Authorization header required", nil)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid Authorization header format (expected Bearer token)", nil)
			return
		}

		tokenString := parts[1]
		claims := &Claims{} // Claims определены в handlers.go (или можно вынести в models.go)

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Убедимся, что метод подписи тот, который мы ожидаем (HS256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return h.JWTSecret, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token signature", nil)
				return
			}
			// Другие ошибки парсинга или истекший токен
			utils.RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Invalid token: %v", err.Error()), nil)
			return
		}

		if !token.Valid {
			utils.RespondWithError(w, http.StatusUnauthorized, "Token is not valid", nil)
			return
		}

		// Токен валиден, добавляем информацию о пользователе в контекст запроса
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		ctx = context.WithValue(ctx, "userRole", claims.Role)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminRequiredMiddleware - это обертка вокруг AuthMiddleware, которая дополнительно проверяет,
// что у пользователя роль 'admin'.
func (h *APIHandler) AdminRequiredMiddleware(next http.Handler) http.Handler {
	return h.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userRole, ok := r.Context().Value("userRole").(models.UserRole)
		if !ok || userRole != models.RoleAdmin {
			utils.RespondWithError(w, http.StatusForbidden, "Access denied: Admin role required", nil)
			return
		}
		next.ServeHTTP(w, r)
	}))
}

// UserRequiredMiddleware - это обертка вокруг AuthMiddleware, которая проверяет,
// что пользователь аутентифицирован (любая роль, кроме полностью анонимного).
// В нашем случае AuthMiddleware уже это делает, но это для примера, если бы были разные уровни.
// func (h *APIHandler) UserRequiredMiddleware(next http.Handler) http.Handler {
// 	return h.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// userID, ok := r.Context().Value("userID").(int64)
// 		// if !ok || userID == 0 { // Простая проверка, что userID есть в контексте
// 		// 	utils.RespondWithError(w, http.StatusUnauthorized, "User authentication required", nil)
// 		// 	return
// 		// }
// 		next.ServeHTTP(w, r) // AuthMiddleware уже проверил токен
// 	}))
// }