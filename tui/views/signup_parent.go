package views

import (
	tea "github.com/charmbracelet/bubbletea"
)

type mode int

const (
	choiceMode mode = iota
	signUpMode
	oauthMode
)

type parentModel struct {
	currentModel tea.Model
}

type (
	switchToOAuthMsg  struct{}
	switchToSignUpMsg struct{}
)

// NewParentModel initializes the parent model with the child models.
func NewParentModel() parentModel {
	return parentModel{
		currentModel: NewSignUpOptionsModel(),
	}
}

func (m parentModel) Init() tea.Cmd {
	return m.currentModel.Init()
}

func (m parentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case switchToOAuthMsg:
		m.currentModel = NewOAuthModel()
		return m, m.currentModel.Init()
	case switchToSignUpMsg:
		m.currentModel = NewSignUpModel()
		return m, m.currentModel.Init()
	default:
		var cmd tea.Cmd
		m.currentModel, cmd = m.currentModel.Update(msg)
		return m, cmd
	}
}

func (m parentModel) View() string {
	return m.currentModel.View()
}
