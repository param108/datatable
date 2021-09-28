package widgets

import (
	"context"
	"sync"

	"strings"

	"github.com/jroimartin/gocui"
	"github.com/param108/datatable/data"
	log "github.com/sirupsen/logrus"
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
		//w.Window.View.MoveCursor(0, 0, false)
		w.formatData()
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
	}).Debugf("TopWindow: Layout")
}

func (w *DataWindow) formatData() {
	ctx := context.Background()
	rows, cols := w.d.GetSize(ctx)

	colMaxWidth := []int{}
	for c := 0; c < cols; c++ {
		colData, _ := w.d.GetColumn(context.Background(), c)
		max := int(0)
		for _, d := range colData {
			if len(d) > max {
				max = len(d)
			}
		}
		colMaxWidth = append(colMaxWidth, max+2)
	}

	line := ""

	for r := 0; r < rows; r++ {
		rowData, _ := w.d.GetRow(context.Background(), r)
		log.Debugln(rowData)
		for idx, d := range rowData {
			line += " " + d
			maxWidth := colMaxWidth[idx]
			if len(d) < maxWidth {
				line += strings.Repeat(" ", maxWidth-(len(d)+1))
			}
		}
		log.Debugln(line)
		line += "\n"
	}
	w.Window.View.Write([]byte(line))

}
