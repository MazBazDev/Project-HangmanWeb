package main

import (
	"encoding/csv"
	"fmt"
	hangman "hangman/HangMan"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type sessionData struct {
	Logged   bool
	Usermane string
	Email    string
	Error    string
	Game     hangman.HangmanData
	Win      int
	Loose    int
	//CorrectAttempts int
	//WrongAttempts   int
}

var session = sessionData{}

func main() {
	http.HandleFunc("/", Routing)
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	fmt.Println("Server started : http://127.0.0.1:8080")
	http.ListenAndServe(":8080", nil)
}

func Routing(w http.ResponseWriter, request *http.Request) {
	switch request.URL.Path {
	case "/":
		template.Must(template.ParseFiles("static/pages/index.html", "static/templates/nav.html", "static/templates/head.html")).Execute(w, session)
	case "/login":
		if session.Logged {
			http.Redirect(w, request, "/", http.StatusSeeOther)
		} else {
			if request.Method == "GET" {
				template.Must(template.ParseFiles("static/pages/login.html", "static/templates/nav.html", "static/templates/head.html")).Execute(w, session)
			} else if request.Method == "POST" {
				Login(w, request)
			}
		}
	case "/register":
		if session.Logged {
			http.Redirect(w, request, "/", http.StatusSeeOther)
		} else {
			if request.Method == "GET" {
				template.Must(template.ParseFiles("static/pages/register.html", "static/templates/nav.html", "static/templates/head.html")).Execute(w, session)
			} else if request.Method == "POST" {
				Register(w, request)
			}
		}
	case "/stats":
		template.Must(template.ParseFiles("static/pages/stats.html", "static/templates/nav.html", "static/templates/head.html")).Execute(w, session)
	case "/dictionary":
		if session.Logged && hangman.GameData.CurrentDictionaryPath == "" {
			if request.Method == "GET" {
				template.Must(template.ParseFiles("static/pages/dictionary.html", "static/templates/nav.html", "static/templates/head.html")).Execute(w, session)
			} else if request.Method == "POST" {
				InitGame(w, request)
			}
		} else {
			http.Redirect(w, request, "/hangman", http.StatusSeeOther)
		}
	case "/hangman":
		if !session.Logged {
			http.Redirect(w, request, "/", http.StatusSeeOther)
		} else {
			if hangman.GameData.CurrentDictionaryPath == "" {
				http.Redirect(w, request, "/dictionary", http.StatusSeeOther)
			}
			if request.Method == "GET" {
				fmt.Println(session)
				template.Must(template.ParseFiles("static/pages/game.html", "static/templates/nav.html", "static/templates/head.html")).Execute(w, session)

			} else if request.Method == "POST" {
				Play(w, request)
			}
			if hangman.GameData.WordFinded {
				template.Must(template.ParseFiles("static/pages/win.html", "static/templates/nav.html", "static/templates/head.html")).Execute(w, session)
				session.Win++
				hangman.GameData.CurrentDictionaryPath = ""
				hangman.GameData.Word = ""
				hangman.GameData.CurrentLetter = ""
				hangman.GameData.PlayedLetters = ""
				hangman.GameData.Attempts = 10
				hangman.GameData.Error = ""
			} else if hangman.GameData.Attempts == 0 && !hangman.GameData.WordFinded {
				template.Must(template.ParseFiles("static/pages/end.html", "static/templates/nav.html", "static/templates/head.html")).Execute(w, session)
				session.Loose++
				hangman.GameData.CurrentDictionaryPath = ""
				hangman.GameData.Word = ""
				hangman.GameData.CurrentLetter = ""
				hangman.GameData.PlayedLetters = ""
				hangman.GameData.Attempts = 10
				hangman.GameData.Error = ""
			}
		}
	case "/logout":
		if request.Method == "POST" && session.Logged {
			session = sessionData{
				Logged:   false,
				Usermane: "",
				Email:    "",
			}
			http.Redirect(w, request, "/", http.StatusSeeOther)
		}
	default:
		template.Must(template.ParseFiles("static/pages/error.html", "static/templates/nav.html", "static/templates/head.html")).Execute(w, session)
	}
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

func PasswordCheck(w http.ResponseWriter, request *http.Request) bool {
	if request.FormValue("password") != request.FormValue("password-confirm") {
		session.Error = "Passwords don't match"
		http.Redirect(w, request, "/register", http.StatusSeeOther)
		return false
	}
	return true
}

func Register(w http.ResponseWriter, request *http.Request) {

	if PasswordCheck(w, request) && !RegisterHasAccount(w, request) {
		data := [][]string{
			{
				request.FormValue("name"),
				request.FormValue("email"),
				HashPassword(request.FormValue("password")),
				strconv.Itoa(session.Win),
				strconv.Itoa(session.Loose),
				//strconv.Itoa(session.CorrectAttempts),
				//strconv.Itoa(session.WrongAttempts),
			},
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

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		fmt.Println(err)
	}

	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Login(w http.ResponseWriter, request *http.Request) {
	file, err := os.Open("data/accounts.csv")
	if err != nil {
		fmt.Println(err)
	}
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	for _, line := range records {
		if line[1] == request.FormValue("email") && CheckPasswordHash(request.FormValue("password"), line[2]) {
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

func InitGame(w http.ResponseWriter, request *http.Request) {
	hangman.GameData.WordFinded = false
	hangman.GameData.WordToFind = ""
	hangman.GameData.PaternsPath = "./HangMan/files/hangman.txt"
	hangman.GameData.CurrentDictionaryPath = "./HangMan/files/dictionary/" + request.FormValue("level") + ".txt"
	hangman.GameData.Attempts = 10
	hangman.GameData.WordToFind = hangman.GetRandomWord(hangman.GameData.CurrentDictionaryPath)
	hangman.WordBegining(hangman.GameData.WordToFind)

	session.Game = hangman.GameData
	fmt.Println("attempts", session.Game.Attempts)
	http.Redirect(w, request, "/hangman", http.StatusSeeOther)
}

func Play(w http.ResponseWriter, request *http.Request) {
	hangman.GameData.CurrentLetter = request.FormValue("letter")
	hangman.Play()
	if hangman.GameData.Attempts > 0 && !hangman.GameData.WordFinded {
		if hangman.GameData.Error != "" {
			session.Error = hangman.GameData.Error
		}
		session.Game = hangman.GameData
		http.Redirect(w, request, "/hangman", http.StatusSeeOther)
	}
}
