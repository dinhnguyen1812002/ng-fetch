package utils

import "github.com/fatih/color"

const (
	Green = "green"
	Blue  = "blue"
	Red   = "red"
)

func PrintColored(text, colorType string) {
	switch colorType {
	case Green:
		color.Green(text)
	case Blue:
		color.Blue(text)
	case Red:
		color.Red(text)
	default:
		color.White(text)
	}
}
