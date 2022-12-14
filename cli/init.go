package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/transferwise/klaabu/klaabu"
)

func initCommand() {
	init := flag.NewFlagSet("init", flag.ExitOnError)

	// TODO change to schema.kml
	schemaFileName := init.String("schema", "schema-new.kml", "Schema file path")
	version := init.String("version", "v1", "Schema version")
	labelsString := init.String("labels", "", "(Optional) Schema labels, e.g. 'org=acme,env=staging'")

	init.Parse(os.Args[2:])

	if len(os.Args) < 3 {
		log.Printf("Usage: klaabu init [OPTIONS] \n\n Subcommands: \n")
		init.PrintDefaults()
		os.Exit(1)
	}

	if *version != "v1" {
		log.Printf("%s is not a valid schema version \n\n", *version)
		log.Printf("Usage: klaabu init [OPTIONS] \n\n Subcommands: \n")
		init.PrintDefaults()
	}

	labels := make(map[string]string)
	if *labelsString != "" {
		for _, v := range strings.Split(*labelsString, ",") {
			pair := strings.Split(v, "=")
			if len(pair) != 2 {
				log.Fatalln("error with labels")
			}

			labels[strings.TrimSpace(pair[0])] = strings.TrimSpace(pair[1])
		}
	}

	schema := klaabu.NewSchema(labels)

	err := klaabu.WriteSchemaToFile(schema, schemaFileName)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully created your new schema '%s' \n", *schemaFileName)
}
