package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
)

func parseFile(filename string) (Graph, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Unable to open file: ", err)
	}
	defer file.Close()

	graph := Graph{}
	graph.initGraph()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return graph, errors.New("malformed file")
	}
	nodeCount, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return graph, err
	}

	// initialize nodes
	for i := 0; i < nodeCount; i++ {
		if !scanner.Scan() {
			return graph, errors.New("malformed file")
		}
		currentLine := scanner.Text()
		split := strings.Split(currentLine, " ")
		if len(split) != nodeCount {
			return graph, errors.New("malformed file")
		}
		x, err := strconv.ParseFloat(split[0], 64)
		if err != nil {
			return graph, err
		}
		y, err := strconv.ParseFloat(split[1], 64)
		if err != nil {
			return graph, err
		}
		graph.addNode(split[2], x, y)
	}
	// adjacency matrix
	for i := 0; i < nodeCount; i++ {
		if !scanner.Scan() {
			return graph, errors.New("malformed file")
		}
		currentLine := scanner.Text()
		split := strings.Split(currentLine, " ")
		if len(split) != nodeCount {
			return graph, errors.New("malformed file")
		}
		currentNode := graph.getNodeNameAtIndex(i)
		for index, value := range split {
			currWeight, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return graph, err
			}
			if currWeight < 0 {
				return graph, errors.New("malformed file")
			}
			pairedNode := graph.getNodeNameAtIndex(index)
			if currWeight > 0 {
				graph.addEdge(currentNode, pairedNode, currWeight)
			}
		}
	}
	return graph, nil
}
