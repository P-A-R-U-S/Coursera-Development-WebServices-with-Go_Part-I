package main

import "fmt"

type Basic struct {
	Base string
}
type Person struct {
	Id      int
	Name    string
	Address string
}

type Account struct {
	Id int
	// Name    string
	Cleaner func(string) string
	Owner   Person
	Person
	Basic
}

func main() {
	// полное объявление структуры
	var acc Account = Account{
		Id: 1,
		// Name: "rvasily",
		Person: Person{
			Name:    "Василий",
			Address: "Москва",
		},
		Basic: Basic{ Base: "Base from Basic "},
	}
	fmt.Printf("%#v\n", acc)

	// короткое объявление структуры
	acc.Owner = Person{2, "Romanov Vasily", "Moscow"}

	fmt.Printf("%#v\n", acc)

	acc.Name = "Valentyn"

	fmt.Println(acc.Name)
	fmt.Println(acc.Person.Name)
	fmt.Println(acc.Basic.Base)



}

