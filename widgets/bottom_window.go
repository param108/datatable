package widgets

import (
	"context"
	"strconv"
	"sync"

	"strings"

	"github.com/jroimartin/gocui"
	"github.com/param108/datatable/keybindings"
	"github.com/param108/datatable/messages"
	log "github.com/sirupsen/logrus"
)

const (
	editValueMode = "edit_value"
	saveAsMode    = "save_as"
)

type BottomWindow struct {
	*Window
	sendEvt   chan *messages.Message
	rdEvt     chan *messages.Message
	currDataX int
	currDataY int
	mode      string
	ctx       context.Context
	wg        sync.WaitGroup
	ks        *keybindings.KeyStore
}

func (w *BottomWindow) Wait() {
	w.wg.Wait()
}

func (w *BottomWindow) SetFocus() error {
	w.Window.G.Cursor = true
	if _, err := w.Window.G.SetCurrentView(w.Window.GetView().Name()); err != nil {
		log.Errorf("bottomWindow: Failed to set view %v", err)
		return err
	}
	return nil
}

func (w *BottomWindow) EventHandler() {
	for {
		select {
		case msg := <-w.rdEvt:
			log.Infof("BottomWindow: Message %s", msg.Key)
			switch msg.Key {
			case messages.SetEditModeMsg:
				// Its edit mode now, extract the value and show it
				w.G.Update(func(g *gocui.Gui) error {
					w.Window.View.Clear()
					w.Window.View.SetCursor(0, 0)
					log.Infof("New Edit: %s", string(msg.Data["value"]))
					w.Window.View.Write([]byte(msg.Data["value"]))
					w.currDataX, _ = strconv.Atoi(msg.Data["X"])
					w.currDataY, _ = strconv.Atoi(msg.Data["Y"])
					w.Window.View.SetCursor(len(msg.Data["value"]), 0)
					w.mode = editValueMode
					g.SetCurrentView(w.Window.Name)
					g.Cursor = true
					return nil
				})
			case messages.SetSaveAsModeMsg:
				// Its edit mode now, extract the value and show it
				w.G.Update(func(g *gocui.Gui) error {
					w.Window.View.Clear()
					w.Window.View.SetCursor(0, 0)
					w.Window.View.Write([]byte(msg.Data["value"]))
					w.Window.View.SetCursor(len(msg.Data["value"]), 0)
					w.mode = saveAsMode
					g.SetCurrentView(w.Window.Name)
					g.Cursor = true
					return nil
				})

			}
		case <-w.ctx.Done():
			log.Infof("Exitting Bottom Window")
			return
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

	if key == gocui.KeyArrowLeft {
		x, y := v.Cursor()
		x = x - 1
		if x >= 0 {
			v.SetCursor(x, y)
		}
		return
	}

	if key == gocui.KeyArrowRight {
		x, y := v.Cursor()
		x = x + 1
		if x < len(v.ViewBuffer()) {
			v.SetCursor(x, y)
		}
		return
	}

	if key == gocui.KeyEsc {
		msg := &messages.Message{
			Key: messages.SetExploreModeMsg,
		}
		w.sendEvt <- msg
		w.Window.View.Clear()
		w.Window.View.SetCursor(0, 0)
		return
	}

	if key == gocui.KeyEnter {
		log.Infof("keyEnter %s", strings.TrimSpace(v.Buffer()))
		var msg *messages.Message

		switch w.mode {
		case editValueMode:
			msg = &messages.Message{
				Key: messages.UpdateValueMsg,
				Data: map[string]string{
					"value": strings.TrimSpace(v.Buffer()),
					"X":     strconv.Itoa(w.currDataX),
					"Y":     strconv.Itoa(w.currDataY),
				},
			}
		case saveAsMode:
			msg = &messages.Message{
				Key: messages.SaveAsMsg,
				Data: map[string]string{
					"value": strings.TrimSpace(v.Buffer()),
				},
			}
		default:
			log.Errorf("bottom_window: Unknown mode %s", w.mode)
			return
		}
		w.sendEvt <- msg
		w.View.Clear()
		w.View.SetCursor(0, 0)
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

	w.View.FgColor = gocui.ColorGreen
	w.View.SelBgColor = gocui.ColorGreen

	//w.View.SetCursor(0, 0)
}

func NewBottomWindow(ctx context.Context, g *gocui.Gui, name string, ks *keybindings.KeyStore, cltRd, cltWr chan *messages.Message) *BottomWindow {
	w := &BottomWindow{Window: &Window{Name: name, G: g}, sendEvt: cltWr, rdEvt: cltRd, ks: ks}
	w.ctx = ctx
	w.Layout()

	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		w.EventHandler()
	}()

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

func (w *BottomWindow) CreateEscKeyHandler() func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		log.Infof("Bottom: Esc Key pressed")
		return nil
	}
}

func (w *BottomWindow) SetKeys() error {
	/*err := w.ks.AddKey(w.Window.Name, gocui.KeyEsc, gocui.ModNone, w.CreateEscKeyHandler())
	if err != nil {
		log.Errorf("Failed to add key handler %+v", err)
		return err
	}*/

	return nil
}
