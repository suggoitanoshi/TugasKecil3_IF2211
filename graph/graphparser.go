package graph

import (
	"errors"
	"strconv"
	"strings"
)

func ParseContent(content string) (Graph, error) {
	g := NewGraph()

	lines := strings.Split(content, "\n")
	nodeCount, err := strconv.Atoi(lines[0])
	if err != nil {
		return g, err
	}

	if len(lines) < (2*nodeCount)+1 {
		return g, errors.New("malformed file")
	}

	// initialize nodes
	for i := 0; i < nodeCount; i++ {
		currentLine := lines[i+1]
		split := strings.Split(currentLine, " ")
		if len(split) < 3 {
			return g, errors.New("malformed file")
		}
		x, err := strconv.ParseFloat(split[0], 64)
		if err != nil {
			return g, err
		}
		y, err := strconv.ParseFloat(split[1], 64)
		if err != nil {
			return g, err
		}
		g.AddNode(strings.Join(split[2:], " "), x, y)
	}
	// adjacency matrix
	for i := 0; i < nodeCount; i++ {
		currentLine := lines[i+nodeCount+1]
		split := strings.Split(currentLine, " ")
		if len(split) != nodeCount {
			return g, errors.New("malformed file")
		}
		currentNode := g.GetNodeNameAtIndex(i)
		for index, value := range split {
			currWeight, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return g, err
			}
			if currWeight < 0 {
				return g, errors.New("malformed file")
			}
			pairedNode := g.GetNodeNameAtIndex(index)
			if currWeight > 0 {
				g.AddEdge(currentNode, pairedNode, currWeight)
			}
		}
	}
	return g, nil
}
