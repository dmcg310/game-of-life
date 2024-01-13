package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Config struct {
	Preset          string `json:"preset"`
	CellColor       string `json:"cell-color"`
	BackgroundColor string `json:"backgroundColor"`
	ScaleFactor     int    `json:"scale-factor"`
	FPS             int    `json:"fps"`
}

type Grid struct {
	cells          [][]bool
	needsRefreshed bool
}

type Game struct {
	isRunning bool
	isPaused  bool
	grid      Grid
	screen    tcell.Screen
	turn      int
	FPS       int
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

	newGame(s, c).run()
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
	screen.Clear()

	return screen, nil
}

func newGame(screen tcell.Screen, c *Config) *Game {
	var cells [][]bool

	w, h := screen.Size()
	if c != nil {
		if c.Preset == "random" {
			cells = generateRandomCells(w, h)
		} else {
			cells = generatePatternCells(w, h, c.Preset)
		}
	} else {
		cells = generateRandomCells(w, h)
	}

	return &Game{
		isRunning: true,
		isPaused:  true,
		grid: Grid{
			cells:          cells,
			needsRefreshed: true,
		},
		screen: screen,
		turn:   0,
		FPS:    c.FPS,
	}
}

func (g *Game) run() {
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
	defer quit(g.screen)

	go func() {
		for {
			ev := g.screen.PollEvent()
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
				if !g.isPaused {
					g.progress()
				}
			case <-quitq:
				return
			}
		}
	}()

	for {
		if g.grid.needsRefreshed {
			g.renderGamestate()
			g.screen.Show()
			g.grid.needsRefreshed = false
		}

		select {
		case ev := <-eventq:
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Rune() == 'p' {
					g.isPaused = !g.isPaused
				}

				if ev.Rune() == ' ' && g.isPaused {
					g.progress()
					g.renderGamestate()
					g.screen.Show()
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

func (g *Game) renderGamestate() {
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

	offset := 0
	offset = g.renderInt("FPS", g.FPS, offset)
	offset = g.renderInt("TURN", g.turn, offset)

	if g.isPaused {
		offset = g.renderContent("PAUSED", offset)
	} else {
		offset = g.renderContent("RUNNING", offset)
	}

	g.grid.needsRefreshed = false
	_ = offset
}

func (g *Game) renderInt(msg string, value int, offset int) int {
	str := fmt.Sprintf("%s %d", msg, value)
	gridWidth := len(g.grid.cells)

	for i, rune := range str {
		g.screen.SetContent(gridWidth-len(str)+i, offset, rune, nil,
			tcell.StyleDefault.
				Foreground(tcell.ColorWhite).
				Background(tcell.ColorBlack))
	}

	offset++
	return offset
}

func (g *Game) renderContent(msg string, offset int) int {
	gridWidth := len(g.grid.cells)

	for i, rune := range msg {
		g.screen.SetContent(gridWidth-len(msg)+i, offset, rune, nil,
			tcell.StyleDefault.
				Foreground(tcell.ColorWhite).
				Background(tcell.ColorBlack))
	}

	offset++
	return offset
}

func (g *Game) progress() {
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
	g.turn++
}

func (g *Game) countNeighbors(x int, y int) int {
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

func (g *Game) withinBounds(x int, y int) bool {
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
		points := []struct{ x, y int }{
			{centerX - 1, centerY}, {centerX, centerY},
			{centerX + 1, centerY},
		}

		for _, p := range points {
			cells[p.x][p.y] = true
		}
	case "toad":
		points := []struct{ x, y int }{
			{centerX - 1, centerY}, {centerX, centerY},
			{centerX + 1, centerY}, {centerX, centerY + 1},
			{centerX + 1, centerY + 1}, {centerX + 2, centerY + 1},
		}

		for _, p := range points {
			cells[p.x][p.y] = true
		}
	case "beacon":
		points := []struct{ x, y int }{
			{centerX - 2, centerY - 1}, {centerX - 2, centerY - 2},
			{centerX - 1, centerY - 2}, {centerX + 1, centerY},
			{centerX + 1, centerY + 1}, {centerX, centerY + 1},
		}

		for _, p := range points {
			cells[p.x][p.y] = true
		}
	case "lwss":
		points := []struct{ x, y int }{
			{centerX - 1, centerY + 1}, {centerX + 2, centerY + 1},
			{centerX - 2, centerY}, {centerX - 2, centerY - 1},
			{centerX + 2, centerY - 1}, {centerX - 2, centerY - 2},
			{centerX - 1, centerY - 2}, {centerX, centerY - 2},
			{centerX + 1, centerY - 2},
		}

		for _, p := range points {
			cells[p.x][p.y] = true
		}
	case "gosper glider gun":
		offsetX, offsetY := centerX-18, centerY-5
		points := []struct{ x, y int }{
			{offsetX + 0, offsetY + 4}, {offsetX + 0, offsetY + 5},
			{offsetX + 1, offsetY + 4}, {offsetX + 1, offsetY + 5},
			{offsetX + 10, offsetY + 4}, {offsetX + 10, offsetY + 5},
			{offsetX + 10, offsetY + 6}, {offsetX + 11, offsetY + 3},
			{offsetX + 11, offsetY + 7}, {offsetX + 12, offsetY + 2},
			{offsetX + 12, offsetY + 8}, {offsetX + 13, offsetY + 2},
			{offsetX + 13, offsetY + 8}, {offsetX + 14, offsetY + 5},
			{offsetX + 15, offsetY + 3}, {offsetX + 15, offsetY + 7},
			{offsetX + 16, offsetY + 4}, {offsetX + 16, offsetY + 5},
			{offsetX + 16, offsetY + 6}, {offsetX + 17, offsetY + 5},
			{offsetX + 20, offsetY + 2}, {offsetX + 20, offsetY + 3},
			{offsetX + 20, offsetY + 4}, {offsetX + 21, offsetY + 2},
			{offsetX + 21, offsetY + 3}, {offsetX + 21, offsetY + 4},
			{offsetX + 22, offsetY + 1}, {offsetX + 22, offsetY + 5},
			{offsetX + 24, offsetY + 0}, {offsetX + 24, offsetY + 1},
			{offsetX + 24, offsetY + 5}, {offsetX + 24, offsetY + 6},
			{offsetX + 34, offsetY + 2}, {offsetX + 34, offsetY + 3},
			{offsetX + 35, offsetY + 2}, {offsetX + 35, offsetY + 3},
		}

		for _, p := range points {
			cells[p.x][p.y] = true
		}
	case "glider":
		points := []struct{ x, y int }{
			{1, 0}, {2, 1}, {0, 2}, {1, 2}, {2, 2},
		}

		for _, p := range points {
			cells[p.x][p.y] = true
		}
	case "block":
		points := []struct{ x, y int }{
			{0, 0}, {1, 0}, {0, 1}, {1, 1},
		}

		for _, p := range points {
			cells[p.x][p.y] = true
		}
	}

	return cells
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

func readConfig() (*Config, error) {
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

	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func reportError(msg error) {
	log.Fatalf("[ERROR] '%+v'", msg)
}
