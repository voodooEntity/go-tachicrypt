package prettywriter

import (
	"fmt"
	"strings"
)

// Color represents the available colors for text formatting.
type Color int

const (
	// Foreground colors
	Black   Color = 30
	Red     Color = 31
	Green   Color = 32
	Yellow  Color = 33
	Blue    Color = 34
	Magenta Color = 35
	Cyan    Color = 36
	White   Color = 37

	// Background colors

	BlackBG   Color = 40
	RedBG     Color = 41
	GreenBG   Color = 42
	YellowBG  Color = 43
	BlueBG    Color = 44
	MagentaBG Color = 45
	CyanBG    Color = 46
	WhiteBG   Color = 47
)

// Write writes the given message to the console with the specified foreground and background colors.
func Write(message string, foregroundColor, backgroundColor Color) {
	fmt.Print("\033[", int(foregroundColor), ";", int(backgroundColor), "m", message, "\033[0m")
}

// Writeln writes the given message to the console with the specified foreground and background colors, followed by a newline.
func Writeln(message string, foregroundColor, backgroundColor Color) {
	Write(message, foregroundColor, backgroundColor)
	fmt.Println()
}

// Writef writes the formatted message to the console with the specified foreground and background colors.
func Writef(format string, foregroundColor, backgroundColor Color, a ...interface{}) {
	fmt.Printf(fmt.Sprint("\033[", int(foregroundColor), ";", int(backgroundColor), "m", format, "\033[0m"), a...)
}

// Writefln writes the formatted message to the console with the specified foreground and background colors, followed by a newline.
func Writefln(format string, foregroundColor, backgroundColor Color, a ...interface{}) {
	Writef(format, foregroundColor, backgroundColor, a...)
	fmt.Println()
}

// Style represents the border style of the box.
type Style int

const (
	SingleLine Style = iota
	DoubleLine
)

// WriteInBox prints the message in a box format within the provided shell width.
func WriteInBox(shellWidth int, message string, foregroundColor, backgroundColor Color, style Style) {
	lines := wrapMessage(message, shellWidth-4) // Account for box corners and padding

	// Calculate top and bottom borders
	topBorder := "┌" + getBorder(shellWidth, style) + "┐"
	bottomBorder := "└" + getBorder(shellWidth, style) + "┘"

	// Print top border
	fmt.Println(topBorder)

	// Iterate through wrapped lines
	for _, line := range lines {
		// Calculate padding based on line length
		padding := strings.Repeat(" ", shellWidth-len(line)-3) // Account for corners, spaces

		// Print line with color and padding
		fmt.Printf("|\033[%d;%dm %s %s \033[0m|\n", int(foregroundColor), int(backgroundColor), line, padding)
	}

	// Print bottom border
	fmt.Println(bottomBorder)
}

// wrapMessage splits the message into lines based on the provided width, breaking at spaces.
func wrapMessage(message string, maxWidth int) []string {
	words := strings.Fields(message)
	var lines []string
	line := ""

	for _, word := range words {
		if len(line)+len(word)+1 > maxWidth { // Account for spaces
			lines = append(lines, line)
			line = ""
		}
		line += word + " "
	}

	lines = append(lines, strings.TrimSpace(line)) // Add the last line

	return lines
}

// getBorder generates the border string based on the shell width and style.
func getBorder(width int, style Style) string {
	borderChar := "-"
	if style == DoubleLine {
		borderChar = "="
	}
	return strings.Repeat(borderChar, width)
}
