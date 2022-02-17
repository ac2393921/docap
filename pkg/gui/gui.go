package gui

import (
	"github.com/golang-collections/collections/stack"
	"github.com/jroimartin/gocui"
)

var OverlappingEdges = false

type Gui struct {
	g     *gocui.Gui
	state guiState
}

type servicePanelState struct {
	SelectedLine int
	ContextIndex int
}

type mainPanelState struct {
	ObjectKey string
}

type panelStates struct {
	Services *servicePanelState
	Main     *mainPanelState
}

type guiState struct {
	MenuItemCount    int
	PreviousViews    *stack.Stack
	Panels           *panelStates
	SubProcessOutput string
	// Stats            map[string]command.ContainerStats
	SessionIndex int
}

func NewGui() (*Gui, error) {
	initialState := guiState{
		Panels: &panelStates{
			Services: &servicePanelState{SelectedLine: -1, ContextIndex: 0},
			Main:     &mainPanelState{ObjectKey: ""},
		},
		SessionIndex:  0,
		PreviousViews: stack.New(),
	}

	g := &Gui{
		state: initialState,
	}

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

	if err = gui.keybindings(g); err != nil {
		return err
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}

func (gui *Gui) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
