## Allgemein Context

Der Grundgedanke eines Context ist, dass mehrere Funktionen oder Goroutines, die zu derselben Aufgabe gehören, einen gemeinsamen Kontext teilen.

Bsp:   
Stell dir vor, ein Webserver erhält eine HTTP-Anfrage:
- Der Handler ruft eine Datenbank auf.
- Die Datenbankfunktion ruft einen externen Service auf.
- Mehrere Goroutines arbeiten parallel.

Wenn der Client die Verbindung schließt oder ein Timeout erreicht wird, sollen alle diese Operationen beendet werden.   
Damit nicht unnötig Rechenleistung blockiert wird.   
Genau dafür gibt es context.

Ein Context wird typischerweise von der aufrufenden Funktion erstellt und an alle weiteren Funktionen weitergereicht.

Funktionen (erstellen immer einen neuen Sub-Context):
- Context = Lebenssteuerung einer Aufgabe
- WithCancel    → ich kann abbrechen (cancel() muss manuell aufgerufen werden)
- WithTimeout   → ich werde nach Zeit beendet (relativ: jetzt + 2s)
- WithDeadline  → ich habe ein Enddatum (absolut oder relativ)
- WithValue     → ich trage Daten mit (nur für Request-scoped Daten (UserID, RequestID) — keine Funktionsparameter ersetzen!)

**Hierarchie:**  
In main wird meist der Root Context initialisiert: `ctx := context.Background()`   
Alle weiteren Context werden davon abgeleitet.   
Vorteil: Wenn der parent Context beendet wird, werden alle abgeleiteten Context ebenfalls beendet.   

## Best practices

### 1. Context nicht in struct speichern

**falsch:**
```go
type Service struct {
    ctx context.Context
}
```
**richtig:**
```go
type Service struct{}
func (s Service) Do(ctx context.Context) {
    // ctx wird als Parameter übergeben, nicht gespeichert
}
```

### 2. cancel() bei withCancel oder withTimeout explizit aufrufen

**richtig:**
```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel() // sonst entstehen memory leaks
```

### 3. Context weiterreichen

**falsch:**
```go
func handler(w http.ResponseWriter, r *http.Request) {
    go doWork() // kein ctx! Kein Abbruch möglich
}
```

**richtig:**
```go
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    go doWork(ctx)
}
```

### 4. context.Background() nicht überall verwenden

**falsch:**
```go
go queryDatabase(context.Background())
```
**richtig:**
```go
go queryDatabase(r.Context())
```

### 5. WithValue nicht missbrauchen

**falsch:**
```go
// wird schnell zur Müllhalde
ctx = context.WithValue(ctx, "db", db)
ctx = context.WithValue(ctx, "logger", logger)
```
**richtig:**
```go
// nur für Request Metadata
ctx = context.WithValue(ctx, userKey{}, userID)
```

### 6. Context nicht ignorieren

**falsch:**
```go
// goroutine läuft ewig
// kein Stop bei cancel/timeout
for {
	doWork()
}
```
**richtig:**
```go
for {
	select {
	case <-ctx.Done():
		return
	default:
		doWork()
	}
}
```

### 7. Parent Context nicht ignorieren

**falsch:**
```go
// immer neuer Context anstatt bestehenden verwenden
ctx, _ := context.WithTimeout(context.Background(), time.Second)
```
**richtig:**
```go
ctx, cancel := context.WithTimeout(parentCtx, time.Second)
defer cancel()
```
