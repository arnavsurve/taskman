package views

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/arnavsurve/taskman/tui/shared"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	statusMsg int
	errMsg    struct{ error }
)

type signupModel struct {
	focusIndex int
	inputs     []textinput.Model
}

func NewSignUpModel() tea.Model {
	m := signupModel{
		inputs: make([]textinput.Model, 3),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = shared.CursorStyle
		t.CharLimit = 32
		t.Prompt = "  "

		switch i {
		case 0:
			t.Placeholder = "Email"
			t.CharLimit = 50
			t.PromptStyle = shared.FocusedStyle
			t.TextStyle = shared.FocusedStyle
			t.Prompt = "> "
			t.Focus()
		case 1:
			t.Placeholder = "Username"
		case 2:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}

func (m signupModel) Init() tea.Cmd {
	return nil
}

func (m signupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			if s == "enter" && m.focusIndex == len(m.inputs) {
				if len(m.inputs) < 2 {
					return m, nil
				}
				email := m.inputs[0].Value()
				username := m.inputs[1].Value()
				password := m.inputs[2].Value()
				return m, m.submitForm(username, password, email)
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				if i == m.focusIndex {
					// Set focused state
					m.inputs[i].Prompt = "> "
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = shared.FocusedStyle
					m.inputs[i].TextStyle = shared.FocusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Prompt = "  "
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = shared.NoStyle
				m.inputs[i].TextStyle = shared.NoStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)
	return m, cmd
}

func NewCustomTextInput(prompt string, focused bool) textinput.Model {
	t := textinput.New()
	t.Prompt = prompt

	if focused {
		t.Focus()
	}

	return t
}

func (m signupModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m signupModel) submitForm(username, password, email string) tea.Cmd {
	return func() tea.Msg {
		data := map[string]string{
			"username": username,
			"password": password,
			"email":    email,
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Printf("Error: %v", err)
			return errMsg{err}
		}

		url := fmt.Sprintf("%s/user", shared.URL)

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Error: %v", err)
			return errMsg{err}
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return errMsg{err}
		}

		log.Printf("\nResponse: %v", result)
		log.Printf("\nStatus: %d", resp.StatusCode)
		return statusMsg(resp.StatusCode)
	}
}

func (m signupModel) View() string {
	var b strings.Builder
	b.WriteString("\n\n")

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &shared.BlurredButton
	if m.focusIndex == len(m.inputs) {
		button = &shared.FocusedButton
	}
	fmt.Fprintf(&b, "\n\n %s\n\n", *button)

	b.WriteString(shared.HelpStyle.Render("ctrl+c or esc to quit"))

	return b.String()
}
