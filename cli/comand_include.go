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
	flagReference = &cli.StringSliceFlag{Name: "reference", Aliases: []string{"ref", "m"}, Usage: "Add reference sequence(s) from `PATH`. Can be single .fasta file or whole directory", Required: true}
	flagPool      = &cli.StringSliceFlag{Name: "pool", Aliases: []string{"p"}, Usage: "Add pool sequence(s) from `PATH`. Can be single .fasta file or whole directory", Required: true}

	commandInclude = cli.Command{
		Name:                   "include",
		Aliases:                []string{"i"},
		Usage:                  "Include nontarget k-mers to target sequence",
		CustomHelpTemplate:     commandHelpTemplate + includeExample,
		UseShortOptionHandling: false,
		Flags: []cli.Flag{
			flagReference,
			flagPool,
			flagOutput,
			flagK,
			flagMin,
			flagMax,
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
				ReverseComplement: false,
				Output:            outfile,
				Log:               os.Stdout,
			}

			var kec keccore.Includer
			kec = &kecoriginal.Includer{}

			mapset := keccore.MapSet{
				TargetMap:    kmapbasic.NewTarget(),
				NontargetMap: kmapbasic.NewNontarget(),
				MaskMap:      mmapbasic.NewMaskMap(),
			}

			err = kec.Init(opts, mapset)
			if err != nil {
				return err
			}

			for _, refFn := range ctx.StringSlice(flagReference.Name) {
				fileList, err := utils.FileListExt(refFn, filedriver.FastaExtensions)
				if err != nil {
					return err
				}
				for _, fn := range fileList {
					kec.AddReferenceFilename(fn)
				}
			}

			for _, poolFn := range ctx.StringSlice(flagPool.Name) {
				fileList, err := utils.FileListExt(poolFn, filedriver.FastaExtensions)
				if err != nil {
					return err
				}
				for _, fn := range fileList {
					kec.AddPoolFilename(fn)
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
