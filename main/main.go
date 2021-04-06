package main

import (
	"fmt"
	"net/http"
	"strconv"
	"syscall/js"
	"tanoshi/graph"

	"github.com/davvo/mercator"
	"github.com/serjvanilla/go-overpass"
)

var document, cost, output, fromSelect, toSelect js.Value
var g graph.Graph
var sigma, sigmaGraph js.Value
var theMap js.Value
var fromMap, firstPick bool
var firstID string

const defaultColor = "#4287f5"
const highlightColor = "#f54e42"
const pathColor = "#39e684"
const defaultWeight = 5
const defaultRadius = 5

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

	selectType := document.Call("querySelector", "#type")
	if selectType.Get("value").String() == "globe" {
		g.IsCartes = false
	} else {
		g.IsCartes = true
	}
	selectType.Call("addEventListener", "change", js.FuncOf(func(this js.Value, ev []js.Value) interface{} {
		t := ev[0].Get("target").Get("value")
		if !t.IsUndefined() {
			if t.String() == "globe" {
				g.IsCartes = false
			} else {
				g.IsCartes = true
			}
		}
		return nil
	}))

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
			mapCont.Get("style").Set("display", "")
			sigmaGraph.Call("clear")
			output.Get("style").Set("display", "none")
			g.IsCartes = false
		} else {
			mapCont.Get("style").Set("display", "none")
			output.Get("style").Set("display", "")
			g.IsCartes = selectType.Get("value").String() != "globe"
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
				edges.Index(i).Set("color", defaultColor)
			}
			nodes := sigmaGraph.Call("nodes")
			for i := 0; i < nodes.Length(); i++ {
				nodes.Index(i).Set("color", defaultColor)
			}
			for i, n := range path {
				node := sigmaGraph.Call("nodes", n)
				node.Set("color", highlightColor)
				if i < len(path)-1 {
					edge := sigmaGraph.Call("edges", path[i]+"-"+path[i+1])
					edge.Set("color", pathColor)
					edge = sigmaGraph.Call("edges", path[i+1]+"-"+path[i])
					edge.Set("color", pathColor)
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
		clearAllGraph()
		err := g.ParseContent(ev[0].String())
		if !fromMap {
			createSigmaGraph()
			if err == nil {
				for _, v := range g.NodeNames {
					createOpt(&fromSelect, v)
					createOpt(&toSelect, v)
				}
			}
		} else {
			createMapGraph()
			avg := []interface{}{}
			for _, n := range g.Nodes {
				avg = append(avg, []interface{}{n.Coord.X, n.Coord.Y})
			}
			theMap.Call("fitBounds", avg)
		}
		return nil
	}))
	return nil
}

func createSigmaGraph() {
	sigmaGraph.Call("clear")
	for name, node := range g.Nodes {
		var displayX, displayY float64
		if g.IsCartes {
			displayX, displayY = node.Coord.X, node.Coord.Y
		} else {
			displayX, displayY = mercator.LatLonToMeters(node.Coord.X, node.Coord.Y)
		}
		sigmaGraph.Call("addNode", map[string]interface{}{
			"id":    name,
			"label": name,
			"x":     displayX,
			"y":     displayY,
			"size":  1,
			"color": "#00f",
		})
	}
	for name, node := range g.Nodes {
		for adj, w := range node.Edges {
			sigmaGraph.Call("addEdge", map[string]interface{}{
				"id":     name + "-" + adj,
				"label":  w,
				"source": name,
				"target": adj,
			})
		}
	}
	sigma.Call("refresh")
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
		(node.close_down.stnodes;.st;);
		out skel qt;
	}`, bbox.Call("getSouth").Float(), bbox.Call("getWest").Float(), bbox.Call("getNorth").Float(), bbox.Call("getEast").Float()))
	clearAllGraph()
	g.IsCartes = false
	for id, node := range result.Nodes {
		if node.Lat != 0 && node.Lon != 0 {
			g.AddNode(strconv.Itoa(int(id)), node.Lat, node.Lon)
		}
	}
	for _, way := range result.Ways {
		lastID := ""
		for idx := 0; idx < len(way.Nodes); idx++ {
			cID := strconv.Itoa(int(way.Nodes[idx].ID))
			if _, ok := g.Nodes[cID]; ok {
				if lastID == "" {
					lastID = cID
					continue
				}
				idB := strconv.Itoa(int(way.Nodes[idx].ID))
				dist := g.GetNodeDistance(lastID, idB)
				g.AddEdge(lastID, idB, dist)
				g.AddEdge(idB, lastID, dist)
				lastID = idB
			}
		}
	}
	for name, node := range g.Nodes {
		for e, w := range node.Edges {
			if w < 0.02 {
				for an, aw := range g.Nodes[e].Edges {
					g.AddEdge(name, an, aw)
					g.AddEdge(an, name, aw)
				}
				g.RemoveNode(e)
			}
		}
	}
	createMapGraph()
}

func createMapGraph() {
	for _, n := range g.Nodes {
		if len(n.Edges) == 0 {
			g.RemoveNode(n.Name)
		} else {
			currentMapOverlay.nodesInMap[n.Name] = js.Global().Get("L").Call("circleMarker", []interface{}{n.Coord.X, n.Coord.Y}, map[string]interface{}{"color": defaultColor, "radius": defaultRadius, "weight": defaultWeight})
			currentMapOverlay.nodesInMap[n.Name].Set("nodeID", n.Name)
			for e := range n.Edges {
				fromPair := []interface{}{n.Coord.X, n.Coord.Y}
				toPair := []interface{}{g.Nodes[e].Coord.X, g.Nodes[e].Coord.Y}
				currentMapOverlay.edgesInMap[n.Name+"-"+e] = js.Global().Get("L").Call("polyline", []interface{}{fromPair, toPair}, map[string]interface{}{"color": defaultColor, "weight": defaultWeight})
			}
		}
	}
	for _, edge := range currentMapOverlay.edgesInMap {
		edge.Call("addTo", currentMapOverlay.overlayLayer)
	}
	for _, node := range currentMapOverlay.nodesInMap {
		node.Call("addTo", currentMapOverlay.overlayLayer)
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
				n.Call("setStyle", map[string]interface{}{"color": defaultColor})
				n.Call("setRadius", defaultRadius)
			}
			currentMapOverlay.pathOver.Call("clearLayers")
		}
		currentMapOverlay.nodesInMap[nID.String()].Call("setStyle", map[string]interface{}{"color": highlightColor})
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
					currentMapOverlay.nodesInMap[path[i]].Call("setStyle", map[string]interface{}{"color": highlightColor})
				}
				js.Global().Get("L").Call("polyline", points, map[string]interface{}{"color": pathColor, "weight": defaultWeight}).Call("addTo", currentMapOverlay.pathOver)
			}
			cost.Set("innerText", pathCost)
			firstPick = true
		}
	}
	return nil
}

func clearAllGraph() {
	g.ClearGraph()
	currentMapOverlay.overlayLayer.Call("clearLayers")
	currentMapOverlay.edgesInMap = map[string]js.Value{}
	currentMapOverlay.nodesInMap = map[string]js.Value{}
	fromSelect.Set("innerHTML", "")
	toSelect.Set("innerHTML", "")
}
