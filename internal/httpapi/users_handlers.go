package httpapi

import "net/http"

type UsersHandlers struct{}

func (h *UsersHandlers) Register(w http.ResponseWriter, r *http.Request) {
	// TODO: Register an account.
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *UsersHandlers) List(w http.ResponseWriter, r *http.Request) {
	// TODO: List public player profiles with pagination.
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *UsersHandlers) GetMe(w http.ResponseWriter, r *http.Request) {
	// TODO: Get the signed-in account.
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *UsersHandlers) EditMe(w http.ResponseWriter, r *http.Request) {
	// TODO: Edit the signed-in account.
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *UsersHandlers) Get(w http.ResponseWriter, r *http.Request) {
	// TODO: Get the public profile identified by r.PathValue("userID").
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
