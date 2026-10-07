# Lernprojekt: REST-Webservice zur Buecherverwaltung in Go

## Ziel

Im Verzeichnis `rest-api/` entsteht eine vollstaendige CRUD-REST-API fuer
Buecher, um zentrale Go-Konzepte praxisnah zu lernen:

- gängige Projektstruktur mit `cmd/` und `internal/`
- Routing mit `net/http` (Go-Standardbibliothek)
- Middleware (Logging und Panic-Recovery)
- in-memory Datenspeicher statt echter Datenbank
- Tests fuer Service, Speicher und HTTP-Verhalten

## Umfang

Ressource: `books`

- `GET /books` (Liste)
- `GET /books/{id}` (Eintrag lesen)
- `POST /books` (Eintrag anlegen)
- `PUT /books/{id}` (Eintrag ersetzen)
- `DELETE /books/{id}` (Eintrag loeschen)

Datenmodell:

- `id` (serverseitig vergeben)
- `title` (Pflicht)
- `author` (Pflicht)
- `published_year` (optional, > 0)

## Projektstruktur

```text
rest-api/
  go.mod
  plan.md
  README.md
  cmd/
    api/
      main.go
  internal/
    book/
      book.go
      repository.go
      service.go
      service_test.go
    storage/
      memory/
        books.go
        books_test.go
    httpapi/
      router.go
      books.go
      json.go
      middleware.go
      books_test.go
      middleware_test.go
```

## Architektur

1. **Domain (`internal/book`)**
   Enthält Modell, Validierung, Service-Logik und fachliche Fehler.
2. **Storage (`internal/storage/memory`)**
   Threadsicherer In-Memory-Store mit `sync.RWMutex` und monoton steigenden IDs.
3. **HTTP (`internal/httpapi`)**
   Routing, Handler, JSON-Parsing, Fehlerantworten und Middleware.
4. **Entry Point (`cmd/api`)**
   Verdrahtet Komponenten und startet den Server mit Timeouts und Graceful Shutdown.

## Middleware (Lernfokus)

- **LoggingMiddleware**: protokolliert Methode, Pfad, Statuscode und Dauer.
- **RecoveryMiddleware**: faengt Panics im Request-Lauf ab und liefert kontrolliert
  `500` statt Serverabsturz.

Verkettung:

`Logging(Recovery(Router))`

So ist jeder Request geloggt, und eine Panic wird zentral behandelt.

## Nicht-Ziele

- keine Persistenz in externer Datenbank
- keine Authentifizierung
- keine Suche/Paginierung

## Abnahme

- API startet mit `go run ./cmd/api`
- alle CRUD-Endpunkte funktionieren
- klare Fehlercodes (`400`, `404`, `413`, `415`, `500`)
- Tests laufen mit `go test ./...`
