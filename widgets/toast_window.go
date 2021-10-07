package widgets

import (
	"context"
	"time"

	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/param108/datatable/messages"
	log "github.com/sirupsen/logrus"
)

type ToastWindow struct {
	*Window
	sendEvt, rdEvt chan *messages.Message
	Msg            string
	ctx            context.Context
}

func (w *ToastWindow) SetFocus() error {
	w.Window.G.Cursor = false
	w.Window.G.SetViewOnTop("Toast")
	go func(g *gocui.Gui) {
		<-time.NewTimer(time.Second * 1).C
		msg := &messages.Message{
			Key: messages.HideToastMsg,
		}
		w.sendEvt <- msg
	}(w.G)
	return nil
}

func (w *ToastWindow) Layout() {
	maxX, maxY := w.G.Size()
	w.MinX = maxX / 6
	w.MinY = maxY / 6
	w.MaxX = (maxX * 5) / 6
	w.MaxY = (maxY * 2) / 6
	log.WithFields(log.Fields{
		"MinX": w.MinX,
		"MinY": w.MinY,
		"MaxX": w.MaxX,
		"MaxY": w.MaxY,
	}).Debugf("ToastWindow: Layout")
}

func (w *ToastWindow) Animate(g *gocui.Gui) error {
	w.GetView().Clear()
	w.GetView().Write([]byte(w.Msg))
	return nil
}

func (w *ToastWindow) CustomSetup() {
}

func (w *ToastWindow) EventHandler() {
	for msg := range w.rdEvt {
		log.Infof("DataWindow: %s", msg.Key)

		switch msg.Key {
		case messages.ShowToastMsg:
			w.Msg = msg.Data["msg"]
			w.G.Update(func(g *gocui.Gui) error {
				// show toast for 5 seconds
				return nil
			})
			// Its edit mode now, extract the value and show it
			log.Infof(fmt.Sprintf("show toast: %s", msg.Data["msg"]))
		}
	}
}

func NewToastWindow(ctx context.Context, g *gocui.Gui, name string, rdEvt, sendEvt chan *messages.Message) *ToastWindow {
	w := &ToastWindow{Window: &Window{Name: name, G: g}}
	w.Layout()
	w.sendEvt = sendEvt
	w.rdEvt = rdEvt
	w.ctx = ctx
	go w.EventHandler()

	return w
}

func (w *ToastWindow) SetKeys() error {
	return nil
}
