package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 { // The user typed a command
		commands := os.Args[1:len(os.Args)]
		handle(commands)
		return
	}
	// In case no commands were typed, show the list of commands
	fmt.Printf("%v \n", availableCommands)
}

func handle(commands []string) {
	switch commands[0] {
	case "init":
		handleInit(commands)
	case "add":
		handleAdd(commands)
	case "status":
		handleStatus()
	case "commit":
		handleCommit(commands)
	case "log":
		handleLog()
	case "vict":
		if commands[1] == "--help" {
			fmt.Printf("%v \n", availableCommands)
			return
		}
		fmt.Printf("vict: '%v' is not a vict command. See 'vict --help'. \n", commands)
	default:
		fmt.Printf("vict: '%v' is not a vict command. See 'vict --help'. \n", commands)
	}
}
