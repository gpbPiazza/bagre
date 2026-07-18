package main

import (
	"bytes"
	"image"
	"log/slog"
	"math"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	// How far a jelly senses neighbors. Keep it LOCAL: too big and every jelly
	// sees the whole population, so cohesion's "average position" becomes the
	// screen center and all of them collapse into one static blob. Small radius
	// = many sub-groups that drift, split, and re-form (a living school).
	// Lower bound: must exceed the drawn sprite size (48*0.5 = 24px) so a jelly
	// senses a neighbor BEFORE their images overlap.
	jellyViewRadius      = 20
	adjustRateAligment   = 0.15
	adjustRateCohesion   = 0.09
	adjustRateSeparation = 0.30
	jellysCount          = 250
)

var (
	unitsByPositions = [screenWidth + 1][screenHeight + 1]int{}
	// units is the collective noun for jellyFish
	units = make(map[int]Unit, 0)

	jellyWalkImg  *ebiten.Image
	jellyDeathImg *ebiten.Image
)

type JellyFish struct {
	position     Vector2D
	velocity     Vector2D
	nextPosition Vector2D
	nextVelocity Vector2D

	tickWhenDied, id int
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
	upperView := j.position.AddVal(jellyViewRadius)
	lowerView := j.position.AddVal(-jellyViewRadius)
	// all variables with prefix all here mean all elements inside of viewBox, inside of JellyFish View Radius
	allJellyFishsVelocity := Vector2D{x: 0, y: 0}
	allJellyFishsPosition := Vector2D{x: 0, y: 0}
	allJellyFishsSeparation := Vector2D{x: 0, y: 0}
	allFleeWes := Vector2D{x: 0, y: 0}

	jellysCount := 0.0
	wesSeen := false

	for i := math.Max(lowerView.x, 0); i <= math.Min(upperView.x, screenWidth); i++ {
		for k := math.Max(lowerView.y, 0); k <= math.Min(upperView.y, screenHeight); k++ {
			seenUnitID := unitsByPositions[int(i)][int(k)]

			if seenUnitID == -1 || j.id == seenUnitID {
				continue
			}

			seenUnit := units[seenUnitID]
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
		accel = accel.Add(allFleeWes)
	}

	return accel
}

// Quanto mais próximo da borda mais rápido será o bounce
func (j *JellyFish) borderBounce(pos, border float64) float64 {

	// Está próximo da bater na borda, passou do limite de vistualização
	// ou seja o passarinho viu a parede e irá mudar de direção
	if pos < jellyViewRadius {
		return 1 / pos
	}
	// Is the same thing but in the other side of the screenView
	// o primeiro If é para o jelly que está próximo da parede em que X é muito pequeno
	// Aqui o x é grande, o mesmo para Y.
	if pos > border-jellyViewRadius {
		return 1 / (pos - border)
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
}

func loadJellyImg() error {
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

	return nil
}

func (j *JellyFish) writeMove() {
	j.position = j.nextPosition
	j.velocity = j.nextVelocity
}
