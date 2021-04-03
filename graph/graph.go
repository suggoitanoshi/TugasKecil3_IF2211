package graph

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

func NewCoord(x, y float64) Coordinate {
	coord := Coordinate{}
	coord.initCoord(x, y)
	return coord
}

func (coord *Coordinate) initCoord(x, y float64) {
	coord.X = x
	coord.Y = y
}

func (coord Coordinate) DistanceTo(other Coordinate) float64 {
	return math.Hypot(coord.X-other.X, coord.Y-other.Y)
}

type Node struct {
	Name  NodeName
	coord Coordinate
	Edges map[NodeName]float64
}

func NewNode(name NodeName, x, y float64) Node {
	node := Node{}
	node.initNode(name, NewCoord(x, y))
	return node
}

func (node *Node) initNode(name string, coord Coordinate) {
	node.Name = name
	node.coord = coord
	node.Edges = make(map[NodeName]float64)
}

func (node *Node) AddEdge(otherNode Node, weight float64) {
	if _, ok := node.Edges[otherNode.Name]; !ok {
		// add not-yet-here edge
		node.Edges[otherNode.Name] = weight
	}
}

func (node *Node) RemoveEdge(otherNode Node) {
	delete(node.Edges, otherNode.Name)
}

type Graph struct {
	NodeNames []NodeName
	Nodes     map[NodeName]Node
}

func NewGraph() Graph {
	return (&Graph{}).initGraph()
}

func (graph *Graph) initGraph() Graph {
	graph.NodeNames = make([]NodeName, 0)
	graph.Nodes = make(map[NodeName]Node)
	return *graph
}

func (graph *Graph) AddNode(name string, x, y float64) {
	if _, ok := graph.Nodes[name]; !ok {
		graph.NodeNames = append(graph.NodeNames, name)
		coord := Coordinate{}
		coord.initCoord(x, y)
		node := Node{}
		node.initNode(name, coord)
		graph.Nodes[name] = node
	}
}

func (graph *Graph) ClearGraph() {
	graph.Nodes = make(map[NodeName]Node)
	graph.NodeNames = []NodeName{}
}

func (graph Graph) GetNodeNameAtIndex(i int) NodeName {
	return graph.NodeNames[i]
}

func (graph *Graph) AddEdge(nameA, nameB string, weight float64) error {
	if _, exist := graph.Nodes[nameA]; !exist {
		return errors.New("nodeA doesn't exist")
	}
	if _, exist := graph.Nodes[nameB]; !exist {
		return errors.New("nodeB doesn't exist")
	}
	if nameA != nameB {
		nodeA := graph.Nodes[nameA]
		nodeB := graph.Nodes[nameB]
		nodeA.AddEdge(nodeB, weight)
		nodeB.AddEdge(nodeA, weight)
	}
	return nil
}

func (graph Graph) ToString() string {
	var str string
	for name, node := range graph.Nodes {
		str += fmt.Sprintf("%s: ", name)
		for adjacent, weight := range node.Edges {
			str += fmt.Sprintf("[%s,%.4f] ", adjacent, weight)
		}
		str += "\n"
	}
	return str
}

func (graph Graph) ShowGraph() {
	fmt.Println(graph.ToString())
}

func (graph Graph) GetNodeEuclideanDistance(nodeA, nodeB NodeName) float64 {
	return graph.Nodes[nodeA].coord.DistanceTo(graph.Nodes[nodeB].coord)
}
