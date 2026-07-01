package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

/*
Der Grundgedanke eines Context ist,
dass mehrere Funktionen oder Goroutines, die zu derselben Aufgabe gehören, einen gemeinsamen Kontext teilen.

Bsp:
Stell dir vor, ein Webserver erhält eine HTTP-Anfrage:
- Der Handler ruft eine Datenbank auf.
- Die Datenbankfunktion ruft einen externen Service auf.
- Mehrere Goroutines arbeiten parallel.

Wenn der Client die Verbindung schließt oder ein Timeout erreicht wird, sollen alle diese Operationen beendet werden.
Damit nicht unnötig Rechenleistung blockiert wird.
Genau dafür gibt es context.

Ein Context wird typischerweise von der aufrufenden Funktion erstellt und an alle weiteren Funktionen weitergereicht.

Funktionen:
Context = Lebenssteuerung einer Aufgabe
WithCancel    → ich kann abbrechen (cancel() muss manuell aufgerufen werden)
WithTimeout   → ich werde nach Zeit beendet (relativ: jetzt + 2s)
WithDeadline  → ich habe ein Enddatum (absolut: time.Now().Add(2s) ist dasselbe wie Timeout — der Unterschied ist nur wie man die Zeit angibt)
WithValue     → ich trage Daten mit (nur für Request-scoped Daten (UserID, RequestID) — keine Funktionsparameter ersetzen!)

Hierarchie:
In main wird meist der Root Context initialisiert: ctx := context.Background()
Alle weiteren Context werden davon abgeleitet.
Vorteil: Wenn der parent Context beendet wird, werden alle abgeleiteten Context ebenfalls beendet.

Hier im Beispiel-Code:
Background()
    │
    └── WithValue(userID="philipp")      ← main()
            │
            └── WithTimeout(2s)          ← fetchAll()
                    │
                    ├── Goroutine 1: queryDB("users",  1s) ✓
                    └── Goroutine 2: queryDB("orders", 3s) ✗
*/

type contextKey string

const keyUserID contextKey = "userID"

func main() {
	ctx := context.Background()                        // root Context erstellen
	ctx = context.WithValue(ctx, keyUserID, "philipp") // neuer Context mit Wert-Speicher erstellen

	fetchAll(ctx)
	/*
		Ausgabe:
		[db:orders] Query für 'philipp' gestartet...
		[db:users] Query für 'philipp' gestartet...
		[db:users] ✓ Fertig
		[db:orders] ✗ Abgebrochen: context deadline exceeded
		Alle Goroutinen beendet
	*/
}

// fetchAll setzt den Timeout und startet alle Queries parallel
// Insgesamt darf fetchAll nur 2 Sekunden dauern (context.WithTimeout)
func fetchAll(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second) // neuer context erstellen
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Starte Goroutine 1")
		queryDB(ctx, "users", 1*time.Second)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Starte Goroutine 2")
		queryDB(ctx, "orders", 3*time.Second)
	}()

	wg.Wait()
	fmt.Println("\nAlle Goroutinen beendet")
}

// queryDB simuliert eine Datenbankabfrage mit konfigurierbarer Latenz
func queryDB(ctx context.Context, table string, latency time.Duration) error {
	userID, ok := ctx.Value(keyUserID).(string)
	if !ok {
		fmt.Printf("[db:%s] kein UserID im Context\n", table)
		return fmt.Errorf("kein UserID im Context")
	}

	fmt.Printf("[db:%s] Query für '%s' gestartet...\n", table, userID)

	select {
	case <-time.After(latency): // Goroutine 1: 1s < 2s → kommt hier an
		fmt.Printf("[db:%s] ✓ Fertig\n", table)
		return nil

	case <-ctx.Done(): // Goroutine 2: 3s > 2s → kommt hier an ✗
		fmt.Printf("[db:%s] ✗ Abgebrochen: %v\n", table, ctx.Err())
		return ctx.Err()
	}
}
