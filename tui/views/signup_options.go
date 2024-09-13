package views

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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
				if i.title == "Sign up" {
					return NewSignUpModel(), nil
				} else if i.title == "Sign in with GitHub" {
					return NewOAuthModel(), nil
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
	return m.list.View()
}

func NewSignUpOptionsModel() signupOptionsModel {
	items := []list.Item{
		item{title: "Sign up", desc: "Sign up with your email, create a username and password"},
		item{title: "Sign in with GitHub", desc: "Sign in using your GitHub account"},
	}

	m := signupOptionsModel{list: list.New(items, list.NewDefaultDelegate(), 60, 20)}
	m.list.Title = "Sign Up"

	return m
}
