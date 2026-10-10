package httpapi

import "net/http"

type TournamentsHandlers struct{}

func (h *TournamentsHandlers) List(w http.ResponseWriter, r *http.Request) {
	// TODO: List tournaments.
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *TournamentsHandlers) Get(w http.ResponseWriter, r *http.Request) {
	// TODO: Get the tournament identified by r.PathValue("tournamentID").
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *TournamentsHandlers) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: Create a tournament for the signed-in account.
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *TournamentsHandlers) Edit(w http.ResponseWriter, r *http.Request) {
	// TODO: Edit the tournament.
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *TournamentsHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	// TODO: Delete the tournament.
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *TournamentsHandlers) Join(w http.ResponseWriter, r *http.Request) {
	// TODO: Join the tournament as the signed-in player.
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *TournamentsHandlers) Leave(w http.ResponseWriter, r *http.Request) {
	// TODO: Leave the tournament before it starts.
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *TournamentsHandlers) UpdateParticipation(w http.ResponseWriter, r *http.Request) {
	// TODO: Pause, resume, or withdraw the signed-in participant.
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *TournamentsHandlers) ListParticipants(w http.ResponseWriter, r *http.Request) {
	// TODO: List tournament participants.
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *TournamentsHandlers) GetStandings(w http.ResponseWriter, r *http.Request) {
	// TODO: Get tournament scores and rankings.
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
