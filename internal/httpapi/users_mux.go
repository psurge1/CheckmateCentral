package httpapi

import "net/http"

func NewUsersMux(h *UsersHandlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /users", h.Register)
	mux.HandleFunc("GET /users", h.List)
	mux.HandleFunc("GET /users/me", h.GetMe)
	mux.HandleFunc("PATCH /users/me", h.EditMe)
	mux.HandleFunc("GET /users/{userID}", h.Get)

	return mux
}
