package main

import (
	"bytes"
	"image/color"
	"log/slog"
	"os"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	screenWidth, screenHeight = 950, 550
)

var (
	darkGrey = color.RGBA{R: 40, G: 45, B: 60, A: 255}
	// counterFaceSource holds the parsed font; parsing is expensive so it
	// happens once at init, never inside Draw.
	// TODO create a system start up
	textFont = func() *text.GoTextFaceSource {
		s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
		if err != nil {
			panic(err)
		}
		return s
	}()
)

type Game struct {
	ScreenWidth   int
	ScreenHeight  int
	tick          int // grows by 1 every Update (~60/sec); our clock for animation
	units         Units
	counter       *Counter
	evenetManager *EventManager
	DrawHitBox    bool
}

func NewGame(l *slog.Logger) *Game {
	eventManager := NewEventManager()

	counter := NewCounter(jellysCount)
	gUnits := NewUnits(eventManager, l, counter)

	g := &Game{
		ScreenWidth:   screenWidth,
		ScreenHeight:  screenHeight,
		units:         gUnits,
		evenetManager: eventManager,
		DrawHitBox:    os.Getenv("DRAW_HIT_BOX") != "",
		counter:       counter,
	}

	eventManager.subscribe(removeUnit, g)

	return g
}

func (g *Game) Update() error {
	g.tick++

	for _, j := range g.units.smack {
		j.nextMove()
	}
	for _, j := range g.units.smack {
		j.writeMove()
	}

	g.units.wes.move()

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
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

func (g *Game) Draw(s *ebiten.Image) {
	s.Fill(darkGrey)

	g.counter.Draw(s)
	g.units.wes.DrawLife(s)

	for _, u := range units {
		drawUnit(s, u, g.tick)
	}

	if g.DrawHitBox {
		g.units.wes.DrawAttackHitBox(s)
	}
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
	unitsPositions[int(position.x)][int(position.y)] = -1
	g.units.smack = slices.DeleteFunc(g.units.smack, func(e *JellyFish) bool {
		return e.ID() == u.ID()
	})

}

func (g *Game) Layout(_, _ int) (int, int) {
	return g.ScreenWidth, g.ScreenHeight
}

func rebuildGrid() {
	for i := range unitsPositions {
		for k := range unitsPositions[i] {
			unitsPositions[i][k] = -1
		}
	}

	for id, u := range units {
		p := u.VecPosition()
		unitsPositions[int(p.x)][int(p.y)] = id
	}
}
