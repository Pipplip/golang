package myOverview

import (
	"fmt"
	"slices"
	"strings"
)

const NUMBER_DASHES = 20

// Damit diese Funktion in main.go aufgerufen werden kann, muss sie exportiert werden
// d.h. der Name der Funktion muss mit einem Großbuchstaben beginnen
func RunOverview() {
	fmt.Println("Hello Go")

	// Simple data types
	// Strings, Numbers (int, uint, float32, float64), Booleans, Errors
	var myBool bool = true      // Boolean
	var myInt int = 5           // Integer
	var myUint uint = 10        // Unsigned Integer - nur positive Werte
	var myFloat float32 = 3.14  // Floating point number
	var myString string = "Hi!" // String

	fmt.Println(strings.Repeat("-", NUMBER_DASHES), "Simple data types")
	fmt.Println(myBool, myInt, myUint, myFloat, myString)

	// ######################################
	// Strings

	// Initialisierungsmöglichkeiten
	var myName string = "Gogo"
	var myName2 string
	var myName3 = "Gopher"
	myName4 := "Name4"
	name, age := "Test Name", 123

	fmt.Println(strings.Repeat("-", NUMBER_DASHES), "Strings")
	fmt.Println(myName)
	fmt.Println(myName2)
	fmt.Println(myName3)
	fmt.Println(myName4)
	fmt.Println(myName3, myName4)

	fmt.Println(name, age)
	fmt.Printf("NameType: %T, NameValue: %v\n", name, name) // NameType: string, NameValue: Test Name

	// ######################################
	// Pointer
	// & (Adresse-Operator): Wird verwendet, um die Speicheradresse einer Variablen zu erhalten. Beispiel: ptr := &a setzt den Pointer ptr auf die Adresse von a
	// * (Dereferenzierungs-Operator): Wird verwendet, um auf den Wert zuzugreifen, der an der Adresse gespeichert ist, auf die der Pointer zeigt.
	// Beispiel: *ptr gibt den Wert zurück, auf den ptr zeigt, und *ptr = 20 ändert den Wert an dieser Adresse.

	a := "A"
	b := &a // b pointed auf die Adresse von a - anders initialisiert var b *string = &a

	fmt.Println(strings.Repeat("-", NUMBER_DASHES), "Pointer")
	fmt.Println(a, b, *b) // A 0xc000022100 A

	*b = "B"
	fmt.Println(a, b, *b) // B 0xc000022100 B

	a = "A2"
	fmt.Println(a, b, *b) // A2 0xc000022100 A2

	// ######################################
	// Aggregate DataTypes: Arrays, Slices, Maps, Struct

	// Arrays: wie Java
	fmt.Println(strings.Repeat("-", NUMBER_DASHES), "Arrays")
	var a1 [3]int   // array mit 3 ints
	fmt.Println(a1) // [0 0 0]

	a1 = [3]int{1, 2, 3}
	a1a := [...]int{1, 2, 3} // [...] - compiler berechnet die Größe
	fmt.Println(a1a[1])      // 1, 2, 3
	fmt.Println(a1[1])       // 2
	a1[2] = 99
	fmt.Println(a1) // [1 2 99]

	fmt.Println(len(a1)) // 3 Länge des Arrays

	a2 := [3]string{"foo", "bar", "baz"}
	a3 := a2        // arrays are copied by value
	fmt.Println(a3) // {"foo, "bar", "baz"}
	a2[0] = "q"
	fmt.Println(a2) // {"q, "bar", "baz"}
	fmt.Println(a3) // {"foo, "bar", "baz"}

	// Slices
	// Like arrays, slices are also used to store multiple values of the same type in a single variable.
	// However, unlike arrays, the length of a slice can grow and shrink as you see fit.
	// Slices are reference types, so when you assign one slice to another, they both point to the same underlying array.
	// In Go, there are several ways to create a slice:
	// (1) Using the []datatype{values} format
	// (2) Create a slice from an array
	// (3) Using the make() function

	fmt.Println(strings.Repeat("-", NUMBER_DASHES), "Slices")

	// (1) Using the []datatype{values} format
	var slice1 []int    // hier muss man nicht die Größe angeben im Gegensatz zum Array
	fmt.Println(slice1) // [] (nil)

	slice1 = []int{1, 2, 3}
	fmt.Println(slice1[1]) // 2
	// Länge und Kapazität
	fmt.Println(len(slice1), cap(slice1)) // 3

	// (2) Create a slice from an array
	arrSlice := [6]int{1, 2, 3, 4, 5, 6}
	slice2 := arrSlice[1:4]
	fmt.Println(slice2) // [2 3 4]

	// (3) Using the make() function
	slice3 := make([]int, 3, 5) // make([]type, length, capacity)
	fmt.Println(slice3)         // [0 0 0]

	// Arbeiten mit Slices: append, delete
	// mit append noch mehr Elemente hinzufügen
	slice1 = append(slice1, 5, 10, 15)
	fmt.Println(slice1) // [1 2 3 5 10 15]

	// mit Delete Elemente entfernen
	slice1 = slices.Delete(slice1, 1, 3)
	fmt.Println(slice1) // [1 5 10 15]

	// Create copy with only needed numbers
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	neededNumbers := numbers[:len(numbers)-10]
	numbersCopy := make([]int, len(neededNumbers))
	copy(numbersCopy, neededNumbers)
	fmt.Printf("numbersCopy = %v\n", numbersCopy)
	fmt.Printf("length = %d\n", len(numbersCopy))
	fmt.Printf("capacity = %d\n", cap(numbersCopy))

	// Map - key/value
	// die maps sind unordered
	fmt.Println(strings.Repeat("-", NUMBER_DASHES), "Maps")
	var m map[string]int // keys = strings, values = int
	m = map[string]int{"foo": 1, "bar": 2}
	fmt.Println(m) // map[bar:2 foo:1]

	fmt.Println(m["foo"]) // 1
	delete(m, "foo")      // remove entry

	m["baz"] = 418 // add a new value
	fmt.Println(m) // map[bar:2 baz:418]

	fmt.Println(m["foo"]) // 0 - standard int wert, da es foo nicht mehr gibt
	// Lösung
	v, ok := m["foo"]
	fmt.Println(v, ok) // 0 false

	m2 := m // copy by reference, wenn man eine Kopie haben will, benutze clone. Ansonsten wird die Referenz kopiert
	m["bar"] = 99
	fmt.Println(m)  // map[bar:99 baz:418]
	fmt.Println(m2) // map[bar:99 baz:418]

	// Make verwenden
	var m3 = make(map[string]int)
	m3["foo"] = 1
	fmt.Println(m3) // map[foo:1]

	// Prüfe ob ein Key existiert
	_, ok_m3 := m3["foo"]
	fmt.Println(ok_m3) // true

	// Structs
	// is used to create a collection of members of different data types, into a single variable.
	// While arrays are used to store multiple values of the same data type into a single variable,
	// structs are used to store multiple values of different data types into a single variable.
	fmt.Println(strings.Repeat("-", NUMBER_DASHES), "Structs")
	var s struct {
		name string
		id   int
	}
	s.name = "Arthur"
	fmt.Println(s.name) // Arthur

	// custom type
	type myStruct struct {
		name string
		id   int
	}

	var s1 myStruct // declare var with custom type
	s1 = myStruct{name: "Arthur", id: 1}
	fmt.Println(s1) // {Arthur 1}

	var s2 myStruct
	s2.name = "Marie"
	s2.id = 2
	fmt.Println(s2) // {Marie 2}

	// s3 := s2 // Eine Kopie der Werte, jedes Struct hat eine eigene Referenz

	// ######################################
	// Verzweigung
	// if else
	fmt.Println(strings.Repeat("-", NUMBER_DASHES), "Conditions")
	var option string
	var index int
	if option == "1" {
		index = 0
	} else if option == "2" {
		index = 1
	} else {
		index = 2
	}
	fmt.Println(index) // 2

	// switch
	var switchOption string = "1"
	switch switchOption {
	case "1":
		index = 0
	case "2":
		index = 1
	default:
		index = 2
	}
	fmt.Println(index) // 0

	// ######################################
	// Schleifen
	// for {...} Endlosschleife
	// for condition {...} Schleife mit Abbruchbedingung
	// for initializer; test; post clause [...} counter based Schleife
	// for key, value := range collection {...} schleife über Collections (array, slice, map)
	// for key := range collection {...} // nur keys
	// for _, value := range collection {...} // nur values
	i := 1
	// Endlosschleife
	/*for {
		fmt.Println(i)
		i += 1
	}*/

	fmt.Println(strings.Repeat("-", NUMBER_DASHES), "Loops")
	// Mit Bedingung
	for i < 3 {
		fmt.Print(i) // 1 2
		i += 1
	}
	fmt.Println()

	// Mit initializer - counter
	for i := 1; i < 3; i++ {
		fmt.Print(i) // 1 2
	}
	fmt.Println()

	// Loop über Collections
	// ----------
	arr4 := [3]int{101, 102, 103}
	for i, v := range arr4 {
		fmt.Println(i, v) // 0 101, 1 102, 3 103
	}
	for i := range arr4 {
		fmt.Println(i) // 0, 1, 2
	}
	// for _, value := range collection {...} // nur values
	for _, v := range arr4 {
		fmt.Println(v) // 101, 102, 103
	}

	// ######################################
	// Functions
	// func functionName (parameters)(return values)
	fmt.Println(strings.Repeat("-", NUMBER_DASHES), "Functions")
	greet1("Greta")
	greet2("Tik", "Tok")
	greetMultiPara("One", "Two", "Three")

	// Funktion mit einem Return type
	resultOfAdd := add(1, 2)
	fmt.Println(resultOfAdd) // 3

	// Funktion mit mehreren Return types
	resultOfDiv, okDiv := divide(1, 2)
	fmt.Println(resultOfDiv, okDiv) // 0 true

	resultOfDiv2, okDiv2 := divide2(1, 2)
	if okDiv2 {
		fmt.Println(resultOfDiv2) // 0
	}

	// Pointer as parameter
	firstName, lastName := "FirstName", "LastName"
	fmt.Println(firstName, lastName) // FirstName LastName
	pointerTest(firstName, &lastName)
	fmt.Println(firstName, lastName) // FirstName -Other new LastName-

	fmt.Println(factorial_recursion(4)) // 24
}

