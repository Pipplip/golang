package main

import (
	"fmt"
	"sync"
	"time"
)

// Main Erklärung von goroutines und Mutex
func main() {
	fmt.Println("goroutines:")
	// coroutines werden mit dem Schlüsselwort 'go' gestartet und laufen parallel zum Hauptprogramm.
	// Sie können auch mit anderen Goroutines parallel laufen.
	go printString("Hello")
	go printString("World")

	// Über channels können coroutines kommunizieren
	// Goroutine A      Channel        Goroutine B
	//    send  ----->   [  ]   ----->   receive

	// Initialisieren eines channels:
	myChannel := make(chan int) // channel, der nur ein int überträgt
	// myChannel2 := make(chan int, 2) // channel, der erst die Goroutines blockiert, wenn mehr als 2 Werte im channel sind
	// ist eine Art buffer

	// sobald etwas im channel ist, werden die coroutines geblockt

	// senden in den channel
	go func() {
		myChannel <- 10
	}()

	// empfangen vom Channel
	myListeningRoutine := <-myChannel
	fmt.Println(myListeningRoutine)

	// WaitGroups = warte, bis alle goroutines fertig sind
	// z.B. man startet 3 goroutines, main() würde aber sofort beenden, man will aber warten bis alle abgeschlossen sind
	var wg sync.WaitGroup
	wg.Add(1) // 1 heißt soviel wie: Ich starte eine Goroutine, auf die ich warten will.

	go func() {
		defer wg.Done() // signalisiert, dass diese goroutine fertig ist
		fmt.Println("Running...")
	}()

	time.Sleep(2 * time.Second)

	// -----------------------------------------------------
	// MUTEX
	fmt.Println("Mutex:")
	// Problemstellung race condition, da beide routines auf x zugreifen:
	// x := 0
	//go func() { x++ }()
	//go func() { x++ }()

	// Lösung:
	var mu sync.Mutex
	var wg2 sync.WaitGroup
	x := 0

	for i := 0; i < 2; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			mu.Lock()
			x++
			mu.Unlock()
		}()
	}
	wg2.Wait()
	fmt.Println(x)

}

func printString(value string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(value)
	}
}

// --------------------------------------------------------

// Worker: nimmt Job, gibt Ergebnis zurück
func worker(job int, results chan int, wg *sync.WaitGroup) {
	defer wg.Done() // wenn aufgabe im worker abgeschlossen ist, sage der workgroup, dass die coroutine fertig ist

	result := job * 2
	results <- result // schicke Ergebnis zurück an channel
}

// Beispiel mit Worker, waitgroups und channels
func main2() {
	// Beispiel:
	// Worker 1 ─┐
	// Worker 2 ─┼──► Channel ───► main
	// Worker 3 ─|
	// Worker 4 ─┘
	//
	//             +
	//	       WaitGroup
	//     (wartet bis alle fertig)

	jobs := []int{1, 2, 3, 4} // Aufgaben, also coroutines

	var wg sync.WaitGroup
	results := make(chan int, len(jobs)) // channel mit genug result-Ergebnisse für jede coroutine

	// Worker starten
	for _, job := range jobs {
		wg.Add(1)                    // erhöhe die Anzahl der Jobs in der waitgroup jedes mal um 1
		go worker(job, results, &wg) // jeder job wird in einer coroutine parallel ausgeführt
	}

	// Warten bis alle fertig sind
	wg.Wait()
	close(results) // schließe channel

	// Ergebnisse ausgeben
	for r := range results {
		fmt.Println(r)
	}
}
