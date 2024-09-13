package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"net/http"
	"os/exec"
	"time"
)

type oauthModel struct {
	authenticated bool
	waiting       bool
}

func NewOAuthModel() oauthModel {
	return oauthModel{
		authenticated: false,
		waiting:       false,
	}
}

func (m oauthModel) Init() tea.Cmd {
	// Open the OAuth URL in the browser
	return openOAuthURL
}

func (m oauthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Allow quitting with "esc" or "ctrl+c" during waiting/authentication
		switch msg.String() {
		case "esc", "ctrl+c":
			return m, tea.Quit
		}

	case authCompleteMsg:
		// Authentication complete
		m.authenticated = true
		m.waiting = false
		return m, nil

	case tea.Msg:
		if m.waiting {
			return m, checkOAuthStatus()
		}
	}

	return m, nil
}

func (m oauthModel) View() string {
	if m.authenticated {
		return "Successfully authenticated!\n\nPress Esc to quit."
	}

	if m.waiting {
		return "Waiting for authentication...\n\nPress Esc to cancel."
	}

	return "Opening browser for GitHub OAuth...\n\nPlease complete authentication in your browser."
}

// openOAuthURL opens the OAuth URL in the user's default browser.
func openOAuthURL() tea.Msg {
	// Replace this with your actual OAuth endpoint URL
	oauthURL := "https://localhost:8080/login/github"

	// Use os/exec to open the default browser
	cmd := exec.Command("open", oauthURL)

	// Sample code for implementing multiple operating systems (currently assuming user is on darwin/MacOS)
	// var cmd *exec.Cmd
	// switch {
	// case isWindows():
	// 	cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", oauthURL)
	// case isMac():
	// 	cmd = exec.Command("open", oauthURL)
	// case isLinux():
	// 	cmd = exec.Command("xdg-open", oauthURL)
	// }

	if err := cmd.Start(); err != nil {
		log.Printf("Error opening browser: %v", err)
	}

	// While the user is authenticating, switch to the "waiting" mode
	return authWaitingMsg{}
}

// authWaitingMsg indicates that the app is waiting for the user to authenticate.
type authWaitingMsg struct{}

// authCompleteMsg indicates that the authentication process has completed successfully.
type authCompleteMsg struct{}

// checkOAuthStatus polls the server to check whether the authentication has completed.
func checkOAuthStatus() tea.Cmd {
	return func() tea.Msg {
		// Simulate waiting for OAuth completion by polling the backend.
		// Replace this with an actual HTTP request or polling mechanism.
		time.Sleep(2 * time.Second)

		// In a real-world scenario, you'd check the authentication status
		// from your backend here. Simulating success for demo purposes.
		response, err := http.Get("YOUR_BACKEND_OAUTH_STATUS_ENDPOINT")
		if err != nil || response.StatusCode != http.StatusOK {
			// Continue polling by sending the authWaitingMsg again
			return authWaitingMsg{}
		}

		// If authentication is complete, return an authCompleteMsg
		return authCompleteMsg{}
	}
}

// Helper functions to detect OS
// func isWindows() bool {
// 	return "windows" == "YOUR_OS_DETECTION_LOGIC"
// }
//
// func isMac() bool {
// 	return "darwin" == "YOUR_OS_DETECTION_LOGIC"
// }
//
// func isLinux() bool {
// 	return "linux" == "YOUR_OS_DETECTION_LOGIC"
// }
