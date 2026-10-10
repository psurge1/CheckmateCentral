package httpapi

import "net/http"

func NewGamesMux(h *GamesHandlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /games", h.List)
	mux.HandleFunc("GET /games/{gameID}", h.Get)
	mux.HandleFunc("GET /games/{gameID}/moves", h.ListMoves)
	mux.HandleFunc("POST /games/{gameID}/moves", h.SubmitMove)
	mux.HandleFunc("POST /games/{gameID}/resign", h.Resign)
	mux.HandleFunc("GET /games/{gameID}/stream", h.Stream)

	return mux
}
