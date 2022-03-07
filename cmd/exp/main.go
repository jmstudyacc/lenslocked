package main

import (
	// ensure you have html/template package for text when running a template
	"html/template"
	"os"
)

type User struct {
	Name string
	Age  int
	Meta UserMeta
	Bio  string
}

type UserMeta struct {
	Visits int
}

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	user := User{
		Name: "John Smith",
		Bio:  `<script>alert("Haha you have been haxxorxed!");</script>`,
		Age:  99,
		Meta: UserMeta{Visits: 5},
	}

	// this is using an anonymous struct and not creating a type struct
	// this is useful when you want to test something quickly - useful when you are running test functions
	//user2 := struct {
	//	Name string
	//}{
	//	Name: "Susan Smith",
	//}

	// you have to define where the template is output to
	err = t.Execute(os.Stdout, user)

	if err != nil {
		panic(err)
	}
}
