package custype

import (
	"sync"
)

type Node struct {
	Left  *Node
	Right *Node
	Above *Node
	Below *Node
}

type ClmHeader struct {
	id    int
	size  int
	left  *ClmHeader
	right *ClmHeader
	clm   *Node
}

type DancingLink struct {
	root *ClmHeader
}

func NewDancingLink(matrix [][]bool) [][]*Node {
	clms := len(matrix[0])
	var previous *ClmHeader
	for idx := range clms {
		clmHeader := &ClmHeader{
			id:    clms - idx - 1,
			size:  0,
			clm:   nil,
			left:  nil,
			right: previous,
		}

		if previous != nil {
			previous.left = clmHeader
		}

		previous = clmHeader
	}

	nodeMatrix := createNodeMatrix(matrix)

	wg := sync.WaitGroup{}

	for i := 0; i < len(matrix[0]); i++ {
		wg.Add(1)
		go bindClmNodes(i, nodeMatrix, &wg)
	}
	wg.Wait()

	wg.Add(len(nodeMatrix))
	for _, row := range nodeMatrix {
		go bindRowNodes(row, &wg)
	}
	wg.Wait()

	return nodeMatrix
}

func createNodeMatrix(matrix [][]bool) [][]*Node {
	nodeMatrix := [][]*Node{}
	wg := sync.WaitGroup{}
	for _, row := range matrix {
		nodeRow := make([]*Node, len(row))
		nodeMatrix = append(nodeMatrix, nodeRow)
		wg.Add(1)
		go createNodesRow(row, nodeRow, &wg)
	}
	wg.Wait()
	return nodeMatrix
}

func createNodesRow(row []bool, nodesRow []*Node, wg *sync.WaitGroup) {
	defer wg.Done()

	if len(row) != len(nodesRow) {
		return
	}

	for idx, cell := range row {
		if cell {
			nodesRow[idx] = &Node{Left: nil, Right: nil, Below: nil, Above: nil}
		} else {
			nodesRow[idx] = nil
		}
	}
}

func bindClmNodes(clmIdx int, matrix [][]*Node, wg *sync.WaitGroup) {
	defer wg.Done()

	if len(matrix) > 0 && len(matrix[0]) < clmIdx {
		return
	}

	firstNodeIdx := 0
	for firstNodeIdx < len(matrix) && matrix[firstNodeIdx][clmIdx] == nil {
		firstNodeIdx++
	}

	if firstNodeIdx == len(matrix) {
		return
	}

	above := matrix[firstNodeIdx][clmIdx]

	for i := firstNodeIdx + 1; i < len(matrix); i++ {
		if matrix[i][clmIdx] != nil {
			node := matrix[i][clmIdx]
			node.Above = above
			above.Below = node
			above = node
		}
	}
}

func bindRowNodes(row []*Node, wg *sync.WaitGroup) {
	defer wg.Done()

	if len(row) == 0 {
		return
	}

	firstNodeIdx := 0
	for firstNodeIdx < len(row) && row[firstNodeIdx] == nil {
		firstNodeIdx++
	}

	if firstNodeIdx == len(row) {
		return
	}

	left := row[firstNodeIdx]
	for i := firstNodeIdx + 1; i < len(row); i++ {
		if row[i] != nil {
			node := row[i]
			node.Left = left
			left.Right = node
			left = node
		}
	}
}
