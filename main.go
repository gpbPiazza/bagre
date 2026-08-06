package main

import (
	"fmt"
	_ "image/png"
	"log/slog"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/gpbPiazza/bagre/pkg/log"
)

func main() {
	os.Exit(run())
}

// run exists so its defers execute before main reaches [os.Exit],
// which would skip them.
func run() int {
	logger, lClose, err := log.InitLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		return 1
	}
	defer func() { _ = lClose() }()

	return initGame(logger)
}

func initGame(l *slog.Logger) int {
	err := loadJellyImg()
	if err != nil {
		l.Error("failed to load jelly imgs", log.Err(err))
		return 1
	}

	game := NewGame(l)

	ebiten.SetWindowTitle("Wes bagre")
	ebiten.SetFullscreen(true)
	if err = ebiten.RunGame(game); err != nil {
		l.Error("failed to run game", log.Err(err))
		return 1
	}

	return 0
}
