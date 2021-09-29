package widgets

import (
	"github.com/jroimartin/gocui"
)

type Widget interface {
	Animate(g *gocui.Gui) error
	SetView() error
	GetName() string
	Layout()
	SetKeys() error
}
