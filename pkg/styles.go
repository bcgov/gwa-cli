package pkg

import (
	"github.com/charmbracelet/lipgloss"
)

var ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#d8292f"))
var SuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
var WarningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#fcba19"))

func Indeterminate() string {
	return WarningStyle.Render("-")
}

func Checkmark() string {
	return SuccessStyle.Render("✓")
}

func Times() string {
	return ErrorStyle.Render("x")
}

func PrintSuccess(output string) string {
	return SuccessStyle.Render(output)
}

func PrintError(output string) string {
	return ErrorStyle.Render(output)
}

func PrintWarning(output string) string {
	return WarningStyle.Render(output)
}
