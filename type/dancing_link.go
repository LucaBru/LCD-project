package custype

import (
	"sync"
)

type Node struct {
	Left   *Node
	Right  *Node
	Above  *Node
	Below  *Node
	header *ClmHeader
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

func NewDancingLink(matrix [][]bool) *DancingLink {
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
	root := previous

	size := 0
	iter := root
	for iter != nil {
		size++
		iter = iter.right
	}

	//create the node matrix
	nodeMatrix := [][]*Node{}
	wg := sync.WaitGroup{}
	clmHeaderIterator := root
	for _, row := range matrix {
		nodeRow := make([]*Node, len(row))
		nodeMatrix = append(nodeMatrix, nodeRow)
		wg.Add(1)
		go createNodesRow(row, nodeRow, clmHeaderIterator, &wg)
		//clmHeaderIterator = clmHeaderIterator.right
	}
	wg.Wait()

	//binds columns nodes
	for i := 0; i < len(matrix[0]); i++ {
		wg.Add(1)
		go bindClmNodes(i, nodeMatrix, &wg)
	}
	wg.Wait()

	//binds rows nodes
	wg.Add(len(nodeMatrix))
	for _, row := range nodeMatrix {
		go bindRowNodes(row, &wg)
	}
	wg.Wait()

	//binds header to its column list
	iterator := root
	iteratorIdx := 0
	for iterator != nil {
		if iterator.size > 0 {
			i := 0
			for ; nodeMatrix[i][iteratorIdx] == nil; i++ {
			}
			nodeMatrix[i][iteratorIdx].header.clm = nodeMatrix[i][iteratorIdx]
		}
		iteratorIdx++
		iterator = iterator.right
	}

	return &DancingLink{
		root,
	}
}

func createNodesRow(row []bool, nodesRow []*Node, clmHeader *ClmHeader, wg *sync.WaitGroup) {
	defer wg.Done()

	if len(row) != len(nodesRow) {
		return
	}

	for idx, cell := range row {
		if cell {
			nodesRow[idx] = &Node{Left: nil, Right: nil, Below: nil, Above: nil, header: clmHeader}
			clmHeader.size++
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
