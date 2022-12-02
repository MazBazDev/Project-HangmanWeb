package hangman

import (
	"strings"

	"github.com/nsf/termbox-go"
)

func OldMain() {
	GameData.PaternsPath = "./files/hangman.txt"
	GameData.SavesPath = "./files/saves/"
	GameData.DictionaryPath = "./files/dictionary/"
	GameData.AsciiPath = "./files/ascii/"
	GameData.Attempts = 10

	Selector("saves")
}

func Selector(what string) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	var Title string
	var Heigth int
	var Rows []string

	if what == "dictionary" {
		Files := ListFilesInFolder(GameData.DictionaryPath)
		for _, v := range Files {
			Rows = append(Rows, strings.Replace(v, ".txt", "", -1))
		}
		Title = "Select your dictionary"
		Heigth = len(Files)

	} else if what == "saves" {
		Rows = []string{""}

		Files := ListFilesInFolder(GameData.SavesPath)
		for _, v := range Files {
			Rows = append(Rows, strings.Replace(v, ".json", "", -1))
		}
		Title = "Select a save / New game"
		Heigth = len(Files) + 1

	} else if what == "ascii" {
		Rows = []string{"Oui", "Non"}
		Title = "Use Ascii art design ?"
		Heigth = 2
	} else if what == "ascii2" {
		Files := ListFilesInFolder(GameData.AsciiPath)
		for _, v := range Files {
			Rows = append(Rows, strings.Replace(v, ".txt", "", -1))
		}
		Title = "Select your Ascii theme"
		Heigth = len(Files)

	}

	if len(Rows) == 0 {
		Selector("dictionary")
	} else {
		var Selectindex int
		err := termbox.Init()
		if err != nil {
			panic(err)
		}

		defer NextSelector(what)

	mainloop:
		for {
			CreateBox(Heigth+2, 94, 0, 0, "white", "black", Title, "white", Rows, "white", 4)

			if what == "saves" {
				TbPrint(5, 1, "cyan", "black", "Start a new game")
			}

			TbPrint(2, Selectindex+1, "white", "black", ">>")
			termbox.Flush()

			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyArrowDown:
					if Selectindex < len(Rows)-1 {
						Selectindex++
					}
				case termbox.KeyArrowUp:
					if Selectindex > 0 {
						Selectindex--
					}
				case termbox.KeyEnter:
					if what == "dictionary" {
						GameData.CurrentDictionaryPath = GetPathFromIndex(GameData.DictionaryPath, Selectindex)
					} else if what == "saves" {
						if Selectindex != 0 {
							GameData.CurrentSavesPath = GetPathFromIndex(GameData.SavesPath, Selectindex-1)
						}
					} else if what == "ascii" {
						if Selectindex == 0 {
							GameData.UseAscii = true
						} else {
							GameData.UseAscii = false
						}
					} else if what == "ascii2" {
						GameData.CurrentAsciiPath = GetPathFromIndex(GameData.AsciiPath, Selectindex)
					}
					break mainloop
				}
			}
		}
	}
}

func NextSelector(what string) {
	if what == "saves" {
		if GameData.CurrentSavesPath == "" {
			Selector("dictionary")
		} else {
			LoadSave(GameData.CurrentSavesPath)
			GameMain()
		}
	} else if what == "dictionary" {
		GameData.WordToFind = GetRandomWord(GameData.CurrentDictionaryPath)
		WordBegining(GameData.WordToFind)
		Selector("ascii")

	} else if what == "ascii" {
		if GameData.UseAscii {
			Selector("ascii2")
		} else {
			GameMain()
		}
	} else if what == "ascii2" {
		GameMain()
	}
}
