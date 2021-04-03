package graph

import (
	"errors"
	"math"
)

func Astar(graph Graph, start NodeName, end NodeName) ([]NodeName, float64, error) {
	cost := 0.0
	done := false
	current := start
	path := make([]NodeName, 0)

	if _, exists := graph.Nodes[start]; !exists {
		return path, math.NaN(), errors.New("invalid starting node")
	}
	if _, exists := graph.Nodes[end]; !exists {
		return path, math.NaN(), errors.New("invalid ending node")
	}
	for !done {
		path = append(path, current)
		if current == end {
			done = true
		} else {
			total := math.Inf(1)
			chosenWeight := 0.0
			for adj, weight := range graph.Nodes[current].Edges {
				f := cost + graph.GetNodeEuclideanDistance(current, end)
				if f < total {
					total = f
					current = adj
					chosenWeight = weight
				}
			}
			cost += chosenWeight
		}
	}
	return path, cost, nil
}
