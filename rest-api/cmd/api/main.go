package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rest-api/internal/book"
	"rest-api/internal/httpapi"
	"rest-api/internal/storage/memory"
)

func main() {
	// Es werden golang -eigene Bibliotheken verwendet
	// Alternativen wären z.B. logrus, zap, chi, gorilla/mux, echo, gin, fiber, etc.
	// gorilla/mux: HTTP-Router für URL-Muster, Parameter und verschachtelte Routen.
	// logrus: Bibliothek für strukturierte Logs mit Feldern und verschiedenen Ausgabeformaten.
	// negroni: Middleware-Framework, mit dem sich HTTP-Handler wie Logging oder Recovery verketten lassen.
	// github.com/rs/cors: Middleware für Cross-Origin Resource Sharing (CORS), um den Zugriff von anderen Domains zu erlauben.
	if err := run(); err != nil {
		slog.Error("server terminated with error", slog.Any("error", err))
		os.Exit(1)
	}
}

func run() error {
	// Verbinden der Komponenten: Logger, Repository, Service, HTTP-Router und HTTP-Server.
	// Slog ist ein strukturierter Logger von go, der hier auf die Standardausgabe schreibt.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// Manuelles Dependency Injection: Store -> Service -> HTTP-Router.
	repo := memory.NewRepository()           // Speicher für Bücher im Arbeitsspeicher. Implementiert das Repository Interface.
	service := book.NewService(repo)         // enthält die Business-Logik für Bücher.
	router := httpapi.NewRouter(service)     // HTTP-Router, der die Endpunkte für die REST-API bereitstellt.
	handler := httpapi.Chain(logger, router) // handler verbindet Router mit HTTP-Verarbeitung, hier mit Einbezug des Loggers

	// Server konfigurieren
	server := &http.Server{
		Addr:              ":8080",
		Handler:           handler, // Anfragen werden an den Handler weitergeleitet, der die Anfragen verarbeitet.
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Server in einem separaten Goroutine starten, um nicht den Hauptthread zu blockieren.
	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	// Warten auf Serverfehler oder auf ein Signal zum Beenden (z.B. Ctrl+C).
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
	}

	// Graceful Shutdown: laufende Requests duerfen innerhalb des Timeouts fertiglaufen.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	err := <-errCh
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
