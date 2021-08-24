package main

import (
	"github.com/jroimartin/gocui"
	"log"
)

var (
	TheUI *UI
)

func setup(g *gocui.Gui) {
	CreateUI(g)
	TheUI.AddWidget(NewTopWindow(g, "Top"))
	TheUI.AddWidget(NewBottomWindow(g, "Bottom"))
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	defer g.Close()

	CreateUI(g)

	TheUI
	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

type UI struct {
	W map[string]Widget
	G *gocui.Gui
}

func CreateUI(g *gocui.Gui) *UI {
	TheUI = &UI{
		W: map[string]Widget{},
		G: g,
	}
}

func (ui *UI) AddWidget(w Widget) {
	ui.W[w.GetName()]= w
}

type Window struct {
	Name string
	MinX int
	MaxX int
	MinY int
	MaxY int
	View *gocui.View
	G *gocui.Gui
}

type Widget interface {
	Animate(g *gocui.Gui)
	SetView(v *gocui.View)
	GetName() string
}

type TopWindow struct {
	*Window
}

type BottomWindow struct {
	*Window
}

func NewTopWindow(g *gocui.Gui, name string) (*TopWindow, error) {
	maxX, maxY := g.Size()
	w := &TopWindow {&Window{name, 0, maxX-1, 0, maxY/10, nil, g}}
	return w
}

func NewBottomWindow(g *gocui.Gui, name string) (*BottomWindow, error) {
	maxX, maxY := g.Size()
	w := &BottomWindow {&Window{name, 0, maxX-1, maxY/10 + 1, maxY - 1, nil, g}}
	return w
}

func (w *TopWindow) Animate(g *gocui.Gui) {
}

func (w *BottomWindow) Animate(g *gocui.Gui) {
}

func layout(g *gocui.Gui) error {
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

