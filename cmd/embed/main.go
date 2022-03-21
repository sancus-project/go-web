package main

import (
	"errors"
	"flag"
	"log"
	"os"
)

type Config struct {
	Package string
	Output  string
	Varname string
}

func main() {

	// Config
	c := &Config{}
	flag.StringVar(&c.Package, "p", "${GOPACKAGE}", "package name")
	flag.StringVar(&c.Output, "o", "-", "output file")
	flag.StringVar(&c.Varname, "n", "Files", "variable name")
	flag.Parse()

	// Embedder
	embed, err := NewEmbedder(c)
	if err == nil {
		// Add items
		items := flag.Args()
		if len(items) == 0 {
			items = []string{"."}
		}

		ok := true
		for _, item := range items {
			if err := embed.Add(item); err != nil {
				log.Println(err)
				ok = false
			}
		}

		if !ok {
			// failed to add items
			err = errors.New("Failed to process all inputs")
		} else if c.Output == "" || c.Output == "-" {
			// write to stdout
			_, err = embed.WriteTo(os.Stdout)
		} else {
			// write to file
			err = embed.WriteFile(c.Output)
		}
	}

	if err != nil {
		log.Fatal(err)
	}
}
