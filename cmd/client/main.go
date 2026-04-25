package main

import (
	"fmt"
	"os"

	"yandexgophkeeper/internal/client"
	"yandexgophkeeper/internal/config"
)

func main() {
	cfg := config.New()

	if len(os.Args) < 2 {
		fmt.Println("usage:")
		fmt.Println("  ping - check server")
		return
	}

	app := client.New(cfg.BaseURL)

	cmd := os.Args[1]

	switch cmd {
	case "ping":
		handlePing(app)
	default:
		fmt.Println("unknown command:", cmd)
	}
}

func handlePing(app *client.App) {
	resp, err := app.Ping()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(resp)
}
