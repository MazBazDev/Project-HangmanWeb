package hangman

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

// --
// Parameters | Type
// row, col, marginTop, marginLeft, linesMarginLeft| int
// boxColor, boxBg, title, titleColor, lineColor | string
// lines | []string
// --
// This function create a box and print an array of string in it
// --
func CreateBox(rows, cols, marginTop, marginLeft int, boxColor, boxBg, title, titleColor string, lines []string, lineColor string, linesMarginLeft int) {
	rows, cols = cols, rows
	for col := 0; col < cols; col++ {
		if col == 0 || col == (cols-1) {
			for row := 0; row < rows; row++ {
				if row == 0 {
					if col == 0 {
						termbox.SetCell(row+marginLeft, col+marginTop, '╔', ColorPicker(boxColor), ColorPicker(boxBg))
					} else if col == cols-1 {
						termbox.SetCell(row+marginLeft, col+marginTop, '╚', ColorPicker(boxColor), ColorPicker(boxBg))
					}
				} else if row == (rows - 1) {
					if col == 0 {
						termbox.SetCell(row+marginLeft, col+marginTop, '╗', ColorPicker(boxColor), ColorPicker(boxBg))
					} else if col == cols-1 {
						termbox.SetCell(row+marginLeft, col+marginTop, '╝', ColorPicker(boxColor), ColorPicker(boxBg))
					}
				} else {
					termbox.SetCell(row+marginLeft, col+marginTop, '═', ColorPicker(boxColor), ColorPicker(boxBg))
				}
			}
			if len(title) > 0 {
				sentence := "[ " + title + " ]"
				for i := 0; i < len(sentence); i++ {
					termbox.SetCell(marginLeft+3+i, marginTop, rune(sentence[i]), ColorPicker(titleColor), ColorPicker(boxBg))
				}
			}
		} else {
			for row := 0; row < rows; row++ {
				if row == 0 || row == (rows-1) {
					termbox.SetCell(row+marginLeft, col+marginTop, '║', ColorPicker(boxColor), ColorPicker(boxBg))
				} else {
					termbox.SetCell(row+marginLeft, col+marginTop, ' ', ColorPicker(boxColor), ColorPicker(boxBg))
				}
			}
		}
	}

	for i, line := range lines {
		for e, letter := range line {
			termbox.SetCell(linesMarginLeft+marginLeft+1+e, marginTop+1+i, letter, ColorPicker(lineColor), ColorPicker(boxBg))
		}
	}
}

// --
// Parameters | Type
// row, col| int
// FontColor, BackgroundColor, text | string
// --
// This prints a string at the coordinates given in the settings
// --
func TbPrint(row, col int, FontColor, BackGroundColor, text string) {
	for _, letter := range text {
		termbox.SetCell(row, col, letter, ColorPicker(FontColor), ColorPicker(BackGroundColor))
		row += runewidth.RuneWidth(letter)
	}
}

// --
// Parameters | Type
// color | string
// --
// Return Type : termbox.Attribute
// --
// This function converts a written color to a termBox color
// --
func ColorPicker(color string) termbox.Attribute {
	switch color {
	case "white":
		return termbox.ColorWhite
	case "red":
		return termbox.ColorRed
	case "bleu":
		return termbox.ColorBlue
	case "cyan":
		return termbox.ColorCyan
	case "green":
		return termbox.ColorGreen
	default:
		return termbox.ColorDefault
	}
}
