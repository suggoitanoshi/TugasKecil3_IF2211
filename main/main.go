package main

import (
	"fmt"
	"net/http"
	"strconv"
	"syscall/js"
	"tanoshi/graph"

	"github.com/serjvanilla/go-overpass"
)

var document, cost, output, fromSelect, toSelect js.Value
var g graph.Graph
var sigma, sigmaGraph js.Value
var theMap js.Value
var fromMap, firstPick bool
var firstID string

type mapOverlay struct {
	nodesInMap   map[string]js.Value
	edgesInMap   map[string]js.Value
	overlayLayer js.Value
	pathOver     js.Value
}

var currentMapOverlay mapOverlay

func main() {
	g = graph.NewGraph()

	fromMap = false
	c := make(chan bool)
	document = js.Global().Get("document")
	file := document.Call("querySelector", "#fileIn")
	file.Call("addEventListener", "change", js.FuncOf(handleFile))
	search := document.Call("querySelector", "#go")
	search.Call("addEventListener", "click", js.FuncOf(startAlgo))
	fromSelect = document.Call("querySelector", "#from")
	toSelect = document.Call("querySelector", "#to")
	output = document.Call("querySelector", "#file-io")

	theMap = js.Global().Get("L").Call("map", "map")
	theMap.Call("setView", []interface{}{0.7893, 113.9213}, 5)
	js.Global().Get("L").Call("tileLayer", "https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png").Call("addTo", theMap)

	theMap.Call("on", "click", js.FuncOf(func(this js.Value, ev []js.Value) interface{} {
		latlng := ev[0].Get("latlng")
		latlng.Get("lat")
		return nil
	}))

	mapCont := document.Call("querySelector", "#map-cont")
	mapCont.Get("style").Set("display", "none")

	toggle := document.Call("querySelector", "#toggle")
	toggle.Call("addEventListener", "click", js.FuncOf(func(this js.Value, ev []js.Value) interface{} {
		fromMap = !fromMap
		if fromMap {
			mapCont.Get("style").Set("display", "block")
			sigmaGraph.Call("clear")
			output.Get("style").Set("display", "none")
		} else {
			mapCont.Get("style").Set("display", "none")
			output.Get("style").Set("display", "block")
		}
		return nil
	}))

	getInt := document.Call("querySelector", "#get-int")
	getInt.Call("addEventListener", "click", js.FuncOf(func(this js.Value, ev []js.Value) interface{} {
		go getGraphFromMap()
		return nil
	}))
	sigma = js.Global().Get("sigma").New("output")
	sigmaGraph = sigma.Get("graph")

	currentMapOverlay = mapOverlay{
		make(map[string]js.Value),
		make(map[string]js.Value),
		js.Global().Get("L").Call("featureGroup"),
		js.Global().Get("L").Call("layerGroup"),
	}

	currentMapOverlay.overlayLayer.Call("on", "click", js.FuncOf(mapNodeClick))
	currentMapOverlay.overlayLayer.Call("addTo", theMap)
	currentMapOverlay.pathOver.Call("addTo", theMap)

	cost = document.Call("querySelector", "#cost")

	firstPick = true

	<-c
}

func startAlgo(this js.Value, ev []js.Value) interface{} {
	if fromSelect.Get("value").String() == "" || toSelect.Get("value").String() == "" {
		return nil
	}
	path, pathCost, err := graph.Astar(g, fromSelect.Get("value").String(), toSelect.Get("value").String())
	if err == nil {
		if !fromMap {
			edges := sigmaGraph.Call("edges")
			for i := 0; i < edges.Length(); i++ {
				edges.Index(i).Set("color", "#00f")
			}
			nodes := sigmaGraph.Call("nodes")
			for i := 0; i < nodes.Length(); i++ {
				nodes.Index(i).Set("color", "#00f")
			}
			for i, n := range path {
				node := sigmaGraph.Call("nodes", n)
				node.Set("color", "#f0f")
				if i < len(path)-1 {
					edge := sigmaGraph.Call("edges", path[i]+"-"+path[i+1])
					edge.Set("color", "#f0f")
					edge = sigmaGraph.Call("edges", path[i+1]+"-"+path[i])
					edge.Set("color", "#f0f")
				}
			}
			sigmaGraph.Call("nodes", path[0]).Set("color", "#0f0")
			sigmaGraph.Call("nodes", path[len(path)-1]).Set("color", "#0f0")
			sigma.Call("refresh")
			cost.Set("innerText", pathCost)
		}
	} else {
		output.Set("innerText", err)
	}
	return nil
}

func handleFile(this js.Value, ev []js.Value) interface{} {
	files := ev[0].Get("target").Get("files")
	files.Index(0).Call("text").Call("then", js.FuncOf(func(this js.Value, ev []js.Value) interface{} {
		var err error
		g.ClearGraph()
		g.IsCartes = true
		err = g.ParseContent(ev[0].String())
		if err == nil {
			fromSelect.Set("innerHTML", "")
			toSelect.Set("innerHTML", "")
			for _, v := range g.NodeNames {
				createOpt(&fromSelect, v)
				createOpt(&toSelect, v)
			}
			sigmaGraph.Call("clear")
			for name, node := range g.Nodes {
				sigmaGraph.Call("addNode", map[string]interface{}{
					"id":    name,
					"label": name,
					"x":     node.Coord.X,
					"y":     node.Coord.Y,
					"size":  1,
					"color": "#00f",
				})
			}
			for name, node := range g.Nodes {
				for adj, _ := range node.Edges {
					sigmaGraph.Call("addEdge", map[string]interface{}{
						"id":     name + "-" + adj,
						"source": name,
						"target": adj,
					})
				}
			}
			sigma.Call("refresh")
		}
		return nil
	}))
	return nil
}

