package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	defaultFPS      = 23
	defaultPattern  = "random"
	defaultFilename = "gol-config.json"
)

type Config struct {
	Preset          string `json:"preset"`
	CellColor       string `json:"cell-color"`
	BackgroundColor string `json:"background-color"`
	ScaleFactor     int    `json:"scale-factor"`
	FPS             int    `json:"fps"`
}

type Colors struct {
	cellStyle       tcell.Style
	backgroundStyle tcell.Style
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
	colors    Colors
}

func main() {
	c := readConfig()
	s := initScreen()

	newGame(s, c).run()
}

func initScreen() tcell.Screen {
	screen, err := tcell.NewScreen()
	if err != nil {
		newAppError(err, "Cannot create a new terminal screen.",
			"Please try to re-run the program.").showAppErrorFatal()
	}

	if err := screen.Init(); err != nil {
		newAppError(err, "Cannot initalise the terminal screen.",
			"Please try to re-run the program, and maybe reset the terminal using `$ reset`").showAppErrorFatal()
	}

	return screen
}

func newGame(screen tcell.Screen, c *Config) *Game {
	w, h := screen.Size()

	var (
		colors Colors
		fps    int
		cells  [][]bool
	)

	if c != nil {
		colors = customColors(c)
		fps = c.FPS
		cells = generatePatternCells(w, h, c.Preset)
	} else {
		colors = defaultColors()
		fps = defaultFPS
		cells = generatePatternCells(w, h, "random")
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
		FPS:    fps,
		colors: colors,
	}
}

func defaultColors() Colors {
	return Colors{
		cellStyle:       tcell.StyleDefault,
		backgroundStyle: tcell.StyleDefault,
	}
}

func customColors(c *Config) Colors {
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

	return Colors{
		cellStyle:       cellStyle,
		backgroundStyle: backgroundStyle,
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
					g.renderGamestate()
					g.screen.Show()
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
				g.screen.SetContent(x, y, cellChar, nil, g.colors.cellStyle)
			} else {
				g.screen.SetContent(x, y, ' ', nil, g.colors.backgroundStyle)
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

func (g *Game) renderInt(msg string, value int, offset int) int {
	str := fmt.Sprintf("%s %d", msg, value)
	return g.renderContent(str, offset)
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
	case "random":
		fallthrough
	default:
		cells = make([][]bool, w)
		for i := range cells {
			cells[i] = make([]bool, h)
		}

		rand.New(rand.NewSource(time.Now().UnixNano()))
		for x := 0; x < len(cells); x++ {
			for y := 0; y < len(cells[x]); y++ {
				cells[x][y] = rand.Float32() < 0.25
			}
		}
	}

	return cells
}

func readConfig() *Config {
	file, err := os.Open(defaultFilename)
	if err != nil {
		// if the file isnt found that is ok, just continue with defaults
		msg := fmt.Sprintf("Cannot open config file: '%s'. Continued with defaults.",
			defaultFilename)
		newAppWarning(msg, "Make sure the file exists and is accessible by the program.").
			showAppWarning()

		return nil
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		msg := fmt.Sprintf("Cannot read config file: '%s'. Continued with defaults.",
			defaultFilename)
		newAppWarning(msg, "Please try re-running the program.").
			showAppWarning()

		return nil
	}

	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		msg := fmt.Sprintf("Cannot parse JSON in config file: '%s'. Continued with defaults.",
			defaultFilename)
		newAppWarning(msg, "Please ensure that the JSON contains no syntactical errors.").
			showAppWarning()

		return nil
	}

	return &config
}
