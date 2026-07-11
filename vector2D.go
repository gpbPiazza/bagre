package main

import "math"

type Vector2D struct {
	x float64
	y float64
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

func (v1 Vector2D) SubtractVal(val float64) Vector2D {
	return Vector2D{
		x: v1.x - val,
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
