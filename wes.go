package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log/slog"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/gpbPiazza/bagre/pkg/log"
)

// conforme o wes for comendo os peixes ele fica maior
// seu tamanho determina se  o inimigo do wes te medo dele ou não
// inimigo do wes quer matar ele. Wes quer comer maximo de peixer possíveis
// Ai nivel de dificuldade pro player
// - o cardume vai ter sprint com CD por tempo
// - ⁠se o predador ficar envolto pot muito jelly fish, eles vai soltar choque
// - ⁠no meio do jogo vai vir o inimigo do wes, um polvo vermelho puto, querendo matar o wes
// - ⁠o wes toda vez que come um  peixe ele cresce
// - ⁠o wes ficar mt grande faz o inimigo do wes ter medo dele, logo, inimigo do wes nao mata ele
// - ⁠wes ganha o jogo de comer todos os peixes
// - ⁠wes perde de ele morrer
// - ⁠wes perde se acabar o tempo
// - ⁠antes do inimigo do wes aparecer, vai ter tartarugas, elas nao dao dano no wes mas elas comem os peixes dele, logo, ele fica menor. Ele compete com elas e elas empurram ele tomando stun

// TODOs para codar
// jelly fish eaten count
// vida do wes
//
// jelly fish aleatorias ficam em estado de attack
// caso wes coma elas ele toma dano.

type Wes struct {
	position              Vector2D
	velocity              Vector2D
	id                    int
	tickesWhenAttackState int
	tickesWhenHurtState   int
	counter               *Counter
	life                  int

	imgs map[unitState]*ebiten.Image

	state  unitState
	logger *slog.Logger
}

func NewWes(
	id int,
	l *slog.Logger,
	eventManager *EventManager,
	c *Counter,
) *Wes {
	imgs, err := loadWesImg()
	if err != nil {
		l.Error("failed to load wes imgs", log.Err(err))
		panic(err)
	}

	w := &Wes{
		position:              Vector2D{x: screenWidth / 2, y: screenHeight / 2},
		velocity:              Vector2D{x: 1.0, y: 1.0},
		id:                    id,
		tickesWhenAttackState: 0,
		tickesWhenHurtState:   0,
		counter:               c,
		state:                 unitStateWalk,
		logger:                l,
		life:                  wesStartingLife,
		imgs:                  imgs,
	}

	eventManager.subscribe(attackAnimationEnded, w)
	eventManager.subscribe(wesTakeDMG, w)
	return w
}

func (w *Wes) Draw() (img *ebiten.Image, ticksWhenStateChanged, tickCountPerPose int, frameCount int) {
	switch w.state {
	case unitStateAttack:
		return w.imgs[w.state], w.tickesWhenAttackState, wesAttackTicksPerPose, wesAttackFrameCount
	case unitStateWalk:
		return w.imgs[w.state], 0, wesWalkTicksPerPose, wesWalkFrameCount
	case unitStateHurt:
		return w.imgs[w.state], 0, wesHurtTicksPerPose, wesHurtFrameCount
	default:
		return w.imgs[w.state], 0, wesWalkTicksPerPose, wesWalkFrameCount
	}
}

func (w *Wes) Position() (float64, float64) {
	return w.position.x, w.position.y
}

func (w *Wes) Handle(et EventType, payload any) {
	switch et {
	case attackAnimationEnded:
		w.state = unitStateWalk
	case wesTakeDMG:
		w.takeDMG(payload)
	default:
		return
	}
}

func (w *Wes) takeDMG(payload any) {
	if w.state == unitStateHurt {
		return
	}

	ticks, ok := payload.(int)
	if !ok {
		return
	}

	w.tickesWhenHurtState = ticks
	w.state = unitStateHurt
	w.life--

	if w.life == 0 {
		w.logger.Warn("WES should be dead IMPLMENT death animation and restart the stage")
		// w.state = unitStateDead
		// TODO implement animation of wes dieing
		// currently when wes die he just vanish from the screen.
	}
}

func (w *Wes) checkState(tick int) {
	if w.state == unitStateHurt {
		tickDiff := tick - w.tickesWhenHurtState
		hasPassedTwoSeconds := tickDiff >= wesHurtDurationTicks

		if hasPassedTwoSeconds {
			w.state = unitStateWalk
			w.tickesWhenHurtState = 0
		}
	}
}

// TODO: implement Draw user UI where will receive wes life
// and counter as depency and draw a betiuful ui of wes spells life and icons

