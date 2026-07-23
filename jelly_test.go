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
