package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type App struct {
	serverAddr string
	client     *http.Client
}

func New(addr string) *App {
	return &App{
		serverAddr: addr,
		client:     &http.Client{},
	}
}

func (a *App) Ping() (string, error) {
	resp, err := a.client.Get(a.serverAddr + "/ping")
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// AddData
func (a *App) AddData(t string, value any, meta string) error {
	token, err := LoadToken()
	if err != nil {
		return err
	}

	reqBody := map[string]any{
		"type":  t,
		"value": value,
		"meta":  meta,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, a.serverAddr+"/data", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("add failed: %d", resp.StatusCode)
	}

	return nil
}

// ListData
func (a *App) ListData() ([]map[string]any, error) {
	token, err := LoadToken()
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodGet, a.serverAddr+"/data", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data []map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&data)

	return data, nil
}

// GetData
func (a *App) GetData(id string) (string, error) {
	token, err := LoadToken()
	if err != nil {
		return "", err
	}

	req, _ := http.NewRequest("GET", a.serverAddr+"/data/"+id, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := a.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}

// UpdateData
func (a *App) UpdateData(id, t, value, meta string) error {
	token, err := LoadToken()
	if err != nil {
		return err
	}

	body, _ := json.Marshal(map[string]string{
		"type":  t,
		"value": value,
		"meta":  meta,
	})

	req, _ := http.NewRequest("PUT", a.serverAddr+"/data/"+id, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	_, err = a.client.Do(req)
	return err
}

// DeleteData
func (a *App) DeleteData(id string) error {
	token, err := LoadToken()
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("DELETE", a.serverAddr+"/data/"+id, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	_, err = a.client.Do(req)
	return err
}
