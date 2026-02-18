package main

import (
	CLI "expense-tracker/cmd/cli"
	"expense-tracker/internal/config"
	"os"
)

func main() {
	cli, err := CLI.NewCLI(os.Args[1:])
	if err != nil {
		panic(err)
	}
	if err := cli.Start(); err != nil {
		panic(err)
	}
	config.SaveState(cli.ExpenseTracker)
}
