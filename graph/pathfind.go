package graph

import (
	"errors"
	"math"
	"tanoshi/prioqueue"
)

type nodePath struct {
	name       NodeName
	before     *nodePath
	cost, heur float64
}

func Astar(graph Graph, start NodeName, end NodeName) ([]NodeName, float64, error) {
	cost := 0.0
	queue := prioqueue.NewPQ()
	path := make([]NodeName, 0)

	if _, exists := graph.Nodes[start]; !exists {
		return path, math.NaN(), errors.New("invalid starting node")
	}
	if _, exists := graph.Nodes[end]; !exists {
		return path, math.NaN(), errors.New("invalid ending node")
	}
	heur := graph.GetNodeDistance(start, end)
	startPath := &nodePath{start, nil, 0, heur}
	queue.Enqueue(startPath, heur)
	visited := make(map[string]*nodePath)
	for !queue.IsEmpty() {
		curr := queue.Dequeue().(*nodePath)
		if curr.name == end {
			cost = curr.cost
			for curr != nil {
				path = append([]NodeName{curr.name}, path...)
				curr = curr.before
			}
			return path, cost, nil
		}
		for adj, weight := range graph.Nodes[curr.name].Edges {
			lastF := math.Inf(1)
			f := curr.cost + weight + graph.GetNodeDistance(adj, end)
			if v, ok := visited[adj]; !ok {
				visited[adj] = &nodePath{adj, curr, curr.cost + weight, f}
			} else {
				lastF = v.heur
			}
			if f < lastF {
				visited[adj].cost = curr.cost + weight
				visited[adj].before = curr
				visited[adj].heur = f
				idx := -1
				for i, n := range queue.GetElements() {
					if n.(*nodePath).name == adj {
						idx = i
						break
					}
				}
				if idx == -1 {
					queue.Enqueue(visited[adj], f)
				} else {
					last := queue.DequeueElementAt(idx).(*nodePath)
					queue.Enqueue(last, f)
				}
			}
		}
	}
	return path, cost, nil
}
