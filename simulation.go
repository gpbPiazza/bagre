package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth, screenHeight = 640, 360
)

var (
	green = color.RGBA{R: 10, G: 255, B: 50}
)

type Simulation struct {
	ScreenWidth  int
	ScreenHeight int
}

func NewSimulation() *Simulation {
	return &Simulation{
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
	}
}

func (s *Simulation) Update() error {
	return nil
}

func (s *Simulation) Draw(screen *ebiten.Image) {
	for _, boid := range flock {
		screen.Set(int(boid.position.x+1), int(boid.position.y+1), green)
		screen.Set(int(boid.position.x-1), int(boid.position.y-1), green)
		screen.Set(int(boid.position.x), int(boid.position.y+1), green)
		screen.Set(int(boid.position.x), int(boid.position.y-1), green)
	}

}

func (s *Simulation) Layout(_, _ int) (int, int) {
	return s.ScreenWidth, s.ScreenHeight
}
