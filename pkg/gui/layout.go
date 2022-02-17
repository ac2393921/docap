package gui

import (
	"github.com/jroimartin/gocui"
)

func (gui *Gui) layout(g *gocui.Gui) error {
	g.Highlight = true
	width, height := g.Size()

	minwidth := 9
	minheight := 10

	if width < minwidth || height < minheight {
		v, err := g.SetView("limit", 0, 0, width-1, height-1)
		if err != nil {
			if err.Error() != "unknown view" {
				return err
			}
			v.Title = "Not Enough Space"
			v.Wrap = true
			_, _ = g.SetViewOnTop("limit")
		}
		return nil
	}

	leftSideWidth := width / 3
	mainSideWidth := width / 6

	var vHeights map[string]int
	usableSpace := height - 4
	tallPanels := 3

	vHeights = map[string]int{
		"project":    3,
		"containers": usableSpace/tallPanels + usableSpace%tallPanels,
		"images":     usableSpace / tallPanels,
		"volumes":    usableSpace / tallPanels,
		"options":    1,
	}

	v, err := g.SetView("main", mainSideWidth+1, vHeights["containers"], width-1, height-2)
	if err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		// v.Wrap = gui.Config.UserConfig.Gui.WrapMainPanel
		v.FgColor = gocui.ColorDefault
	}

	if v, err := g.SetView("services", 0, 0, leftSideWidth, vHeights["containers"]-1); err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		v.Title = "services"
	}

	if v, err := g.SetView("containers", leftSideWidth+1, 0, leftSideWidth*2, vHeights["containers"]-1); err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		v.Title = "containers"
	}

	if v, err := g.SetView("images", leftSideWidth*2+1, 0, width-1, vHeights["containers"]-1); err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		v.Title = "images"
	}

	return nil
}
