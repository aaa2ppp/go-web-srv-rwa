package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"rwa/internal/dto"
)

const (
	defaultLimit  = 20
	defaultOffset = 0
)

type httpError struct {
	Status  int
	Message string
}

func (e httpError) Error() string {
	return e.Message
}

type tErrors struct {
	Body []string `json:"body"`
}

type tErrorResponse struct {
	Errors tErrors `json:"errors"`
}

func writeError(w http.ResponseWriter, err error) {
	var (
		status  = http.StatusInternalServerError
		message string
	)

	switch {
	case errors.Is(err, dto.ErrUnauthorized):
		status = http.StatusUnauthorized
	case errors.Is(err, dto.ErrNotFound):
		status = http.StatusNotFound
	default:
		switch err := err.(type) {
		case httpError:
			status = err.Status
			message = err.Message
		case *httpError:
			status = err.Status
			message = err.Message
		}
	}

	if message == "" {
		message = http.StatusText(status)
	}

	writeJSONResponse(w, status, &tErrorResponse{
		Errors: tErrors{
			Body: []string{message},
		},
	})
}

func writeJSONResponse(w http.ResponseWriter, status int, resp interface{}) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("can't write response: %v", err)
	}
}

func getAuthorizationToken(r *http.Request) (string, error) {
	v := r.Header.Get("Authorization")
	if v == "" || !strings.HasPrefix(v, "Token") {
		return "", httpError{http.StatusUnauthorized, ""}
	}
	return strings.TrimLeft(v[5:], ": "), nil
}
