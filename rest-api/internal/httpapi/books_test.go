package httpapi

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"rest-api/internal/book"
	"rest-api/internal/storage/memory"
)

func newTestHandler() http.Handler {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	service := book.NewService(memory.NewRepository())
	return Chain(logger, NewRouter(service))
}

func TestBooksCRUDFlow(t *testing.T) {
	handler := newTestHandler()

	createBody := `{"title":"Clean Code","author":"Robert Martin","published_year":2008}`
	createReq := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	handler.ServeHTTP(createRec, createReq)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", createRec.Code)
	}

	var created book.Book
	if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if created.ID != 1 {
		t.Fatalf("expected id 1, got %d", created.ID)
	}
	if got := createRec.Header().Get("Location"); got != "/books/1" {
		t.Fatalf("unexpected location header: %q", got)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/books/1", nil)
	getRec := httptest.NewRecorder()
	handler.ServeHTTP(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", getRec.Code)
	}

	updateBody := `{"title":"Clean Code 2","author":"Robert Martin"}`
	updateReq := httptest.NewRequest(http.MethodPut, "/books/1", strings.NewReader(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateRec := httptest.NewRecorder()
	handler.ServeHTTP(updateRec, updateReq)
	if updateRec.Code != http.StatusOK {
		t.Fatalf("expected 200 on update, got %d", updateRec.Code)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/books", nil)
	listRec := httptest.NewRecorder()
	handler.ServeHTTP(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("expected 200 on list, got %d", listRec.Code)
	}
	var listed []book.Book
	if err := json.NewDecoder(listRec.Body).Decode(&listed); err != nil {
		t.Fatalf("decode list response: %v", err)
	}
	if len(listed) != 1 || listed[0].Title != "Clean Code 2" {
		t.Fatalf("unexpected list payload: %+v", listed)
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/books/1", nil)
	deleteRec := httptest.NewRecorder()
	handler.ServeHTTP(deleteRec, deleteReq)
	if deleteRec.Code != http.StatusNoContent {
		t.Fatalf("expected 204 on delete, got %d", deleteRec.Code)
	}
}

func TestBooksInputValidationAndErrors(t *testing.T) {
	handler := newTestHandler()

	noTypeReq := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(`{"title":"X","author":"Y"}`))
	noTypeRec := httptest.NewRecorder()
	handler.ServeHTTP(noTypeRec, noTypeReq)
	if noTypeRec.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("expected 415, got %d", noTypeRec.Code)
	}

	unknownFieldReq := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(`{"title":"X","author":"Y","extra":1}`))
	unknownFieldReq.Header.Set("Content-Type", "application/json")
	unknownFieldRec := httptest.NewRecorder()
	handler.ServeHTTP(unknownFieldRec, unknownFieldReq)
	if unknownFieldRec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for unknown field, got %d", unknownFieldRec.Code)
	}

	invalidIDReq := httptest.NewRequest(http.MethodGet, "/books/abc", nil)
	invalidIDRec := httptest.NewRecorder()
	handler.ServeHTTP(invalidIDRec, invalidIDReq)
	if invalidIDRec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid id, got %d", invalidIDRec.Code)
	}

	notFoundReq := httptest.NewRequest(http.MethodGet, "/books/9999", nil)
	notFoundRec := httptest.NewRecorder()
	handler.ServeHTTP(notFoundRec, notFoundReq)
	if notFoundRec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for missing book, got %d", notFoundRec.Code)
	}
}

func TestBooksRequestBodyTooLarge(t *testing.T) {
	handler := newTestHandler()

	payload := map[string]string{
		"title":  strings.Repeat("A", int(maxBodyBytes)),
		"author": "Author",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected 413, got %d", rec.Code)
	}
}

func TestBooksMethodNotAllowedOnItem(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(http.MethodPost, "/books/1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}

	allow := rec.Header().Get("Allow")
	if allow == "" {
		t.Fatal("expected Allow header")
	}
	for _, method := range []string{http.MethodGet, http.MethodPut, http.MethodDelete} {
		if !strings.Contains(allow, method) {
			t.Fatalf("allow header %q does not include %s", allow, method)
		}
	}
}

func TestBooksHeadUsesGetRoute(t *testing.T) {
	handler := newTestHandler()

	createReq := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(`{"title":"Go","author":"Team"}`))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	handler.ServeHTTP(createRec, createReq)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", createRec.Code)
	}

	headReq := httptest.NewRequest(http.MethodHead, "/books/"+strconv.Itoa(1), nil)
	headRec := httptest.NewRecorder()
	handler.ServeHTTP(headRec, headReq)
	if headRec.Code != http.StatusOK {
		t.Fatalf("expected 200 for HEAD, got %d", headRec.Code)
	}
}
