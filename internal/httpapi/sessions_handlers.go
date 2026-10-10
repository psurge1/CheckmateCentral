package httpapi

import "net/http"

type SessionsHandlers struct{}

func (h *SessionsHandlers) Login(w http.ResponseWriter, r *http.Request) {
	// TODO: Authenticate the account and create a session.
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *SessionsHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	// TODO: Invalidate the current session.
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