func (w *Wes) DrawLife(screen *ebiten.Image) {
	msg := fmt.Sprintf("Vida: %d/%d", w.life, wesStartingLife)

	face := &text.GoTextFace{
		Source: textFont,
		Size:   counterFontSize,
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(wesLifePosX, wesLifePosY)
	op.ColorScale.ScaleWithColor(color.White)

	text.Draw(screen, msg, face, op)
}

func (w *Wes) State() unitState {
	return w.state
}

func (w *Wes) VecPosition() Vector2D {
	return w.position
}

func (w *Wes) ID() int {
	return w.id
}

func (w *Wes) VecVelocity() Vector2D {
	return w.velocity
}

func (w *Wes) Scale() (float64, float64) {
	switch w.state {
	case unitStateAttack:
		return 1, 1
	default:
		return 1, 1
	}
}

func (w *Wes) move() {
	newPosition := w.position
	if isKeyPressed(ebiten.KeyW) {
		newPosition.y = w.position.SubY(wesSpeedRation).y
	}
	if isKeyPressed(ebiten.KeyS) {
		newPosition.y = w.position.AddY(wesSpeedRation).y
	}
	if isKeyPressed(ebiten.KeyA) {
		newPosition.x = w.position.SubX(wesSpeedRation).x
	}
	if isKeyPressed(ebiten.KeyD) {
		newPosition.x = w.position.AddX(wesSpeedRation).x
	}

	if newPosition.x > screenWidth-wesBorderMargin {
		newPosition.x = screenWidth - wesBorderMargin
	}
	if newPosition.x < wesBorderMargin {
		newPosition.x = wesBorderMargin
	}
	if newPosition.y > screenHeight-wesBorderMargin {
		newPosition.y = screenHeight - wesBorderMargin
	}
	if newPosition.y < wesBorderMargin {
		newPosition.y = wesBorderMargin
	}

	w.position = newPosition
}

func (w *Wes) Die(_ int) {}

func (w *Wes) IsPlayer() bool { return true }

func (w *Wes) Attack(tick int) []Unit {
	if w.state == unitStateAttack {
		return nil
	}

	w.state = unitStateAttack
	w.tickesWhenAttackState = tick

	const n = 50
	ax, ay := w.position.x+n, w.position.y-n
	by := w.position.y + n
	cx, cy := w.position.x, w.position.y

	var unitsEaten []Unit

	// Upper half of the triangle: from C's row up to A's row, each row
	// starting one column further right so the edge C->A stays diagonal.
	x0 := cx
	for k := math.Min(cy, screenHeight); k >= math.Max(ay, 0); k-- {
		unitsEaten = append(unitsEaten, w.unitsInRow(k, x0, ax)...)
		x0++
	}

	// Lower half: from the row below C down to B's row, same shrinking rows.
	x0 = cx + 1
	for k := math.Max(cy+1, 0); k <= math.Min(by, screenHeight); k++ {
		unitsEaten = append(unitsEaten, w.unitsInRow(k, x0, ax)...)
		x0++
	}

	w.counter.Add(len(unitsEaten))

	return unitsEaten
}

// unitsInRow collects every unit on grid row k between columns fromX and toX
// (clamped to the screen), skipping empty cells and wes himself.
func (w *Wes) unitsInRow(k, fromX, toX float64) []Unit {
	var found []Unit
	for i := math.Max(fromX, 0); i <= math.Min(toX, screenWidth); i++ {
		seenUnitID := unitsPositions[int(i)][int(k)]
		if seenUnitID == -1 || w.id == seenUnitID {
			continue
		}
		if seenUnit, ok := units[seenUnitID]; ok {
			found = append(found, seenUnit)
		}
	}

	return found
}

// DrawAttackHitBox is a TEMP debug helper: it outlines Wes's attack range as a
// green square (2*wesViewRadius on a side, centered on Wes). Delete when done.
func (w *Wes) DrawAttackHitBox(screen *ebiten.Image) {
	green := hitBoxGreen
	// Square hit box
	// vector.StrokeRect(
	// 	screen,
	// 	float32(w.position.x-wesViewRadius),
	// 	float32(w.position.y-wesViewRadius),
	// 	wesViewRadius*2,
	// 	wesViewRadius*2,
	// 	1,
	// 	green,
	// 	false,
	// )
	//
	// Triangle hit box: C (Wes) -> A (right-up) -> B (right-down) -> back to C.
	const n = 25
	ax, ay := float32(w.position.x+n), float32(w.position.y-n)
	bx, by := float32(w.position.x+n), float32(w.position.y+n)
	cx, cy := float32(w.position.x), float32(w.position.y)

	vector.StrokeLine(screen, cx, cy, ax, ay, 1, green, false) // C -> A
	vector.StrokeLine(screen, ax, ay, bx, by, 1, green, false) // A -> B
	vector.StrokeLine(screen, bx, by, cx, cy, 1, green, false) // B -> C
}

func loadImage(path string) (*ebiten.Image, error) {
	rawImg, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	imgPNG, _, err := image.Decode(bytes.NewReader(rawImg))
	if err != nil {
		return nil, err
	}

	ebiImg := ebiten.NewImageFromImage(imgPNG)

	return ebiImg, nil
}

func loadWesImg() (map[unitState]*ebiten.Image, error) {
	imgsByState := make(map[unitState]*ebiten.Image)

	walkWesImg, err := loadImage("./assets/wes/Walk.png")
	if err != nil {
		return nil, err
	}

	attackWesImg, err := loadImage("./assets/wes/Attack.png")
	if err != nil {
		return nil, err
	}

	hurtWesImg, err := loadImage("./assets/wes/Hurt.png")
	if err != nil {
		return nil, err
	}

	imgsByState[unitStateWalk] = walkWesImg
	imgsByState[unitStateAttack] = attackWesImg
	imgsByState[unitStateHurt] = hurtWesImg

	return imgsByState, nil
}
