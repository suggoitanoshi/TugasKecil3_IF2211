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

func (coord Coordinate) EuclideanDistanceTo(other Coordinate) float64 {
	return math.Hypot(coord.X-other.X, coord.Y-other.Y)
}

func degToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

func (coord Coordinate) toRad() (float64, float64) {
	return degToRad(coord.X), degToRad(coord.Y)
}

func haversine(rad float64) float64 {
	return math.Pow(math.Sin(rad/2), 2)
}

func (coord Coordinate) HaversineDistanceTo(other Coordinate) float64 {
	aLat, aLon := coord.toRad()
	bLat, bLon := other.toRad()
	h := haversine(bLat-aLat) + math.Cos(aLat)*math.Cos(bLat)*haversine(bLon-aLon)
	return 2 * 6371.0 * math.Asin(math.Sqrt(h))
}

type Node struct {
	Name  NodeName
	Coord Coordinate
	Edges map[NodeName]float64
}

func NewNode(name NodeName, x, y float64) Node {
	node := Node{}
	node.initNode(name, NewCoord(x, y))
	return node
}

func (node *Node) initNode(name string, coord Coordinate) {
	node.Name = name
	node.Coord = coord
	node.Edges = make(map[NodeName]float64)
}

func (node *Node) AddEdge(otherNode Node, weight float64) {
	if _, ok := node.Edges[otherNode.Name]; !ok {
		// add not-yet-here edge
		node.Edges[otherNode.Name] = weight
	}
}

type Graph struct {
	NodeNames []NodeName
	Nodes     map[NodeName]Node
	IsCartes  bool
}

func NewGraph() Graph {
	return (&Graph{}).initGraph()
}

func (graph *Graph) initGraph() Graph {
	graph.NodeNames = make([]NodeName, 0)
	graph.Nodes = make(map[NodeName]Node)
	graph.IsCartes = false
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

func (graph *Graph) RemoveNode(name string) {
	for e := range graph.Nodes[name].Edges {
		delete(graph.Nodes[e].Edges, name)
	}
	for i, n := range graph.NodeNames {
		if n == name {
			graph.NodeNames = append(graph.NodeNames[:i], graph.NodeNames[i+1:]...)
			break
		}
	}
	delete(graph.Nodes, name)
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

func (graph Graph) GetNodeDistance(nodeA, nodeB NodeName) float64 {
	if graph.IsCartes {
		return graph.GetNodeEuclideanDistance(nodeA, nodeB)
	} else {
		return graph.GetNodeHaversineDistance(nodeA, nodeB)
	}
}

func (graph Graph) GetNodeEuclideanDistance(nodeA, nodeB NodeName) float64 {
	return graph.Nodes[nodeA].Coord.EuclideanDistanceTo(graph.Nodes[nodeB].Coord)
}

func (graph Graph) GetNodeHaversineDistance(nodeA, nodeB NodeName) float64 {
	return graph.Nodes[nodeA].Coord.HaversineDistanceTo(graph.Nodes[nodeB].Coord)
}
