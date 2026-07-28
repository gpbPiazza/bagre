package main

// running a real game populates the package-global units/unitsPositions;
// restore them so this test can't leak state into others (e.g. TestNextMove).
// t.Cleanup(func() {
// 	units = make(map[int]Unit, 0)
// 	unitsPositions = [screenWidth + 2][screenHeight + 2]int{}
// })
//
// g := NewGame(nil)
//
// // pass 1 second in the game
// for range 60 {
// 	err := g.Update()
// 	require.NoError(t, err)
// }
//
// TODOS
// - I need some way to simulate keys pressed under test

// test arch rules
// Every test should start a entire game
// We must change the state of the game using game.Update method or
// somulating user inputs
// smack positions or wes positions it's ok to manipulate to test other cases.

// asserts

// func TestHandleEleticJellyFishs(t *testing.T) {
// 	t.Run("every second jelly of one random chunk attacks, the rest keeps walking", func(t *testing.T) {
// 		var smack []*JellyFish
// 		for range 16 {
// 			smack = append(smack, &JellyFish{state: unitStateWalk})
// 		}
// 		g := &Game{units: Units{smack: smack}}
//
// 		for _, j := range smack {
// 			require.Equal(t, unitStateWalk, j.state, "every jelly must start walking")
// 		}
//
// 		g.handleEleticJellyFishs(nil)
//
// 		chunkLen := len(smack) / 4
// 		var attackerChunks []int
// 		walking := 0
// 		for i, j := range smack {
// 			switch j.state {
// 			case unitStateAttack:
// 				attackerChunks = append(attackerChunks, i/chunkLen)
// 			case unitStateWalk:
// 				walking++
// 			}
// 		}
//
// 		assert.Len(t, attackerChunks, chunkLen/2, "half of one chunk must be attacking")
// 		assert.Equal(t, len(smack)-chunkLen/2, walking, "all other jellies must keep walking")
// 		for _, c := range attackerChunks {
// 			assert.Equal(t, attackerChunks[0], c, "all attackers must come from the same chunk")
// 		}
// 	})
//
// 	t.Run("attackers are capped at 25 even in a big chunk", func(t *testing.T) {
// 		var smack []*JellyFish
// 		for range 240 {
// 			smack = append(smack, &JellyFish{state: unitStateWalk})
// 		}
// 		g := &Game{units: Units{smack: smack}}
//
// 		for _, j := range smack {
// 			require.Equal(t, unitStateWalk, j.state, "every jelly must start walking")
// 		}
//
// 		g.handleEleticJellyFishs(nil)
//
// 		attacking := 0
// 		walking := 0
// 		for _, j := range smack {
// 			switch j.state {
// 			case unitStateAttack:
// 				attacking++
// 			case unitStateWalk:
// 				walking++
// 			}
// 		}
//
// 		assert.Equal(t, 20, attacking, "a chunk of 60 has 30 candidates, but the limit caps attackers at 20")
// 		assert.Equal(t, len(smack)-20, walking, "all other jellies must keep walking")
// 	})
// }
