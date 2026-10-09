# Golang

IDE (VSC) und Go vorbereiten:
1) Go unter go.dev downloaden und installieren
2) Visual Studio Code installieren und "Go" Extension installieren
3) Neuen Workspace anlegen und Directory in VSC laden
4) Command Pallete öffnen in VSC (view -> command Palette oder `Str+Shift+P`)
   und nach "Go: install/update Tools suchen"
   Alle Tools selektieren und installieren

In Intellij IDEA:  
1) Go Plugin installieren
2) Settings -> Go -> GOPATH (Global und Project Path) setzen (z.B. `C:\Users\<user>\go` bzw. `usr/local/go`)
3) Settings -> Go -> GOROOT setzen (z.B. `C:\Users\<user>\go` bzw. `usr/local/go`)
4) Settings -> Go -> Go Modules (vgo) integration aktivieren

Übersicht über Standardpackages:
`pkg.go.dev/std`

Gutes GitHub: https://github.com/avelino/awesome-go   
Offizielle Tour of go: https://go.dev/tour/list

Workspace anlegen:
1) Neuen Workspace im Explorer oder Konsole erstellen (`mkdir base-project`)
2) Ins Verzeichnis gehen (`cd base-project`) und Haupt-Modul initialisieren (`go mod init base-project`)
   -> Dadurch entsteht eine `go.mod` Datei
3) `main.go` Datei erstellen (Im Explorer oder VSC (Über File - new File) mit Name `main.go`) - Import package main, wenn dort die main Funktion sein soll
4) Programm ausführen: In Konsole `go run ./main.go`

Neues Package anlegen:
1) Erstelle im Workspace einen neuen Ordner (`mkdir myOverview`)
2) Erstelle darin eine neue Datei `overview.go`

## Wichtige Befehle

| Befehl                             | Beschreibung                                         |
|------------------------------------|------------------------------------------------------|
| `go run ./main.go` oder `go run .` | Programm ausführen                                   |
| `go build`                         | Programm kompilieren und bauen                       |
| `go test ./...`                    | Alle Tests ausführen                                 |
| `go mod tidy`                      | Ergänzt benötigte Dependencies und entfernt unnötige |
| `go get <package>`                 | Fügt ein Package als Dependency hinzu                |
| `go fmt ./...`                     | Formatiert Go-Code                                   |
| `go vet ./...`                     | Prüft auf typische Fehler                            |
| `go mod download`                  | Läft Modul Dependencies herunter                     |
| `go version`                       | Zeigt die installierte Go Version                    |
| `go list ./...`                    | listet Go-Packages im Projekt auf                    |


## Datenstrukturen

1) Slices = Listen
2) Maps = Key-Value
3) Structs = eigene Datentypen
4) Arrays = feste Listen
5) Pointer
6) Interfaces

> Slice + Map + Struct = 90% von Go

Go code ist in Packages gruppiert, und Packages sind in Modulen gruppiert. Dein Modul gibt die Abhängigkeiten an, die zum Ausführen deines Codes benötigt werden, einschließlich der Go-Version und der Menge anderer Module, die es benötigt.

Für welche Themen ist Go relevant?
- Backend Development:
Go ist eine beliebte Programmiersprache für die Backend-Entwicklung, insbesondere für Webanwendungen und APIs.
- Cloud Computing:
cloud infrastructure.
- System Programming:
Go provides low-level system access.
- Microservices:
Go excels at building microservices.
- DevOps:
Go is popular for DevOps tooling.
- Network Programming:
Go has strong networking capabilities.
- Concurrent Programming:
Go makes concurrent programming simple.
