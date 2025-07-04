package display

import (
	"github.com/dmcg310/game-of-life/internal/errors"
	"github.com/gdamore/tcell/v2"
)

func InitScreen() tcell.Screen {
	screen, err := tcell.NewScreen()
	if err != nil {
		errors.NewAppError(err, "Cannot create a new terminal screen.",
			"Please try to re-run the program.").ShowAppErrorFatal()
	}

	if err := screen.Init(); err != nil {
		errors.NewAppError(err, "Cannot initialise the terminal screen.",
			"Please try to re-run the program, and maybe reset the terminal using `$ reset`").
			ShowAppErrorFatal()
	}

	return screen
}
