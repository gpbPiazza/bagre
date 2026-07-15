package main

import (
	"fmt"
	_ "image/png"
	"log/slog"
	"os"

	"github.com/gpbPiazza/bagre/pkg/log"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	logger, lClose, err := log.InitLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = lClose() }()

	os.Exit(initGame(logger))
}

func initGame(l *slog.Logger) int {
	err := loadJellyImg()
	if err != nil {
		l.Error("failed to load jelly imgs", log.Err(err))
		return 1
	}

	err = loadWesImg()
	if err != nil {
		l.Error("failed to load wes imgs", log.Err(err))
		return 1
	}

	wes = NewWes(jellysCount+2, l)
	NewSmack(wes, l)

	s := NewGame()

	ebiten.SetWindowSize(s.ScreenWidth, s.ScreenWidth)
	ebiten.SetWindowTitle("Wes bagre")
	if err := ebiten.RunGame(s); err != nil {
		l.Error("failed to run game", log.Err(err))
		return 1
	}

	return 0
}
