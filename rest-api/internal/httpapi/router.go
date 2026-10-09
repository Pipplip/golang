package httpapi

import (
	"log/slog"
	"net/http"

	"rest-api/internal/book"
)

// Der Handler ist die Übersetzungsschicht zwischen HTTP und dem Service.
// HTTP liefert Dinge wie URL, Methode, JSON Body...
// Der Service liefert Dinge wie Book, []Book, Fehler...
// Der Handler verbindet beides.
type Handler struct {
	service *book.Service
}

func NewRouter(service *book.Service) *http.ServeMux {
	mux := http.NewServeMux()
	h := &Handler{service: service}

	// Methodenspezifisches Routing ab Go 1.22: Methode + Pfad in einem Pattern.
	// Verbinden einer Route mit einer Handlerfunktion. Die Handlerfunktion ist eine Methode des Handlers, die die Business-Logik aufruft und die HTTP-Antwort erstellt.
	mux.HandleFunc("GET /books", h.listBooks)
	mux.HandleFunc("GET /books/{id}", h.getBook)
	mux.HandleFunc("POST /books", h.createBook)
	mux.HandleFunc("PUT /books/{id}", h.updateBook)
	mux.HandleFunc("DELETE /books/{id}", h.deleteBook)
	mux.HandleFunc("GET /{$}", HomeHandler)

	return mux
}

func Chain(logger *slog.Logger, next http.Handler) http.Handler {
	// Reihenfolge ist absichtlich: Logging(Recovery(Router)).
	// LoggingMiddleware misst die Dauer und protokolliert nach der Anfrage Methode, Pfad und Statuscode. Ein kleiner statusRecorder merkt sich dafür, welchen HTTP-Status die Handler geschrieben haben
	// RecoveryMiddleware fängt unerwartete Panics ab, protokolliert den Fehler samt Stacktrace und sendet, wenn noch keine Antwort begonnen hat, eine HTTP-500-Antwort
	// Bei einer Anfrage läuft zuerst das Logging, dann Recovery und zuletzt der Router
	return LoggingMiddleware(logger)(
		RecoveryMiddleware(logger)(
			CORSMiddleware(next),
		),
	)
}
