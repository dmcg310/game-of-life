package main

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

type config struct {
	Preset          string `json:"preset"`
	CellColor       string `json:"cell-color"`
	BackgroundColor string `json:"backgroundColor"`
}

type grid struct {
	cells          [][]bool
	needsRefreshed bool
}

type game struct {
	isRunning bool
	grid      grid
	screen    tcell.Screen
}

func main() {
	c, err := readConfig()
	if err != nil {
		reportError(err)
	}

	s, err := initScreen()
	if err != nil {
		reportError(err)
	}

	newGame(s, c).gameLoop()
}

func initScreen() (tcell.Screen, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err := screen.Init(); err != nil {
		return nil, err
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).
		Foreground(tcell.ColorReset)

	screen.SetStyle(defStyle)
	screen.EnableMouse()
	screen.Clear()

	return screen, nil
}

func newGame(screen tcell.Screen, c *config) *game {
	var cells [][]bool

	w, h := screen.Size()
	// w, h := 50, 50
	if c != nil {
		if c.Preset == "random" {
			cells = generateRandomCells(w, h)
		} else {
			cells = generatePatternCells(w, h, c.Preset)
		}
	} else {
		cells = generateRandomCells(w, h)
	}

	return &game{
		isRunning: true,
		grid: grid{
			cells:          cells,
			needsRefreshed: true,
		},
		screen: screen,
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
	cellChar := '█'

	for x := 0; x < len(g.grid.cells); x++ {
		for y := 0; y < len(g.grid.cells[x]); y++ {
			if g.grid.cells[x][y] {
				g.screen.SetContent(x, y, cellChar, nil, tcell.StyleDefault)
			} else {
				g.screen.SetContent(x, y, ' ', nil, tcell.StyleDefault)
			}
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

func generatePatternCells(w int, h int, pattern string) [][]bool {
	cells := make([][]bool, w)
	for i := range cells {
		cells[i] = make([]bool, h)
	}

	centerX, centerY := w/2, h/2

	switch pattern {
	// oscillators
	case "blinker":
		points := []struct{ x, y int }{{centerX - 1, centerY}, {centerX, centerY}, {centerX + 1, centerY}}
		for _, p := range points {
			cells[p.x][p.y] = true
		}
	case "toad":
		points := []struct{ x, y int }{{centerX - 1, centerY}, {centerX, centerY}, {centerX + 1, centerY},
			{centerX, centerY + 1}, {centerX + 1, centerY + 1}, {centerX + 2, centerY + 1}}
		for _, p := range points {
			cells[p.x][p.y] = true
		}
	case "beacon":
		points := []struct{ x, y int }{
			{centerX - 2, centerY - 1}, {centerX - 2, centerY - 2},
			{centerX - 1, centerY - 2}, {centerX + 1, centerY},
			{centerX + 1, centerY + 1}, {centerX, centerY + 1},
		}

		setCells(points, cells)
	case "lwss":
		points := []struct{ x, y int }{
			{centerX - 1, centerY + 1}, {centerX + 2, centerY + 1},
			{centerX - 2, centerY},
			{centerX - 2, centerY - 1}, {centerX + 2, centerY - 1},
			{centerX - 2, centerY - 2}, {centerX - 1, centerY - 2}, {centerX, centerY - 2}, {centerX + 1, centerY - 2},
		}
		for _, p := range points {
			cells[p.x][p.y] = true
		}
	case "gosper glider gun":
		points := []struct{ x, y int }{
			// todo
		}
		for _, p := range points {
			cells[p.x][p.y] = true
		}

	case "glider":
		points := []struct{ x, y int }{{1, 0}, {2, 1}, {0, 2}, {1, 2}, {2, 2}}
		for _, p := range points {
			cells[p.x][p.y] = true
		}
	case "block":
		points := []struct{ x, y int }{{0, 0}, {1, 0}, {0, 1}, {1, 1}}
		for _, p := range points {
			cells[p.x][p.y] = true
		}
	}

	return cells
}

func setCells(points []struct{ x, y int }, cells [][]bool) {
	for _, p := range points {
		if p.x >= 0 && p.x < len(cells) && p.y >= 0 && p.y < len(cells[0]) {
			cells[p.x][p.y] = true
		}
	}
}

func generateRandomCells(w int, h int) [][]bool {
	cells := make([][]bool, w)
	for i := range cells {
		cells[i] = make([]bool, h)
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))
	for x := 0; x < len(cells); x++ {
		for y := 0; y < len(cells[x]); y++ {
			cells[x][y] = rand.Float32() < 0.25
		}
	}

	return cells
}

func readConfig() (*config, error) {
	file, err := os.Open("gol-config.json")
	// if the file isnt found that is ok, just continue with defaults
	if err != nil {
		return nil, nil
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config config
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func reportError(msg error) {
	log.Fatalf("[ERROR] '%+v'", msg)
}
