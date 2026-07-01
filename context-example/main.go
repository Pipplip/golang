package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

/*
Hierarchy im Beispiel-Code:
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
	// withTimeout erstellt automatisch eine cancel funktion. Nach 2 Sek wird diese aufgerufen
	defer cancel() // stellt sicher, dass cancel spätestens aufgerufen wird, wenn fetchAll fertig ist

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

	select { // select = reagiert auf mehrere mögliche Ereignisse
	case <-time.After(latency): // Goroutine 1: 1s < 2s → kommt hier an
		fmt.Printf("[db:%s] ✓ Fertig\n", table)
		return nil

		// ctx.Done = Context wurde beendet
		// ctx.Err = warum wurde beendet
	case <-ctx.Done(): // Abbruchsignal
		fmt.Printf("[db:%s] ✗ Abgebrochen: %v\n", table, ctx.Err())
		return ctx.Err()
	}
}
