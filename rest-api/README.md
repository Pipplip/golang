# REST API Lernprojekt (Buecherverwaltung)

Dieses Projekt zeigt eine kleine, aber vollstaendige REST-API in Go mit
Standardbibliothek (`net/http`) und In-Memory-Speicher.

## Lernziele

- gängige Go-Projektstruktur mit `cmd/` und `internal/`
- Routing mit Go 1.22+ ServeMux-Mustern (`GET /books/{id}`)
- Middleware verstehen und verketten
- saubere Trennung von Domain, HTTP-Schicht und Speicher
- threadsicherer In-Memory-Store ohne echte Datenbank

## Projektstruktur

```text
rest-api/
  cmd/api/main.go                  # Einstiegspunkt, Serverstart, Shutdown
  internal/book/                   # Domainmodell, Validierung, Service, Repository-Interface
  internal/storage/memory/         # threadsicherer In-Memory-Store - implementiert Repository-Interface
  internal/httpapi/                # Routing, Handler, JSON, Middleware
```

## Datenmodell

```json
{
  "id": 1,
  "title": "Clean Code",
  "author": "Robert Martin",
  "published_year": 2008
}
```

- `id`: serverseitig vergeben
- `title`: Pflichtfeld
- `author`: Pflichtfeld
- `published_year`: optional, > 0

## Routing

Das Routing laeuft mit `http.NewServeMux()` und methodenspezifischen Mustern:

- `GET /books`
- `GET /books/{id}`
- `POST /books`
- `PUT /books/{id}`
- `DELETE /books/{id}`

`{id}` wird mit `r.PathValue("id")` gelesen und in eine positive Ganzzahl geparst.
Go behandelt `HEAD` automatisch ueber passende `GET`-Routen.

## Middleware erklaert

Eine Middleware in Go ist vereinfacht gesagt eine Funktion, die vor und/oder nach deinem eigentlichen HTTP-Handler ausgeführt wird.  
Die Middleware kann Anfragen modifizieren, zusätzliche Logik ausführen oder Fehler abfangen, bevor die Anfrage an den eigentlichen Handler weitergeleitet wird.  

Sie sitzt also zwischen dem HTTP-Server und deinem Handler und kann so das Verhalten der Anfrage beeinflussen.

**Beispiel:** Logging, Authentifizierung, CORS, Rate Limiting, Recovery von Panics.

Die Middleware ist als klassische Go-Decorator-Funktion aufgebaut:

```go
func(http.Handler) http.Handler
```

Verkettung in diesem Projekt:

`Logging(Recovery(Router))`

1. **LoggingMiddleware**
   - misst Dauer und protokolliert Methode, Pfad und Statuscode
2. **RecoveryMiddleware**
   - faengt Panics im Request-Lauf ab
   - liefert, falls noch nichts geschrieben wurde, einen kontrollierten
     `500`-JSON-Fehler statt Serverabsturz

Wichtig: Recovery ersetzt keine regulaere Fehlerbehandlung, sondern ist ein
Sicherheitsnetz fuer unerwartete Panics.

## API-Endpunkte

| Methode | Pfad | Beschreibung |
|---|---|---|
| GET | `/books` | Alle Buecher lesen |
| GET | `/books/{id}` | Ein Buch lesen |
| POST | `/books` | Neues Buch anlegen |
| PUT | `/books/{id}` | Buch vollstaendig ersetzen |
| DELETE | `/books/{id}` | Buch loeschen |

## Fehlerverhalten

- `400` ungueltige ID / ungueltiges JSON / unbekannte Felder / Validierung
- `404` Buch nicht gefunden
- `413` Request-Body groesser als 1 MiB
- `415` fehlender oder falscher `Content-Type` bei POST/PUT
- `500` unerwarteter interner Fehler

Fehlerantworten sind einheitlich:

```json
{"error":"..."}
```

## Starten

```bash
cd rest-api
go run ./cmd/api
```

Server: `http://localhost:8080`

## Beispiele mit curl

Liste:

```bash
curl http://localhost:8080/books
```

Anlegen:

```bash
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{"title":"Clean Code","author":"Robert Martin","published_year":2008}'
```

Eintrag lesen:

```bash
curl http://localhost:8080/books/1
```

Ersetzen:

```bash
curl -X PUT http://localhost:8080/books/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Clean Code 2","author":"Robert Martin"}'
```

Loeschen:

```bash
curl -X DELETE http://localhost:8080/books/1
```

## Tests

```bash
go test ./...
```

## Grenzen des In-Memory-Ansatzes

- Daten gehen bei Neustart verloren.
- Kein Shared State ueber mehrere Prozesse.
- Ideal fuer Lernen und schnelle Prototypen, nicht fuer produktive Persistenz.
