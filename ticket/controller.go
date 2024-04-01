package ticket

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jinzhu/copier"
)

type FilmService interface {
	CreateNewFilm(context.Context, Film) (*Film, error)
}

var filmSvc FilmService = nil

func filmCreateHandler(w http.ResponseWriter, r *http.Request) {
	var payload DTOFilmCreate
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Failed to parse JSON payload", http.StatusBadRequest)
		return
	}

	var req Film
	if err := copier.Copy(&req, payload); err != nil {
		http.Error(w, "Failed to parse FILM model from payload", http.StatusInternalServerError)
		return
	}

	result, err := filmSvc.CreateNewFilm(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to create new film", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Failed to parse film JSON payload", http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(resp); err != nil {
		http.Error(w, "Failed to write response data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
