package api

import (
	"net/http"
	"github.com/gorilla/mux"
)

func (h *APIHandler) RegisterRoutes(router *mux.Router) {
	apiV1 := router.PathPrefix("/api/v1").Subrouter()

	apiV1.HandleFunc("/health", h.HealthCheck).Methods(http.MethodGet)

	authRouter := apiV1.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", h.RegisterUser).Methods(http.MethodPost)
	authRouter.HandleFunc("/login", h.LoginUser).Methods(http.MethodPost)

	estRouterPublic := apiV1.PathPrefix("/establishments").Subrouter()
	estRouterPublic.HandleFunc("", h.GetEstablishments).Methods(http.MethodGet)
	estRouterPublic.HandleFunc("/{id:[0-9]+}", h.GetEstablishmentByIDPublic).Methods(http.MethodGet)

	// Admin routes
	adminRouter := apiV1.PathPrefix("/admin").Subrouter()
	// adminRouter.Use(h.AdminRequiredMiddleware) 

	adminRouter.HandleFunc("/my-establishment", h.GetMyEstablishmentAdminMVP).Methods(http.MethodGet)
	adminRouter.HandleFunc("/my-establishment", h.UpdateMyEstablishmentAdminMVP).Methods(http.MethodPut)

	// Маршруты для управления залами конкретного заведения
	adminEstHallsRouter := adminRouter.PathPrefix("/establishments/{establishment_id:[0-9]+}/halls").Subrouter()
	adminEstHallsRouter.HandleFunc("", h.AdminCreateHallHandler).Methods(http.MethodPost)
	adminEstHallsRouter.HandleFunc("/{hall_id:[0-9]+}", h.AdminDeleteHallHandler).Methods(http.MethodDelete)
	adminEstHallsRouter.HandleFunc("/{hall_id:[0-9]+}", h.AdminUpdateHallHandler).Methods(http.MethodPut)

	// НОВЫЕ МАРШРУТЫ: Управление местами в конкретном зале
	adminHallPlacesRouter := adminEstHallsRouter.PathPrefix("/{hall_id:[0-9]+}/places").Subrouter()
	adminHallPlacesRouter.HandleFunc("", h.AdminCreatePlaceHandler).Methods(http.MethodPost)          // POST .../halls/{hall_id}/places
	adminHallPlacesRouter.HandleFunc("/{place_id:[0-9]+}", h.AdminUpdatePlaceHandler).Methods(http.MethodPut) // PUT .../halls/{hall_id}/places/{place_id}
	adminHallPlacesRouter.HandleFunc("/{place_id:[0-9]+}", h.AdminDeletePlaceHandler).Methods(http.MethodDelete) // DELETE .../halls/{hall_id}/places/{place_id}
    // GET для мест уже приходит с данными о заведении/зале
}