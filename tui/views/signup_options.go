package views

import (
	"github.com/arnavsurve/taskman/tui/shared"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type signupOptionsModel struct {
	list     list.Model
	selected list.Item
}

func NewSignUpOptionsModel() signupOptionsModel {
	items := []list.Item{
		item{title: "Create an account", desc: "Sign up using your email, create a username and password"},
		item{title: "Sign in with GitHub", desc: "Sign in using your GitHub account"},
	}

	m := signupOptionsModel{list: list.New(items, list.NewDefaultDelegate(), 60, 20)}
	m.list.Title = "Sign Up"
	m.list.Styles.Title = shared.TitleStyle

	return m
}

func (m signupOptionsModel) Init() tea.Cmd {
	return nil
}

func (m signupOptionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				if i.title == "Create an account" {
					return parentModel{currentModel: NewSignUpModel()}, tea.Batch(
						func() tea.Msg { return switchToSignUpMsg{} },
					)
				} else if i.title == "Sign in with GitHub" {
					// return NewOAuthModel(), nil
					return parentModel{currentModel: NewOAuthModel()}, tea.Batch(
						func() tea.Msg { return switchToOAuthMsg{} },
					)
				}
			}
		}

		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m signupOptionsModel) View() string {
	listView := m.list.View()

	styledListView := shared.FocusedStyle.Render(listView)
	// styledListView := lipgloss.NewStyle().
	// 	Foreground(lipgloss.Color("240")).
	// 	Render(listView)

	return styledListView
}
