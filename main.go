package main

import (
	"log"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/param108/datatable/data"
	mylog "github.com/param108/datatable/log"
	"github.com/param108/datatable/widgets"
	"github.com/sirupsen/logrus"
	"sync"
)

var (
	TheUI     *UI
	TheSource data.DataSource
)

func setup(g *gocui.Gui, filename string) {
	if src, err := data.NewCSV(filename); err != nil {
		panic(err)
	} else {
		TheSource = src
	}

	CreateUI(g)
	TheUI.AddWidget(widgets.NewDataWindow(g, "Data", TheSource))
	TheUI.AddWidget(widgets.NewBottomWindow(g, "Bottom"))
	TheUI.D = TheSource

}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	defer mylog.Close()

	defer g.Close()

	logrus.Infof("Created gui %p", g)

	setup(g, "data.csv")

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	var wg sync.WaitGroup

	quit := make(chan int)

	go func() {
		wg.Add(1)
		defer wg.Done()
		animate(g, quit)
	}()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	close(quit)
	wg.Wait()

}

type UI struct {
	W map[string]widgets.Widget
	G *gocui.Gui
	D data.DataSource
}

func CreateUI(g *gocui.Gui) *UI {
	TheUI = &UI{
		W: map[string]widgets.Widget{},
		G: g,
	}
	return TheUI
}

func (ui *UI) AddWidget(w widgets.Widget) {
	ui.W[w.GetName()] = w
}

func layout(g *gocui.Gui) error {
	for _, w := range TheUI.W {
		logrus.Infof("Layout for view %s %p", w.GetName(), g)
		w.Layout()
		if err := w.SetView(); err != nil {
			logrus.Errorf("Failed to setview %+v", err)
			return err
		}
	}
	return nil
}

func animate(g *gocui.Gui, quit chan int) {
	t := time.NewTicker(time.Second)
	for {
		select {
		case <-t.C:
			for _, w := range TheUI.W {
				g.Update(w.Animate)
			}
		case <-quit:
			return
		}
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
