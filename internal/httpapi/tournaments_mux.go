package httpapi

import "net/http"

func NewTournamentsMux(h *TournamentsHandlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tournaments", h.List)
	mux.HandleFunc("GET /tournaments/{tournamentID}", h.Get)
	mux.HandleFunc("POST /tournaments", h.Create)
	mux.HandleFunc("PATCH /tournaments/{tournamentID}", h.Edit)
	mux.HandleFunc("DELETE /tournaments/{tournamentID}", h.Delete)
	mux.HandleFunc("PUT /tournaments/{tournamentID}/participants/me", h.Join)
	mux.HandleFunc("DELETE /tournaments/{tournamentID}/participants/me", h.Leave)
	mux.HandleFunc("PATCH /tournaments/{tournamentID}/participants/me", h.UpdateParticipation)
	mux.HandleFunc("GET /tournaments/{tournamentID}/participants", h.ListParticipants)
	mux.HandleFunc("GET /tournaments/{tournamentID}/standings", h.GetStandings)

	return mux
}
