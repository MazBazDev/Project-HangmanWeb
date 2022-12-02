package hangman

import "github.com/nsf/termbox-go"

// --
// This function closes the termbox and calls the DeleteSaveIfWinOrLoose() function
// --
func EndGame() {
	termbox.Close()
	DeleteSaveIfWinOrLoose()
}
