package main

import (
	"bytes"
	"image/color"
	"log/slog"
	"math/rand/v2"
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

	logger                *slog.Logger
	tickEletricJellyStart int
}

func NewGame(l *slog.Logger) *Game {
	eventManager := NewEventManager()

	counter := NewCounter(jellysCount, eventManager)
	gUnits := NewUnits(eventManager, l, counter)

	g := &Game{
		ScreenWidth:   screenWidth,
		ScreenHeight:  screenHeight,
		units:         gUnits,
		evenetManager: eventManager,
		DrawHitBox:    os.Getenv("DRAW_HIT_BOX") != "",
		counter:       counter,
		logger:        l,
	}

	eventManager.subscribe(removeUnit, g)
	eventManager.subscribe(startEletricJellyFishs, g)

	return g
}

func (g *Game) Update() error {
	g.tick++

	for _, j := range g.units.smack {
		j.nextMove()
	}
	for _, j := range g.units.smack {
		j.writeMove()
		j.checkState(g.tick)
	}

	var underAttack []*JellyFish
	for _, j := range g.units.smack {
		if j.state == unitStateAttack {
			underAttack = append(underAttack, j)
		}
	}
	if len(underAttack) > 10 {
		panic("we have more jelly under attack them expected")
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

	// TODO create a struct reponsiable to deal with
	// wes enemies stages and how they will come into the game.
	g.stage(g.tick)

	rebuildGrid()
	return nil
}

func (g *Game) Draw(s *ebiten.Image) {
	s.Fill(darkGrey)

	g.counter.Draw(s)
	g.units.wes.DrawLife(s)

	// NOTE: if we draw wes after smack they will not see wes. order here matters
	// because smack need calc next position then write their move
	// drawUnit(s, g.units.wes, g.tick)
	// PQ O WES TA SUMINDO DO MAPA?
	// Quando eu parei de iterar no map de units as smacks simplemente pararm de considerar a posição do wes.
	// TODO: tenho que fazer tests para:
	// smack ver o wes e fugir dele
	// wes come jelly e assert em tudo que deve acontecer
	//  - jelly morrendo animação e parar de se mover
	//  - wes esperar animação de atack acabar para atacar denovo
	//  - counter ++
	for _, u := range units {
		drawUnit(s, u, g.tick)
	}

	if g.DrawHitBox {
		g.units.wes.DrawAttackHitBox(s)
		for _, u := range g.units.smack {
			u.DrawAttackHitBox(s)
		}
	}
}

func (g *Game) Handle(et EventType, payload any) {
	switch et {
	case removeUnit:
		g.handleRemoveUnit(payload)
	case startEletricJellyFishs:
		g.tickEletricJellyStart = g.tick
	default:
		return
	}
}

// stage checks what ever should happend given the stage that we are
// currently we have stage 1 - eletricJellies,
//
// map 1
// stages of map 1
//
//	0 - nothing passive jelly
//	1 - eletric jelly
//	2 - turtles in the map
//	3 - shark attack
//	4 - octopus hold
//	5 - octopus boss
func (g *Game) stage(tick int) {
	tickDiff := tick - g.tickEletricJellyStart
	// Fire ONCE per 5s boundary. tickDiff/60 holds each value for a whole
	// second, so `% 5` on it would be true for all 60 ticks of that second and
	// call smackAttack 60x per boundary — each pass selecting a new random
	// chunk, blowing past the 10-attacker cap. Modulo the RAW tick count
	// instead: only one tick in every 300 (5*60) satisfies this.
	hasPassedFiveSeconds := tickDiff > 0 && tickDiff%(5*60) == 0

	if !(hasPassedFiveSeconds && g.tickEletricJellyStart != 0) {
		return
	}

	g.logger.Info("passed 5 seconds trigger eletric Smack")
	g.smackAttack(tick)
}

func (g *Game) smackAttack(tick int) {
	var chunks [][]*JellyFish
	chunkLen := len(g.units.smack) / 4

	var chunk []*JellyFish
	for _, s := range g.units.smack {
		chunk = append(chunk, s)

		if len(chunk) != chunkLen {
			continue
		}

		chunks = append(chunks, chunk)
		chunk = nil
	}

	index := rand.IntN(len(chunks))
	count := 0
	jellyAttackLimit := 10
	for _, j := range chunks[index] {
		if count < jellyAttackLimit {
			count++
			j.Attack(tick)
		}
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
