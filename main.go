package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	setEmptyFlocksMapsPosition()
	NewFlock()
	setEachBoidPositionIntoFlockMap()
	s := NewSimulation()

	ebiten.SetWindowSize(s.ScreenWidth, s.ScreenWidth)
	ebiten.SetWindowTitle("Boids simulation")
	if err := ebiten.RunGame(s); err != nil {
		log.Fatal(err)
	}
}
