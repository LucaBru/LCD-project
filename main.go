package main

import (
	custype "sudoku-solver/type"
)

func main() {
	matrix := [][]bool{{true, false, false}, {true, false, true}, {true, false, false}}
	nodeMatrix := custype.NewDancingLink(matrix)

}
