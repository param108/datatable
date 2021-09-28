package widgets

import (
	"github.com/jroimartin/gocui"
	"github.com/param108/datatable/data"
	log "github.com/sirupsen/logrus"
	"sync"
)

type DataWindow struct {
	*Window
	d       data.DataSource
	mx      sync.Mutex
	changed bool
}

func NewDataWindow(g *gocui.Gui, name string, d data.DataSource) *DataWindow {
	w := &DataWindow{Window: &Window{Name: name, G: g}, d: d}
	w.Layout()
	w.changed = true
	return w
}

func (w *DataWindow) Animate(g *gocui.Gui) error {
	w.mx.Lock()
	defer w.mx.Unlock()
	if w.changed {
		w.Window.View.Write([]byte("Hello"))
		w.changed = false
	}
	return nil
}

func (w *DataWindow) Layout() {

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
