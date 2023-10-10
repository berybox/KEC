package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:                   "KEC",
		Version:                "1.1",
		Usage:                  "Search for unique sequences by k-mer exclusion",
		UseShortOptionHandling: false,
		CustomAppHelpTemplate:  appHelpTemplate,
		HideHelp:               true,
		Commands: []*cli.Command{
			&commandExclude,
			&commandInclude,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
