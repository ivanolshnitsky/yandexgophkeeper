package client

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		return fmt.Errorf("login failed: %d", resp.StatusCode)
	}

	var res map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	return SaveToken(res["token"])
}
