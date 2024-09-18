package views

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strings"

	"github.com/arnavsurve/taskman/tui/shared"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	authWaitingMsg  struct{}
	authCompleteMsg struct{}
)

type oauthModel struct {
	authenticated bool
}

func NewOAuthModel() oauthModel {
	return oauthModel{
		authenticated: false,
	}
}

func (m oauthModel) Init() tea.Cmd {
	// Open the OAuth URL in the browser
	return openOAuthURL()
}

func (m oauthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Allow quitting with "esc" or "ctrl+c"
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			return m, tea.Quit
		}

	case authWaitingMsg:
		// Waiting for authentication
		m.authenticated = false
		return m, nil

	case authCompleteMsg:
		// Authentication complete
		m.authenticated = true
		return m, nil

	case tea.Msg:
		if m.authenticated == false {
			return m, checkOAuthStatus()
		}
	}

	return m, nil
}

func (m oauthModel) View() string {
	var b strings.Builder

	// b.WriteString(shared.TitleStyle.Render("Sign Up"))
	b.WriteString("\n\n")
	b.WriteString(" Opening browser for GitHub OAuth...\n\n Please complete the sign up process in your browser.")

	if m.authenticated {
		b.WriteString("Successfully authenticated!\n\nPress Esc to quit.")
	}

	if m.authenticated == false {
		b.WriteString(shared.BlurredStyle.Render("\n\n Waiting for authentication...\n\n Press esc/q to cancel."))
	}

	return b.String()
}

// openOAuthURL opens the OAuth URL in the user's default browser.
func openOAuthURL() tea.Cmd {
	return func() tea.Msg {
		fmt.Println("openOAuthURL called")
		oauthURL := "http://localhost:8080/oauth2/github"

		var cmd *exec.Cmd
		switch {
		case isWindows():
			cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", oauthURL)
		case isMac():
			cmd = exec.Command("open", oauthURL)
		case isLinux():
			cmd = exec.Command("xdg-open", oauthURL)
		}

		if err := cmd.Start(); err != nil {
			log.Printf("Error opening browser: %v", err)
			return nil
		}

		return authWaitingMsg{}
	}
}

// checkOAuthStatus polls the server to check whether the authentication has completed.
func checkOAuthStatus() tea.Cmd {
	return func() tea.Msg {
		// Simulate waiting for OAuth completion by polling the backend.
		// Replace this with an actual HTTP request or polling mechanism.
		// time.Sleep(2 * time.Second)

		// In a real-world scenario, you'd check the authentication status
		// from your backend here. Simulating success for demo purposes.
		response, err := http.Get("YOUR_BACKEND_OAUTH_STATUS_ENDPOINT")
		if err != nil || response.StatusCode != http.StatusOK {
			return authWaitingMsg{}
		}

		return authCompleteMsg{}
	}
}

// Helper functions to detect OS
func isMac() bool {
	return strings.Contains(strings.ToLower(runtime.GOOS), "darwin")
}

func isLinux() bool {
	return strings.Contains(strings.ToLower(runtime.GOOS), "linux")
}

func isWindows() bool {
	return strings.Contains(strings.ToLower(runtime.GOOS), "windows")
}
