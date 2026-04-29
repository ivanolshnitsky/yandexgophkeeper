package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *App) Register(username, password string) error {
	body, err := json.Marshal(credentials{username, password})
	if err != nil {
		return err
	}

	resp, err := http.Post(a.serverAddr+"/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("register failed: %d", resp.StatusCode)
	}

	return nil
}

func (a *App) Login(username, password string) error {
	body, err := json.Marshal(credentials{username, password})
	if err != nil {
		return err
	}

	resp, err := http.Post(a.serverAddr+"/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d: %s", resp.StatusCode, body)
	}

	var res map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	return SaveToken(res["token"])
}
