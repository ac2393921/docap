package gui

import (
	"log"

	"github.com/jroimartin/gocui"
)

var OverlappingEdges = false

type Gui struct {
	g *gocui.Gui
}

func NewGui() (*Gui, error) {
	g := &Gui{}
	return g, nil
}

func (gui *Gui) Run() error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	defer g.Close()

	gui.g = g

	g.SetManagerFunc(gui.layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
