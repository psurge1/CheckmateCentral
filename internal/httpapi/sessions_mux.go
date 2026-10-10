package httpapi

import "net/http"

func NewSessionsMux(h *SessionsHandlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /sessions", h.Login)
	mux.HandleFunc("DELETE /sessions/current", h.Logout)

	return mux
}