func createOpt(el *js.Value, v string) {
	opt := document.Call("createElement", "option")
	opt.Set("value", v)
	opt.Set("innerText", v)
	el.Call("appendChild", opt)
}

func getGraphFromMap() {
	bbox := theMap.Call("getBounds")

	client := overpass.NewWithSettings("https://lz4.overpass-api.de/api/interpreter", 1, http.DefaultClient)
	result, _ := client.Query(fmt.Sprintf(`[out:json][timeout:240][bbox:%.5f,%.5f,%.5f,%.5f];
	way[highway];
	foreach->.st{
	  .st > -> .stnodes;
	  .stnodes < -> .allclose;
	  way.allclose[highway]->.allclose;
	  (.allclose; - .st;)->.close;
	  .close>->.close_down;
	  node.close_down.stnodes;
	  out skel;
	  way.close;
	  out skel;
	}`, bbox.Call("getSouth").Float(), bbox.Call("getWest").Float(), bbox.Call("getNorth").Float(), bbox.Call("getEast").Float()))
	g.ClearGraph()
	currentMapOverlay.overlayLayer.Call("clearLayers")
	g.IsCartes = false
	for id, node := range result.Nodes {
		if node.Lat != 0 && node.Lon != 0 {
			g.AddNode(strconv.Itoa(int(id)), node.Lat, node.Lon)
		}
	}
	for _, way := range result.Ways {
		for _, nodeA := range way.Nodes {
			idA := strconv.Itoa(int(nodeA.ID))
			if _, ok := g.Nodes[idA]; ok {
				for _, nodeB := range way.Nodes {
					idB := strconv.Itoa(int(nodeB.ID))
					if idA != idB {
						if _, ok := g.Nodes[idB]; ok {
							dist := g.GetNodeDistance(idA, idB)
							if dist > .025 {
								g.AddEdge(idA, idB, dist)
								g.AddEdge(idB, idA, dist)
							} else {
								for e, w := range g.Nodes[idA].Edges {
									if _, ok := g.Nodes[idB].Edges[e]; !ok {
										g.AddEdge(idB, e, w)
										g.AddEdge(e, idB, w)
									}
								}
								g.RemoveNode(idA)
							}
						}
					}
				}
			}
		}
	}
	for _, n := range g.Nodes {
		if len(n.Edges) == 0 {
			g.RemoveNode(n.Name)
		} else {
			currentMapOverlay.nodesInMap[n.Name] = js.Global().Get("L").Call("circle", []interface{}{n.Coord.X, n.Coord.Y}, map[string]interface{}{"color": "skyblue", "radius": 5})
			currentMapOverlay.nodesInMap[n.Name].Set("nodeID", n.Name)
			currentMapOverlay.nodesInMap[n.Name].Call("addTo", currentMapOverlay.overlayLayer)
			for e, _ := range n.Edges {
				fromPair := []interface{}{n.Coord.X, n.Coord.Y}
				toPair := []interface{}{g.Nodes[e].Coord.X, g.Nodes[e].Coord.Y}
				currentMapOverlay.edgesInMap[n.Name+"-"+e] = js.Global().Get("L").Call("polyline", []interface{}{fromPair, toPair}, map[string]interface{}{"color": "skyblue"})
				currentMapOverlay.edgesInMap[n.Name+"-"+e].Call("addTo", currentMapOverlay.overlayLayer)
			}
		}
	}
}

func mapNodeClick(this js.Value, ev []js.Value) interface{} {
	pFrom := ev[0].Get("sourceTarget")
	if pFrom.IsUndefined() {
		return nil
	}
	nID := pFrom.Get("nodeID")
	if !nID.IsUndefined() {
		if firstPick {
			for _, n := range currentMapOverlay.nodesInMap {
				n.Call("setStyle", map[string]interface{}{"color": "skyblue"})
			}
			currentMapOverlay.pathOver.Call("clearLayers")
		}
		currentMapOverlay.nodesInMap[nID.String()].Call("setStyle", map[string]interface{}{"color": "red"})
		if firstPick {
			firstID = nID.String()
			firstPick = false
		} else {
			path, pathCost, _ := graph.Astar(g, firstID, nID.String())
			if len(path) > 0 {
				points := []interface{}{}
				for i := 0; i < len(path); i++ {
					fromP := g.Nodes[path[i]].Coord
					points = append(points, []interface{}{fromP.X, fromP.Y})
				}
				js.Global().Get("L").Call("polyline", points, map[string]interface{}{"color": "red"}).Call("addTo", currentMapOverlay.pathOver)
			}
			cost.Set("innerText", pathCost)
			firstPick = true
		}
	}
	return nil
}
