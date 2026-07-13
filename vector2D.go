package main

import "math"

type Vector2D struct {
	x float64
	y float64
}

func NewVector(x, y float64) Vector2D {
	return Vector2D{x, y}
}

func (v1 Vector2D) Add(v2 Vector2D) Vector2D {
	return Vector2D{
		x: v1.x + v2.x,
		y: v1.y + v2.y,
	}
}

func (v1 Vector2D) Subtract(v2 Vector2D) Vector2D {
	return Vector2D{
		x: v1.x - v2.x,
		y: v1.y - v2.y,
	}
}

func (v1 Vector2D) Multiply(v2 Vector2D) Vector2D {
	return Vector2D{
		x: v1.x * v2.x,
		y: v1.y * v2.y,
	}
}

func (v1 Vector2D) AddVal(val float64) Vector2D {
	return Vector2D{
		x: v1.x + val,
		y: v1.y + val,
	}
}

func (v1 Vector2D) SubVal(val float64) Vector2D {
	return Vector2D{
		x: v1.x - val,
		y: v1.y - val,
	}
}

func (v1 Vector2D) AddX(val float64) Vector2D {
	return Vector2D{
		x: v1.x + val,
		y: v1.y,
	}
}

func (v1 Vector2D) AddY(val float64) Vector2D {
	return Vector2D{
		x: v1.x,
		y: v1.y + val,
	}
}

func (v1 Vector2D) SubX(val float64) Vector2D {
	return Vector2D{
		x: v1.x - val,
		y: v1.y,
	}
}

func (v1 Vector2D) SubY(val float64) Vector2D {
	return Vector2D{
		x: v1.x,
		y: v1.y - val,
	}
}

func (v1 Vector2D) MultiplyVal(val float64) Vector2D {
	return Vector2D{
		x: v1.x * val,
		y: v1.y * val,
	}
}

func (v1 Vector2D) DivisionVal(val float64) Vector2D {
	return Vector2D{
		x: v1.x / val,
		y: v1.y / val,
	}
}

func (v1 Vector2D) LimitVal(lower, uppper float64) Vector2D {
	return Vector2D{
		x: math.Min(math.Max(v1.x, lower), uppper),
		y: math.Min(math.Max(v1.y, lower), uppper),
	}
}

func (v1 Vector2D) Distance(v2 Vector2D) float64 {
	xDif := math.Pow(v1.x-v2.x, 2)
	yDif := math.Pow(v1.y-v2.y, 2)
	return math.Sqrt(xDif + yDif)
}

func (v1 Vector2D) Pythagoras() float64 {
	return math.Sqrt(math.Pow(v1.x, 2) + math.Pow(v1.y, 2))
}

// ScaleToLength keeps the vector's DIRECTION but sets its LENGTH to target.
// It normalizes (divide by current magnitude -> unit vector) then scales up.
// A zero vector has no direction, so we leave it untouched to avoid div-by-zero.
func (v1 Vector2D) ScaleToLength(target float64) Vector2D {
	mag := v1.Pythagoras()
	if mag == 0 {
		return v1
	}
	return v1.DivisionVal(mag).MultiplyVal(target)
}
