package main

import "github.com/gdamore/tcell/v2"

type Colors struct {
	cellStyle       tcell.Style
	backgroundStyle tcell.Style
}

func DefaultColors() *Colors {
	return &Colors{
		cellStyle:       tcell.StyleDefault,
		backgroundStyle: tcell.StyleDefault,
	}
}

func CustomColors(c *Config) *Colors {
	var cellStyle, backgroundStyle tcell.Style

	if c.CellColor != "" {
		cellColor := tcell.GetColor(c.CellColor)
		cellStyle = tcell.StyleDefault.Foreground(cellColor)
	} else {
		cellStyle = tcell.StyleDefault
	}

	if c.BackgroundColor != "" {
		backgroundColor := tcell.GetColor(c.BackgroundColor)
		backgroundStyle = tcell.StyleDefault.Background(backgroundColor)
	} else {
		backgroundStyle = tcell.StyleDefault
	}

	return &Colors{
		cellStyle:       cellStyle,
		backgroundStyle: backgroundStyle,
	}
}
