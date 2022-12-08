package main

import (
	"encoding/csv"
	"fmt"
	hangman "hangman/HangMan"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type oldData struct {
	Attemps    int
	WordToFind string
	Word       string
}
type sessionData struct {
	Logged     bool
	Usermane   string
	Email      string
	Win        int
	Loose      int
	Played     int
	Error      string
	Game       hangman.HangmanData
	OldDatas   oldData
	Scoreboard []ScorePlayer
}

type ScorePlayer struct {
	Name  string
	Win   int
	loose int
}

var session = sessionData{}

func main() {
	Scoreboard()
	http.HandleFunc("/", Routing)
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	fmt.Println("Server started : http://127.0.0.1:8080")
	http.ListenAndServe(":8080", nil)
}

func Routing(w http.ResponseWriter, request *http.Request) {
	switch request.URL.Path {
	case "/":
		template.Must(template.ParseFiles("static/templates/layout.html", "static/pages/index.html")).Execute(w, session)
	case "/login":
		if session.Logged {
			http.Redirect(w, request, "/", http.StatusSeeOther)
		} else {
			if request.Method == "GET" {
				template.Must(template.ParseFiles("static/templates/layout.html", "static/pages/login.html")).Execute(w, session)
			} else if request.Method == "POST" {
				Login(w, request)
			}
		}
	case "/register":
		if session.Logged {
			http.Redirect(w, request, "/", http.StatusSeeOther)
		} else {
			if request.Method == "GET" {
				template.Must(template.ParseFiles("static/templates/layout.html", "static/pages/register.html")).Execute(w, session)
			} else if request.Method == "POST" {
				Register(w, request)
			}
		}
	case "/stats":
		template.Must(template.ParseFiles("static/templates/layout.html", "static/pages/stats.html")).Execute(w, session)
	case "/dictionary":
		if session.Logged && hangman.GameData.WordToFind == "" {
			if request.Method == "GET" {
				template.Must(template.ParseFiles("static/templates/layout.html", "static/pages/dictionary.html")).Execute(w, session)
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
				template.Must(template.ParseFiles("static/templates/layout.html", "static/pages/game.html")).Execute(w, session)

			} else if request.Method == "POST" {
				Play(w, request)
			}
		}
	case "/win":
		template.Must(template.ParseFiles("static/templates/layout.html", "static/pages/win.html")).Execute(w, session)
	case "/loose":
		template.Must(template.ParseFiles("static/templates/layout.html", "static/pages/loose.html")).Execute(w, session)
	case "/endGame":
		if request.Method == "POST" && session.Logged {
			ClearGameStruct()
		}
		http.Redirect(w, request, "/", http.StatusSeeOther)
	case "/logout":
		if request.Method == "POST" && session.Logged {
			ClearGameStruct()
			session = sessionData{
				Logged:   false,
				Usermane: "",
				Email:    "",
				Game:     hangman.HangmanData{},
			}
			http.Redirect(w, request, "/", http.StatusSeeOther)
		}
	default:
		template.Must(template.ParseFiles("static/templates/layout.html", "static/pages/error.html")).Execute(w, session)
	}
}

func RegisterHasAccount(w http.ResponseWriter, request *http.Request) bool {
	file, err := os.Open("data/accounts.csv")
	if err != nil {
		Log("Function: RegisterHasAccount", err)
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
		var newData string
		data := [][]string{
			{
				request.FormValue("name"),
				request.FormValue("email"),
				HashPassword(request.FormValue("password")),
				strconv.Itoa(session.Win),
				strconv.Itoa(session.Loose),
				strconv.Itoa(session.Played),
			},
		}
		file, err := os.Open("data/accounts.csv")
		if err != nil {
			Log("Func Register", err)
		}
		reader := csv.NewReader(file)
		records, _ := reader.ReadAll()

		records = append(records, data[0])
		newData = strings.Join(records[0], ",")
		count := 1
		for count < len(records) {
			newData = newData + "\n" + strings.Join(records[count], ",")
			count++
		}
		os.WriteFile("data/accounts.csv", []byte(newData), 0644)

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
		Log("Func Login", err)
	}
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	for _, line := range records {
		if line[1] == request.FormValue("email") && CheckPasswordHash(request.FormValue("password"), line[2]) {
			played, _ := strconv.Atoi(line[5])
			session = sessionData{
				Logged:   true,
				Usermane: line[0],
				Email:    line[1],
				Played:   played,
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
	hangman.GameData = hangman.HangmanData{
		WordFinded:            false,
		WordToFind:            "",
		PaternsPath:           "./HangMan/files/hangman.txt",
		CurrentDictionaryPath: "./HangMan/files/dictionary/" + request.FormValue("level") + ".txt",
		Attempts:              10,
	}
	hangman.GameData.WordToFind = hangman.GetRandomWord(hangman.GameData.CurrentDictionaryPath)
	hangman.WordBegining(hangman.GameData.WordToFind)

	session.Game = hangman.GameData
	http.Redirect(w, request, "/hangman", http.StatusSeeOther)
}

func Play(w http.ResponseWriter, request *http.Request) {
	hangman.GameData.CurrentLetter = request.FormValue("letter")
	hangman.Play()

	if hangman.GameData.Error != "" {
		session.Error = hangman.GameData.Error
	}

	session.Game = hangman.GameData

	if hangman.GameData.WordFinded {
		session.Win++
		session.Played++
		UpdateStats()
		SetOldDatas()

		http.Redirect(w, request, "/win", http.StatusSeeOther)

	} else if hangman.GameData.Attempts <= 0 && !hangman.GameData.WordFinded {
		session.Loose++
		session.Played++
		UpdateStats()
		SetOldDatas()

		http.Redirect(w, request, "/loose", http.StatusSeeOther)
	} else {
		http.Redirect(w, request, "/hangman", http.StatusSeeOther)
	}
}

func SetOldDatas() {
	session.OldDatas = oldData{
		Attemps:    session.Game.Attempts,
		WordToFind: hangman.GameData.WordToFind,
		Word:       hangman.GameData.Word,
	}
	ClearGameStruct()
}

func ClearGameStruct() {
	hangman.GameData = hangman.HangmanData{
		CurrentDictionaryPath: "",
		Word:                  "",
		CurrentLetter:         "",
		PlayedLetters:         "",
		Attempts:              10,
		Error:                 "",
	}
	session.Error = ""
	session.Game = hangman.GameData
}

func Log(Error string, err1 error) {
	f, err := os.OpenFile("./data/latest.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Log error opening file: %v", err)
	}

	defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)
	log.Println(Error, err1, "-------------")
}

func UpdateStats() {
	var newData string

	file, err := os.Open("data/accounts.csv")
	if err != nil {
		Log("Func Login", err)
	}
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	for index, line := range records {
		if strings.Compare(line[1], session.Email) == 0 {
			var newData []string = []string{line[0], line[1], line[2], strconv.Itoa(session.Win), strconv.Itoa(session.Loose), strconv.Itoa(session.Played)}
			records[index] = newData
		}
	}

	csvFile, err := os.OpenFile("data/accounts.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		Log("Func Register", err)
	}
	newData = strings.Join(records[0], ",")
	count := 1
	for count < len(records) {
		newData = newData + "\n" + strings.Join(records[count], ",")
		count++
	}

	os.WriteFile("data/accounts.csv", []byte(newData), 0644)

	csvFile.Close()

}

func Scoreboard() {
	file, err := os.Open("data/accounts.csv")
	if err != nil {
		Log("Func Login", err)
	}
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()
	for _, line := range records {
		twin, _ := strconv.Atoi(line[3])
		tloose, _ := strconv.Atoi(line[4])
		session.Scoreboard = append(session.Scoreboard, ScorePlayer{
			Name:  line[0],
			Win:   twin,
			loose: tloose,
		})
	}
	session = sessionData{}
}
