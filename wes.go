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
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// wesViewRadius is how far Wes's attack reaches. Shared by Attack (the real
// hitbox) and DrawAttackHitBox (the debug visualization) so they never drift.
const wesViewRadius = 100

var (
	attackWesImg *ebiten.Image
	walkWesImg   *ebiten.Image
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

type Wes struct {
	position Vector2D
	velocity Vector2D
	id       int

	state  unitState
	logger *slog.Logger
}

func NewWes(id int, l *slog.Logger) *Wes {
	w := &Wes{
		position: Vector2D{x: screenWidth / 2, y: screenHeight / 2},
		velocity: Vector2D{x: 1.0, y: 1.0},
		id:       id,
		state:    unitStateWalk,
		logger:   l,
	}
	return w
}

func (w *Wes) Draw() (img *ebiten.Image, tickesWhenDead, tickCountPerPose int, frameCount int) {
	switch w.state {
	case unitStateAttack:
		return attackWesImg, 0, 3, 6
	case unitStateWalk:
		return walkWesImg, 0, 5, 4
	default:
		return walkWesImg, 0, 5, 4
	}
}

func (w *Wes) Position() (float64, float64) {
	return w.position.x, w.position.y
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
		return 1.3, 1.3
	default:
		return 1, 1
	}
}

func (w *Wes) move() {
	const speed = 2.0

	newPosition := w.position
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		newPosition.y = w.position.SubY(speed).y
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		newPosition.y = w.position.AddY(speed).y
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		newPosition.x = w.position.SubX(speed).x
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		newPosition.x = w.position.AddX(speed).x
	}

	// 10 px visually good
	if newPosition.x > screenWidth-10 {
		newPosition.x = screenWidth - 10
	}
	if newPosition.x < 10 {
		newPosition.x = 10
	}
	if newPosition.y > screenHeight-10 {
		newPosition.y = screenHeight - 10
	}
	if newPosition.y < 10 {
		newPosition.y = 10
	}

	w.position = newPosition

}

func (w *Wes) Die(tick int) {}

func (w *Wes) IsPlayer() bool { return true }

// TODO
// 1. fix attack hit box to be a trapézio and not a square
// 2. fix animation to
// 2.1 find the units to  attack
// 2.2 set animation of atack
// 2.3 wait animation of atack to finish
// 2.4 put animation of runing back
// 3. We should only suport press and release action to attack
//

func (w *Wes) Attack() []Unit {
	w.state = unitStateAttack

	// upperView := w.position.AddVal(wesViewRadius)
	// lowerView := w.position.AddVal(-wesViewRadius)

	const n = 20
	ax, _ := float64(w.position.x+n), float64(w.position.y-n)
	_, by := float64(w.position.x+n), float64(w.position.y+n)
	cx, cy := float64(w.position.x), float64(w.position.y)

	var unitsEaten []Unit
	for i := cx; i <= math.Min(ax, screenWidth); i++ {
		for k := cy; k <= math.Min(by, screenHeight); k++ {
			fmt.Println("scan sequence", "x", i, "y", k)

			seenUnitID := unitsByPositions[int(i)][int(k)]
			if seenUnitID == -1 || w.id == seenUnitID {
				continue
			}
			seenUnit := units[seenUnitID]

			unitsEaten = append(unitsEaten, seenUnit)
		}
	}

	return unitsEaten
}

// DrawAttackHitBox is a TEMP debug helper: it outlines Wes's attack range as a
// green square (2*wesViewRadius on a side, centered on Wes). Delete when done.
func (w *Wes) DrawAttackHitBox(screen *ebiten.Image) {
	// Square hit box
	green := color.RGBA{R: 0, G: 255, B: 0, A: 255}
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

	// Triangle hit box: C (Wes) -> A (right-up) -> B (right-down) -> back to C.
	const n = 20
	ax, ay := float32(w.position.x+n), float32(w.position.y-n)
	bx, by := float32(w.position.x+n), float32(w.position.y+n)
	cx, cy := float32(w.position.x), float32(w.position.y)

	vector.StrokeLine(screen, cx, cy, ax, ay, 1, green, false) // C -> A
	vector.StrokeLine(screen, ax, ay, bx, by, 1, green, false) // A -> B
	vector.StrokeLine(screen, bx, by, cx, cy, 1, green, false) // B -> C
}

func loadWesImg() error {
	walk, err := os.ReadFile("./assets/wes/Walk.png")
	if err != nil {
		return err
	}

	walkImg, _, err := image.Decode(bytes.NewReader(walk))
	if err != nil {
		return err
	}

	walkWesImg = ebiten.NewImageFromImage(walkImg)

	attack, err := os.ReadFile("./assets/wes/Attack.png")
	if err != nil {
		return err
	}

	attackImg, _, err := image.Decode(bytes.NewReader(attack))
	if err != nil {
		return err
	}

	attackWesImg = ebiten.NewImageFromImage(attackImg)

	return nil
}
