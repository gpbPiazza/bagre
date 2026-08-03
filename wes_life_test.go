package main

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWesLife(t *testing.T) {
	t.Run("wes loses a life when touched by an attacking jellyfish", func(t *testing.T) {
		t.Cleanup(func() {
			units = make(map[int]Unit, 0)
		})

		game := NewGame(slog.Default())
		wesX, wesY := game.units.wes.Position()

		// Only wes plus one jellyfish, close enough that the jellyfish attack
		// box touches him.
		const touchOffset = 5.0
		units = make(map[int]Unit, 0)
		units[game.units.wes.ID()] = game.units.wes

		jelly := newJellyFish(0, slog.Default(), game.evenetManager)
		jelly.position = NewVector(wesX+touchOffset, wesY)
		jelly.nextPosition = jelly.position
		units[jelly.ID()] = jelly
		game.units.smack = []*JellyFish{jelly}
		rebuildGrid()

		// Damage only flows while tickStartAttack is nonzero, so the attack
		// starts on the tick the next Update will run.
		jelly.Attack(game.tick + 1)

		require.Equal(t, wesStartingLife, game.units.wes.life, "wes must start at full life")
		require.Equal(t, unitStateWalk, game.units.wes.State(), "wes must start walking")
		require.Equal(t, unitStateAttack, jelly.State(), "the jellyfish must be attacking")
		require.Less(t, jelly.position.Distance(game.units.wes.VecPosition()), float64(jellyAttackViewRadius),
			"the jellyfish must start within touch range of wes")

		err := game.Update()
		require.NoError(t, err)

		assert.Equal(t, wesStartingLife-1, game.units.wes.life, "touching an attacking jellyfish costs one life")
		assert.Equal(t, unitStateHurt, game.units.wes.State(), "taking damage puts wes into the hurt state")
	})

	t.Run("wes cannot lose life while in hurt state", func(t *testing.T) {
		t.Cleanup(func() {
			units = make(map[int]Unit, 0)
		})

		game := NewGame(slog.Default())
		wesX, wesY := game.units.wes.Position()

		const touchOffset = 5.0
		units = make(map[int]Unit, 0)
		units[game.units.wes.ID()] = game.units.wes

		jelly := newJellyFish(0, slog.Default(), game.evenetManager)
		jelly.position = NewVector(wesX+touchOffset, wesY)
		jelly.nextPosition = jelly.position
		units[jelly.ID()] = jelly
		game.units.smack = []*JellyFish{jelly}
		rebuildGrid()

		jelly.Attack(game.tick + 1)

		// First touch puts wes into the hurt state.
		err := game.Update()
		require.NoError(t, err)
		require.Equal(t, wesStartingLife-1, game.units.wes.life, "the first touch must cost one life")
		require.Equal(t, unitStateHurt, game.units.wes.State(), "the first touch must hurt wes")

		// Keep the jellyfish glued to wes for the rest of the hurt window: it
		// publishes damage every tick, but none of it may land.
		for range wesHurtDurationTicks - 1 {
			jelly.position = NewVector(wesX+touchOffset, wesY)
			jelly.nextPosition = jelly.position

			err = game.Update()
			require.NoError(t, err)

			require.Equal(t, unitStateAttack, jelly.State(), "the jellyfish must keep attacking the whole window")
			require.Equal(t, unitStateHurt, game.units.wes.State(), "wes stays hurt for the whole hurt window")
			assert.Equal(t, wesStartingLife-1, game.units.wes.life, "wes cannot lose life while hurt")
		}
	})

	t.Run("after hurtState duration passes wes can take damage again", func(t *testing.T) {
		t.Cleanup(func() {
			units = make(map[int]Unit, 0)
		})

		game := NewGame(slog.Default())
		wesX, wesY := game.units.wes.Position()

		const touchOffset = 5.0
		units = make(map[int]Unit, 0)
		units[game.units.wes.ID()] = game.units.wes

		jelly := newJellyFish(0, slog.Default(), game.evenetManager)
		jelly.position = NewVector(wesX+touchOffset, wesY)
		jelly.nextPosition = jelly.position
		units[jelly.ID()] = jelly
		game.units.smack = []*JellyFish{jelly}
		rebuildGrid()

		jelly.Attack(game.tick + 1)

		// First touch puts wes into the hurt state.
		err := game.Update()
		require.NoError(t, err)
		require.Equal(t, wesStartingLife-1, game.units.wes.life, "the first touch must cost one life")
		require.Equal(t, unitStateHurt, game.units.wes.State(), "the first touch must hurt wes")

		// Ride out the entire hurt window with the jellyfish still in contact.
		for range wesHurtDurationTicks {
			jelly.position = NewVector(wesX+touchOffset, wesY)
			jelly.nextPosition = jelly.position

			err = game.Update()
			require.NoError(t, err)
		}

		require.Equal(t, unitStateWalk, game.units.wes.State(), "wes must leave the hurt state once it expires")
		require.Equal(t, wesStartingLife-1, game.units.wes.life, "riding out the hurt window costs no extra life")
		require.Equal(t, unitStateAttack, jelly.State(), "the jellyfish must still be attacking")

		// With the hurt window over, the very next touch lands again.
		jelly.position = NewVector(wesX+touchOffset, wesY)
		jelly.nextPosition = jelly.position

		err = game.Update()
		require.NoError(t, err)

		assert.Equal(t, wesStartingLife-2, game.units.wes.life, "wes takes damage again after the hurt window ends")
		assert.Equal(t, unitStateHurt, game.units.wes.State(), "the new touch puts wes back into the hurt state")
	})
}
