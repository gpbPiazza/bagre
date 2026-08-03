package main

import (
	"log/slog"
	"math"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWesMoves(t *testing.T) {
	t.Run("wes start the game at the center", func(t *testing.T) {
		game := NewGame(slog.Default())
		require.Equal(t, game.units.wes.position, NewVector(screenWidth/2, screenHeight/2))
	})

	t.Run("wes press W moves up", func(t *testing.T) {
		isKeyPressed = isKeyPressedTest(t, ebiten.KeyW)
		t.Cleanup(func() { isKeyPressed = isKeyPressedEmptyState(t) })

		game := NewGame(slog.Default())

		startingX, startingY := game.units.wes.Position()

		for range testOneSec {
			isKeyPressed(ebiten.KeyW)
			err := game.Update()
			require.NoError(t, err)
		}

		require.Equal(t, testOneSec, game.tick)
		gotX, gotY := game.units.wes.Position()

		require.InDelta(t, startingX, gotX, floatTolerance, "X position")
		require.InDelta(t, startingY-wesSpeedRation*testOneSec, gotY, floatTolerance, "Y position")
	})

	t.Run("wes press S moves down", func(t *testing.T) {
		isKeyPressed = isKeyPressedTest(t, ebiten.KeyS)
		t.Cleanup(func() { isKeyPressed = isKeyPressedEmptyState(t) })

		game := NewGame(slog.Default())

		startingX, startingY := game.units.wes.Position()

		for range testOneSec {
			isKeyPressed(ebiten.KeyS)
			err := game.Update()
			require.NoError(t, err)
		}

		require.Equal(t, testOneSec, game.tick)
		gotX, gotY := game.units.wes.Position()

		require.InDelta(t, startingX, gotX, floatTolerance, "X position")
		require.InDelta(t, startingY+wesSpeedRation*testOneSec, gotY, floatTolerance, "Y position")
	})

	t.Run("wes press A moves left", func(t *testing.T) {
		isKeyPressed = isKeyPressedTest(t, ebiten.KeyA)
		t.Cleanup(func() { isKeyPressed = isKeyPressedEmptyState(t) })

		game := NewGame(slog.Default())

		startingX, startingY := game.units.wes.Position()

		for range testOneSec {
			isKeyPressed(ebiten.KeyA)
			err := game.Update()
			require.NoError(t, err)
		}

		require.Equal(t, testOneSec, game.tick)
		gotX, gotY := game.units.wes.Position()

		require.InDelta(t, startingX-wesSpeedRation*testOneSec, gotX, floatTolerance, "X position")
		require.InDelta(t, startingY, gotY, floatTolerance, "Y position")
	})

	t.Run("wes press D moves right", func(t *testing.T) {
		isKeyPressed = isKeyPressedTest(t, ebiten.KeyD)
		t.Cleanup(func() { isKeyPressed = isKeyPressedEmptyState(t) })

		game := NewGame(slog.Default())

		startingX, startingY := game.units.wes.Position()

		for range testOneSec {
			isKeyPressed(ebiten.KeyD)
			err := game.Update()
			require.NoError(t, err)
		}

		require.Equal(t, testOneSec, game.tick)
		gotX, gotY := game.units.wes.Position()

		require.InDelta(t, startingX+wesSpeedRation*testOneSec, gotX, floatTolerance, "X position")
		require.InDelta(t, startingY, gotY, floatTolerance, "Y position")
	})

	t.Run("wes cannot move beyond the screen limits", func(t *testing.T) {
		t.Run("pressing D stops at the right border", func(t *testing.T) {
			isKeyPressed = isKeyPressedTest(t, ebiten.KeyD)
			t.Cleanup(func() { isKeyPressed = isKeyPressedEmptyState(t) })

			game := NewGame(slog.Default())
			game.units.wes.position = NewVector(screenWidth-20, screenHeight/2)

			_, startingY := game.units.wes.Position()

			for range testOneSec {
				isKeyPressed(ebiten.KeyD)
				err := game.Update()
				require.NoError(t, err)
			}

			require.Equal(t, testOneSec, game.tick)
			gotX, gotY := game.units.wes.Position()

			require.InDelta(t, float64(screenWidth-wesBorderMargin), gotX, floatTolerance, "X clamped at the right border")
			require.InDelta(t, startingY, gotY, floatTolerance, "Y position")
		})

		t.Run("pressing A stops at the left border", func(t *testing.T) {
			isKeyPressed = isKeyPressedTest(t, ebiten.KeyA)
			t.Cleanup(func() {
				isKeyPressed = isKeyPressedEmptyState(t)
			})

			game := NewGame(slog.Default())
			game.units.wes.position = NewVector(20, screenHeight/2)

			_, startingY := game.units.wes.Position()

			for range testOneSec {
				isKeyPressed(ebiten.KeyA)
				err := game.Update()
				require.NoError(t, err)
			}

			require.Equal(t, testOneSec, game.tick)
			gotX, gotY := game.units.wes.Position()

			require.InDelta(t, float64(wesBorderMargin), gotX, floatTolerance, "X clamped at the left border")
			require.InDelta(t, startingY, gotY, floatTolerance, "Y position")
		})

		t.Run("pressing W stops at the top border", func(t *testing.T) {
			isKeyPressed = isKeyPressedTest(t, ebiten.KeyW)
			t.Cleanup(func() {
				isKeyPressed = isKeyPressedEmptyState(t)
			})

			game := NewGame(slog.Default())
			game.units.wes.position = NewVector(screenWidth/2, 20)

			startingX, _ := game.units.wes.Position()

			for range testOneSec {
				isKeyPressed(ebiten.KeyW)
				err := game.Update()
				require.NoError(t, err)
			}

			require.Equal(t, testOneSec, game.tick)
			gotX, gotY := game.units.wes.Position()

			require.InDelta(t, startingX, gotX, floatTolerance, "X position")
			require.InDelta(t, float64(wesBorderMargin), gotY, floatTolerance, "Y clamped at the top border")
		})

		t.Run("pressing S stops at the bottom border", func(t *testing.T) {
			isKeyPressed = isKeyPressedTest(t, ebiten.KeyS)
			t.Cleanup(func() {
				isKeyPressed = isKeyPressedEmptyState(t)
			})

			game := NewGame(slog.Default())
			game.units.wes.position = NewVector(screenWidth/2, screenHeight-20)

			startingX, _ := game.units.wes.Position()

			for range testOneSec {
				isKeyPressed(ebiten.KeyS)
				err := game.Update()
				require.NoError(t, err)
			}

			require.Equal(t, testOneSec, game.tick)
			gotX, gotY := game.units.wes.Position()

			require.InDelta(t, startingX, gotX, floatTolerance, "X position")
			require.InDelta(t, float64(screenHeight-wesBorderMargin), gotY, floatTolerance, "Y clamped at the bottom border")
		})
	})

	t.Run("diagonal cases", func(t *testing.T) {
		t.Run("wes press W+D moves up and right", func(t *testing.T) {
			isKeyPressed = isKeyPressedTest(t, ebiten.KeyW, ebiten.KeyD)
			t.Cleanup(func() {
				isKeyPressed = isKeyPressedEmptyState(t)
			})

			game := NewGame(slog.Default())

			startingX, startingY := game.units.wes.Position()

			for range testOneSec {
				isKeyPressed(ebiten.KeyW)
				isKeyPressed(ebiten.KeyD)
				err := game.Update()
				require.NoError(t, err)
			}

			require.Equal(t, testOneSec, game.tick)
			gotX, gotY := game.units.wes.Position()

			require.InDelta(t, startingX+wesSpeedRation*testOneSec, gotX, floatTolerance, "X position")
			require.InDelta(t, startingY-wesSpeedRation*testOneSec, gotY, floatTolerance, "Y position")
		})

		t.Run("wes press W+A moves up and left", func(t *testing.T) {
			isKeyPressed = isKeyPressedTest(t, ebiten.KeyW, ebiten.KeyA)
			t.Cleanup(func() {
				isKeyPressed = isKeyPressedEmptyState(t)
			})

			game := NewGame(slog.Default())

			startingX, startingY := game.units.wes.Position()

			for range testOneSec {
				isKeyPressed(ebiten.KeyW)
				isKeyPressed(ebiten.KeyA)
				err := game.Update()
				require.NoError(t, err)
			}

			require.Equal(t, testOneSec, game.tick)
			gotX, gotY := game.units.wes.Position()

			require.InDelta(t, startingX-wesSpeedRation*testOneSec, gotX, floatTolerance, "X position")
			require.InDelta(t, startingY-wesSpeedRation*testOneSec, gotY, floatTolerance, "Y position")
		})

		t.Run("wes press S+D moves down and right", func(t *testing.T) {
			isKeyPressed = isKeyPressedTest(t, ebiten.KeyS, ebiten.KeyD)
			t.Cleanup(func() {
				isKeyPressed = isKeyPressedEmptyState(t)
			})

			game := NewGame(slog.Default())

			startingX, startingY := game.units.wes.Position()

			for range testOneSec {
				isKeyPressed(ebiten.KeyS)
				isKeyPressed(ebiten.KeyD)
				err := game.Update()
				require.NoError(t, err)
			}

			require.Equal(t, testOneSec, game.tick)
			gotX, gotY := game.units.wes.Position()

			require.InDelta(t, startingX+wesSpeedRation*testOneSec, gotX, floatTolerance, "X position")
			require.InDelta(t, startingY+wesSpeedRation*testOneSec, gotY, floatTolerance, "Y position")
		})

		t.Run("wes press S+A moves down and left", func(t *testing.T) {
			isKeyPressed = isKeyPressedTest(t, ebiten.KeyS, ebiten.KeyA)
			t.Cleanup(func() {
				isKeyPressed = isKeyPressedEmptyState(t)
			})

			game := NewGame(slog.Default())

			startingX, startingY := game.units.wes.Position()

			for range testOneSec {
				isKeyPressed(ebiten.KeyS)
				isKeyPressed(ebiten.KeyA)
				err := game.Update()
				require.NoError(t, err)
			}

			require.Equal(t, testOneSec, game.tick)
			gotX, gotY := game.units.wes.Position()

			require.InDelta(t, startingX-wesSpeedRation*testOneSec, gotX, floatTolerance, "X position")
			require.InDelta(t, startingY+wesSpeedRation*testOneSec, gotY, floatTolerance, "Y position")
		})
	})
}

func TestWesAttack(t *testing.T) {
	// The attack hit box is a triangle that opens to wes's RIGHT: its apex is on
	// wes and its far edge is vertical. We surround wes with a 9x5 grid of
	// jellyfish and assert that only the ones caught by that triangle are eaten;
	// everything behind wes (to his left) survives.
	//
	//   before attack (w = wes, j = jelly)      after attack (x = eaten)
	//
	//     j j j j j j j j j                        j j j j j j x x x
	//     j j j j j j j j j                        j j j j j x x x x
	//     j j j w j j j j j                        j j j w x x x x x
	//     j j j j j j j j j                        j j j j j x x x x
	//     j j j j j j j j j                        j j j j j j x x x
	t.Run("wes attack has a triangle hit box", func(t *testing.T) {
		// space is "just pressed" for the single tick we run, so Update fires
		// exactly one attack.
		isKeyJustPressed = isKeyJustPressedTest(t, ebiten.KeySpace)
		t.Cleanup(func() {
			isKeyJustPressed = isKeyJustPressedEmptyState(t)
			units = make(map[int]Unit, 0)
		})

		game := NewGame(slog.Default())

		wesX, wesY := game.units.wes.Position()

		// Steps between grid cells. The horizontal:vertical ratio (10:15) is what
		// makes the eaten region narrow by one column per row away from wes.
		const stepX, stepY = 10.0, 15.0

		units = make(map[int]Unit, 0)
		units[game.units.wes.ID()] = game.units.wes
		var smack []*JellyFish
		id := 0
		for gx := -3; gx <= 5; gx++ {
			for gy := -2; gy <= 2; gy++ {
				if gx == 0 && gy == 0 {
					continue // wes sits at the center cell
				}

				jelly := newJellyFish(id, slog.Default(), game.evenetManager)
				jelly.position = NewVector(wesX+float64(gx)*stepX, wesY+float64(gy)*stepY)
				jelly.nextPosition = jelly.position

				units[jelly.ID()] = jelly
				smack = append(smack, jelly)
				id++
			}
		}
		game.units.smack = smack
		rebuildGrid()

		require.Len(t, game.units.smack, 44, "grid must fully surround wes")

		// Advance a single tick with space pressed: Update runs wes.Attack and
		// then Die on everything it caught, all within this one tick.
		err := game.Update()
		require.NoError(t, err)

		eaten := game.units.unitsEaten

		require.Equal(t, unitStateAttack, game.units.wes.State(), "attacking puts wes into the attack state")

		gotEaten := map[[2]int]bool{}
		for _, u := range eaten {
			assert.Equal(t, unitStateDead, u.State(), "an eaten jellyfish must be in the dead state")

			j, ok := u.(*JellyFish)
			require.True(t, ok, "an eaten unit must be a jellyfish")
			assert.Equal(t, game.tick, j.tickWhenDied, "death must record the tick it happened on")

			ux, uy := u.Position()
			assert.Greater(t, ux, wesX, "no jellyfish behind wes may be eaten")

			// The tick's flee/flock step nudges each jelly by <=1.5px before
			// Attack reads the grid, so round back to the original grid cell.
			// Identities are unaffected: Attack reads the grid built before the tick.
			gx := int(math.Round((ux - wesX) / stepX))
			gy := int(math.Round((uy - wesY) / stepY))
			gotEaten[[2]int{gx, gy}] = true
		}

		// The triangle in front of wes: apex column reaches wes, far column (gx=5)
		// is the vertical edge, and the diagonals trim one column per row.
		wantEaten := map[[2]int]bool{
			{1, 0}:  true,
			{2, -1}: true, {2, 0}: true, {2, 1}: true,
			{3, -2}: true, {3, -1}: true, {3, 0}: true, {3, 1}: true, {3, 2}: true,
			{4, -2}: true, {4, -1}: true, {4, 0}: true, {4, 1}: true, {4, 2}: true,
			{5, -2}: true, {5, -1}: true, {5, 0}: true, {5, 1}: true, {5, 2}: true,
		}

		assert.Len(t, eaten, 19, "attack eats exactly the triangle in front of wes")
		assert.Equal(t, wantEaten, gotEaten, "eaten jellyfish form wes's triangular hit box")
	})

	t.Run("wes cannot attack again before the current attack finishes", func(t *testing.T) {
		game := NewGame(slog.Default())

		game.units.wes.Attack(game.tick)
		require.Equal(t, unitStateAttack, game.units.wes.State(), "the first attack must start an attack")

		secondAttack := game.units.wes.Attack(game.tick)

		assert.Nil(t, secondAttack, "wes cannot attack again while the current attack is unfinished")

		_, _, tickCountPerPose, frameCount := game.units.wes.Draw()
		attackAnimationTicks := tickCountPerPose * frameCount

		for tick := 1; tick <= attackAnimationTicks; tick++ {
			err := game.Update()
			require.NoError(t, err)

			if tick < attackAnimationTicks {
				require.Equal(t, unitStateAttack, game.units.wes.State(), "wes stays under attack while the animation plays")
				continue
			}

			assert.Equal(t, unitStateWalk, game.units.wes.State(), "wes returns to walk on the last animation tick")
		}
	})

	t.Run("jellyfish eaten count toward wes total", func(t *testing.T) {
		t.Cleanup(func() {
			units = make(map[int]Unit, 0)
		})

		game := NewGame(slog.Default())
		wesX, wesY := game.units.wes.Position()

		// Only wes plus three jellyfish lined up right in front of him, all inside
		// the attack triangle so every one gets eaten.
		units = make(map[int]Unit, 0)
		units[game.units.wes.ID()] = game.units.wes

		const eaten = 3
		for i := range eaten {
			jelly := newJellyFish(i, slog.Default(), game.evenetManager)
			jelly.position = NewVector(wesX+float64((i+1)*10), wesY)
			units[jelly.ID()] = jelly
		}
		rebuildGrid()

		require.Zero(t, game.counter.unitsEaten, "the counter must start empty")

		game.units.wes.Attack(game.tick)

		assert.Equal(t, eaten, game.counter.unitsEaten, "every eaten jellyfish counts toward wes total")
	})
}
