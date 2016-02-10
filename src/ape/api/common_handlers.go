package api

import (
	"ape/models"
	"encoding/json"
	"log"
	"net/http"
)

func respondOK(rw http.ResponseWriter, body models.Viewer) {
	rw.Header().Add("Content-Type", "application/json")
	respondStatus(rw, body, http.StatusOK)
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
		log.Println(err)
	}
}

func respondCreated(rw http.ResponseWriter, body models.Viewer, location string) {
	if location != "" {
		// TODO: Set location header
	}
	respondStatus(rw, body, http.StatusCreated)
}

func respondStatus(rw http.ResponseWriter, body models.Viewer, status int) {
	view := body.View()
	jsn, err := json.MarshalIndent(view, "", " ")
	if err != nil {
		log.Printf("respondOK failed to marshal json response with error: %v\n", err)
	}

	if len(*view) == 0 {
		rw.WriteHeader(http.StatusNoContent)
		return
	}
	rw.WriteHeader(status)
	rw.Write(jsn)
}

// ErrorResponse encodes errors into http response body
type ErrorResponse struct {
	Errors []string `json:"errors"`
}
