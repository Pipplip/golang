package httpapi

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

// Die Middleware läuft bei einem Request und bei einer Response
// In der Chain im Router werden die Middleware-Funktionen initiert und aufgerufen. Jede Middleware kann entscheiden, ob sie die Anfrage an die nächste Middleware weitergibt oder nicht.
// Normalerweise wird bei einer Response die Middleware in umgekehrter Reihenfolge aufgerufen, d.h. die letzte Middleware wird zuerst aufgerufen, dann die vorletzte usw. bis zur ersten Middleware. Danach wird die Response in umgekehrter Reihenfolge zurückgegeben, d.h. die erste Middleware wird zuerst aufgerufen, dann die zweite usw. bis zur letzten Middleware.
// Anfrage:  Logging → Recovery → Router
// Antwort:  Router → Recovery → Logging

type statusRecorder struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
}

func (w *statusRecorder) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.wroteHeader = true
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *statusRecorder) Write(data []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(data)
}

func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(rec, r)

			logger.Info(
				"http request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", rec.statusCode),
				slog.Duration("duration", time.Since(start)),
			)
		})
	}
}

func RecoveryMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rec := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
			defer func() {
				if recovered := recover(); recovered != nil {
					logger.Error(
						"panic recovered",
						slog.Any("panic", recovered),
						slog.String("stack", string(debug.Stack())),
					)

					// Nur wenn noch nichts geschrieben wurde, senden wir eine saubere JSON-Fehlerantwort.
					if !rec.wroteHeader {
						writeInternalServerError(rec)
					}
				}
			}()

			// Handler ausfuehren; bei Panic greift der defer-Block oben.
			next.ServeHTTP(rec, r)
		})
	}
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Language, Content-type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
