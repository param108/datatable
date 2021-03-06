package widgets

import (
	"github.com/jroimartin/gocui"
)

type Widget interface {
	Animate(g *gocui.Gui) error
	SetView() error
	GetView() *gocui.View
	GetName() string
	Layout()
	CustomSetup()
	SetKeys() error
	SetFocus() error
	// Cancel the context and call wait for shutdown
	Wait()
}
