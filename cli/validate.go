package main

import (
	"flag"
	"log"
	"os"

	"github.com/transferwise/klaabu/klaabu"
)

func validateCommand() {
	validateFlags := flag.NewFlagSet("validate", flag.ExitOnError)

	schemaFileName := validateFlags.String("schema", "schema.kml", "Schema file path")
	validateFlags.Parse(os.Args[2:])

	schema, err := klaabu.LoadSchemaFromKmlFile(*schemaFileName)
	if err != nil {
		log.Fatalln(err)
	}

	err = schema.Validate()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Validation successful")
}
