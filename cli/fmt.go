package main

import (
	"bufio"
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/transferwise/klaabu/klaabu"
)

func fmtCommand() {
	flags := flag.NewFlagSet("fmt", flag.ExitOnError)

	schemaFileName := flags.String("schema", "schema.kml", "Schema file path")
	fmtCheck := flags.Bool("check", false, "Checks the formatting of your schema only, but doesn't actually format it!")
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

	if *fmtCheck {
		fmtBuffer := bytes.Buffer{}
		err := klaabu.MarshalKml(node, &fmtBuffer)
		if err != nil {
			log.Fatalln(err)
		}

		originalScanner := bufio.NewScanner(bytes.NewBuffer(originalBytes))
		fmtScanner := bufio.NewScanner(&fmtBuffer)
		lineNumber := 1
		mismatchLines := false

		for originalScanner.Scan() && fmtScanner.Scan() {
			if originalScanner.Text() != fmtScanner.Text() {
				log.Printf("- %5d: %s\n", lineNumber, originalScanner.Text())
				log.Printf("+ %5d: %s\n\n", lineNumber, fmtScanner.Text())
				mismatchLines = true
			}

			lineNumber++
		}

		if mismatchLines {
			log.Fatalln("Schema is not properly formatted. Run `klaabu fmt`.")
		} else {
			log.Println("Schema is properly formatted, thank you!")
		}

	} else {
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
}
