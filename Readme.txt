IDE und Go vorbereiten:
1) Go unter go.dev downloaden und installieren
2) Visual Studio Code installieren und "Go" Extension installieren
3) Neuen Workspace anlegen und Directory in VSC laden
4) Command Pallete öffnen in VSC (view -> command Palette oder Str+Shift+P)
und nach "Go: install/update Tools suchen"
Alle Tools selektieren und installieren

Übersicht über Standardpackages:
pkg.go.dev/std

Gutes GitHub: https://github.com/avelino/awesome-go

Workspace anlegen:
1) Neuen Workspace im Explorer oder Konsole erstellen (mkdir BasisProject)
2) Ins Verzeichnis gehen (cd BasisProject) und Haupt-Modul initialisieren (go mod init BasisProject)
-> Dadurch entsteht eine go.mod Datei
3) main.go Datei erstellen (Im Explorer oder VSC (Über File - new File) mit Name main.go) - Import package main, wenn dort die main Funktion sein soll
4) Programm ausführen: In Konsole go run ./main.go

Neues Package anlegen:
1) Erstelle im Workspace einen neuen Ordner (mkdir myOverview)
2) Erstelle darin eine neue Datei overview.go

Go code is grouped into packages, and packages are grouped into modules. Your module specifies dependencies needed to run your code, including the Go version and the set of other modules it requires. 

Which Subjects Are Go Relevant For?
    Backend Development:
    	Go is excellent for building server-side applications.
    Cloud Computing:
    	Go is widely used in cloud infrastructure.
    System Programming:
    	Go provides low-level system access.
    Microservices:
    	Go excels at building microservices.
    DevOps:
    	Go is popular for DevOps tooling.
    Network Programming:
    	Go has strong networking capabilities.
    Concurrent Programming:
    	Go makes concurrent programming simple.
