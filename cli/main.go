package main

import (
	"os"

	"github.com/urfave/cli"
)

var (
	target    string
	flgTarget = cli.StringFlag{
		Name:        "t",
		Usage:       "Add target sequence(s) from `PATH` as .fasta file or whole directory",
		Destination: &target,
		Required:    true,
	}

	nonTarget    string
	flgNonTarget = cli.StringFlag{
		Name:        "n",
		Usage:       "Add nontarget sequence(s) from `PATH` as .fasta file or whole directory",
		Destination: &nonTarget,
		Required:    true,
	}

	master    string
	flgMaster = cli.StringFlag{
		Name:        "m",
		Usage:       "Add master sequence(s) from `PATH` as .fasta file or whole directory",
		Destination: &master,
		Required:    true,
	}

	pool    string
	flgPool = cli.StringFlag{
		Name:        "p",
		Usage:       "Add pool sequence(s) from `PATH` as .fasta to include in consensus sequence",
		Destination: &pool,
		Required:    true,
	}

	output    string
	flgOutput = cli.StringFlag{
		Name:        "o",
		Usage:       "Output `PATH` to store resulting .fasta file",
		Destination: &output,
		Required:    true,
	}

	k    = 12
	flgK = cli.IntFlag{
		Name:        "k",
		Usage:       "K-mer size",
		Destination: &k,
		Required:    false,
		Value:       k,
	}

	min    int
	flgMin = cli.IntFlag{
		Name:        "min",
		Usage:       "Minimum size of resulting sequence",
		Destination: &min,
		Required:    false,
		Value:       k + 1,
	}

	max    int
	flgMax = cli.IntFlag{
		Name:        "max",
		Usage:       "Maximum size of resulting sequence (0 = unlimited)",
		Destination: &max,
		Required:    false,
		Value:       0,
	}

	reverse    bool
	flgReverse = cli.BoolFlag{
		Name:        "r",
		Usage:       "Also exclude reverse complements of the sequences",
		Destination: &reverse,
		Required:    false,
		Value:       false,
	}

	cmdExclude = cli.Command{
		CustomHelpTemplate: commandHelpTemplate + excludeExample,
		Name:               "exclude",
		Aliases:            []string{"e"},
		Usage:              "Exclude nontarget k-mers from target sequence",

		Flags: []cli.Flag{&flgTarget, &flgNonTarget, &flgOutput, &flgK, &flgReverse, &flgMin, &flgMax},

		Action: exclude,
	}

	cmdInclude = cli.Command{
		CustomHelpTemplate: commandHelpTemplate + includeExample,
		Name:               "include",
		Aliases:            []string{"i"},
		Usage:              "Create common sequences from master sequence and other sequences to include",

		Flags: []cli.Flag{&flgMaster, &flgPool, &flgOutput, &flgK, &flgMin, &flgMax},

		Action: include,
	}
	cmds = []*cli.Command{&cmdExclude, &cmdInclude}
)

func main() {
	app := cli.NewApp()
	app.CustomAppHelpTemplate = appHelpTemplate
	app.Version = "1.0"
	app.Name = "KEC"
	app.Usage = "Search for unique sequences"
	app.HideHelp = true
	app.Commands = cmds
	app.Run(os.Args)
}
