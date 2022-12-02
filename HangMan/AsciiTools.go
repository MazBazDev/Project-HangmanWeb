package hangman

// --
// Parameters | Type
// toFind | string
// line| int
// --
// This function return one line of a full word printed in ascii art characters
// --
func OneLineAsciiArt(toFind string, line int) string {
	res := ""
	tabToFind := []rune{}
	for _, v := range toFind {
		tabToFind = append(tabToFind, v)
	}

	for i := 0; i < len(tabToFind); i++ {
		res += OneLetterAsciiArt(toFind, tabToFind[i], line)
	}
	return res
}

// --
// Parameters | Type
// toFind | string
// letter | rune
// line| int
// --
// This function returns one letter of a word in ascii art character
// --
func OneLetterAsciiArt(toFind string, letter rune, line int) string {
	tabLetter := GetAsciiPattern(GameData.CurrentAsciiPath, letter)
	j := ""
	for i := 0; i < len(tabLetter); i++ {
		j = tabLetter[line]
		break
	}
	return j
}
