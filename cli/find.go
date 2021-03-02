package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/erikkn/klaabu/klaabu"
)

// The purpose of the 'find' command is more to 'explore' the schema and searching for something in a broad way. 'GET', however, is used if you already know what you need to lookup and do a specific search.

func findCommand() {
	find := flag.NewFlagSet("find", flag.ExitOnError)

	schemaFileName := find.String("schema", "schema.kml", "(Optional) Schema file path (default './schema.kml')")
	termsString := find.String("label", "", "(Optional) Filter by labels. Example: 'az=euc1-az1,type=vpc'. No filtering when empty.")
	err := find.Parse(os.Args[2:])

	parentId := find.Arg(0)

	if err != nil || len(os.Args) < 2 {
		log.Printf("Usage: klaabu find [OPTIONS] [PARENT_ID] \n\n Subcommands: \n")
		find.PrintDefaults()
		os.Exit(1)
	}

	schema, err := klaabu.LoadSchemaFromKmlFile(*schemaFileName)
	if err != nil {
		log.Fatalln(err)
	}

	var parent *klaabu.Prefix
	if parentId == "" {
		parent = schema.Root
	} else {
		parent = schema.PrefixById(parentId)
		if parent == nil {
			log.Fatalln(err)
		}
	}

	terms := make([]*klaabu.LabelSearchTerm, 0)

	if *termsString != "" {
		termStrings := strings.Split(*termsString, ",")
		for _, termString := range termStrings {
			termSlice := strings.Split(termString, "=")
			var term klaabu.LabelSearchTerm
			if len(termSlice) == 2 {
				term.Key = strings.TrimSpace(termSlice[0])
				value := strings.TrimSpace(termSlice[1])
				term.Value = &value
			} else if len(termSlice) == 1 {
				term.Key = strings.TrimSpace(termSlice[0])
				term.Value = nil
			}
			terms = append(terms, &term)
		}
	}

	println(">>> num terms: ", len(terms))
	for _, term := range terms {
		println(">>> term: ", term.Key, term.Value)
	}

	cidrs := parent.FindPrefixesByLabelTerms(terms)

	for _, v := range cidrs {
		labels := make([]string, 0, len(v.Labels))
		for k, v := range v.Labels {
			labels = append(labels, k+"="+v)
		}

		log.Printf("%s: %s [%s]\n", string(v.Cidr), strings.Join(v.Aliases, "|"), strings.Join(labels, ","))
	}
}
