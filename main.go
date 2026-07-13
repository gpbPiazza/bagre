package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	err := loadJellyImg()
	if err != nil {
		log.Fatal(err)
	}

	err = loadWesImg()
	if err != nil {
		log.Fatal(err)
	}

	wes = NewWes(jellysCount + 2)
	setEmptySmacksMapsPosition()
	NewSmack()
	setEachJellyFishPositionIntoSmackMap()

	s := NewGame()

	ebiten.SetWindowSize(s.ScreenWidth, s.ScreenWidth)
	ebiten.SetWindowTitle("Wes bagre")
	if err := ebiten.RunGame(s); err != nil {
		log.Fatal(err)
	}
}
