package hangman

import "github.com/nsf/termbox-go"

// --
// This function display the navigation bar
// --
func NavBar() {
	var selectedIndex int
	switch GameData.CurrentPage {
	case 0:
		selectedIndex = 20
	case 1:
		selectedIndex = 42
	case 2:
		selectedIndex = 65
	}

	body := []string{"Welcome               Game                   Help"}
	CreateBox(3, 94, 0, 0, "white", "black", "Welcome", "white", body, "white", 21)
	if !(GameData.CurrentPage > 2) {
		TbPrint(selectedIndex-1, 1, "white", "black", ">>")
	}
}

// --
// Parameters | Type
// Page | string
// --
// This function is used to switch from the current page to another page passed in parameter
// --
func NavigateTo(page int) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	switch page {
	case 0:
		GameData.CurrentPage = page
		PageWelcome()
	case 1:
		GameData.CurrentPage = page
		PageGame(GameData.Attempts, GameData.Word, GameData.PlayedLetters, GameData.CurrentLetter, GameData.PaternsPath)
	case 2:
		GameData.CurrentPage = page
		PageHelp()
	case 3:
		GameData.CurrentPage = page
		PageFinal(GameData.WordFinded, GameData.Attempts, GameData.Word, GameData.WordToFind, GetHangPatern(GameData.PaternsPath, GameData.Attempts))
	}
	NavBar()
	termbox.Flush()
}
