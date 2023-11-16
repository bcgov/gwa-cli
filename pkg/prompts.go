package pkg

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	// "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// type FieldType int
//
// const (
//
//	TextInput FieldType = iota
//	RadioInput
//
// )
//
//	type Component interface {
//		View() string
//	}
var (
	PromptSymbolStyle = SuccessStyle.Copy().Bold(true)
	PromptErrorStyle  = ErrorStyle.Copy().Bold(true)
	PromptStyle       = lipgloss.NewStyle().Bold(true)
	PromptBulletStyle = SuccessStyle.Copy().Bold(true).Render("?")
	BoldStyle         = lipgloss.NewStyle().Bold(true)
	FocusedStyle      = SuccessStyle.Copy().Bold(true)
	InputStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))

	ItemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	SelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("6"))

	ButtonBlurred = BoldStyle.Copy().Render("? Submit")
	ButtonFocused = FocusedStyle.Copy().Render("> Submit")
)

type PromptType int

const (
	TextInput PromptType = iota
	ListInput
)

type PromptField struct {
	PromptType
	Label      string
	IsRequired bool
	Value      string
	TextInput  textinput.Model
	List       list.Model
	Validator  func(input string) error
}

func NewPromptLabel(label string) string {
	return fmt.Sprintf("%s %s: ", PromptBulletStyle, label)
}

func NewPromptError(err error) tea.Cmd {
	return func() tea.Msg {
		return err
	}
}

//	func NewTextInput(prompt string, placeholder string, focus bool) PromptField {
//		input := textinput.New()
//		input.Prompt = prompt
//		input.PromptStyle = PromptStyle
//		input.Placeholder = placeholder
//		if focus {
//
//			input.Focus()
//		}
//
//		return PromptField{
//			Input: input,
//			Type:  TextInput,
//		}
//	}
type ListItem string

func (i ListItem) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(ListItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return SelectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func NewList(items []list.Item) list.Model {
	l := list.New(items, itemDelegate{}, 20, len(items)+2)
	l.SetShowStatusBar(false)
	l.SetShowPagination(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	return l
}
