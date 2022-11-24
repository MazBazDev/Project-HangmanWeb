package main

import (
	"fmt"
	"net/http"
)

type Student struct {
	Name  string
	Age   string
	Quote string
	Hobby string
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to Hangman Web!")
	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)

}