func greet1(name string) {
	fmt.Println("Hello", name)
}

func greet2(name1, name2 string) { // oder (name1 string, name2 string)
	fmt.Println("Hello", name1, "and", name2)
}

func greetMultiPara(names ...string) {
	for _, v := range names {
		fmt.Println("Hello", v)
	}
}

// function mit return type int
func add(l, r int) int {
	return l + r
}

// function mit mehreren return types
// return type int, bool
func divide(l, r int) (int, bool) {
	if r == 0 {
		return 0, false
	}
	return l / r, true
}

func divide2(l, r int) (result int, ok bool) {
	if r == 0 {
		return
	}
	result = l / r
	ok = true
	return
}

func pointerTest(name string, otherLastName *string) {
	// aus dem ersten Wert "name string", wird hier eine Kopie angelegt (Pass by value)
	// und mit der Kopie gearbeitet. Deshalb wird in der aufrufenden Funktion der Wert nicht geändert.
	// Der Wert wird hier lokal kopiert und damit gearbeitet.
	// "otherLastName" wird als Pointer übergeben, deshalb wird der
	// Wert in der aufrufenden Funktion überschrieben
	// Es wird der Wert geteilt
	name = "NewFirstName"
	*otherLastName = "-Other new LastName-"

	// Allgemein: Benutze Pointer um Speicher zu teilen, ansonsten Werte
}

func factorial_recursion(x float64) (y float64) {
	if x > 0 {
		//fmt.Println("x", x)
		y = x * factorial_recursion(x-1)
		fmt.Println("y", y)
	} else {
		y = 1
	}
	return
}
