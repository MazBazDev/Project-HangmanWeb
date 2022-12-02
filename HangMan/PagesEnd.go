package hangman

import "strconv"

// This function print a box
func PageWelcome() {
	body := []string{
		"",
		"Welcome to HangMan Termbox version.",
		"Good luck to you !",
		"/!\\ PLEASE NOTE: each letter/word entered is final!",
		"Press \"Enter\" to confirm your choice.",
		"By ANNEG Noemie & YAKOUBEN Mazigh"}

	CreateBox(9, 94, 4, 0, "white", "black", "Welcome", "white", body, "white", 5)
}

// This function print a box
func PageHelp() {
	body := []string{
		"",
		"1. \"ESC\" to quit.",
		"2. \"ENTER\" to confirm your choice.",
		"3. \"BACKSPACE\" or \"DEL\" to delete the last letter.",
	}
	CreateBox(7, 94, 4, 0, "white", "black", "Help", "white", body, "white", 5)
}

// --
// Parameters | Type
// Status | bool
// Atempts | int
// Word | string
// WordToFind | string
// Hangman | []string
// --
// Return type : []string
// --
// This function displays the final page with the current progress.
// --
func PageFinal(Status bool, Attempts int, Word, WordToFind string, HangMan []string) {
	if Status {
		body := []string{
			" __     __                                 _ ",
			" \\ \\   / /                                | |",
			"  \\ \\_/ /__  _   _  __      _____  _ __   | |",
			"   \\   / _ \\| | | | \\ \\ /\\ / / _ \\| '_ \\  | |",
			"    | | (_) | |_| |  \\ V  V / (_) | | | | |_|",
			"    |_|\\___/ \\__,_|   \\_/\\_/ \\___/|_| |_| (_)",
			"",
			"",
			"",
		}

		for _, v := range HangMan {
			body = append(body, "                   "+v)
		}

		end := []string{
			"",
			"         Press any key on the keyboard",
			"                to close the game",
		}

		for _, v := range end {
			body = append(body, ""+v)
		}
		CreateBox(22, 94, 4, 0, "white", "black", "You won !", "white", body, "white", 22)

		finalSentence := "You find \"" + GameData.Word + "\" with " + strconv.Itoa(Attempts) + " attempts left"

		TbPrint((94/2)-(len(finalSentence)/2)-1, 12, "white", "black", finalSentence)

	} else {
		body := []string{
			"__     __           _           _",
			"\\ \\   / /          | |         | |  ",
			" \\ \\_/ /__  _   _  | | ___  ___| |_ ",
			"  \\   / _ \\| | | | | |/ _ \\/ __| __|",
			"   | | (_) | |_| | | | (_) \\__ \\ |_ ",
			"   |_|\\___/ \\__,_| |_|\\___/|___/\\__|",
			"",
			"",
		}

		for _, v := range HangMan {
			body = append(body, "             "+v)
		}

		end := []string{
			"",
			"      The word to find was:",
			"      \"" + WordToFind + "\"",
			"",
			"      You finded:",
			"      \"" + Word + "\"",
			"",
			"    Press any key on the keyboard",
			"         to close the game",
		}

		for _, v := range end {
			body = append(body, ""+v)
		}
		CreateBox(26, 94, 4, 0, "white", "black", "You lost !", "white", body, "white", 28)
	}
}
