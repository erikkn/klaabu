package main

import (
	"flag"
	"github.com/erikkn/klaabu/klaabu"
	"github.com/erikkn/klaabu/klaabu/iputil"
	"log"
	"os"
	"sort"
)

func spaceCommand() {
	flagSet := flag.NewFlagSet("space", flag.ExitOnError)

	schemaFileName := flagSet.String("schema", "schema.kml", "(Optional) Schema file path (default './schema.kml')")

	err := flagSet.Parse(os.Args[2:])

	parentId := flagSet.Arg(0)

	if err != nil || len(os.Args) < 2 {
		log.Printf("Usage: klaabu space [OPTIONS] [PARENT_ID] \n\n Subcommands: \n")
		flagSet.PrintDefaults()
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
			log.Fatalf("Parent not found: %s", parentId)
		}
	}

	min, max, err := iputil.MinMaxIP(string(parent.Cidr))
	if err != nil {
		log.Fatal(err)
	}

	childCidrs := make([]klaabu.Cidr, 0, len(parent.Children))
	for _, child := range parent.Children {
		childCidrs = append(childCidrs, child.Cidr)
	}

	sort.Slice(childCidrs, func(i1, i2 int) bool {
		min1, _, _ := iputil.MinMaxIP(string(childCidrs[i1]))
		min2, _, _ := iputil.MinMaxIP(string(childCidrs[i2]))
		cmp, _ := iputil.CompareIPs(min1, min2)
		return cmp < 0
	})

	for _, cidr := range childCidrs {
		childMin, childMax, err := iputil.MinMaxIP(string(cidr))
		if err != nil {
			log.Fatal(err)
		}

		childMinCmp, err := iputil.CompareIPs(min, childMin)
		if err != nil {
			log.Fatal(err)
		}

		childStartsAtMin := childMinCmp == 0

		if !childStartsAtMin {
			// detected a gap between min and the child start, grab it
			childStartMinusOne, err := iputil.PreviousIP(childMin)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("From: %s, To: %s", min, childStartMinusOne)
		}

		// continue right after the child ends
		min, err = iputil.NextIP(childMax)
		if err != nil {
			log.Fatal(err)
		}
	}

	minMaxCmp, err := iputil.CompareIPs(min, max)
	if err != nil {
		log.Fatal(err)
	}

	if minMaxCmp < 0 {
		// have some IP space remaining
		log.Printf("From: %s, To: %s", min, max)
	}
}
