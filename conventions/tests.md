# Testing Conventions

## HARD RULE: tests never touch production code

Implementing a test means changing `_test.go` files **only**. Never add getters, helpers, seams, or any other production code to make a test easier to write. Tests live in `package main`.

## Start a whole game, drive it through ticks

Every test starts a real game and changes its state the way the game does — by advancing ticks and simulating input — never by hand-building internals.

- Build the game: `game := NewGame(slog.Default())`.
- Advance state by looping `game.Update()` (one call = one tick, ~1/60s). `testOneSec` = 60 ticks.

Setting up a scenario by manipulating package state directly is fine (tests live in `package main`): assign `game.units.wes.position`, reset the global `units` map, place jellyfish, call `rebuildGrid()`. Anything you mutate, restore in `t.Cleanup` (see below).

## Target the highest-level function possible

Test the outer function that encapsulates the behavior and assert on its observable output. This keeps tests stable across internal refactors — inner functions can be renamed, split, or inlined and the test still holds.

`game.Update()` is usually that top level: it drives `move`, `Attack`, `Die`, and `checkState`. Assert on what a player would observe (`wes.State()`, `wes.Position()`), never on a private helper (`borderBounce`) or an intermediate field (`nextPosition`).

## No hardcoded math — derive expected values from source

Never bake a computed literal into a test. Read tunable values from the same place production reads them, and let the test do the arithmetic:

```go
// good — tracks the real animation length
_, _, tickCountPerPose, frameCount := game.units.wes.Draw()
attackAnimationTicks := tickCountPerPose * frameCount

// bad — breaks on every balance tweak
const attackAnimationTicks = 6 * 6
```

Speeds, cooldowns, and animation lengths are balancing knobs. A test must break only when the *logic* changes, not when a balance value moves. (The example below derives the expected distance from the production constant `wesSpeedRation`, not a literal.)

## One `t.Run` per case, no table tests

- Do **not** write table tests.
- Every case is its own `t.Run` and has everything it needs: its own setup, execution, and assertions.
- Each `t.Run` must not depend on any other. If a case needs many steps, they all live inside that one `t.Run`.

## Simulate input through mockable seams; restore globals in `t.Cleanup`

External input is wrapped in swappable package vars (`isKeyPressed`, `isKeyJustPressed`) so tests can drive it. Swap with the helpers and always restore:

```go
isKeyPressed = isKeyPressedTest(t, ebiten.KeyW)
t.Cleanup(func() { isKeyPressed = isKeyPressedEmptyState(t) })
```

Only mock what the case actually uses: a test that never simulates a key press must not swap `isKeyPressed`/`isKeyJustPressed` at all — the swap and restore lines are pure noise there.

Those helpers (`isKeyPressedTest` / `isKeyPressedEmptyState`, `isKeyJustPressedTest` / `isKeyJustPressedEmptyState`) and shared constants like `testOneSec` live in `helpers_test.go`. Put anything reused across test files there, not inline in a single `_test.go`.

Rule: any package global a test mutates — input funcs, the `units` map, positions — is restored in `t.Cleanup` so cases stay independent.

## Assertions: testify

Use `github.com/stretchr/testify`:

- `require` — for **preconditions** that must hold before the test can run (global state, starting positions). A failed `require` stops the subtest immediately.
- `assert` — for the **outcomes** being tested.

## Example

```go
func TestWesMoves(t *testing.T) {
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

		require.Equal(t, game.tick, testOneSec)
		gotX, gotY := game.units.wes.Position()

		require.Equal(t, startingX, gotX, "X position")
		require.Equal(t, startingY-wesSpeedRation*testOneSec, gotY, "Y position")
	})
}
```
