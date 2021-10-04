package widgets

import (
	"sync"

	"strings"

	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/mitchellh/colorstring"
	"github.com/param108/datatable/data"
	"github.com/param108/datatable/keybindings"
	"github.com/param108/datatable/messages"
	log "github.com/sirupsen/logrus"
	"strconv"
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
	sendEvt chan *messages.Message
	rdEvt   chan *messages.Message
}

func NewDataWindow(g *gocui.Gui, name string, d data.DataSource, ks *keybindings.KeyStore,
	rdEvt, sendEvt chan *messages.Message) *DataWindow {
	w := &DataWindow{Window: &Window{Name: name, G: g}, d: d}
	w.Layout()
	w.changed = true
	w.crsr = &Cursor{X: 0, Y: 1}
	w.ks = ks
	w.sendEvt = sendEvt
	w.rdEvt = rdEvt

	go w.EventHandler()

	return w
}

func (w *DataWindow) EventHandler() {
	for msg := range w.rdEvt {
		log.Infof("DataWindow: %s", msg.Key)

		switch msg.Key {
		case messages.UpdateValueMsg:
			// Its edit mode now, extract the value and show it
			w.G.Update(func(g *gocui.Gui) error {
				log.Infof("DataWindow: updateValue %s %s %s", msg.Data["X"], msg.Data["Y"])
				value := msg.Data["value"]
				X, _ := strconv.Atoi(msg.Data["X"])
				Y, _ := strconv.Atoi(msg.Data["Y"])
				err := w.d.Set(Y, X, value)
				if err != nil {
					log.Errorf("DataWindow: failed set %v", err)
				}
				return nil
			})
		case messages.SaveAsMsg:
			if err := w.d.SaveAs(msg.Data["value"]); err != nil {
				// FIXME: Add a toast to notify the user
				log.Errorf("data_window: Failed to save as %v", err)
			}
		}
	}
}

func (w *DataWindow) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	log.Infof("Edit: %s", string(ch))

	switch ch {
	case 'e':
		value, err := w.d.Get(w.crsr.Y, w.crsr.X)
		if err != nil {
			log.Errorf("invalid ordinates col: %d row: %d", w.crsr.Y, w.crsr.X)
			return
		}

		msg := &messages.Message{
			Key: messages.SetEditModeMsg,
			Data: map[string]string{
				"value": value.(string),
				"X":     strconv.Itoa(w.crsr.X),
				"Y":     strconv.Itoa(w.crsr.Y),
			},
		}
		w.sendEvt <- msg
	case 's':
		err := w.d.Save()
		if err != nil {
			log.Errorf("failed to save %v", err)
		}
	case 'w':
		msg := &messages.Message{
			Key: messages.SetSaveAsModeMsg,
			Data: map[string]string{
				"value": w.d.Source(),
			},
		}
		w.sendEvt <- msg

	}
	return
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
		w.d.ClearChanged()
	}
	return nil
}

func (w *DataWindow) CustomSetup() {
	w.View.Editor = w
	w.View.Editable = true
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
		log.Infoln(line)
		line += "\n"
	}

	return []byte(line)

}
