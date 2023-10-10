package main

import (
	"os"

	"github.com/berybox/KEC/pkg/filedriver"
	"github.com/berybox/KEC/pkg/keccore"
	"github.com/berybox/KEC/pkg/keccore/kecoriginal"
	"github.com/berybox/KEC/pkg/maps/kmapbasic"
	"github.com/berybox/KEC/pkg/maps/mmapbasic"
	"github.com/berybox/KEC/pkg/utils"
	"github.com/urfave/cli/v2"
)

var (
	defaultK = 12

	flagTarget    = &cli.StringSliceFlag{Name: "target", Aliases: []string{"t"}, Usage: "Add target sequence(s) from `PATH`. Can be single .fasta file or whole directory", Required: true}
	flagNontarget = &cli.StringSliceFlag{Name: "nontarget", Aliases: []string{"n"}, Usage: "Add nontarget sequence(s) from `PATH`. Can be single .fasta file or whole directory", Required: true}
	flagOutput    = &cli.StringFlag{Name: "output", Aliases: []string{"o"}, Usage: "Output `PATH` to store resulting .fasta file", Required: true}

	flagK       = &cli.IntFlag{Name: "k", Usage: "K-mer size", Value: defaultK}
	flagMin     = &cli.IntFlag{Name: "min", Usage: "Minimum size of resulting sequence", Value: defaultK + 1}
	flagMax     = &cli.IntFlag{Name: "max", Usage: "Maximum size of resulting sequence (0 = unlimited)", Value: 0}
	flagReverse = &cli.BoolFlag{Name: "r", Usage: "Also exclude reverse complements of the sequences", Value: false}

	commandExclude = cli.Command{
		Name:                   "exclude",
		Aliases:                []string{"e"},
		Usage:                  "Exclude nontarget k-mers from target sequence",
		CustomHelpTemplate:     commandHelpTemplate + excludeExample,
		UseShortOptionHandling: false,
		Flags: []cli.Flag{
			flagTarget,
			flagNontarget,
			flagOutput,
			flagK,
			flagMin,
			flagMax,
			flagReverse,
		},
		Action: func(ctx *cli.Context) error {
			outfile, err := os.Create(ctx.String(flagOutput.Name))
			if err != nil {
				return err
			}
			defer outfile.Close()

			opts := keccore.Options{
				K:                 ctx.Int(flagK.Name),
				MinSize:           ctx.Int(flagMin.Name),
				MaxSize:           ctx.Int(flagMax.Name),
				ReverseComplement: ctx.Bool(flagReverse.Name),
				Output:            outfile,
				Log:               os.Stdout,
			}

			var kec keccore.Excluder
			kec = &kecoriginal.Excluder{}

			mapset := keccore.MapSet{
				TargetMap:    kmapbasic.NewTarget(),
				NontargetMap: kmapbasic.NewNontarget(),
				MaskMap:      mmapbasic.NewMaskMap(),
			}

			err = kec.Init(opts, mapset)
			if err != nil {
				return err
			}

			for _, targetFn := range ctx.StringSlice(flagTarget.Name) {
				fileList, err := utils.FileListExt(targetFn, filedriver.FastaExtensions)
				if err != nil {
					return err
				}
				for _, fn := range fileList {
					kec.AddTargetFilename(fn)
				}
			}

			for _, nontargetFn := range ctx.StringSlice(flagNontarget.Name) {
				fileList, err := utils.FileListExt(nontargetFn, filedriver.FastaExtensions)
				if err != nil {
					return err
				}
				for _, fn := range fileList {
					kec.AddNontargetFilename(fn)
				}
			}

			err = kec.Run()
			if err != nil {
				return err
			}

			return nil
		},
	}
)
