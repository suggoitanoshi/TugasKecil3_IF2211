package main

import (
	"errors"
	"fmt"
	"math"
)

type (
	NodeName = string
)

type Coordinate struct {
	X, Y float64
}

func (coord *Coordinate) initCoord(x, y float64) {
	coord.X = x
	coord.Y = y
}

func (coord *Coordinate) distanceTo(other Coordinate) float64 {
	return math.Hypot(coord.X-other.X, coord.Y-other.Y)
}

type Node struct {
	Name  NodeName
	coord Coordinate
	Edges map[NodeName]float64
}

func (node *Node) initNode(name string, coord Coordinate) {
	node.Name = name
	node.coord = coord
	node.Edges = make(map[NodeName]float64)
}

func (node *Node) addEdge(otherNode Node, weight float64) {
	if _, ok := node.Edges[otherNode.Name]; !ok {
		// add not-yet-here edge
		node.Edges[otherNode.Name] = weight
	}
}

func (node *Node) removeEdge(otherNode Node) {
	if _, ok := node.Edges[otherNode.Name]; ok {
		// delete only if exists
		delete(node.Edges, otherNode.Name)
	}
}

type Graph struct {
	NodeNames []NodeName
	Nodes     map[NodeName]Node
}

func (graph *Graph) initGraph() Graph {
	graph.NodeNames = make([]NodeName, 0)
	graph.Nodes = make(map[NodeName]Node)
	return *graph
}

func (graph *Graph) addNode(name string, x, y float64) {
	if _, ok := graph.Nodes[name]; !ok {
		graph.NodeNames = append(graph.NodeNames, name)
		coord := Coordinate{}
		coord.initCoord(x, y)
		node := Node{}
		node.initNode(name, coord)
		graph.Nodes[name] = node
	}
}

func (graph *Graph) clearGraph() {
	graph.Nodes = make(map[NodeName]Node)
	graph.NodeNames = []NodeName{}
}

func (graph *Graph) getNodeNameAtIndex(i int) NodeName {
	return graph.NodeNames[i]
}

func (graph *Graph) addEdge(nameA, nameB string, weight float64) error {
	if _, exist := graph.Nodes[nameA]; !exist {
		return errors.New("nodeA doesn't exist")
	}
	if _, exist := graph.Nodes[nameB]; !exist {
		return errors.New("nodeB doesn't exist")
	}
	if nameA != nameB {
		nodeA := graph.Nodes[nameA]
		nodeB := graph.Nodes[nameB]
		nodeA.addEdge(nodeB, weight)
		nodeB.addEdge(nodeA, weight)
	}
	return nil
}

func (graph *Graph) showGraph() {
	for name, node := range graph.Nodes {
		fmt.Printf("%s: ", name)
		for adjacent, weight := range node.Edges {
			fmt.Printf("[%s,%.4f] ", adjacent, weight)
		}
		fmt.Printf("\n")
	}
}
