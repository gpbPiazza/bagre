package main

import (
	"log/slog"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// rebuildGrid indexes unitsPositions by position, so nextMove must never
// produce a nextPosition outside [0, screenWidth] x [0, screenHeight].
// Regression for: panic index out of range [-2] in rebuildGrid — a jelly at
// y=-2.18 got a bounce (1/pos with pos<0) pushing it further off-screen
// every frame instead of back in.
func TestNextMove(t *testing.T) {
	t.Run("out of screen above recovers (panic log frame)", func(t *testing.T) {
		require.Empty(t, units, "no other unit may influence the acceleration")

		// Mid-screen x (clear of the side borders) and a y pushed past the top
		// edge, so only the top bounce acts.
		j := &JellyFish{
			position: NewVector(screenWidth/2, -jellyViewRadius),
			velocity: NewVector(0.5, -1.4),
			state:    unitStateWalk,
		}

		j.nextMove()

		assert.Positive(t, j.nextVelocity.y, "bounce must turn the jelly back down into the screen")
		assert.GreaterOrEqual(t, j.nextPosition.x, 0.0)
		assert.LessOrEqual(t, j.nextPosition.x, float64(screenWidth))
		assert.GreaterOrEqual(t, j.nextPosition.y, 0.0)
		assert.LessOrEqual(t, j.nextPosition.y, float64(screenHeight))
	})

	t.Run("out of screen on the left recovers", func(t *testing.T) {
		require.Empty(t, units, "no other unit may influence the acceleration")

		// x pushed past the left edge, mid-screen y (clear of top/bottom borders).
		j := &JellyFish{position: NewVector(-jellyViewRadius, screenHeight/2), velocity: NewVector(-1.5, 0), state: unitStateWalk}

		j.nextMove()

		assert.Positive(t, j.nextVelocity.x, "bounce must turn the jelly back right into the screen")
		assert.GreaterOrEqual(t, j.nextPosition.x, 0.0)
		assert.LessOrEqual(t, j.nextPosition.x, float64(screenWidth))
		assert.GreaterOrEqual(t, j.nextPosition.y, 0.0)
		assert.LessOrEqual(t, j.nextPosition.y, float64(screenHeight))
	})

	t.Run("out of screen beyond bottom-right corner recovers", func(t *testing.T) {
		require.Empty(t, units, "no other unit may influence the acceleration")

		// Pushed past both the right and bottom edges.
		j := &JellyFish{
			position: NewVector(screenWidth+jellyViewRadius, screenHeight+jellyViewRadius),
			velocity: NewVector(1, 1),
			state:    unitStateWalk,
		}

		j.nextMove()

		assert.Negative(t, j.nextVelocity.x, "bounce must turn the jelly back left into the screen")
		assert.Negative(t, j.nextVelocity.y, "bounce must turn the jelly back up into the screen")
		assert.GreaterOrEqual(t, j.nextPosition.x, 0.0)
		assert.LessOrEqual(t, j.nextPosition.x, float64(screenWidth))
		assert.GreaterOrEqual(t, j.nextPosition.y, 0.0)
		assert.LessOrEqual(t, j.nextPosition.y, float64(screenHeight))
	})

	t.Run("exactly at top-left corner stays finite and on screen", func(t *testing.T) {
		require.Empty(t, units, "no other unit may influence the acceleration")

		// A raw 1/pos bounce at pos=0 is +Inf, which turns the position into
		// NaN after ScaleToLength.
		j := &JellyFish{position: NewVector(0, 0), velocity: NewVector(-1, -1), state: unitStateWalk}

		j.nextMove()

		assert.False(t, math.IsNaN(j.nextPosition.x), "position must stay a real number")
		assert.False(t, math.IsNaN(j.nextPosition.y), "position must stay a real number")
		assert.GreaterOrEqual(t, j.nextPosition.x, 0.0)
		assert.LessOrEqual(t, j.nextPosition.x, float64(screenWidth))
		assert.GreaterOrEqual(t, j.nextPosition.y, 0.0)
		assert.LessOrEqual(t, j.nextPosition.y, float64(screenHeight))
	})

	t.Run("exactly at bottom-right corner stays finite and on screen", func(t *testing.T) {
		require.Empty(t, units, "no other unit may influence the acceleration")

		j := &JellyFish{position: NewVector(screenWidth, screenHeight), velocity: NewVector(1, 1), state: unitStateWalk}

		j.nextMove()

		assert.False(t, math.IsNaN(j.nextPosition.x), "position must stay a real number")
		assert.False(t, math.IsNaN(j.nextPosition.y), "position must stay a real number")
		assert.GreaterOrEqual(t, j.nextPosition.x, 0.0)
		assert.LessOrEqual(t, j.nextPosition.x, float64(screenWidth))
		assert.GreaterOrEqual(t, j.nextPosition.y, 0.0)
		assert.LessOrEqual(t, j.nextPosition.y, float64(screenHeight))
	})

	t.Run("near top-left border moving outward is pushed back in", func(t *testing.T) {
		require.Empty(t, units, "no other unit may influence the acceleration")

		// On-screen but inside the push-cap zone (0 < pos < borderPushCap), so the
		// bounce clamps the divisor instead of blowing up.
		j := &JellyFish{position: NewVector(borderPushCap/2, borderPushCap/2), velocity: NewVector(-1, -1), state: unitStateWalk}

		j.nextMove()

		assert.Positive(t, j.nextVelocity.x, "bounce must turn the jelly back right into the screen")
		assert.Positive(t, j.nextVelocity.y, "bounce must turn the jelly back down into the screen")
		assert.GreaterOrEqual(t, j.nextPosition.x, 0.0)
		assert.LessOrEqual(t, j.nextPosition.x, float64(screenWidth))
		assert.GreaterOrEqual(t, j.nextPosition.y, 0.0)
		assert.LessOrEqual(t, j.nextPosition.y, float64(screenHeight))
	})
}

func TestSmackFlocking(t *testing.T) {
	t.Run("jellyfish group up, match direction, and keep distance from close neighbors", func(t *testing.T) {
		t.Skip("yet I didn't have a regression here and make a good test to assert this behavior it's" +
			"always resulting into something, hard to maintain and not necessaraly is testing what I really want," +
			" we need more time thinking on this one")

		game := NewGame(slog.Default())

		const windows = 4
		twoSeconds := 2 * testOneSec

		for range windows {
			for range twoSeconds {
				err := game.Update()
				require.NoError(t, err)
			}
		}
	})

	t.Run("jellyfish never stop moving", func(t *testing.T) {
		game := NewGame(slog.Default())

		for range testOneSec {
			err := game.Update()
			require.NoError(t, err)

			for _, j := range game.units.smack {
				speed := j.VecVelocity().Pythagoras()
				assert.GreaterOrEqual(t, speed, jellyMinSpeed-floatTolerance,
					"every jellyfish must keep cruising at or above the minimum speed")
			}
		}
	})

	t.Run("jellyfish never move faster than their cruising speed", func(t *testing.T) {
		game := NewGame(slog.Default())

		for range testOneSec {
			err := game.Update()
			require.NoError(t, err)

			for _, j := range game.units.smack {
				speed := j.VecVelocity().Pythagoras()
				assert.LessOrEqual(t, speed, jellyMaxSpeed+floatTolerance,
					"no jellyfish may move faster than the maximum cruising speed")
			}
		}
	})
}

func TestSmackAttack(t *testing.T) {
	// countAttacking reports how many of the school are currently attacking.
	countAttacking := func(game *Game) int {
		attacking := 0
		for _, j := range game.units.smack {
			if j.State() == unitStateAttack {
				attacking++
			}
		}
		return attacking
	}

	t.Run("jellyfish can enter an attack state", func(t *testing.T) {
		game := NewGame(slog.Default())

		require.Zero(t, countAttacking(game), "no jellyfish attacks before the electric stage")

		// The electric stage only records its start on a nonzero tick, so advance
		// one tick before crossing the "eaten enough" threshold that starts it.
		require.NoError(t, game.Update())
		game.counter.Add(startEletricJellyAt + 1)
		require.NotZero(t, game.tickEletricJellyStart, "eating past the threshold must start the electric stage")

		// The stage fires its first wave one interval after it starts.
		for range eletricAttackIntervalTicks {
			require.NoError(t, game.Update())
		}

		assert.Positive(t, countAttacking(game), "the electric stage must put jellyfish into the attack state")
	})

	t.Run("touching wes while attacking damages him", func(t *testing.T) {
		t.Cleanup(func() { units = make(map[int]Unit, 0) })

		game := NewGame(slog.Default())
		wesX, wesY := game.units.wes.Position()

		// Only wes and one jellyfish, placed inside the jelly's attack radius.
		units = make(map[int]Unit, 0)
		units[game.units.wes.ID()] = game.units.wes

		jelly := newJellyFish(0, slog.Default(), game.evenetManager)
		const touchOffset = 5.0
		jelly.position = NewVector(wesX+touchOffset, wesY)
		jelly.nextPosition = jelly.position
		units[jelly.ID()] = jelly
		game.units.smack = []*JellyFish{jelly}
		rebuildGrid()

		// Damage only flows while tickStartAttack is nonzero, so start the attack
		// on the tick the next Update will process.
		jelly.Attack(game.tick + 1)

		require.Equal(t, wesStartingLife, game.units.wes.life, "wes must start at full life")
		require.Equal(t, unitStateAttack, jelly.State(), "the jelly must be attacking")
		require.Less(t, jelly.VecPosition().Distance(game.units.wes.VecPosition()), float64(jellyAttackViewRadius),
			"the jelly must start within touch range of wes")

		require.NoError(t, game.Update())

		assert.Equal(t, wesStartingLife-1, game.units.wes.life, "touching an attacking jellyfish costs wes a life")
		assert.Equal(t, unitStateHurt, game.units.wes.State(), "taking damage puts wes into the hurt state")
	})

	t.Run("attack state reverts to normal after a while", func(t *testing.T) {
		t.Cleanup(func() { units = make(map[int]Unit, 0) })

		game := NewGame(slog.Default())

		// A single jelly off on its own, so nothing but the passage of time acts
		// on its attack state.
		units = make(map[int]Unit, 0)
		jelly := newJellyFish(0, slog.Default(), game.evenetManager)
		jelly.position = NewVector(screenWidth/2, screenHeight/2)
		jelly.nextPosition = jelly.position
		units[jelly.ID()] = jelly
		game.units.smack = []*JellyFish{jelly}
		rebuildGrid()

		attackStartTick := game.tick + 1
		jelly.Attack(attackStartTick)
		require.Equal(t, unitStateAttack, jelly.State(), "the jelly must start in the attack state")

		revertTick := attackStartTick + jellyAttackDurationTicks
		for game.tick < revertTick {
			require.NoError(t, game.Update())

			if game.tick < revertTick {
				require.Equal(t, unitStateAttack, jelly.State(), "the jelly stays attacking until the duration elapses")
				continue
			}
			assert.Equal(t, unitStateWalk, jelly.State(), "the jelly returns to walk once the attack duration passes")
		}
	})

	t.Run("no more than a limited number of jellyfish attack at the same time", func(t *testing.T) {
		game := NewGame(slog.Default())

		require.NoError(t, game.Update())
		game.counter.Add(startEletricJellyAt + 1)
		require.NotZero(t, game.tickEletricJellyStart, "the electric stage must have started")

		// Drive through a couple of attack waves.
		maxSeen := 0
		for range 2 * eletricAttackIntervalTicks {
			require.NoError(t, game.Update())

			attacking := countAttacking(game)
			require.LessOrEqual(t, attacking, maxAttackingJellys,
				"never more than the cap may attack at the same time")
			if attacking > maxSeen {
				maxSeen = attacking
			}
		}

		require.Positive(t, maxSeen, "at least one wave must have fired, otherwise the cap is untested")
	})
}

func TestSmackDeath(t *testing.T) {
	t.Skip("TODO implement test")

	t.Run("eaten jellyfish stop walking and play a death animation before disappearing", func(_ *testing.T) {

	})
}
