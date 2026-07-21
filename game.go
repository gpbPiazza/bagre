package main

import (
	"fmt"
	"image/color"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth, screenHeight = 950, 550
)

var (
	darkGrey = color.RGBA{R: 40, G: 45, B: 60, A: 255}
)

type Game struct {
	ScreenWidth   int
	ScreenHeight  int
	tick          int // grows by 1 every Update (~60/sec); our clock for animation
	units         Units
	evenetManager *EventManager
}

func NewGame(gUnits Units, evenEventManager *EventManager) *Game {
	return &Game{
		ScreenWidth:   screenWidth,
		ScreenHeight:  screenHeight,
		units:         gUnits,
		evenetManager: evenEventManager,
	}
}

func (g *Game) Update() error {
	g.tick++

	g.units.wes.move()

	for _, j := range g.units.smack {
		j.nextMove()
	}
	for _, j := range g.units.smack {
		j.writeMove()
	}

	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		fmt.Println("RELEASED SPACE")
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		fmt.Println("HITTED SPACE")
		g.units.unitsEaten = append(g.units.unitsEaten, g.units.wes.Attack(g.tick)...)
	}

	for _, u := range g.units.unitsEaten {
		u.Die(g.tick)
	}

	for _, u := range g.units.unitsEaten {
		checkState(u, g.tick, g.evenetManager)
	}

	checkState(g.units.wes, g.tick, g.evenetManager)

	rebuildGrid()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(darkGrey)
	for _, u := range units {
		drawUnit(screen, u, g.tick)
	}

	// TEMP: visualize Wes's attack hitbox. Delete when done debugging.
	g.units.wes.DrawAttackHitBox(screen)
}

func (g *Game) Handle(et EventType, payload any) {
	switch et {
	case removeUnit:
		g.handleRemoveUnit(payload)
	default:
		return
	}
}

func (g *Game) handleRemoveUnit(payload any) {
	u, ok := payload.(Unit)
	if !ok {
		panic("tu ta fazendo merda gabriel")
	}

	delete(units, u.ID())
	position := u.VecPosition()
	unitsByPositions[int(position.x)][int(position.y)] = -1
	g.units.smack = slices.DeleteFunc(g.units.smack, func(e *JellyFish) bool {
		return e.ID() == u.ID()
	})

}

func (g *Game) Layout(_, _ int) (int, int) {
	return g.ScreenWidth, g.ScreenHeight
}

func rebuildGrid() {
	for i := range unitsByPositions {
		for k := range unitsByPositions[i] {
			unitsByPositions[i][k] = -1
		}
	}

	for id, u := range units {
		p := u.VecPosition()
		unitsByPositions[int(p.x)][int(p.y)] = id
	}
}
