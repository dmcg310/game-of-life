package main

import (
	"github.com/gdamore/tcell/v2"
	"log"
	"math/rand"
	"time"
)

type config struct {
    termW int
    termH int
    randomnessValue float32
}

type grid struct {
	cells          [][]bool
	needsRefreshed bool
}

type game struct {
	isRunning bool
	grid      grid
	screen    tcell.Screen
    config config
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
    config := config{
        termW: w,
        termH: h,
        randomnessValue: 0.17,
    }

	cells := config.generateCells()

	return &game{
		isRunning: true,
		grid: grid{
			cells:          cells,
			needsRefreshed: true,
		},
		screen: screen,
        config: config,
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

	for g.isRunning {
		g.screen.Show()
		ev := g.screen.PollEvent()

		g.progress()

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

func (g *game) progress() {
	currentGrid := g.grid.cells
	tempGrid := make([][]bool, len(currentGrid))

	for x := range currentGrid {
		tempGrid[x] = make([]bool, len(currentGrid[x]))
		for y := range currentGrid[x] {
			count := g.countNeighbors(x, y)

			if currentGrid[x][y] {
				tempGrid[x][y] = count == 2 || count == 3
			} else {
				tempGrid[x][y] = count == 3
			}
		}
	}

	g.grid.cells = tempGrid
	g.grid.needsRefreshed = true
}

func (g *game) countNeighbors(x int, y int) int {
	count := 0

	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}

			nx := x + dx
			ny := y + dy

			if g.withinBounds(nx, ny) {
				if g.grid.cells[nx][ny] {
					count++
				}
			}
		}
	}

	return count
}

func (g *game) withinBounds(x int, y int) bool {
	return x >= 0 && x < len(g.grid.cells) && y >= 0 && y < len(g.grid.cells[x])
}

func (c *config) generateCells() [][]bool {
	cells := make([][]bool, c.termW)
	for i := range cells {
		cells[i] = make([]bool, c.termH)
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))
	for x := 0; x < len(cells); x++ {
		for y := 0; y < len(cells[x]); y++ {
			cells[x][y] = rand.Float32() < c.randomnessValue
		}
	}

	return cells
}

func reportError(msg error) {
	log.Fatalf("[ERROR] '%+v'", msg)
}
