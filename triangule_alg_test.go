package main

import (
	"fmt"
	"testing"
)

func Test_run_triangle(t *testing.T) {

	// ┼─────→ +X
	// │   (right)
	// │
	// │
	// ↓
	//  +Y (down)

	grid := [11][11]int{
		//     =0  1  2  3  4  5  6  7  8  9  10
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, // 0
		{1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1}, // 1
		{1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1}, // 2
		{1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1}, // 3
		{1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1}, // 4
		{1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1}, // 5
		{1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1}, // 6
		{1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1}, // 7
		{1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1}, // 8
		{1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1}, // 9
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, // 10
	}

	x0 := 3
	y0 := 5 //bin/→  the tip (leftmost 0 in the whole grid)
	n := 4

	xa := x0 + n
	ya := y0 - n

	x_a := 7
	y_a := 1 //→  top-right corner (topmost 0)

	assertPoint := func(x, y, x_, y_ int) {
		if x != x_ {
			t.Error("x it's not with the expected value")
		}
		if y != y_ {
			t.Error("y it's not with the expected value")
		}
	}

	assertPoint(xa, ya, x_a, y_a)

	x_b := 7
	y_b := 9 //→  bottom-right corner (bottommost 0)

	xb := x0 + n
	yb := y0 + n

	assertPoint(xb, yb, x_b, y_b)

	expectedZeroCount := 25
	var zeroGot []int

	fmt.Print("getting the upper part\n")
	newX0 := x0
	for k := y0; k >= ya; k-- {
		for i := newX0; i <= xa; i++ {
			val := grid[k][i]

			fmt.Printf("visited x:%d y:%d - val:%d\n", i, k, val)

			zeroGot = append(zeroGot, val)

		}
		fmt.Print("===/===\n")

		newX0++
	}

	fmt.Print("getting the lower part\n")
	newY0 := y0 + 1
	newX02 := x0 + 1
	for k := newY0; k <= yb; k++ {
		for i := newX02; i <= xa; i++ {
			val := grid[k][i]

			fmt.Printf("visited x:%d y:%d - val:%d\n", i, k, val)

			zeroGot = append(zeroGot, val)

		}
		fmt.Print("===/===\n")

		newX02++
	}

	if len(zeroGot) != expectedZeroCount {
		t.Errorf("expected to found %d got - %d", expectedZeroCount, len(zeroGot))
	}
}
