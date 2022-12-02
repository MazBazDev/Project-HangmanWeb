package hangman

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

// --
// Parameters | Type
// Atempts | int
// word, PlayedLetters, CurrentLetter, hangmanPaternsPath | string
// --
// This function calls the function to create the model of the game page
// --
var asciiUpBoxes int
var asciiUpBoxesHeigh int

func PageGame(Attempts int, word, PlayedLetters, CurrentLetter, hangmanPaternsPath string) {
	if GameData.UseAscii {
		asciiUpBoxesHeigh = 7
		for _, v := range OneWordAsciiArt(GameData.WordToFind) {
			if len(v) > asciiUpBoxes {
				asciiUpBoxes = len(v)
			}
		}
	}

	if GameData.Error != "" {
		CreateBox(4, 70+asciiUpBoxes, 4, 0, "white", "black", "Info", "white", []string{"", GameData.Error}, "red", ((70/2)-(len(GameData.Error)/2)-1)+(asciiUpBoxes/2))
	} else {

		body := []string{
			"You have 10 attempts to find this word",
			"       Good luck and Have Fun ;)",
		}

		DisplayInfo(body, "white")
	}

	AttemptsBox(Attempts)
	HangBox(GetHangPatern(GameData.PaternsPath, Attempts))
	DisplayWord(word)
	DisplayPlayedLetters(PlayedLetters)
	DisplayCurrentLetter(CurrentLetter)
}

func DisplayInfo(msg []string, TextColor string) {
	CreateBox(4, 70+asciiUpBoxes, 4, 0, "white", "black", "Info", "white", msg, TextColor, ((70/2)-(len(msg[0])/2)-1)+(asciiUpBoxes/2))
}

// --
// Parameters | Type
// path | string
// step | int
// --
// This function returns one hangman pattern at a time gives step formatted as []string
// --
func GetHangPatern(path string, step int) []string {
	step = (9) - step

	var fileLines []string

	readFile, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	i := 0
	start := step*7 + step
	end := step*7 + 6 + step

	for fileScanner.Scan() {
		if i >= start && i <= end {
			fileLines = append(fileLines, fileScanner.Text())
		}
		i++
	}

	readFile.Close()
	return fileLines
}

// --
// Parameters | Type
// hangman | []string
// --
// This function returns display the current state of the hangman
// --
func HangBox(hangman []string) {
	CreateBox(9, 19, 8, 75+asciiUpBoxes, "white", "black", "HangMan", "white", hangman, "white", 4)
}

// --
// Parameters | Type
// attempts | int
// --
// This function display the current attempts left
// --
func AttemptsBox(attempts int) {
	CreateBox(3, 19, 4, 75+asciiUpBoxes, "white", "black", "Attempts", "white", []string{strconv.Itoa(attempts)}, "white", (19/2)-(len(strconv.Itoa(attempts))/2)-1)
}

// --
// Parameters | Type
// word | string
// --
// This function display the word to find, with letters or with ascii art
// --
func DisplayWord(word string) {
	if GameData.UseAscii {
		CreateBox(5+asciiUpBoxesHeigh, 70+asciiUpBoxes, 9, 0, "white", "black", "Word", "white", OneWordAsciiArt(GameData.Word), "white", 28)
	} else {
		CreateBox(5, 70, 9, 0, "white", "black", "Word", "white", []string{"", word}, "white", (70/2)-(len(word)/2)-2)
	}
}

// --
// Parameters | Type
// PlayedLetters | string
// --
// This function display all played letters
// --
func DisplayPlayedLetters(PlayedLetters string) {
	CreateBox(5, 19, 18, 75+asciiUpBoxes, "white", "black", "Letters", "white", []string{"", PlayedLetters}, "white", 2)
}

// --
// Parameters | Type
// PlayedLetters | string
// --
// This function display the current letters
// --
func DisplayCurrentLetter(CurrentLetter string) {
	if GameData.UseAscii {
		CreateBox(5+asciiUpBoxesHeigh, 70+asciiUpBoxes, 15+asciiUpBoxesHeigh, 0, "white", "black", "Press \"ENTER\" to try your letter/word", "white", OneWordAsciiArt(CurrentLetter), "white", 10)
	} else {
		CreateBox(5+asciiUpBoxesHeigh, 70+asciiUpBoxes, 15+asciiUpBoxesHeigh, 0, "white", "black", "Press \"ENTER\" to try your letter/word", "white", []string{"", CurrentLetter}, "white", (70/2)+asciiUpBoxes-(len(CurrentLetter)/2)-2)
	}
}
