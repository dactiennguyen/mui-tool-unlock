package colors

import "fmt"

// ANSI Color codes
const (
	// Reset
	Reset = "\033[0m"

	// Basic Colors
	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	// Bright Colors
	BrightBlack   = "\033[90m"
	BrightRed     = "\033[91m"
	BrightGreen   = "\033[92m"
	BrightYellow  = "\033[93m"
	BrightBlue    = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan    = "\033[96m"
	BrightWhite   = "\033[97m"

	// Text Styles
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"
	Blink     = "\033[5m"
	Reverse   = "\033[7m"

	// Background Colors
	BgBlack   = "\033[40m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"
)

// Themed color functions for better semantics
func Success(text string) string {
	return BrightGreen + "✅ " + text + Reset
}

func Error(text string) string {
	return BrightRed + "❌ " + text + Reset
}

func Warning(text string) string {
	return BrightYellow + "⚠️  " + text + Reset
}

func Info(text string) string {
	return BrightCyan + "ℹ️  " + text + Reset
}

func Progress(text string) string {
	return BrightBlue + "⏳ " + text + Reset
}

func Download(text string) string {
	return BrightMagenta + "⬇️  " + text + Reset
}

func Upload(text string) string {
	return BrightMagenta + "📤 " + text + Reset
}

func Browser(text string) string {
	return BrightCyan + "🌐 " + text + Reset
}

func Device(text string) string {
	return BrightYellow + "📱 " + text + Reset
}

func Unlock(text string) string {
	return BrightGreen + "🔓 " + text + Reset
}

func Lock(text string) string {
	return BrightRed + "🔒 " + text + Reset
}

func Key(text string) string {
	return BrightYellow + "🔑 " + text + Reset
}

func Email(text string) string {
	return BrightBlue + "📧 " + text + Reset
}

func Notice(text string) string {
	return BrightCyan + "📢 " + text + Reset
}

func Rocket(text string) string {
	return BrightMagenta + "🚀 " + text + Reset
}

func Package(text string) string {
	return BrightGreen + "📦 " + text + Reset
}

func Tool(text string) string {
	return BrightYellow + "🔧 " + text + Reset
}

func Save(text string) string {
	return BrightGreen + "💾 " + text + Reset
}

func Trophy(text string) string {
	return BrightYellow + "🎉 " + text + Reset
}

// Utility functions
func BoldText(text string) string {
	return Bold + text + Reset
}

func DimText(text string) string {
	return Dim + text + Reset
}

func UnderlineText(text string) string {
	return Underline + text + Reset
}

// Header creates a beautiful header with borders
func Header(title string) string {
	border := BrightCyan + "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" + Reset
	titleLine := BrightWhite + Bold + title + Reset
	return fmt.Sprintf("\n%s\n%s\n%s", border, titleLine, border)
}

// Section creates a section divider
func Section(title string) string {
	border := BrightBlue + "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" + Reset
	titleLine := BrightWhite + title + Reset
	return fmt.Sprintf("\n%s\n%s", titleLine, border)
}

// Gradient creates a gradient effect (simple version)
func Gradient(text string) string {
	// Simple gradient effect using different shades
	return BrightMagenta + text + Reset
}

// Rainbow creates rainbow text effect
func Rainbow(text string) string {
	colors := []string{BrightRed, BrightYellow, BrightGreen, BrightCyan, BrightBlue, BrightMagenta}
	result := ""
	for i, char := range text {
		color := colors[i%len(colors)]
		result += color + string(char) + Reset
	}
	return result
}

// Create styled prompts
func Prompt(text string) string {
	return BrightWhite + Bold + text + Reset
}

func Input(text string) string {
	return BrightCyan + "🔗 " + text + Reset
}
