package utils

import "fmt"

const (
	ColorReset   = "\x1b[0m"
	ColorBold    = "\x1b[1m"
	ColorRegular = "\x1b[2m"
	ColorRed     = "\x1b[31m"
	ColorGreen   = "\x1b[32m"
	ColorYellow  = "\x1b[33m"
	ColorBlue    = "\x1b[34m"
)

func PrintError(err error) {
	fmt.Println(ColorRed + err.Error() + ColorReset)
}

func PrintSuccess(msg string) {
	fmt.Println(ColorGreen + msg + ColorReset)
}

func PrintWarning(msg string) {
	fmt.Println(ColorYellow + msg + ColorReset)
}

func PrintInfo(msg string) {
	fmt.Println(ColorBlue + msg + ColorReset)
}

func PrintBold(msg string) {
	fmt.Println(ColorBold + msg + ColorReset)
}

func Bold(msg string) string {
	return ColorBold + msg + ColorRegular
}

func Red(msg string) string {
	return ColorRed + msg
}

func Green(msg string) string {
	return ColorGreen + msg
}

func Yellow(msg string) string {
	return ColorYellow + msg
}

func Blue(msg string) string {
	return ColorBlue + msg
}

func Warning(msg string) string {
	return ColorYellow + msg
}

func Info(msg string) string {
	return ColorBlue + msg
}

func Error(msg string) string {
	return ColorRed + msg
}

func Success(msg string) string {
	return ColorGreen + msg
}
