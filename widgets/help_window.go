package widgets

import (
	"context"

	"github.com/jroimartin/gocui"
	"github.com/mitchellh/colorstring"
	"github.com/param108/datatable/messages"
	log "github.com/sirupsen/logrus"
)

type HelpWindow struct {
	*Window
	sendEvt chan *messages.Message
	rdEvt   chan *messages.Message
	ctx     context.Context
}

func (w *HelpWindow) Wait() {
	// nothing required for help
}

func (w *HelpWindow) SetFocus() error {
	w.Window.G.Cursor = false
	w.Window.G.SetViewOnTop(w.Window.GetView().Name())
	if _, err := w.Window.G.SetCurrentView(w.Window.GetView().Name()); err != nil {
		log.Errorf("bottomWindow: Failed to set view %v", err)
		return err
	}
	return nil
}

func (w *HelpWindow) Layout() {
	maxX, maxY := w.G.Size()
	w.MinX = maxX / 6
	w.MinY = maxY / 6
	w.MaxX = (maxX * 5) / 6
	w.MaxY = (maxY * 5) / 6
	log.WithFields(log.Fields{
		"MinX": w.MinX,
		"MinY": w.MinY,
		"MaxX": w.MaxX,
		"MaxY": w.MaxY,
	}).Debugf("HelpWindow: Layout")
}

func (w *HelpWindow) EventHandler() {
	for msg := range w.rdEvt {
		log.Infof("DataWindow: %s", msg.Key)

		switch msg.Key {
		case messages.UpdateValueMsg:
			// Its edit mode now, extract the value and show it
			log.Infof("Update value seen")
		}
	}
}

func (w *HelpWindow) Animate(g *gocui.Gui) error {
	return nil
}

func (w *HelpWindow) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	log.Infof("Edit: %s %d", string(ch), key)

	switch ch {
	case 'q':
		msg := &messages.Message{
			Key: messages.CloseHelpWindow,
		}
		w.sendEvt <- msg
		return
	}

	return
}

func (w *HelpWindow) CustomSetup() {
	w.View.Editor = w
	w.View.Editable = true
	w.View.Clear()
	w.View.Write([]byte(colorstring.Color(`
    [underline]DataTable[reset]

    Choose the cell you want to edit using the arrow keys.
    Selected cell looks like [invert]this[reset].
    Once you find the cell, use the letter [invert]e[reset] to edit.

    The value of the cell will be seen in the bottom window.
    Use the alphanumeric keys to edit. [invert]Backspace[reset] to delete.
    Press [invert]Enter[reset] to finish editting.
          [invert]Esc[reset] to cancel edit.

    After editting, press [invert]s[reset] to save.

    press [invert]w[reset] to save as.
       filename will appear in the bottom window for editting. Enter to save.

    press [invert]a[reset] to add a column.

    Press [invert]q[reset] to close this.
    Press [invert]ctrl-h[reset] to see this again.
`)))
}

func (w *HelpWindow) SetKeys() error {
	return nil
}

func NewHelpWindow(ctx context.Context, g *gocui.Gui, name string, rdEvt, sendEvt chan *messages.Message) *HelpWindow {
	w := &HelpWindow{Window: &Window{Name: name, G: g}}
	w.Layout()
	w.sendEvt = sendEvt
	w.rdEvt = rdEvt
	w.ctx = ctx
	return w
}
