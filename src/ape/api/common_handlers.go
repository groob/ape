package api

import (
	"ape/models"
	"encoding/json"
	"net/http"
)

func respondOK(rw http.ResponseWriter, body models.Viewer) {
	rw.Header().Add("Content-Type", "application/json")
	respond(rw, body, http.StatusOK)
}

func respondError(rw http.ResponseWriter, status int, err error) {
	// TODO accept variadic error param
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)

	resp := &ErrorResponse{Errors: make([]string, 0, 1)}
	if err != nil {
		resp.Errors = append(resp.Errors, err.Error())
	}

	if err := json.NewEncoder(rw).Encode(resp); err != nil {
		resp.Errors = append(resp.Errors, err.Error())
	}
}

func respondCreated(rw http.ResponseWriter, body models.Viewer, location string) {
	if location != "" {
		// TODO: Set location header
	}
	respond(rw, body, http.StatusCreated)
}

func respond(rw http.ResponseWriter, body models.Viewer, status int) {
	view, err := body.View()
	switch err {
	case nil:
		break
	case models.ErrNotFound:
		respondError(rw, http.StatusNotFound, err)
		return
	default:
		respondError(rw, http.StatusInternalServerError, err)
		return
	}

	rw.WriteHeader(status)
	rw.Write(view.Data)
}

// ErrorResponse encodes errors into http response body
type ErrorResponse struct {
	Errors []string `json:"errors"`
}
