package widgets

import (
	"sync"

	"strings"

	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/mitchellh/colorstring"
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
	if w.changed || w.d.Changed() {
		w.Window.View.Write(w.formatData())
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

func (w *DataWindow) formatData() []byte {
	rows, cols := w.d.GetSize()

	colMaxWidth := []int{}
	for c := 0; c < cols; c++ {
		colData, _ := w.d.GetColumn(c)
		max := int(0)
		for _, d := range colData {
			if len(d.Value) > max {
				max = len(d.Value)
			}
		}
		colMaxWidth = append(colMaxWidth, max+2)
	}

	line := ""

	for r := 0; r < rows; r++ {
		rowData, _ := w.d.GetRow(r)
		log.Debugln(rowData)
		for idx, d := range rowData {
			tl := ""

			tl += " " + d.Value
			maxWidth := colMaxWidth[idx]
			if len(d.Value) < maxWidth {
				tl += strings.Repeat(" ", maxWidth-(len(d.Value)+1))
			}
			if r == 0 {
				tl = colorstring.Color(fmt.Sprintf("%s%s%s", d.BgColor, d.FgColor, d.Attr) + tl)
			}
			line += tl
		}
		log.Infoln(line)
		line += "\n"
	}

	return []byte(line)

}
