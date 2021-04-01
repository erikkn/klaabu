package main

import (
	"log"
	"os"
)

var commands = map[string]func(){
	"export-terraform": exportTerraformCommand,
	"find":             findCommand,
	"fmt":              fmtCommand,
	"get":              getCommand,
	"init":             initCommand,
	"space":            spaceCommand,
	"validate":         validateCommand,
	"version":          versionCommand,
}

func helpCommand(c string) {
	log.Printf("Usage: klaabu COMMAND [OPTIONS] [ARGS] \n\n Common commands: ")
	for key := range commands {
		log.Printf("    %s \n", key)
	}
	os.Exit(1)
}

func main() {
	// Remove the timestamp prefix
	log.SetFlags(0)

	if len(os.Args) < 2 {
		helpCommand("")
	}

	commandName := os.Args[1]
	c, ok := commands[commandName]
	if !ok {
		helpCommand(commandName)
	}

	c()
}
