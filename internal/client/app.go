package client

import (
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
