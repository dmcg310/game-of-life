package main

import (
	"github.com/gdamore/tcell/v2"
	"log"
	"time"
)

type grid struct {
	cells          [][]bool
	needsRefreshed bool
}

type terminal struct {
	width  int
	height int
}

type game struct {
	isRunning bool
	grid      grid
	screen    tcell.Screen
	terminal  terminal
}

func main() {
	initScreen().gameLoop()
}

func initScreen() *game {
	screen, err := tcell.NewScreen()
	if err != nil {
		reportError(err)
	}

	if err := screen.Init(); err != nil {
		reportError(err)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).
		Foreground(tcell.ColorReset)

	screen.SetStyle(defStyle)
	screen.EnableMouse()
	screen.Clear()

	w, h := screen.Size()

	cells := make([][]bool, w)
	for i := range cells {
		cells[i] = make([]bool, h)
	}

	return &game{
		isRunning: true,
		grid: grid{
			cells:          cells,
			needsRefreshed: true,
		},
		screen: screen,
		terminal: terminal{
			width:  w,
			height: h,
		},
	}
}

func (g *game) gameLoop() {
	quit := func(screen tcell.Screen) {
		maybePanic := recover()
		screen.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit(g.screen)

	delay := time.Millisecond * 100 // 100ms

	for g.isRunning {
		g.screen.Show()
		ev := g.screen.PollEvent()

		if g.grid.needsRefreshed {
			g.renderGamestate()
			g.screen.Show()
		}

		switch ev := ev.(type) {
		case *tcell.EventResize:
			g.screen.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape ||
				ev.Key() == tcell.KeyCtrlC ||
				ev.Rune() == 'q' {
				g.isRunning = false
			}
		}

		time.Sleep(delay)
	}
}

func (g *game) renderGamestate() {
	for x := 0; x < len(g.grid.cells); x++ {
		for y := 0; y < len(g.grid.cells[x]); y++ {
			char := ' '
			if g.grid.cells[x][y] {
				char = 'x'
			}

			g.screen.SetContent(x, y, char, nil, tcell.StyleDefault)
		}
	}

	g.grid.needsRefreshed = false
}

func reportError(msg error) {
	log.Fatalf("[ERROR] '%+v'", msg)
}
