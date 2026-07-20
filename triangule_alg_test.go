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
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1},
		{1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1},
		{1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1},
		{1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1},
		{1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1},
		{1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	}

	x0 := 4
	y0 := 6
	n := 4

	xa := x0 + n
	// ya := y0 - n

	// xb := x0 + n
	// yb := ya + n

	expectedZeroCount := 24
	var zeroGot []int

	// Algoritmo que eu pensei em fazer
	// Como eu espero sempre produzir um triangulo equilatero, que possui lados eguais eu posso assumir uma condição dele
	// em uma matrix, o triangulo sempre será a representado como o ponto -1, na próxima iteração.
	// 1. começa no ponto x0 and y0
	// 1.1 itera até o limite te xa ou xb(são iguais)
	// note -> aqui você terá pego todos os pontos do meio do triangulo.
	// 2. Iremo iterar na parte superior do triangulo
	// 2.1 o "range de x" muda de x0 até ax, vai ser x0-1 até ax
	// 2.2 o y vaiser y0-1
	// 2.3 o y fica constante, itera até o final do range de x
	// 2.4 a gente repete este processo até chegar no ponto xa e ya
	// note -> aqui nós teremos varrido todos aos pontos do meio para cima.
	//
	// 3. Nós faremos a mesma operação da seção 2 só que para parte de baixo do triangulo.
	// Para isso, basta invez de diminuirmos o y, nós iremos aumentalo, iterando no mesmo range the x

	for i := x0; i < xa; i++ {
		// for k := y0; k < y0; k++ {
		val := grid[i][y0]
		fmt.Printf("visited x:%d y:%d - val:%d\n", i, y0, val)

		zeroGot = append(zeroGot, val)
		// }
	}

	if len(zeroGot) != expectedZeroCount {
		t.Errorf("expected to found %d got - %d", expectedZeroCount, len(zeroGot))
	}
}
