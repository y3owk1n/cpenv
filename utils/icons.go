package utils

import "github.com/fatih/color"

const (
	Checkmark = "✓"
	Cross     = "✖"
	Info      = "ℹ"
	Warn      = "⚠"
)

func ColoredIcon(icon string, fgColor color.Attribute) string {
	return color.New(fgColor).Sprint(icon)
}

func SuccessIcon() string {
	return ColoredIcon(Checkmark, color.FgGreen)
}

func ErrorIcon() string {
	return ColoredIcon(Cross, color.FgRed)
}

func WarningIcon() string {
	return ColoredIcon(Warn, color.FgYellow)
}

func InfoIcon() string {
	return ColoredIcon(Info, color.FgBlue)
}

func WhiteText(text string) string {
	return color.New(color.FgWhite).Sprint(text)
}

func CyanText(text string) string {
	return color.New(color.FgCyan).Sprint(text)
}

func GreenText(text string) string {
	return color.New(color.FgGreen).Sprint(text)
}
