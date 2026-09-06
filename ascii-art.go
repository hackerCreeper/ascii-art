package ascii-art

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	Reset       = "\033[0m"
	ColorText   = "\033[1;37m"
	ColorRain   = "\033[32m"  
	ColorDim    = "\033[2;32m"  
	MoveCursor  = "\033[H"   
	ClearScreen = "\033[H\033[2J"
	HideCursor  = "\033[?25l"
	ShowCursor  = "\033[?25h"
)

const (
	Width  = 80
	Height = 14
)

const E = ' '

var Font = map[rune][5][3]rune{
	' ': {{E, E, E}, {E, E, E}, {E, E, E}, {E, E, E}, {E, E, E}},
	'A': {{E, '-', E}, {'/', ' ', '\\'}, {'|', '-', '|'}, {'|', ' ', '|'}, {'|', ' ', '|'}},
	'B': {{'-', '-', E}, {'|', ' ', '\\'}, {'-', '-', E}, {'|', ' ', '\\'}, {'-', '-', E}},
	'C': {{E, '_', '_'}, {'/', ' ', E}, {'|', ' ', E}, {'|', ' ', E}, {'\\', '_', '/'}},
	'D': {{'-', '-', E}, {'|', ' ', '\\'}, {'|', ' ', '|'}, {'|', ' ', '|'}, {'-', '-', E}},
	'E': {{'-', '-', '-'}, {'|', '_', '_'}, {'|', ' ', E}, {'|', '_', '_'}, {'-', '-', '-'}},
	'F': {{'-', '-', '-'}, {'|', '_', '_'}, {'|', ' ', E}, {'|', ' ', E}, {'|', ' ', E}},
	'G': {{E, '_', '_'}, {'/', ' ', E}, {'|', ' ', '-'}, {'|', ' ', '|'}, {'\\', '_', '/'}},
	'H': {{'|', ' ', '|'}, {'|', ' ', '|'}, {'|', '-', '|'}, {'|', ' ', '|'}, {'|', ' ', '|'}},
	'I': {{'-', '-', '-'}, {E, '|', E}, {E, '|', E}, {E, '|', E}, {'-', '-', '-'}},
	'J': {{E, E, '-'}, {E, E, '|'}, {E, E, '|'}, {'|', ' ', '|'}, {'\\', '_', '/'}},
	'K': {{'|', ' ', '/'}, {'|', '_', '\\'}, {'|', ' ', '\\'}, {'|', ' ', '\\'}, {'|', ' ', '\\'}},
	'L': {{'|', ' ', E}, {'|', ' ', E}, {'|', ' ', E}, {'|', ' ', E}, {'|', '_', '_'}},
	'M': {{'|', '\\', '|'}, {'|', '|', '|'}, {'|', ' ', '|'}, {'|', ' ', '|'}, {'|', ' ', '|'}},
	'N': {{'|', '\\', ' '}, {'|', ' ', '|'}, {'|', ' ', '|'}, {'|', ' ', '|'}, {'|', ' ', '|'}},
	'O': {{E, '_', E}, {'/', ' ', '\\'}, {'|', ' ', '|'}, {'|', ' ', '|'}, {'\\', '_', '/'}},
	'P': {{'-', '-', E}, {'|', ' ', '|'}, {'-', '-', E}, {'|', ' ', E}, {'|', ' ', E}},
	'Q': {{E, '_', E}, {'/', ' ', '\\'}, {'|', ' ', '|'}, {'\\', '_', '/'}, {E, E, '\\'}},
	'R': {{'-', '-', E}, {'|', ' ', '|'}, {'-', '-', E}, {'|', ' ', '\\'}, {'|', ' ', '\\'}},
	'S': {{E, '_', '_'}, {'/', '_', '_'}, {'\\', '_', '_'}, {E, E, '\\'}, {'_', '_', '/'}},
	'T': {{'-', '-', '-'}, {E, '|', E}, {E, '|', E}, {E, '|', E}, {E, '|', E}},
	'U': {{'|', ' ', '|'}, {'|', ' ', '|'}, {'|', ' ', '|'}, {'|', ' ', '|'}, {'\\', '_', '/'}},
	'V': {{'|', ' ', '|'}, {'|', ' ', '|'}, {'\\', ' ', '/'}, {'\\', ' ', '/'}, {E, 'V', E}},
	'W': {{'|', ' ', '|'}, {'|', ' ', '|'}, {'|', ' ', '|'}, {'|', '\\', '|'}, {E, 'V', E}},
	'X': {{'\\', ' ', '/'}, {'\\', ' ', '/'}, {E, 'X', E}, {'/', ' ', '\\'}, {'/', ' ', '\\'}},
	'Y': {{'|', ' ', '|'}, {'\\', ' ', '/'}, {E, 'V', E}, {E, '|', E}, {E, '|', E}},
	'Z': {{'-', '-', '-'}, {E, E, '/'}, {E, '/', E}, {'/', ' ', E}, {'-', '-', '-'}},
	'0': {{E, '_', E}, {'/', ' ', '\\'}, {'|', ' ', '|'}, {'|', ' ', '|'}, {'\\', '_', '/'}},
	'1': {{E, '/', E}, {E, '|', E}, {E, '|', E}, {E, '|', E}, {E, '|', E}},
	'2': {{'-', '-', E}, {E, E, '\\'}, {E, '-', '/'}, {'/', ' ', E}, {'-', '-', '-'}},
	'3': {{'-', '-', E}, {E, E, '\\'}, {'-', '-', E}, {E, E, '\\'}, {'-', '-', 'E'}},
	'4': {{'|', ' ', '|'}, {'|', ' ', '|'}, {'-', '-', '-'}, {E, E, '|'}, {E, E, '|'}},
	'5': {{'-', '-', '-'}, {'|', '_', '_'}, {E, '_', '\\'}, {E, E, '|'}, {'_', '_', '/'}},
	'6': {{E, '_', E}, {'/', ' ', E}, {'|', '_', E}, {'|', ' ', '|'}, {'\\', '_', '/'}},
	'7': {{'-', '-', '-'}, {E, E, '/'}, {E, '/', E}, {E, '/', E}, {E, '/', E}},
	'8': {{E, '_', E}, {'/', ' ', '\\'}, {E, '=', E}, {'/', ' ', '\\'}, {'\\', '_', '/'}},
	'9': {{E, '_', E}, {'/', ' ', '|'}, {'\\', '_', '|'}, {E, E, '|'}, {'_', '_', '/'}},
}

