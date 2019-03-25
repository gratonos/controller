package controller

import (
	"fmt"
)

type colorID int

const (
	black colorID = iota + 30
	red
	green
	yellow
	blue
	magenta
	cyan
	white
)

func colorize(text string, color colorID) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", color, text)
}
