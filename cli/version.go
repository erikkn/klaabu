package main

import (
	"flag"
	"log"
	"os"
)

var Version string

func versionCommand() {
	commandName := "version"
	flagSet := flag.NewFlagSet(commandName, flag.ExitOnError)
	flagSet.Parse(os.Args[2:])

	if len(os.Args) != 2 {
		log.Printf("Usage: klaabu %s\n", flagSet.Name())
		flagSet.PrintDefaults()
		os.Exit(1)
	}

	if Version != "" {
		log.Println(Version)
	} else {
		log.Println("development")
	}
}
