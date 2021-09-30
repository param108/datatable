package widgets

import (
	"github.com/jroimartin/gocui"
	"github.com/param108/datatable/messages"
	log "github.com/sirupsen/logrus"
)

type BottomWindow struct {
	*Window
	sendEvt chan *messages.Message
	rdEvt   chan *messages.Message
}

func (w *BottomWindow) EventHandler() {
	for msg := range w.rdEvt {
		switch msg.Key {
		case messages.SetEditModeMsg:
			// Its edit mode now, extract the value and show it
			w.G.Update(func(g *gocui.Gui) error {
				w.Window.View.Clear()
				w.Window.View.Write([]byte(msg.Data["value"]))
				return nil
			})
		}
	}
}

func (w *BottomWindow) CustomSetup() {

}

func NewBottomWindow(g *gocui.Gui, name string, cltRd, cltWr chan *messages.Message) *BottomWindow {
	w := &BottomWindow{Window: &Window{Name: name, G: g}, sendEvt: cltWr, rdEvt: cltRd}
	w.Layout()
	go w.EventHandler()
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
	}).Debugf("BottomWindow: Layout")

}

func (w *BottomWindow) Animate(g *gocui.Gui) error {
	return nil
}

func (w *BottomWindow) SetKeys() error {
	return nil
}
