package main

import (
	"slices"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

var isKeyPressedTest = func(t *testing.T, pressedKeys ...ebiten.Key) func(key ebiten.Key) bool {
	t.Helper()
	return func(key ebiten.Key) bool {
		return slices.Contains(pressedKeys, key)
	}
}

var isKeyPressedEmptyState = func(t *testing.T) func(key ebiten.Key) bool {
	t.Helper()
	return func(key ebiten.Key) bool {
		return false
	}
}

var isKeyJustPressedTest = func(t *testing.T, pressedKeys ...ebiten.Key) func(key ebiten.Key) bool {
	t.Helper()
	return func(key ebiten.Key) bool {
		return slices.Contains(pressedKeys, key)
	}
}

var isKeyJustPressedEmptyState = func(t *testing.T) func(key ebiten.Key) bool {
	t.Helper()
	return func(key ebiten.Key) bool {
		return false
	}
}

const testOneSec = 60
