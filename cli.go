package main

import (
	"strconv"

	"github.com/urfave/cli/v2"
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

func PrepareConfigAndColors(args *Arguments) (*Config, *Colors) {
	configFile := ReadConfig()
	config := NewConfig(configFile)

	if config == nil {
		config = NewConfigWithDefaults()
	}

	if args.Preset != "" {
		config.Preset = args.Preset
	}

	if args.FPS != 0 {
		config.FPS = args.FPS
	}

	var colors *Colors
	if configFile != nil {
		colors = CustomColors(config)
	} else {
		colors = DefaultColors()
	}

	return config, colors
}
