package main

import "github.com/gdamore/tcell/v2"

func InitScreen() tcell.Screen {
	screen, err := tcell.NewScreen()
	if err != nil {
		newAppError(err, "Cannot create a new terminal screen.",
			"Please try to re-run the program.").showAppErrorFatal()
	}

	if err := screen.Init(); err != nil {
		newAppError(err, "Cannot initalise the terminal screen.",
			"Please try to re-run the program, and maybe reset the terminal using `$ reset`").
			showAppErrorFatal()
	}

	return screen
}
