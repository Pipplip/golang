package main

// Interfaces beschreiben Verhalten
// Wenn etwas diese Methode hat, dann gehört es dazu
// Nicht wie in Java, wo man implements braucht und den Vertrag damit sozusagen unterschreibt
// Go sagt einfach. Der Typ passt zufällig, weil er die Methode implementiert.

// Java: Dog → unterschreibt Vertrag (implements)
// Go: Dog → passt einfach zufällig in den Vertrag rein

type Speaker interface {
	Speak() string
}

type Dog struct{}
type Cat struct{}

func (d Dog) Speak() string {
	return "Wau"
}

func saySomething(s Speaker) {
	println(s.Speak())
}

func main() {
	saySomething(Dog{}) // Der Dog kann Speak()en, deshalb ist er ein Speaker
	//saySomething(Cat{})
}
