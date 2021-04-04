package graph

import (
	"errors"
	"strconv"
	"strings"
)

func (g *Graph) ParseContent(content string) error {
	lines := strings.Split(content, "\n")
	nodeCount, err := strconv.Atoi(lines[0])
	if err != nil {
		return err
	}

	if len(lines) < (2*nodeCount)+1 {
		return errors.New("malformed file")
	}

	// initialize nodes
	for i := 0; i < nodeCount; i++ {
		currentLine := lines[i+1]
		split := strings.Split(currentLine, " ")
		if len(split) < 3 {
			return errors.New("malformed file")
		}
		x, err := strconv.ParseFloat(split[0], 64)
		if err != nil {
			return err
		}
		y, err := strconv.ParseFloat(split[1], 64)
		if err != nil {
			return err
		}
		g.AddNode(strings.Join(split[2:], " "), x, y)
	}
	// adjacency matrix
	for i := 0; i < nodeCount; i++ {
		currentLine := lines[i+nodeCount+1]
		split := strings.Split(currentLine, " ")
		if len(split) != nodeCount {
			return errors.New("malformed file")
		}
		currentNode := g.GetNodeNameAtIndex(i)
		for index, value := range split {
			currWeight, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			if currWeight < 0 {
				return errors.New("malformed file")
			}
			pairedNode := g.GetNodeNameAtIndex(index)
			if currWeight > 0 {
				g.AddEdge(currentNode, pairedNode, g.GetNodeDistance(currentNode, pairedNode))
			}
		}
	}
	return nil
}
