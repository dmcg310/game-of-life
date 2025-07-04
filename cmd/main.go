package main

import (
	"fmt"
	"github.com/dmcg310/game-of-life/internal/display"
	"github.com/dmcg310/game-of-life/internal/errors"
	"github.com/dmcg310/game-of-life/internal/game"
	"os"

	_cli "github.com/dmcg310/game-of-life/internal/cli"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:      "Game of Life",
		Usage:     "Conway's Game of Life in the Terminal!",
		UsageText: "game-of-life [global options] <pattern> <fps>",
		Description: `This program runs Conway's Game of Life. You can optionally specify
        a pattern and frames per second (fps) as arguments, or create a configuration file 
        which additionally includes color options.`,
		ArgsUsage: "[pattern] [fps]",
		Action: func(ctx *cli.Context) error {
			args := _cli.ParseCLIArgs(ctx)
			config, colors := _cli.PrepareConfigAndColors(args)

			s := display.InitScreen()
			w, h := s.Size()
			g := game.NewGrid(w, h, config.Preset)

			game.NewGame(config).Run(s, g, colors)

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "patterns",
				Aliases: []string{"p"},
				Usage:   "List available patterns",
				Action: func(c *cli.Context) error {
					fmt.Println("Available patterns: blinker, toad, beacon, lwss, gosper-glider-gun, glider, block, random")
					return nil
				},
			},
			{
				Name:    "config-location",
				Aliases: []string{"cl"},
				Usage:   "Echo configuration directory path",
				Action: func(c *cli.Context) error {
					s, _ := os.UserConfigDir()
					fmt.Printf("%s/gol/\n", s)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		errors.NewAppError(err, "Cannot run the current application.",
			"Please try to re-run the program.").ShowAppErrorFatal()
	}
}
