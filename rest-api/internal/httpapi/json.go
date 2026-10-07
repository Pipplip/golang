package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
)

const maxBodyBytes int64 = 1 << 20

type apiError struct {
	status  int
	message string
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeAPIError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

func writeInternalServerError(w http.ResponseWriter) {
	writeAPIError(w, http.StatusInternalServerError, "internal server error")
}

func requireJSONContentType(r *http.Request) error {
	raw := r.Header.Get("Content-Type")
	if raw == "" {
		return errors.New("content-type must be application/json")
	}

	mediaType, _, err := mime.ParseMediaType(raw)
	if err != nil {
		return errors.New("invalid content-type header")
	}
	if mediaType != "application/json" {
		return errors.New("content-type must be application/json")
	}
	return nil
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) *apiError {
	// Hard limit fuer den Request-Body, um Speicherverbrauch kontrolliert zu halten.
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	dec := json.NewDecoder(r.Body)
	// Unbekannte Felder werden als Fehler behandelt, damit Tippfehler sichtbar sind.
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		var maxBytesErr *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxErr):
			return &apiError{status: http.StatusBadRequest, message: "request body contains invalid JSON"}
		case errors.Is(err, io.ErrUnexpectedEOF):
			return &apiError{status: http.StatusBadRequest, message: "request body contains invalid JSON"}
		case errors.As(err, &typeErr):
			msg := fmt.Sprintf("request body has wrong type for field %q", typeErr.Field)
			return &apiError{status: http.StatusBadRequest, message: msg}
		case errors.As(err, &maxBytesErr):
			return &apiError{status: http.StatusRequestEntityTooLarge, message: "request body must not be larger than 1 MiB"}
		case errors.Is(err, io.EOF):
			return &apiError{status: http.StatusBadRequest, message: "request body must not be empty"}
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			msg := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return &apiError{status: http.StatusBadRequest, message: "request body contains unknown field " + msg}
		default:
			return &apiError{status: http.StatusBadRequest, message: "request body is not valid"}
		}
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		// Zweiter Decode darf nur EOF liefern, sonst kamen mehrere JSON-Werte.
		return &apiError{status: http.StatusBadRequest, message: "request body must contain only one JSON object"}
	}

	return nil
}
