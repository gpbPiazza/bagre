package main

import (
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
	t.Skip("I'm refactoring the architecture")

	t.Run("out of screen above recovers (panic log frame)", func(t *testing.T) {
		require.Empty(t, units, "no other unit may influence the acceleration")

		j := &JellyFish{position: NewVector(771.77, -2.18), velocity: NewVector(0.5, -1.4), state: unitStateWalk}

		j.nextMove()

		assert.Positive(t, j.nextVelocity.y, "bounce must turn the jelly back down into the screen")
		assert.GreaterOrEqual(t, j.nextPosition.x, 0.0)
		assert.LessOrEqual(t, j.nextPosition.x, float64(screenWidth))
		assert.GreaterOrEqual(t, j.nextPosition.y, 0.0)
		assert.LessOrEqual(t, j.nextPosition.y, float64(screenHeight))
	})

	t.Run("out of screen on the left recovers", func(t *testing.T) {
		require.Empty(t, units, "no other unit may influence the acceleration")

		j := &JellyFish{position: NewVector(-1, 300), velocity: NewVector(-1.5, 0), state: unitStateWalk}

		j.nextMove()

		assert.Positive(t, j.nextVelocity.x, "bounce must turn the jelly back right into the screen")
		assert.GreaterOrEqual(t, j.nextPosition.x, 0.0)
		assert.LessOrEqual(t, j.nextPosition.x, float64(screenWidth))
		assert.GreaterOrEqual(t, j.nextPosition.y, 0.0)
		assert.LessOrEqual(t, j.nextPosition.y, float64(screenHeight))
	})

	t.Run("out of screen beyond bottom-right corner recovers", func(t *testing.T) {
		require.Empty(t, units, "no other unit may influence the acceleration")

		j := &JellyFish{position: NewVector(screenWidth+1, screenHeight+1), velocity: NewVector(1, 1), state: unitStateWalk}

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

		j := &JellyFish{position: NewVector(0.3, 0.3), velocity: NewVector(-1, -1), state: unitStateWalk}

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
	t.Skip("TODO implement test")

	t.Run("jellyfish group up, match direction, and keep distance from close neighbors", func(_ *testing.T) {

	})

	t.Run("jellyfish never stop moving", func(_ *testing.T) {

	})

	t.Run("jellyfish never move faster than their cruising speed", func(_ *testing.T) {

	})
}

func TestSmackFleeWes(t *testing.T) {
	t.Skip("TODO implement test")

	t.Run("jellyfish flee from wes when he gets close", func(_ *testing.T) {

	})
}

func TestSmackAttack(t *testing.T) {
	t.Skip("TODO implement test")

	t.Run("jellyfish can enter an attack state", func(_ *testing.T) {

	})

	t.Run("touching wes while attacking damages him", func(_ *testing.T) {

	})

	t.Run("attack state reverts to normal after a while", func(_ *testing.T) {

	})

	t.Run("no more than a limited number of jellyfish attack at the same time", func(_ *testing.T) {

	})
}

func TestSmackDeath(t *testing.T) {
	t.Skip("TODO implement test")

	t.Run("eaten jellyfish play a death animation before disappearing", func(_ *testing.T) {

	})
}
