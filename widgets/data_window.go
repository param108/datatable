package widgets

import (
	"sync"

	"strings"

	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/mitchellh/colorstring"
	"github.com/param108/datatable/data"
	"github.com/param108/datatable/keybindings"
	log "github.com/sirupsen/logrus"
)

const (
	arrowUp    = "up"
	arrowDown  = "down"
	arrowLeft  = "left"
	arrowRight = "right"
)

type Cursor struct {
	X int
	Y int
}

type DataWindow struct {
	*Window
	d       data.DataSource
	mx      sync.Mutex
	changed bool
	crsr    *Cursor
	ks      *keybindings.KeyStore
}

func NewDataWindow(g *gocui.Gui, name string, d data.DataSource, ks *keybindings.KeyStore) *DataWindow {
	w := &DataWindow{Window: &Window{Name: name, G: g}, d: d}
	w.Layout()
	w.changed = true
	w.crsr = &Cursor{X: 0, Y: 1}
	w.ks = ks
	return w
}
func (w *DataWindow) SetKeys() error {
	err := w.ks.AddKey(w.Window.Name, gocui.KeyArrowUp, gocui.ModNone, w.CreateArrowKeyHandler(arrowUp))
	if err != nil {
		log.Errorf("Failed to add key handler %+v", err)
		return err
	}

	err = w.ks.AddKey(w.Window.Name, gocui.KeyArrowDown, gocui.ModNone, w.CreateArrowKeyHandler(arrowDown))
	if err != nil {
		log.Errorf("Failed to add key handler %+v", err)
		return err
	}

	err = w.ks.AddKey(w.Window.Name, gocui.KeyArrowRight, gocui.ModNone, w.CreateArrowKeyHandler(arrowRight))
	if err != nil {
		log.Errorf("Failed to add key handler %+v", err)
		return err
	}

	err = w.ks.AddKey(w.Window.Name, gocui.KeyArrowLeft, gocui.ModNone, w.CreateArrowKeyHandler(arrowLeft))
	if err != nil {
		log.Errorf("Failed to add key handler %+v", err)
		return err
	}
	return nil
}

func (w *DataWindow) CreateArrowKeyHandler(dir string) func(*gocui.Gui, *gocui.View) error {
	log.Infof("CreateArrowKeyHandler Called %s", dir)

	return func(g *gocui.Gui, v *gocui.View) error {
		log.Infof("ArrowKey Called %s", dir)
		changed := false
		sizeY, sizeX := w.d.GetSize()
		switch dir {
		case arrowUp:
			newY := w.crsr.Y - 1
			// cant go above the header
			if newY > 0 {
				w.crsr.Y = newY
				changed = true
			}
		case arrowDown:
			newY := w.crsr.Y + 1
			if newY < sizeY {
				w.crsr.Y = newY
				changed = true
			}
		case arrowRight:
			newX := w.crsr.X + 1
			if newX < sizeX {
				w.crsr.X = newX
				changed = true
			}
		case arrowLeft:
			newX := w.crsr.X - 1

			if newX >= 0 {
				w.crsr.X = newX
				changed = true
			}
		}

		if changed {
			w.mx.Lock()
			defer w.mx.Unlock()
			w.changed = true
		}
		return nil
	}
}

func (w *DataWindow) Animate(g *gocui.Gui) error {
	w.mx.Lock()
	defer w.mx.Unlock()
	if w.changed || w.d.Changed() {
		w.Window.View.Clear()
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
			} else if idx == w.crsr.X && r == w.crsr.Y {
				tl = colorstring.Color("[invert]" + tl)
			}
			line += tl
		}
		log.Debugln(line)
		line += "\n"
	}

	return []byte(line)

}
