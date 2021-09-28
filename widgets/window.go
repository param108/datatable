package widgets

import (
	"github.com/jroimartin/gocui"
)

type Window struct {
	Name string
	MinX int
	MaxX int
	MinY int
	MaxY int
	View *gocui.View
	G    *gocui.Gui
}

func (w *Window) SetView() error {
	v, err := w.G.SetView(w.Name, w.MinX, w.MinY, w.MaxX, w.MaxY)
	if err != gocui.ErrUnknownView {
		return err
	}
	w.View = v
	return nil
}

func (w *Window) GetName() string {
	return w.Name
}

func (w *Window) Refresh(g *gocui.Gui) {
}
