package widgets

import (
	"github.com/jroimartin/gocui"
)

type Widget interface {
	Animate(g *gocui.Gui)
	SetView() error
	GetName() string
	Layout()
}
