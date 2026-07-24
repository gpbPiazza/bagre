package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleEleticJellyFishs(t *testing.T) {
	t.Run("every second jelly of one random chunk attacks, the rest keeps walking", func(t *testing.T) {
		var smack []*JellyFish
		for range 16 {
			smack = append(smack, &JellyFish{state: unitStateWalk})
		}
		g := &Game{units: Units{smack: smack}}

		for _, j := range smack {
			require.Equal(t, unitStateWalk, j.state, "every jelly must start walking")
		}

		g.handleEleticJellyFishs(nil)

		chunkLen := len(smack) / 4
		var attackerChunks []int
		walking := 0
		for i, j := range smack {
			switch j.state {
			case unitStateAttack:
				attackerChunks = append(attackerChunks, i/chunkLen)
			case unitStateWalk:
				walking++
			}
		}

		assert.Len(t, attackerChunks, chunkLen/2, "half of one chunk must be attacking")
		assert.Equal(t, len(smack)-chunkLen/2, walking, "all other jellies must keep walking")
		for _, c := range attackerChunks {
			assert.Equal(t, attackerChunks[0], c, "all attackers must come from the same chunk")
		}
	})

	t.Run("attackers are capped at 25 even in a big chunk", func(t *testing.T) {
		var smack []*JellyFish
		for range 240 {
			smack = append(smack, &JellyFish{state: unitStateWalk})
		}
		g := &Game{units: Units{smack: smack}}

		for _, j := range smack {
			require.Equal(t, unitStateWalk, j.state, "every jelly must start walking")
		}

		g.handleEleticJellyFishs(nil)

		attacking := 0
		walking := 0
		for _, j := range smack {
			switch j.state {
			case unitStateAttack:
				attacking++
			case unitStateWalk:
				walking++
			}
		}

		assert.Equal(t, 20, attacking, "a chunk of 60 has 30 candidates, but the limit caps attackers at 20")
		assert.Equal(t, len(smack)-20, walking, "all other jellies must keep walking")
	})
}
