package cli

import (
	"github.com/urfave/cli/v2"
	"strconv"

	configuration "github.com/dmcg310/game-of-life/internal/config"
)

type Arguments struct {
	Preset string
	FPS    int
}

func ParseCLIArgs(ctx *cli.Context) *Arguments {
	args := &Arguments{}
	if ctx.NArg() > 0 {
		args.Preset = ctx.Args().Get(0)
	}

	if ctx.NArg() > 1 {
		fps, err := strconv.Atoi(ctx.Args().Get(1))
		if err == nil {
			args.FPS = fps
		}
	}

	return args
}

func PrepareConfigAndColors(args *Arguments) (*configuration.Config, *configuration.Colors) {
	configFile := configuration.ReadConfig()
	config := configuration.NewConfig(configFile)

	if config == nil {
		config = configuration.NewConfigWithDefaults()
	}

	if args.Preset != "" {
		config.Preset = args.Preset
	}

	if args.FPS != 0 {
		config.FPS = args.FPS
	}

	var colors *configuration.Colors
	if configFile != nil {
		colors = configuration.CustomColors(config)
	} else {
		colors = configuration.DefaultColors()
	}

	return config, colors
}
