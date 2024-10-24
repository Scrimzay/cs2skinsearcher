package textinput

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().Background(lipgloss.Color("#01FAC6")).Foreground(lipgloss.Color("#030303")).Bold(true).Padding(0, 1, 0)
)

type (
	errMsg error
)

type Output struct {
	IsKnife bool
	IsStatTrak bool
	Weapon string
	Skin string
	Condition string
}

// func (o *Output) update(val string) {
// 	o.Output = val
// }

type model struct {
	textInput textinput.Model
	err error
	output *Output
	step int // Added to keep track of the current input step
	headers []string // List of prompts
}

func InitialTextInputModel(output *Output) model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40

	return model{
		textInput: ti,
		err: nil,
		output: output,
		step: 0,
		headers: []string{
			"Is it a StatTrak™? (Yes/No)",
			"Is it a ★ knife? (Yes/No)",
			"What weapon do you want to look up? If it uses a -, then please include it.",
			"What skin do you want to look up? If it uses a -, then please include it.",
			"What condition do you want? If it uses a -, then please include it.",
		},
	}
}

func (m *model) nextPrompt() {
	m.step++
	if m.step < len(m.headers) {
		m.textInput.SetValue("")
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			value := m.textInput.Value()
			if len(value) > 0 {
				// Store the input based on the current step
				switch m.step {
				case 0:
					lowerValue := strings.ToLower(value)
					if lowerValue == "yes" || lowerValue == "y" {
						m.output.IsStatTrak = true
					} else {
						m.output.IsStatTrak = false
					}
				case 1:
					lowerValue := strings.ToLower(value)
					if lowerValue == "yes" || lowerValue == "y" {
						m.output.IsKnife = true
					} else {
						m.output.IsKnife = false
					}
				case 2:
					m.output.Weapon = value
				case 3:
					m.output.Skin = value
				case 4:
					m.output.Condition = value
					return m, tea.Quit // All inputs collected
				}
				m.nextPrompt()
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	header := titleStyle.Render(m.headers[m.step])
	return fmt.Sprintf("%s\n\n%s\n\n",
		header,
		m.textInput.View(),
	) + "\n"
}