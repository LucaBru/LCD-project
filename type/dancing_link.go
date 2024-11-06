package custype

import (
	"sync"
)

const CLMSIZE int = 324
const ROWSIZE int = 729

type Node struct {
	Left   *Node
	Right  *Node
	Above  *Node
	Below  *Node
	header *ClmHeader
}

type ClmHeader struct {
	id   int
	size int
	node Node
}

type DancingLink struct {
	root *ClmHeader
}

/*
given the boolean matrix
	=> create the node matrix, which creates the headers too
	=> binds the rows
	=> binds the columns
*/

func NewDancingLink(matrix [ROWSIZE][CLMSIZE]bool) *DancingLink {
	nodeMatrix := getNodeMatrix(&matrix)
	root := bindColumns(nodeMatrix)

	return &DancingLink{
		root,
	}
}

func getHeaders() *ClmHeader {
	clms := CLMSIZE

	root := &ClmHeader{id: 0, size: 0, node: Node{}}
	root.node.header = root
	clms--
	iter := root

	for i := 0; i < clms; i++ {
		header := &ClmHeader{id: i, size: 0, node: Node{Left: &iter.node}}
		header.node.header = header
		iter.node.Right = &header.node
		iter = iter.node.header
	}

	iter.node.Right = &root.node
	root.node.Left = &iter.node

	return root
}

func bindColumns(matrix *[ROWSIZE][CLMSIZE]*Node) *ClmHeader {
	header := getHeaders()
	headerIter := header
	wg := sync.WaitGroup{}
	wg.Add(CLMSIZE)

	for i := 0; i < CLMSIZE; i++ {
		go bindsSingleColumn(matrix, i, headerIter, &wg)
		headerIter = headerIter.node.Right.header
	}
	wg.Add(CLMSIZE)

	return header
}

func bindsSingleColumn(matrix *[ROWSIZE][CLMSIZE]*Node, clmIdx int, header *ClmHeader, wg *sync.WaitGroup) {
	if clmIdx < 0 || clmIdx >= CLMSIZE {
		return
	}

	idx := 0
	for ; idx < ROWSIZE && matrix[idx][clmIdx] == nil; idx++ {
	}

	if idx == ROWSIZE {
		return

	}

	previous := matrix[idx][clmIdx]
	header.size++
	header.node.Below = previous
	previous.Above = &header.node
	previous.header = header

	for i := idx; i < ROWSIZE; i++ {
		if matrix[i][clmIdx] != nil {
			node := matrix[i][clmIdx]
			previous.Below = node
			node.Above = previous
			node.header = header
			header.size++
			previous = node
		}
	}

	previous.Below = &header.node
	header.node.Above = previous
}

func getNodeMatrix(matrix *[ROWSIZE][CLMSIZE]bool) *[ROWSIZE][CLMSIZE]*Node {
	nodeMatrix := [ROWSIZE][CLMSIZE]*Node{}
	wg := sync.WaitGroup{}
	wg.Add(ROWSIZE)
	for idx, row := range matrix {
		nodeRow := [CLMSIZE]*Node{}
		go createNodesRow(&row, &nodeRow, &wg)
		nodeMatrix[idx] = nodeRow
	}
	wg.Wait()
	return &nodeMatrix
}

func createNodesRow(row *[CLMSIZE]bool, nodesRow *[CLMSIZE]*Node, wg *sync.WaitGroup) {
	defer wg.Done()

	i := 0
	for ; i < CLMSIZE && !row[i]; i++ {
	}

	previous := &Node{}

	if i == CLMSIZE {
		return

	}

	for idx, cell := range row[i+1:] {
		if cell {
			node := &Node{}
			node.Left = previous
			previous.Right = node
			nodesRow[idx] = node
			previous = node
		}
	}
}

/* func algorithmX(root *ClmHeader, partialSol []*Node) []*Node {
	if root == nil {
		return partialSol
	}

	clm := root
	for iter := root.right; iter != nil; iter = iter.right {
		if iter.size > clm.size {
			clm = iter
		}
	}

	//getting the matrix & the clm updates the matrix
	//to create a new column view => easy

	//include r in the partial solution,
	//take a modified view of the matrix
	//proceed recursively

}

func cover(node *Node, dancingLink *DancingLink) *DancingLink {
	for rowIter := node.Right; rowIter != node; rowIter = rowIter.Right {
		//row iter point to all ones in the same row
		for clm
	}
} */
