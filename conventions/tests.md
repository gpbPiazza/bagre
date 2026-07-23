# Testing Conventions

## Target the highest-level function possible

If a function encapsulates two or more inner functions, write the test at that top level. This keeps tests stable across internal refactors — we can restructure inner functions without adjusting tests, as long as the output stays correct.

Example: `nextMove` calls `calcAcceleration`, which calls `borderBounce`. Tests target `nextMove` and assert on the observable outcome (`nextPosition`, `nextVelocity`), never on `borderBounce` directly. If `borderBounce` is renamed, split, or inlined, the tests still hold.

## Structure: one `t.Run` per case, no table tests

- Do **not** write table tests.
- Every case is its own `t.Run` and has everything it needs to execute: its own setup, its own execution, its own assertions.
- Each `t.Run` must not depend on any other `t.Run`. If a case needs many functions to execute, all of those executions go inside that one `t.Run`.

## Assertions: testify

Use `github.com/stretchr/testify`:

- `require` — for **preconditions** that must be true before the test can run (e.g. global maps or units state). This shows the reader that what must be true in fact is true before the test executes. A failed `require` stops the subtest immediately.
- `assert` — for the **outcomes** being tested.

## Example

```go
func TestNextMove(t *testing.T) {
	t.Run("out of screen above recovers", func(t *testing.T) {
		require.Empty(t, units, "no other unit may influence the acceleration")

		j := &JellyFish{position: NewVector(771.77, -2.18), velocity: NewVector(0.5, -1.4), state: unitStateWalk}

		j.nextMove()

		assert.Positive(t, j.nextVelocity.y, "bounce must turn the jelly back down into the screen")
		assert.GreaterOrEqual(t, j.nextPosition.y, 0.0)
		assert.LessOrEqual(t, j.nextPosition.y, float64(screenHeight))
	})

	t.Run("near top-left border moving outward is pushed back in", func(t *testing.T) {
		require.Empty(t, units, "no other unit may influence the acceleration")

		j := &JellyFish{position: NewVector(0.3, 0.3), velocity: NewVector(-1, -1), state: unitStateWalk}

		j.nextMove()

		assert.Positive(t, j.nextVelocity.x)
		assert.Positive(t, j.nextVelocity.y)
	})
}
```
