package main

import (
	custype "sudoku-solver/type"
	"sync"
)

const ROWSIZE = 729
const CLMSIZE = 324

func clmNumbers(clm int, matrix [9][9]int) []int {
	res := []int{}
	for idx := range matrix {
		if n := matrix[idx][clm]; n != -1 {
			res = append(res, n)
		}
	}
	return res
}

func boxNumbers(cellRowIdx, cellClmIdx int, matrix [9][9]int) []int {
	res := make([]int, 0)
	rowIdx := cellRowIdx / 3
	clmIdx := cellClmIdx / 3
	for i := rowIdx; i < rowIdx+3; i++ {
		for j := 0; j < clmIdx+3; j++ {
			if n := matrix[i][j]; n != -1 {
				res = append(res, n)
			}
		}
	}
	return res
}

func rowNumbers(row int, matrix [9][9]int) []int {
	res := make([]int, 0)
	for i := 0; i < 9; i++ {
		if n := matrix[row][i]; n != -1 {
			res = append(res, n)
		}
	}
	return res
}

/*
one go routine per each constraint
Row-column constraint
	guarantee that a cell has exactly one value
	R1C1 = { R1C1#1, R1C1#2, R1C1#3, R1C1#4, R1C1#5, R1C1#6, R1C1#7, R1C1#8, R1C1#9 } means that cell [1][1] can have exactly one value in [1..9]

Row-number constraint
	guarantee that a row has all [1..9] numbers
	R1#1 = { R1C1#1, R1C2#1, R1C3#1, R1C4#1, R1C5#1, R1C6#1, R1C7#1, R1C8#1, R1C9#1 }

Column-number constraint
	guarantee that a column has all [1..9] numbers
	C1#1 = { R1C1#1, R2C1#1, R3C1#1, R4C1#1, R5C1#1, R6C1#1, R7C1#1, R8C1#1, R9C1#1 }.

Box-number constraint
	as rows and columns :)
	B1#1 = { R1C1#1, R1C2#1, R1C3#1, R2C1#1, R2C2#1, R2C3#1, R3C1#1, R3C2#1, R3C3#1 }.

each row is named R<value>C<value>#<value> (order is from right to left, so firstly change number, than clm*9 then row*81)
creating a go routine for each constraint means that it has to put 729 cell to true
can be split further to 9 go routines per constraints that set 81 cell to true
*/

func rowClmConstraint(matrix [ROWSIZE][CLMSIZE]bool, wg *sync.WaitGroup) {
	// R1C1 = { R1C1#1, R1C1#2, R1C1#3, R1C1#4, R1C1#5, R1C1#6, R1C1#7, R1C1#8, R1C1#9 }
	defer wg.Done()
	for i := 0; i < 9; i++ {
		rowIdx := i * 81
		for j := 0; j < 9; j++ {
			clmIdx := j * 9
			for h := 0; h < 9; h++ {
				matrix[rowIdx+clmIdx+h][i*9+j] = true
			}
		}
	}
}

func rowNumConstraint(matrix [ROWSIZE][CLMSIZE]bool, wg *sync.WaitGroup) {
	//R1#1 = { R1C1#1, R1C2#1, R1C3#1, R1C4#1, R1C5#1, R1C6#1, R1C7#1, R1C8#1, R1C9#1 }
	defer wg.Done()
	for i := 0; i < 9; i++ {
		rowIdx := i * 81
		for j := 0; j < 9; j++ {
			for h := 0; h < 9; h++ {
				clmIdx := h * 9
				matrix[rowIdx+clmIdx+j][81+j+i*9] = true
			}
		}
	}
}

func clmNumConstraint(matrix [ROWSIZE][CLMSIZE]bool, wg *sync.WaitGroup) {
	//C1#1 = { R1C1#1, R2C1#1, R3C1#1, R4C1#1, R5C1#1, R6C1#1, R7C1#1, R8C1#1, R9C1#1 }
	defer wg.Done()
	for i := 0; i < 9; i++ {
		clmIdx := i * 9
		for j := 0; j < 9; j++ {
			numIdx := j
			for h := 0; h < 9; h++ {
				rowIdx := h * 81
				matrix[rowIdx+clmIdx+numIdx][81*2+j+i*9] = true
			}
		}
	}
}

func boxNumConstraint(matrix [ROWSIZE][CLMSIZE]bool, wg *sync.WaitGroup) {
	//B1#1 = { R1C1#1, R1C2#1, R1C3#1, R2C1#1, R2C2#1, R2C3#1, R3C1#1, R3C2#1, R3C3#1 }
	defer wg.Done()
	for i := 0; i < 9; i++ {
		boxIdx := i
		startingRow := boxIdx / 3 * 3
		startingClm := boxIdx % 3 * 3
		for j := 0; j < 9; j++ {
			numIdx := j
			for rowIdx := startingRow; rowIdx < startingRow+3; rowIdx++ {
				for clmIdx := startingClm; clmIdx < startingClm+3; clmIdx++ {
					matrix[rowIdx*81+clmIdx*9+numIdx][81*3+boxIdx*9+numIdx] = true
				}
			}
		}
	}
}

func main() {
	//matrix := [9][9]int{}
	constraintMatrix := [ROWSIZE][CLMSIZE]bool{}
	wg := sync.WaitGroup{}

	constraints := []func([ROWSIZE][CLMSIZE]bool, *sync.WaitGroup){rowClmConstraint, rowNumConstraint, clmNumConstraint, boxNumConstraint}

	for _, constraint := range constraints {
		wg.Add(1)
		go constraint(constraintMatrix, &wg)
	}
	wg.Wait()

	custype.NewDancingLink(constraintMatrix)

}
