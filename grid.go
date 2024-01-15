package main

import (
	"math/rand"
	"time"
)

type Grid struct {
	Cells          [][]bool
	NeedsRefreshed bool
}

func NewGrid(w int, h int, pattern string) *Grid {
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
	case "gosper-glider-gun":
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
		NewAppWarning("Unknown pattern found. Continued with random as a default.",
			"Check the pattern exists or any typos.").ShowAppWarning()

		cells = make([][]bool, w)
		for i := range cells {
			cells[i] = make([]bool, h)
		}

		rand.New(rand.NewSource(time.Now().UnixNano()))
		for x := 0; x < len(cells); x++ {
			for y := 0; y < len(cells[x]); y++ {
				cells[x][y] = rand.Float32() < 0.25 // % chance of a cell being alive
			}
		}
	}

	return &Grid{
		Cells:          cells,
		NeedsRefreshed: true,
	}
}
