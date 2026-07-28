package main

import (
	"bytes"
	"image"
	"image/color"
	"log/slog"
	"math"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	// How far a jelly senses neighbors. Keep it LOCAL: too big and every jelly
	// sees the whole population, so cohesion's "average position" becomes the
	// screen center and all of them collapse into one static blob. Small radius
	// = many sub-groups that drift, split, and re-form (a living school).
	// Lower bound: must exceed the drawn sprite size (48*0.5 = 24px) so a jelly
	// senses a neighbor BEFORE their images overlap.
	jellyViewRadius       = 20
	jellyAttackViewRadius = 25

	adjustRateAligment      = 0.15
	adjustRateCohesion      = 0.09
	adjustRateSeparation    = 0.30
	adjustRateSeparationWes = 0.45
	jellysCount             = 400
)

var (
	jellyWalkImg   *ebiten.Image
	jellyDeathImg  *ebiten.Image
	jellyAttackImg *ebiten.Image
)

type JellyFish struct {
	position     Vector2D
	velocity     Vector2D
	nextPosition Vector2D
	nextVelocity Vector2D

	tickWhenDied, id int
	tickStartAttack  int
	state            unitState
	logger           *slog.Logger
	events           *EventManager
}

func newJellyFish(id int, l *slog.Logger, e *EventManager) *JellyFish {
	b := &JellyFish{
		position:     Vector2D{x: rand.Float64() * screenWidth, y: rand.Float64() * screenHeight},
		velocity:     Vector2D{x: (rand.Float64() * 2) - 1.0, y: (rand.Float64() * 2) - 1.0},
		id:           id,
		state:        unitStateWalk,
		logger:       l,
		tickWhenDied: 0,
		events:       e,
	}
	b.nextVelocity = b.velocity
	b.nextPosition = b.position

	return b
}

func (j *JellyFish) Draw() (img *ebiten.Image, ticksWhenDead, tickCountPerPose int, frameCount int) {
	switch j.state {
	case unitStateWalk:
		return jellyWalkImg, 0, 5, 4
	case unitStateDead:
		return jellyDeathImg, j.tickWhenDied, 10, 6
	case unitStateAttack:
		return jellyAttackImg, 0, 5, 4
	default:
		return jellyWalkImg, 0, 5, 4
	}
}

func (j *JellyFish) Position() (float64, float64) {
	return j.position.x, j.position.y
}

func (j *JellyFish) Scale() (float64, float64) {
	return 1, 1
}

func (j *JellyFish) VecPosition() Vector2D {
	return j.position
}

func (j *JellyFish) ID() int {
	return j.id
}

func (j *JellyFish) VecVelocity() Vector2D {
	return j.velocity
}

func (j *JellyFish) State() unitState {
	return j.state
}

func (j *JellyFish) Attack(tick int) {
	if j.state == unitStateAttack {
		return
	}

	j.logger.Info("SETTING JELLY TO ATTACK")
	j.tickStartAttack = tick
	j.state = unitStateAttack
}

// TODO adicionar tests para essa parada como um todo desde
// comeu suficinete pro stage -> seta pro attack random -> acada 5 secs seleciona mais 10 random
// por quatro segundos elas ficam em estado de atack
func (j *JellyFish) checkState(tick int) {
	tickDiff := tick - j.tickStartAttack
	hasPassedFourSeconds := tickDiff >= 4*60

	if j.state == unitStateAttack && j.tickStartAttack != 0 {
		for _, unit := range j.unitsSeen(jellyAttackViewRadius) {
			if unit.IsPlayer() {
				j.events.Publish(wesTakeDMG, nil)
			}
		}
	}

	if hasPassedFourSeconds && j.tickStartAttack != 0 {
		j.state = unitStateWalk
		j.tickStartAttack = 0
	}
}

func (j *JellyFish) DrawAttackHitBox(s *ebiten.Image) {
	green := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	vector.StrokeRect(
		s,
		float32(j.position.x-jellyAttackViewRadius),
		float32(j.position.y-jellyAttackViewRadius),
		jellyAttackViewRadius*2,
		jellyAttackViewRadius*2,
		1,
		green,
		false,
	)
}

// unitSeen will iterate over a box in the grind where viewRadius val will create
// a upper point and a lower point from the the jelly position. Iteration start at first x element all smmaler y -> bigger, go to next x repeat y.
func (j *JellyFish) unitsSeen(viewRadius float64) []Unit {
	var seen []Unit

	upperView := j.position.AddVal(viewRadius)
	lowerView := j.position.AddVal(-viewRadius)

	for i := math.Max(lowerView.x, 0); i <= math.Min(upperView.x, screenWidth); i++ {
		for k := math.Max(lowerView.y, 0); k <= math.Min(upperView.y, screenHeight); k++ {
			seenUnitID := unitsPositions[int(i)][int(k)]
			if seenUnitID == -1 || j.id == seenUnitID {
				continue
			}
			if seenUnit, ok := units[seenUnitID]; ok {
				seen = append(seen, seenUnit)
			}
		}
	}

	return seen
}

func (j *JellyFish) Die(tick int) {
	if j.state == unitStateDead {
		return
	}

	j.tickWhenDied = tick
	j.state = unitStateDead
}

