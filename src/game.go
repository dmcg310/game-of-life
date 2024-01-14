package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Game struct {
	IsRunning bool
	IsPaused  bool
	Turn      int
	FPS       int
}

func NewGame(c *Config) *Game {
	return &Game{
		IsRunning: true,
		IsPaused:  true,
		Turn:      0,
		FPS:       c.FPS,
	}
}

func (g *Game) Run(screen tcell.Screen, grid *Grid, colors *Colors) {
	eventq := make(chan tcell.Event)
	quitq := make(chan struct{})
	ticker := time.NewTicker(time.Second / time.Duration(g.FPS))
	defer ticker.Stop()

	quit := func(screen tcell.Screen) {
		maybePanic := recover()
		screen.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit(screen)

	go func() {
		for {
			ev := screen.PollEvent()
			if ev == nil {
				return
			}

			eventq <- ev
		}
	}()

	go func() {
		for {
			select {
			case <-ticker.C:
				if !g.IsPaused {
					g.progress(grid)
				}
			case <-quitq:
				return
			}
		}
	}()

	for {
		if grid.needsRefreshed {
			g.RenderGamestate(grid, screen, colors)
			screen.Show()
			grid.needsRefreshed = false
		}

		select {
		case ev := <-eventq:
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Rune() == 'p' {
					g.IsPaused = !g.IsPaused
					g.RenderGamestate(grid, screen, colors)
					screen.Show()
				}

				if ev.Rune() == ' ' && g.IsPaused {
					g.progress(grid)
					g.RenderGamestate(grid, screen, colors)
					screen.Show()
				}

				if ev.Key() == tcell.KeyEscape ||
					ev.Key() == tcell.KeyCtrlC ||
					ev.Rune() == 'q' {
					close(quitq)
					return
				}
			}
		case <-quitq:
			break
		default:
			time.Sleep(time.Millisecond * 10) // 10ms
		}
	}
}

func (g *Game) RenderGamestate(grid *Grid, screen tcell.Screen, colors *Colors) {
	cellChar := '█'

	for x := 0; x < len(grid.cells); x++ {
		for y := 0; y < len(grid.cells[x]); y++ {
			if grid.cells[x][y] {
				screen.SetContent(x, y, cellChar, nil, colors.cellStyle)
			} else {
				screen.SetContent(x, y, ' ', nil, colors.backgroundStyle)
			}
		}
	}

	offset := 0
	offset = g.RenderInt(grid, screen, "FPS", g.FPS, offset)
	offset = g.RenderInt(grid, screen, "TURN", g.Turn, offset)

	if g.IsPaused {
		offset = g.RenderContent(grid, screen, "PAUSED", offset)
	} else {
		offset = g.RenderContent(grid, screen, "RUNNING", offset)
	}

	grid.needsRefreshed = false
	_ = offset
}

func (g *Game) progress(grid *Grid) {
	currentGrid := grid.cells
	tempGrid := make([][]bool, len(currentGrid))

	for x := range currentGrid {
		tempGrid[x] = make([]bool, len(currentGrid[x]))
		for y := range currentGrid[x] {
			count := g.countNeighbors(x, y, grid)

			if currentGrid[x][y] {
				tempGrid[x][y] = count == 2 || count == 3
			} else {
				tempGrid[x][y] = count == 3
			}
		}
	}

	grid.cells = tempGrid
	grid.needsRefreshed = true
	g.Turn++
}

func (g *Game) countNeighbors(x int, y int, grid *Grid) int {
	count := 0

	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}

			nx := x + dx
			ny := y + dy

			if g.withinBounds(nx, ny, grid) {
				if grid.cells[nx][ny] {
					count++
				}
			}
		}
	}

	return count
}

func (g *Game) withinBounds(x int, y int, grid *Grid) bool {
	return x >= 0 && x < len(grid.cells) && y >= 0 && y < len(grid.cells[x])
}

func (g *Game) RenderContent(
	grid *Grid, screen tcell.Screen, msg string, offset int,
) int {
	gridWidth := len(grid.cells)

	for i, rune := range msg {
		screen.SetContent(gridWidth-len(msg)+i, offset, rune, nil,
			tcell.StyleDefault.
				Foreground(tcell.ColorWhite).
				Background(tcell.ColorBlack))
	}

	offset++
	return offset

}

func (g *Game) RenderInt(
	grid *Grid, screen tcell.Screen, msg string, value int, offset int,
) int {
	str := fmt.Sprintf("%s %d", msg, value)
	return g.RenderContent(grid, screen, str, offset)
}
