package widgets

import (
	"strconv"

	"github.com/jroimartin/gocui"
	"github.com/param108/datatable/messages"
	log "github.com/sirupsen/logrus"
)

type BottomWindow struct {
	*Window
	sendEvt   chan *messages.Message
	rdEvt     chan *messages.Message
	currDataX int
	currDataY int
}

func (w *BottomWindow) EventHandler() {
	for msg := range w.rdEvt {
		switch msg.Key {
		case messages.SetEditModeMsg:
			// Its edit mode now, extract the value and show it
			w.G.Update(func(g *gocui.Gui) error {
				w.Window.View.Clear()
				w.Window.View.Write([]byte(msg.Data["value"]))
				w.currDataX, _ = strconv.Atoi(msg.Data["X"])
				w.currDataY, _ = strconv.Atoi(msg.Data["Y"])
				return nil
			})
		}
	}
}

func (w *BottomWindow) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	log.Infof("BottomWindow Edit: %s %d %d", string(ch), key, mod)

	if key == gocui.KeyBackspace || key == gocui.KeyBackspace2 {
		log.Infof("backspace")
		v.EditDelete(true)
		return
	}

	if key == gocui.KeyEnter {
		msg := &messages.Message{
			Key: messages.UpdateValueMsg,
			Data: map[string]string{
				"value": v.Buffer(),
				"X":     strconv.Itoa(w.currDataX),
				"Y":     strconv.Itoa(w.currDataY),
			},
		}
		w.sendEvt <- msg
		return
	}
	if ch == 0 {
		return
	}

	v.EditWrite(ch)
}

func (w *BottomWindow) CustomSetup() {
	w.View.Editor = w
	w.View.Editable = true
	w.View.Overwrite = true
	w.View.SetCursor(0, 0)
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
