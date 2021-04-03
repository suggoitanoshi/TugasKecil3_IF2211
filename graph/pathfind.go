package graph

import (
	"errors"
	"math"
)

type nodePath struct {
	name       NodeName
	before     *nodePath
	cost, heur float64
}

func Astar(graph Graph, start NodeName, end NodeName) ([]NodeName, float64, error) {
	cost := 0.0
	queue := make([]nodePath, 0)
	path := make([]NodeName, 0)

	if _, exists := graph.Nodes[start]; !exists {
		return path, math.NaN(), errors.New("invalid starting node")
	}
	if _, exists := graph.Nodes[end]; !exists {
		return path, math.NaN(), errors.New("invalid ending node")
	}
	queue = append(queue, nodePath{start, nil, 0, graph.GetNodeEuclideanDistance(start, end)})
	for len(queue) != 0 {
		smallestF := math.Inf(1)
		var curr *nodePath
		var idx int
		for i, v := range queue {
			f := v.cost + v.heur
			if f < smallestF {
				curr = &queue[i]
				idx = i
				smallestF = f
			}
		}
		queue = append(queue[idx+1:], queue[:idx]...)
		if curr.name == end {
			for curr != nil {
				path = append([]NodeName{curr.name}, path...)
				cost += curr.cost
				curr = curr.before
			}
			return path, cost, nil
		}
		for adj, weight := range graph.Nodes[curr.name].Edges {
			lastF := math.Inf(1)
			idx := -1
			for i, n := range queue {
				if n.name == adj {
					lastF = n.heur
					idx = i
				}
			}
			f := curr.cost + weight + graph.GetNodeEuclideanDistance(adj, end)
			if f < lastF {
				if idx == -1 {
					idx = len(queue)
					queue = append(queue, nodePath{})
					queue[idx].name = adj
				}
				queue[idx].heur = f
				queue[idx].cost = curr.cost + weight
				queue[idx].before = curr
			}
		}
	}
	return path, cost, nil
}
