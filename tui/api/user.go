package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) submitForm(username, password, email string) tea.Cmd {
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

		log.Printf("Response: %v", result)
		log.Printf("Status: %d", resp.StatusCode)
		return statusMsg(resp.StatusCode)
	}
}
