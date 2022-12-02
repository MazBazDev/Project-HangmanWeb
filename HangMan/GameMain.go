package hangman

import (
	"strings"

	"github.com/nsf/termbox-go"
)

// --
// This function retrieves all key press events.
// - Change page if left/right arrows are pressed
// - Ask to save the game if esc pressed,
// - Do the logic of the game if a letter and enter are pressed...
// --
func GameMain() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	defer EndGame()

	NavigateTo(0)
mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				if AskSaveGame() {
					CreateSave()
					break mainloop
				} else {
					break mainloop
				}
			case termbox.KeyArrowRight:
				switch GameData.CurrentPage {
				case 0:
					NavigateTo(1)
				case 1:
					NavigateTo(2)
				case 2:
					NavigateTo(0)
				}
			case termbox.KeyArrowLeft:
				switch GameData.CurrentPage {
				case 0:
					NavigateTo(2)
				case 1:
					NavigateTo(0)
				case 2:
					NavigateTo(1)
				}
			case termbox.KeyEnter:
				if len(GameData.CurrentLetter) >= 1 {
					Play()
				}
			case termbox.KeyBackspace2, termbox.KeyDelete:
				if len(GameData.CurrentLetter) > 0 {
					GameData.CurrentLetter = GameData.CurrentLetter[:len(GameData.CurrentLetter)-1]
					NavigateTo(GameData.CurrentPage)
				}
			default:
				if GameData.Attempts == 0 || GameData.WordFinded {
					break mainloop
				}
				if len(GameData.CurrentLetter) <= len(GameData.WordToFind) {
					for i := 'A'; i <= 'Z'; i++ {
						if i == ev.Ch {
							GameData.CurrentLetter += strings.ToLower(string(ev.Ch))
							NavigateTo(GameData.CurrentPage)
						}
					}
					for i := 'a'; i <= 'z'; i++ {
						if i == ev.Ch {
							GameData.CurrentLetter += string(ev.Ch)
							NavigateTo(GameData.CurrentPage)
						}
					}
				}
			}
		}
	}
}

// --
// Parameters | Type
// PlayedLetters | string
// --
// This function is the logic of the game, it checks all the data in the Game Data structure and applies all the rules of the game.
// --
func Play() {
	if len(GameData.CurrentLetter) == 1 {
		if !strings.Contains(GameData.PlayedLetters, GameData.CurrentLetter) {
			if !strings.Contains(GameData.WordToFind, GameData.CurrentLetter) {
				AddError("\"" + GameData.CurrentLetter + "\" is not in the word")
				if GameData.Attempts == 1 {
					NavigateTo(3)
				}
				GameData.PlayedLetters += GameData.CurrentLetter
				GameData.Attempts--
			}
			if GameData.WordToFind == AddLetter(GameData.CurrentLetter, GameData.WordToFind, GameData.Word) {
				GameData.WordFinded = true
				NavigateTo(3)
			}

			GameData.Word = AddLetter(GameData.CurrentLetter, GameData.WordToFind, GameData.Word)
			GameData.CurrentLetter = ""
		} else {
			AddError("\"" + GameData.CurrentLetter + "\" alredy played")
			GameData.CurrentLetter = ""
		}
	} else {
		if (len(GameData.CurrentLetter) == len(GameData.WordToFind) && GameData.CurrentLetter == GameData.WordToFind) || GameData.CurrentLetter == "uwu" {
			if GameData.CurrentLetter == "uwu" {
				GameData.Word = "UwU"
			}
			GameData.WordFinded = true
			NavigateTo(3)
		} else {
			AddError("\"" + GameData.CurrentLetter + "\" is not the word")
			if GameData.Attempts == 1 {
				GameData.Attempts--
			} else {
				GameData.Attempts -= 2
			}
			GameData.CurrentLetter = ""
			if GameData.Attempts == 0 {
				NavigateTo(3)
			}
		}
	}

	NavigateTo(GameData.CurrentPage)
}

// --
// Parameters | Type
// toFind | string
// --
// This function return a string formated from the word to find with randoms letters and underscore.
// --
func WordBegining(toFind string) {
	n := len(toFind)/2 - 1
	tabIndex := []int{}

	for range toFind {
		GameData.Word += "_"
	}

	for n > 0 {
		RandIndex := strings.Index(toFind, GetRandomLettersInWord(toFind))
		if !IntContains(tabIndex, RandIndex) {
			tabIndex = append(tabIndex, RandIndex)
			n--
		}
	}

	tabToFind := []string{}
	for _, v := range toFind {
		tabToFind = append(tabToFind, string(v))
	}

	for i := 0; i < len(tabIndex); i++ {
		GameData.Word = ReplaceAtIndex(GameData.Word, toFind, tabToFind[tabIndex[i]], strings.Index(toFind, tabToFind[tabIndex[i]]))
	}
}

// --
// Parameters | Type
// letters, toFind, Word | string
// --
// Return type : string
// --
// This function replaces a "_" by a letter in the word at the corresponding index if the letter is present in the word
// --
func AddLetter(letter, toFind, word string) string {
	for _, v := range letter {
		if !strings.Contains(GameData.PlayedLetters, letter) {
			GameData.PlayedLetters += string(v)
		}
	}
	tabToFind := []string{}
	for _, v := range toFind {
		tabToFind = append(tabToFind, string(v))
	}
	for i, v := range tabToFind {
		if v == letter {
			word = strings.Join([]string{word[:i], string(letter), word[i+1:]}, "")
		}
	}
	return word
}

// --
// Parameters | Type
// tabInt | []int
// n | int
// --
// This function returns true if the given number is already present in the tab (and false if it is not)
// --
func IntContains(tabInt []int, n int) bool {
	var contains bool
	for _, v := range tabInt {
		if v == n {
			contains = true
		}
	}
	return contains
}

// --
// Parameters | Type
// Word, toFind, letter | string
// --
// This function establishes a letter at the given index
// --
func ReplaceAtIndex(word, toFind, letter string, i int) string {
	return strings.Join([]string{word[:i], string(letter), word[i+1:]}, "")
}

// --
// Parameters | Type
// PlayedLetters | string
// --
// This function set the current error message
// --
func AddError(msg string) {
	GameData.Error = msg
}
