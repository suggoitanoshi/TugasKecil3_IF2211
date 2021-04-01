package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var app *gtk.Application

func readFile(filename string) {
	ba, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("Could not open file: ", err)
	}
	fileContent := string(ba)
	fmt.Println(fileContent)
}

func chooseFile() {
	// do something
	dialog, err := gtk.FileChooserNativeDialogNew("Choose file", nil, gtk.FILE_CHOOSER_ACTION_OPEN, "Open", "Cancel")
	if err != nil {
		log.Fatal("Fail opening file: ", err)
	}
	res := dialog.Run()
	if res == int(gtk.RESPONSE_ACCEPT) {
		chooser := dialog.FileChooser
		readFile(chooser.GetFilename())
	}
}

var signals = map[string]interface{}{
	"chooseFile": chooseFile,
}

func main() {
	const appID = "fun.suggoi"
	app, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Fatal("Failed to create app: ", err)
	}
	app.Connect("activate", func() { startup(app) })
	os.Exit(app.Run(os.Args))
}

func startup(app *gtk.Application) {
	builder, err := gtk.BuilderNewFromFile("layout.glade")
	if err != nil {
		log.Fatal("Could not create builder: ", err)
	}
	builder.ConnectSignals(signals)
	// obj, err := builder.GetObject("chooseFile")
	// if err != nil {
	// 	log.Fatal("Could not access button: ", err)
	// }
	// chooseFile := obj.(*gtk.Button)
	// chooseFile.Connect("clicked")
	obj, err := builder.GetObject("window")
	if err != nil {
		log.Fatal("Could not access window: ", err)
	}
	window := obj.(*gtk.Window)
	window.ShowAll()
	app.AddWindow(window)
}
