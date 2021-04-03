package main

import (
	"fmt"
	"strings"
	"syscall/js"
	"tanoshi/graph"
)

var document, output, fromSelect, toSelect js.Value
var g graph.Graph

func main() {
	c := make(chan bool)
	document = js.Global().Get("document")
	file := document.Call("createElement", "input")
	file.Set("type", "file")
	file.Call("addEventListener", "change", js.FuncOf(handleFile))
	document.Get("body").Call("appendChild", file)
	search := document.Call("createElement", "button")
	search.Set("value", "Go!")
	search.Set("innerText", "Go!")
	search.Call("addEventListener", "click", js.FuncOf(startAlgo))
	fromSelect = document.Call("createElement", "select")
	toSelect = document.Call("createElement", "select")
	output = document.Call("createElement", "p")
	document.Get("body").Call("appendChild", output)
	document.Get("body").Call("appendChild", fromSelect)
	document.Get("body").Call("appendChild", toSelect)
	document.Get("body").Call("appendChild", search)
	<-c
}

func startAlgo(this js.Value, ev []js.Value) interface{} {
	path, cost, err := graph.Astar(g, fromSelect.Get("value").String(), toSelect.Get("value").String())
	if err == nil {
		output.Set("innerText", strings.Join(path, "->")+fmt.Sprintf("%.5f", cost))
	}
	return nil
}

func handleFile(this js.Value, ev []js.Value) interface{} {
	files := ev[0].Get("target").Get("files")
	files.Index(0).Call("text").Call("then", js.FuncOf(func(this js.Value, ev []js.Value) interface{} {
		var err error
		g, err = graph.ParseContent(ev[0].String())
		if err == nil {
			output.Set("innerText", g.ToString())
			for _, v := range g.NodeNames {
				createOpt(&fromSelect, v)
				createOpt(&toSelect, v)
			}
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
