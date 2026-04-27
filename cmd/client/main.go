package main

import (
	"fmt"
	"os"
	"yandexgophkeeper/internal/config"

	"yandexgophkeeper/internal/client"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: register|login|add|list")
		return
	}
	cfg := config.New()

	app := client.New(cfg.ClientURL)

	switch os.Args[1] {

	case "register":
		if len(os.Args) < 4 {
			fmt.Println("usage: register <user> <pass>")
			return
		}
		err := app.Register(os.Args[2], os.Args[3])
		fmt.Println("register:", err)

	case "login":
		if len(os.Args) < 4 {
			fmt.Println("usage: login <user> <pass>")
			return
		}
		err := app.Login(os.Args[2], os.Args[3])
		fmt.Println("login:", err)

	case "add":
		if len(os.Args) < 5 {
			fmt.Println("usage: add <type> <value> <meta>")
			return
		}
		err := app.AddData(os.Args[2], os.Args[3], os.Args[4])
		fmt.Println("add:", err)

	case "list":
		data, err := app.ListData()
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		fmt.Println(data)
	case "get":
		if len(os.Args) < 3 {
			fmt.Println("usage: get <id>")
			return
		}
		res, err := app.GetData(os.Args[2])
		fmt.Println(res, err)

	case "update":
		if len(os.Args) < 6 {
			fmt.Println("usage: update <id> <type> <value> <meta>")
			return
		}
		err := app.UpdateData(os.Args[2], os.Args[3], os.Args[4], os.Args[5])
		fmt.Println("update:", err)

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("usage: delete <id>")
			return
		}
		err := app.DeleteData(os.Args[2])
		fmt.Println("delete:", err)

	default:
		fmt.Println("unknown command")
	}
}
