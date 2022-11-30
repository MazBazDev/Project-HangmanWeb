package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type sessionData struct {
	Logged   bool
	Usermane string
	Email    string
	Error    string
}

var session = sessionData{}

func main() {
	http.HandleFunc("/", Routing)

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)
}

func Routing(w http.ResponseWriter, request *http.Request) {
	switch request.URL.Path {
	case "/":
		template.Must(template.ParseFiles("static/pages/index.html")).Execute(w, session)
	case "/login":
		if session.Logged {
			http.Redirect(w, request, "/", http.StatusSeeOther)
		} else {
			if request.Method == "GET" {
				template.Must(template.ParseFiles("static/pages/login.html")).Execute(w, session)
			} else if request.Method == "POST" {
				Login(w, request)
			}
		}
	case "/register":
		if session.Logged {
			http.Redirect(w, request, "/", http.StatusSeeOther)
		} else {
			if request.Method == "GET" {
				template.Must(template.ParseFiles("static/pages/register.html")).Execute(w, session)
			} else if request.Method == "POST" {
				Register(w, request)
			}
		}
	case "/stats":
		template.Must(template.ParseFiles("static/pages/stats.html")).Execute(w, session)
	case "/play":
		template.Must(template.ParseFiles("static/pages/game.html")).Execute(w, session)
	case "/logout":
		if request.Method == "POST" && session.Logged {
			session = sessionData{
				Logged:   false,
				Usermane: "",
				Email:    "",
			}
			http.Redirect(w, request, "/", http.StatusSeeOther)
		}
	}
}

func PasswordCheck(w http.ResponseWriter, request *http.Request) bool {
	if request.FormValue("password") != request.FormValue("password-confirm") {
		session.Error = "Passwords don't match"
		http.Redirect(w, request, "/register", http.StatusSeeOther)
		return false
	}
	return true
}

func RegisterHasAccount(w http.ResponseWriter, request *http.Request) bool {
	file, err := os.Open("data/accounts.csv")
	if err != nil {
		fmt.Println(err)
	}
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	for _, line := range records {
		if line[1] == request.FormValue("email") {
			session.Error = "You already have an account !"
			http.Redirect(w, request, "/register", http.StatusSeeOther)
			return true
		}
	}

	return false
}

func Register(w http.ResponseWriter, request *http.Request) {
	if PasswordCheck(w, request) && !RegisterHasAccount(w, request) {
		data := [][]string{
			{request.FormValue("name"), request.FormValue("email"), request.FormValue("password")},
		}

		//create a file
		csvFile, err := os.OpenFile("data/accounts.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

		if err != nil {
			log.Fatalf("Failed to create file,: %", err)
		}
		//initialize csv writer
		csvWriter := csv.NewWriter(csvFile)
		for _, value := range data {
			csvWriter.Write(value)
		}

		csvWriter.Flush()
		csvFile.Close()

		session = sessionData{
			Logged:   true,
			Usermane: request.FormValue("name"),
			Email:    request.FormValue("email"),
		}

		http.Redirect(w, request, "/", http.StatusSeeOther)
	}
}

func Login(w http.ResponseWriter, request *http.Request) {
	file, err := os.Open("data/accounts.csv")
	if err != nil {
		fmt.Println(err)
	}
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	for _, line := range records {
		if line[1] == request.FormValue("email") && line[2] == request.FormValue("password") {
			session = sessionData{
				Logged:   true,
				Usermane: line[0],
				Email:    line[1],
			}
			http.Redirect(w, request, "/", http.StatusSeeOther)
		}
	}

	if !session.Logged {
		session.Error = "Bad credentials, retry or create an account"
		http.Redirect(w, request, "/login", http.StatusSeeOther)
	}
}
