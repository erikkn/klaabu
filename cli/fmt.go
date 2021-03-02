package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/erikkn/klaabu/klaabu"
)

func fmtCommand() {
	flags := flag.NewFlagSet("fmt", flag.ExitOnError)

	schemaFileName := flags.String("schema", "schema.kml", "Schema file path")
	flags.Parse(os.Args[2:])

	node, err := klaabu.LoadKmlFromFile(*schemaFileName)
	if err != nil {
		log.Fatalln(err)
	}

	schema, err := klaabu.KmlToSchema(node)
	if err != nil {
		log.Fatalln(err)
	}

	err = schema.Validate()
	if err != nil {
		log.Fatalln(err)
	}

	originalBytes, err := ioutil.ReadFile(*schemaFileName)
	if err != nil {
		log.Fatalln(err)
	}

	wf, err := os.OpenFile(*schemaFileName, os.O_WRONLY, 0)
	if err != nil {
		log.Fatalln(err)
	}

	err = klaabu.MarshalKml(node, wf)
	if err != nil {
		log.Fatalln(err)
	}

	updatedBytes, err := ioutil.ReadFile(*schemaFileName)
	if err != nil {
		log.Fatalln(err)
	}

	if string(originalBytes) != string(updatedBytes) {
		log.Println(*schemaFileName)
	}
}
