package httpapi

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"rest-api/internal/book"
)

func (h *Handler) listBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.service.List(r.Context())
	if err != nil {
		writeInternalServerError(w)
		return
	}
	// API-Vertrag: Leere Listen als [] statt null ausgeben.
	if books == nil {
		books = []book.Book{}
	}

	writeJSON(w, http.StatusOK, books)
}

func (h *Handler) getBook(w http.ResponseWriter, r *http.Request) {
	// PathValue liest den {id}-Teil direkt aus dem ServeMux-Routenmuster.
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "id must be a positive integer")
		return
	}

	found, err := h.service.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, book.ErrNotFound) {
			writeAPIError(w, http.StatusNotFound, "book not found")
			return
		}
		writeInternalServerError(w)
		return
	}

	writeJSON(w, http.StatusOK, found)
}

func (h *Handler) createBook(w http.ResponseWriter, r *http.Request) {
	if err := requireJSONContentType(r); err != nil {
		writeAPIError(w, http.StatusUnsupportedMediaType, err.Error())
		return
	}

	var input book.Input
	if err := decodeJSONBody(w, r, &input); err != nil {
		writeAPIError(w, err.status, err.message)
		return
	}

	created, err := h.service.Create(r.Context(), input)
	if err != nil {
		var validationErr *book.ValidationError
		if errors.As(err, &validationErr) {
			writeAPIError(w, http.StatusBadRequest, validationErr.Error())
			return
		}
		writeInternalServerError(w)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/books/%d", created.ID))
	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) updateBook(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "id must be a positive integer")
		return
	}
	if err := requireJSONContentType(r); err != nil {
		writeAPIError(w, http.StatusUnsupportedMediaType, err.Error())
		return
	}

	var input book.Input
	if err := decodeJSONBody(w, r, &input); err != nil {
		writeAPIError(w, err.status, err.message)
		return
	}

	updated, err := h.service.Update(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, book.ErrNotFound) {
			writeAPIError(w, http.StatusNotFound, "book not found")
			return
		}
		var validationErr *book.ValidationError
		if errors.As(err, &validationErr) {
			writeAPIError(w, http.StatusBadRequest, validationErr.Error())
			return
		}
		writeInternalServerError(w)
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) deleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "id must be a positive integer")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, book.ErrNotFound) {
			writeAPIError(w, http.StatusNotFound, "book not found")
			return
		}
		writeInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	if id <= 0 {
		return 0, errors.New("id must be positive")
	}
	return id, nil
}
