package widgets

import (
	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
)

type TopWindow struct {
	*Window
}

func NewTopWindow(g *gocui.Gui, name string) *TopWindow {
	w := &TopWindow{&Window{Name: name, G: g}}
	w.Layout()
	return w
}

func (w *TopWindow) Animate(g *gocui.Gui) {
}

func (w *TopWindow) Layout() {

	maxX, maxY := w.G.Size()
	w.MinX = 0
	w.MinY = 0
	w.MaxX = maxX - 1
	w.MaxY = (maxY * 9) / 10
	log.WithFields(log.Fields{
		"MinX": w.MinX,
		"MinY": w.MinY,
		"MaxX": w.MaxX,
		"MaxY": w.MaxY,
	}).Infof("TopWindow: Layout")
}