// // The behavior you actually want
//
// "Come together, keep moving, separate, re-form" = a living school. That requires:
//
// - Constant motion — fish always cruise; forces change direction, not whether they move.
// - Local awareness — each fish reacts to a few near neighbors, not the whole tank.
// - Tension, not equilibrium — separation slightly wins up close (no overlap), cohesion pulls loosely at medium range, alignment keeps sub-groups coherent — but they never perfectly cancel.
func (j *JellyFish) calcAcceleration() Vector2D {
	// all variables with prefix all here mean all elements inside of viewBox, inside of JellyFish View Radius
	allJellyFishsVelocity := Vector2D{x: 0, y: 0}
	allJellyFishsPosition := Vector2D{x: 0, y: 0}
	allJellyFishsSeparation := Vector2D{x: 0, y: 0}
	allFleeWes := Vector2D{x: 0, y: 0}

	jellysCount := 0.0
	wesSeen := false

	unistsSeen := j.unitsSeen(jellyViewRadius)

	for _, seenUnit := range unistsSeen {
		seenPosition := seenUnit.VecPosition()
		seenVelocity := seenUnit.VecVelocity()

		if seenUnit.IsPlayer() {
			wesSeen = true
			dist := seenPosition.Distance(j.position)
			separation := j.position.Subtract(seenPosition).DivisionVal(dist - dist/2) // push aways too close
			allFleeWes = allFleeWes.Add(separation)                                    // push aways too close
		}

		dist := seenPosition.Distance(j.position)
		if dist < jellyViewRadius {
			jellysCount++
			allJellyFishsVelocity = allJellyFishsVelocity.Add(seenVelocity)   // Direction
			allJellyFishsPosition = allJellyFishsPosition.Add(seenPosition)   // move to the center
			separation := j.position.Subtract(seenPosition).DivisionVal(dist) // push aways too close
			allJellyFishsSeparation = allJellyFishsSeparation.Add(separation) // push aways too close
		}
	}

	borderBounceX := j.borderBounce(j.position.x, screenWidth)
	borderBouncey := j.borderBounce(j.position.y, screenHeight)
	accel := Vector2D{x: borderBounceX, y: borderBouncey}

	if jellysCount > 0 {
		avgVelocity := allJellyFishsVelocity.DivisionVal(jellysCount)
		avgPosition := allJellyFishsPosition.DivisionVal(jellysCount)
		accelAligment := avgVelocity.Subtract(j.velocity).MultiplyVal(adjustRateAligment)
		accelCohesion := avgPosition.Subtract(j.position).MultiplyVal(adjustRateCohesion)
		accelSepartion := allJellyFishsSeparation.MultiplyVal(adjustRateSeparation)
		accel = accel.Add(accelAligment).Add(accelCohesion).Add(accelSepartion)
	}

	if wesSeen {
		accel = accel.Add(allFleeWes.MultiplyVal(adjustRateSeparationWes))
	}

	return accel
}

// Quanto mais próximo da borda mais rápido será o bounce
func (j *JellyFish) borderBounce(pos, border float64) float64 {

	// Está próximo da bater na borda, passou do limite de vistualização
	// ou seja o passarinho viu a parede e irá mudar de direção
	if pos < jellyViewRadius {
		// Cap pos at 0.5: at pos <= 0 (already past the border) a raw 1/pos
		// flips sign (or hits +Inf) and shoves the jelly further off-screen,
		// so keep the push finite and always pointing back inside.
		return 1 / math.Max(pos, 0.5)
	}
	// Is the same thing but in the other side of the screenView
	// o primeiro If é para o jelly que está próximo da parede em que X é muito pequeno
	// Aqui o x é grande, o mesmo para Y.
	if pos > border-jellyViewRadius {
		return 1 / math.Min(pos-border, -0.5)
	}

	return 0
}

func (j *JellyFish) IsPlayer() bool {
	return false
}

func (j *JellyFish) nextMove() {
	if j.state == unitStateDead {
		return
	}

	minSpeed := 0.5
	maxSpeed := 1.5

	accel := j.calcAcceleration()

	j.nextVelocity = j.velocity.Add(accel)

	// Cruising speed: only change the LENGTH of velocity, never its direction.
	// Floor (minSpeed) = never stall into a frozen blob.
	// Ceiling (maxSpeed) = never blast off across the screen.
	velocityMag := j.nextVelocity.Pythagoras()
	if velocityMag < minSpeed {
		j.nextVelocity = j.nextVelocity.ScaleToLength(minSpeed)
	}
	if velocityMag > maxSpeed {
		j.nextVelocity = j.nextVelocity.ScaleToLength(maxSpeed)
	}

	j.nextPosition = j.position.Add(j.nextVelocity)

	// Hard wall: rebuildGrid indexes unitsPositions by position, so a jelly
	// outside [0, screenWidth] x [0, screenHeight] panics the game. The
	// bounce steers away from borders, but flee/flock forces can overpower
	// it in a single frame — clamp as the last line of defense.
	j.nextPosition.x = math.Min(math.Max(j.nextPosition.x, 0), screenWidth)
	j.nextPosition.y = math.Min(math.Max(j.nextPosition.y, 0), screenHeight)
}

func loadJellyImg() (err error) {
	jellyBytesPng, err := os.ReadFile("./assets/jellyfish/Walk.png")
	if err != nil {
		return err
	}

	jellyImg, _, err := image.Decode(bytes.NewReader(jellyBytesPng))
	if err != nil {
		return err
	}
	jellyWalkImg = ebiten.NewImageFromImage(jellyImg)

	diePng, err := os.ReadFile("./assets/jellyfish/Death.png")
	if err != nil {
		return err
	}
	dieImg, _, err := image.Decode(bytes.NewReader(diePng))
	if err != nil {
		return err
	}
	jellyDeathImg = ebiten.NewImageFromImage(dieImg)

	attackPng, err := os.ReadFile("./assets/jellyfish/Attack.png")
	if err != nil {
		return err
	}
	attackImg, _, err := image.Decode(bytes.NewReader(attackPng))
	if err != nil {
		return err
	}
	jellyAttackImg = ebiten.NewImageFromImage(attackImg)

	return nil
}

func (j *JellyFish) writeMove() {
	j.position = j.nextPosition
	j.velocity = j.nextVelocity
}
