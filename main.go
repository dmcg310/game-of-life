package main

import (
	"github.com/gdamore/tcell/v2"
	"log"
)

type game struct {
	isRunning bool
	terminal  terminal
}

type terminal struct {
	width  int
	height int
}

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		reportError(err)
	}

	if err := screen.Init(); err != nil {
		reportError(err)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	screen.SetStyle(defStyle)
	screen.EnableMouse()
	screen.EnablePaste()
	screen.Clear()

	quit := func() {
		maybePanic := recover()
		screen.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	w, h := screen.Size()

	g := &game{
		isRunning: true,
		terminal: terminal{
			width:  w,
			height: h,
		},
	}

	for g.isRunning {
		screen.Show()
		ev := screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
				return
			}
		}
	}
}

func reportError(msg error) {
	log.Fatalf("[ERROR] '%+v'", msg)
}