var RainDrops [Width]int

func InitRain() {
	for i := 0; i < Width; i++ {
		rainDrops[i] = rand.Intn(Height)
	}
}

func GetTextBox(text string) ([5][Width]rune, [5][Width]bool) {
	var textMask [5][Width]rune
	var shadowMask [5][Width]bool
	
	text = strings.ToUpper(text)
	
	charStride := 3 + 2 
	textWidth := len(text) * charStride
	
	startX := (Width - textWidth) / 2
	if startX < 0 {
		startX = 0
	}
	
	for row := 0; row < 5; row++ {
		for c := 0; c < Width; c++ {
			textMask[row][c] = ' '
		}
	}
	
	for row := 0; row < 5; row++ {
		currentX := startX
		for _, char := range text {
			matrix, exists := font[char]
			if !exists {
				matrix = font[' ']
			}
			
			for c := 0; c < 3; c++ {
				symbol := matrix[row][c]
				if symbol != ' ' && currentX < Width {
					textMask[row][currentX] = symbol
					
					for shadowX := currentX - 2; shadowX <= currentX+2; shadowX++ {
						if shadowX >= 0 && shadowX < Width {
							shadowMask[row][shadowX] = true
						}
					}
				}
				currentX++
			}
			currentX += 2 
		}
	}
	return textMask, shadowMask
}

func Start() {
	inputText := "go 2026"

	if len(os.Args) > 1 {
		inputText = strings.Join(os.Args[1:], " ")
	}

	rand.Seed(time.Now().UnixNano())
	initRain()

	specialChars := []rune{'!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '-', '+', '/', '<', '>', '?', ':', ';'}
	textMask, shadowMask := getTextBox(inputText)

	fmt.Print(HideCursor)
	fmt.Print(ClearScreen)
	
	for {
		fmt.Print(MoveCursor)
		
		var output strings.Builder

		for y := 0; y < Height; y++ {
			for x := 0; x < Width; x++ {
				
				isTextZone := y >= 4 && y < 9
				
				if isTextZone {
					textY := y - 4
					symbol := textMask[textY][x]
					
					if symbol != ' ' {
						output.WriteString(ColorText + string(symbol) + Reset)
						continue
					}
					
					if shadowMask[textY][x] {
						output.WriteString(" ")
						continue
					}
				}

				if rainDrops[x] == y {
					char := specialChars[rand.Intn(len(specialChars))]
					output.WriteString(ColorRain + string(char) + Reset)
				} else if rainDrops[x] == y-1 || rainDrops[x] == y-2 {
					char := specialChars[rand.Intn(len(specialChars))]
					output.WriteString(ColorDim + string(char) + Reset)
				} else {
					output.WriteString(" ")
				}
			}
			output.WriteString("\n")
		}

		fmt.Print(output.String())

		for i := 0; i < Width; i++ {
			rainDrops[i] = (rainDrops[i] + 1) % Height
		}

		time.Sleep(80 * time.Millisecond)
	}
}
