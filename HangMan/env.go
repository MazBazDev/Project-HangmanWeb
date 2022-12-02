package hangman

type HangmanData struct {
	CurrentPage           int
	Attempts              int
	WordToFind            string
	Word                  string
	CurrentLetter         string
	PlayedLetters         string
	WordFinded            bool
	SavesPath             string
	CurrentSavesPath      string
	DictionaryPath        string
	CurrentDictionaryPath string
	PaternsPath           string
	UseAscii              bool
	AsciiPath             string
	CurrentAsciiPath      string
	Error                 string
}

var GameData HangmanData
