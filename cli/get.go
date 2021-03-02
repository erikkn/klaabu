package main

import (
	"flag"
	"github.com/erikkn/klaabu/klaabu"
	"log"
	"os"
)

func getCommand() {
	get := flag.NewFlagSet("get", flag.ExitOnError)
	schemaFileName := get.String("schema", "schema.kml", "(Optional) Schema file path (defaults to ./schema.kml)")
	err := get.Parse(os.Args[2:])

	prefixId := get.Arg(0)

	if err != nil || len(os.Args) < 3 {
		log.Printf("Usage: klaabu get [OPTIONS] PREFIX_ID \n\n Subcommands: \n")
		get.PrintDefaults()
		os.Exit(1)
	}

	schema, err := klaabu.LoadSchemaFromKmlFile(*schemaFileName)
	if err != nil {
		log.Fatalln(err)
	}

	prefix := schema.PrefixById(prefixId)
	if prefix == nil {
		log.Fatalf("not found: %s", prefixId)
	}

	// TODO: Output all the individual fields, might change again so a `todo` for later.
	log.Println(prefix)
}
