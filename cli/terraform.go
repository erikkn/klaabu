package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/erikkn/klaabu/klaabu"
	"github.com/erikkn/klaabu/klaabu/terraform"
)

func exportTerraformCommand() {
	export := flag.NewFlagSet("export-terraform", flag.ExitOnError)

	schemaFileName := export.String("schema", "schema.kml", "Schema file path")
	export.Parse(os.Args[2:])

	if len(os.Args) < 2 {
		log.Printf("Usage: klaabu export-terraform [OPTIONS] \n\n Subcommands: \n")
		export.PrintDefaults()
		os.Exit(1)
	}

	schema, err := klaabu.LoadSchemaFromKmlFile(*schemaFileName)
	if err != nil {
		log.Fatalln(err)
	}

	terraformJson, err := terraform.Generate(schema)
	if err != nil {
		log.Fatalf("error while generating JSON of your schema with error message: %s \n", err)
	}

	fmt.Println(string(terraformJson))
	//log.Println(string(terraformJson))

}
