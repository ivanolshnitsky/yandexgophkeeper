package main

import (
	"encoding/json"
	"fmt"
	"os"

	"yandexgophkeeper/internal/client"
	"yandexgophkeeper/internal/config"
)

var (
	buildVersion = "dev"
	buildDate    = "unknown"
)

func printUsage() {
	fmt.Println(`
	Usage:
	  register <user> <pass>
	  login <user> <pass>
	  add <type> <value> <meta>
	  list
	  get <id>
	  update <id> <type> <value> <meta>
	  delete <id>
	  version
	`)
}

func printVersion() {
	fmt.Println("version:", buildVersion)
	fmt.Println("date:", buildDate)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	cfg, err := config.New()
	if err != nil {
		fmt.Println("config error:", err)
		return
	}
	app := client.New(cfg.ClientURL)

	switch os.Args[1] {

	case "register":
		if len(os.Args) < 4 {
			fmt.Println("usage: register <user> <pass>")
			return
		}
		if err := app.Register(os.Args[2], os.Args[3]); err != nil {
			fmt.Println("register error:", err)
			return
		}
		fmt.Println("ok")

	case "login":
		if len(os.Args) < 4 {
			fmt.Println("usage: login <user> <pass>")
			return
		}
		if err := app.Login(os.Args[2], os.Args[3]); err != nil {
			fmt.Println("login error:", err)
			return
		}
		fmt.Println("ok")

	case "add":
		if len(os.Args) < 5 {
			fmt.Println("usage: add <type> <json_value> <meta>")
			return
		}

		if err := app.AddData(os.Args[2], json.RawMessage(os.Args[3]), os.Args[4]); err != nil {
			fmt.Println("error:", err)
			return
		}
		fmt.Println("ok")

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
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		fmt.Println(res)

	case "update":
		if len(os.Args) < 6 {
			fmt.Println("usage: update <id> <type> <value> <meta>")
			return
		}

		if err := app.UpdateData(os.Args[2], os.Args[3], os.Args[4], os.Args[5]); err != nil {
			fmt.Println("update error:", err)
			return
		}
		fmt.Println("ok")

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("usage: delete <id>")
			return
		}

		if err := app.DeleteData(os.Args[2]); err != nil {
			fmt.Println("delete error:", err)
			return
		}
		fmt.Println("ok")

	case "version":
		printVersion()

	default:
		fmt.Println("unknown command")
		printUsage()
	}
}
