package pkg

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// This module has some basic defaults for a multistep prompt form TUI
//
// Initialize in a Cobra Run/RunE like so:
//
//       model := initGenerateModel(ctx)
//       if _, err := tea.NewProgram(model).Run(); err != nil {
//       	return err
//       }
//       return nil
//
// Example bootstrap function:
//
//       const (
//         name = iota
//         email
//       )
//
//       func initGenerateModel(ctx *pkg.AppContext) pkg.GenerateModel {
//         var prompts = make([]pkg.PromptField, 2)
//
//         // Use incrementing consts so it's easy to reference a field in a slice
//         prompts[name] = pkg.NewTextInput("Name", "Users' name", true)
//         ...
//         return pkg.GenerateModel{
//           Action: actionCallback,
//           Ctx: ctx,
//           Prompts: prompts,
//         }
//       }

//       func runGenerateConfig(m pkg.GenerateModel) tea.Cmd {
//         return func() tea.Msg {
//           err := someServerAction(struct{
//             m.Prompts[name].TextInput().Value
//           })
//           if err != nil {
//             return pkg.PromptOutputErr{Err: err}
//           }
//           return pkg.PromptCompleteEvent("string to display success")
//         }
//       }

// Prompt only styles
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

// Form Component
type GenerateModel struct {
	// Like an HTML form action, the command will implement this callback to submit the user input
	Action   func(GenerateModel) tea.Cmd
	Ctx      *AppContext
	ErrorMsg string
	// Keeps track of which input is focused
	focusIndex int
	// Optional text rendered above the form
	Header string
	// Internal request state
	isRequesting bool

	Prompts []PromptField
	// TODO: The module is responsible for creating the spinner, would be nice to make this internal
	Spinner spinner.Model
}

func (m GenerateModel) Init() tea.Cmd {
	return m.Spinner.Tick
}

func (m GenerateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	// Form validation
	case PromptValidationErr:
		m.ErrorMsg = msg.Error()
		return m, nil
	case PromptOutputErr:
		fmt.Println(fmt.Sprintf("%s %s %s", Times(), BoldStyle.Render("[ERROR]"), msg))
		return m, tea.Quit
	case PromptFieldValidEvent:
		m.Prompts[m.focusIndex].Value = string(msg)
		m.focusIndex++

		if m.focusIndex < len(m.Prompts) {
			promptType := m.Prompts[m.focusIndex].PromptType
			if promptType == TextInput {
				return m, m.Prompts[m.focusIndex].TextInput.Focus()
			}
		}

		return m, nil
	case PromptCompleteEvent:
		m.isRequesting = false
		fmt.Println(fmt.Sprintf("\n\n%s %s %s", Checkmark(), SuccessStyle.Render("[Success]"), msg))
		return m, tea.Quit

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd

	// Keyboard events
	case tea.KeyMsg:
		m.ErrorMsg = ""
		key := msg.String()
		if key == "esc" || key == "ctrl+c" {
			return m, tea.Quit
		}

		switch key {
		case "enter":
			totalPrompts := len(m.Prompts)
			if totalPrompts == m.focusIndex {
				m.isRequesting = true
				return m, m.Action(m)
			}

			if m.focusIndex < totalPrompts {
				return m, ValidateField(m.Prompts[m.focusIndex])
			}
		}
	}

	// Update the currently focused input
	if m.focusIndex < len(m.Prompts) {
		current := m.Prompts[m.focusIndex]
		switch current.PromptType {
		case TextInput:
			m.Prompts[m.focusIndex].TextInput, cmd = current.TextInput.Update(msg)
		case ListInput:
			m.Prompts[m.focusIndex].List, cmd = current.List.Update(msg)
		}
	}

	return m, cmd
}

func (m GenerateModel) View() string {
	var b strings.Builder
	b.WriteRune('\n')

	if m.Header != "" {
		b.WriteString(m.Header)
	}

	for i, p := range m.Prompts {
		if i > m.focusIndex {
			continue
		}

		// Render the actual input
		if i == m.focusIndex {
			if m.ErrorMsg != "" {
				b.WriteString(fmt.Sprintf("%s %s\n", ErrorStyle.Copy().Bold(true).Render("#"), m.ErrorMsg))
			}

			var s string
			switch p.PromptType {
			case TextInput:
				s = p.TextInput.View()
			case ListInput:
				b.WriteString(NewPromptLabel(p.Label))
				s = p.List.View()
			}
			b.WriteString(s)
		}

		// Render the entered value
		if i < m.focusIndex {
			switch p.PromptType {
			case TextInput:
				b.WriteString(p.TextInput.Prompt)
			case ListInput:
				b.WriteString(NewPromptLabel(p.Label))
			}
			b.WriteString(InputStyle.Render(p.Value))
		}
		b.WriteRune('\n')
	}

	if len(m.Prompts) == m.focusIndex {
		buttonText := fmt.Sprintf("%s Submit", PromptBulletStyle)
		b.WriteString(buttonText)
	}

	// Request results
	if m.isRequesting {
		s := fmt.Sprintf("\n%s Requesting", m.Spinner.View())
		b.WriteString(s)
	}

	return b.String()
}

// Form Actions
type PromptFieldValidEvent string
type PromptCompleteEvent string

type PromptValidationErr struct{ Err error }

func (e PromptValidationErr) Error() string {
	return e.Err.Error()
}

type PromptOutputErr struct{ Err error }

func (e PromptOutputErr) Error() string {
	return e.Err.Error()
}

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

// Convenience function to make a consistently-styled textinput component
func NewTextInput(prompt string, placeholder string, required bool) PromptField {
	input := textinput.New()
	input.Prompt = NewPromptLabel(prompt)
	input.Placeholder = placeholder

	if !required {
		input.Placeholder = strings.Join([]string{placeholder, "(optional)"}, " ")
	}

	return PromptField{
		PromptType: TextInput,
		TextInput:  input,
		IsRequired: required,
	}
}

// L
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

// Convenience function to make a single-choice, consistently styled list component
func NewList(label string, items []string) PromptField {
	listItems := []list.Item{}
	for _, i := range items {
		listItems = append(listItems, ListItem(i))
	}
	l := list.New(listItems, itemDelegate{}, 20, len(items)+2)
	l.SetShowStatusBar(false)
	l.SetShowPagination(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false)

	return PromptField{
		PromptType: ListInput,
		List:       l,
		Label:      label,
	}
}

// TODO: if needed a multi-choice list component

// Only runs if a `PromptField` has a `Validator` callback
func ValidateField(p PromptField) tea.Cmd {
	return func() tea.Msg {
		value := p.TextInput.Value()
		switch p.PromptType {
		case TextInput:
			if p.IsRequired && value == "" {
				return PromptValidationErr{fmt.Errorf("Field is required")}
			}
		case ListInput:
			v, ok := p.List.SelectedItem().(ListItem)
			if ok {
				return PromptFieldValidEvent(v)
			}
		}
		if p.Validator != nil {
			err := p.Validator(value)
			if err != nil {
				return PromptValidationErr{err}
			}
		}
		return PromptFieldValidEvent(value)
	}
}
