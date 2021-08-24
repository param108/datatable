package main

import (
	"github.com/jroimartin/gocui"
	"github.com/param108/datatable/widgets"
	"log"
	mylog "github.com/param108/datatable/log"
	"github.com/sirupsen/logrus"

)

var (
	TheUI *UI
)

func setup(g *gocui.Gui) {
	CreateUI(g)
	TheUI.AddWidget(widgets.NewTopWindow(g, "Top"))
	TheUI.AddWidget(widgets.NewBottomWindow(g, "Bottom"))
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	defer mylog.Close()

	defer g.Close()

	logrus.Infof("Created gui %p", g)

	setup(g)

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

type UI struct {
	W map[string]widgets.Widget
	G *gocui.Gui
}

func CreateUI(g *gocui.Gui) *UI {
	TheUI = &UI{
		W: map[string]widgets.Widget{},
		G: g,
	}
	return TheUI
}

func (ui *UI) AddWidget(w widgets.Widget) {
	ui.W[w.GetName()]= w
}

func layout(g *gocui.Gui) error {
	for _, w := range (TheUI.W) {
		logrus.Infof("Layout for view %s %p", w.GetName(), g)
		w.Layout()
		if err := w.SetView(); err != nil {
			logrus.Errorf("Failed to setview %+v", err)
			return err
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

