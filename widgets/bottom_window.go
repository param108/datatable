package widgets

import (
	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
)

type BottomWindow struct {
	*Window
}

func NewBottomWindow(g *gocui.Gui, name string) *BottomWindow {
	w := &BottomWindow{&Window{Name: name, G: g}}
	w.Layout()
	return w
}

func (w *BottomWindow) Layout() {
	maxX, maxY := w.G.Size()
	w.MinX = 0
	w.MaxX = maxX - 1
	w.MinY = ((maxY * 9) / 10) + 1
	w.MaxY = maxY - 1
	log.WithFields(log.Fields{
		"MinX": w.MinX,
		"MinY": w.MinY,
		"MaxX": w.MaxX,
		"MaxY": w.MaxY,
	}).Infof("BottomWindow: Layout")

}

func (w *BottomWindow) Animate(g *gocui.Gui) error {
	return nil
}
