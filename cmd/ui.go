package cmd

import (
	"github.com/jroimartin/gocui"
	"github.com/param108/datatable/data"
	"github.com/param108/datatable/eventbus"
	"github.com/param108/datatable/keybindings"
	"github.com/param108/datatable/messages"
	"github.com/param108/datatable/widgets"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"time"
)

type UI struct {
	W     map[string]widgets.Widget
	G     *gocui.Gui
	D     data.DataSource
	KS    *keybindings.KeyStore
	CV    *gocui.View
	EB    *eventbus.EventBus
	F     string
	Stack []string
	quit  chan int
}

// push a view name on stack
func (ui *UI) Push(v string) {
	ui.Stack = append(ui.Stack, v)
	logrus.Infof("Pushing %s", v)

}

func (ui *UI) Pop(def string) string {
	if len(ui.Stack) == 0 {
		logrus.Infof("Popping %s", def)

		return def
	}
	v := ui.Stack[0]
	logrus.Infof("Popping %s", def)

	ui.Stack = ui.Stack[1:]
	return v
}

func (ui *UI) CentralCommand(CNCrd, CNCwr chan *messages.Message) {
	for msg := range CNCrd {
		logrus.Infof("CNC: message %s", msg.Key)

		switch msg.Key {
		case messages.SetEditModeMsg:
			ui.G.Update(func(g *gocui.Gui) error {
				err := ui.W["Bottom"].SetFocus()
				if err != nil {
					logrus.Errorf("CNC: Failed to set view: Bottom")
					return err
				}
				ui.CV = ui.W["Bottom"].GetView()
				return nil
			})
		case messages.UpdateValueMsg:
			ui.G.Update(func(g *gocui.Gui) error {
				err := ui.W["Data"].SetFocus()
				if err != nil {
					logrus.Errorf("CNC: Failed to set view: Bottom")
					return err
				}
				ui.CV = ui.W["Data"].GetView()
				return nil
			})
		case messages.SetSaveAsModeMsg:
			ui.G.Update(func(g *gocui.Gui) error {
				err := ui.W["Bottom"].SetFocus()
				if err != nil {
					logrus.Errorf("CNC: Failed to set view: Bottom")
					return err
				}
				ui.CV = ui.W["Bottom"].GetView()
				return nil
			})
		case messages.SaveAsMsg:
			ui.G.Update(func(g *gocui.Gui) error {
				err := ui.W["Data"].SetFocus()
				if err != nil {
					logrus.Errorf("CNC: Failed to set view: Bottom")
					return err
				}
				ui.CV = ui.W["Data"].GetView()
				return nil
			})
		case messages.ShowHelpWindow:
			ui.G.Update(func(g *gocui.Gui) error {
				ui.Push(ui.CV.Name())
				ui.W["Help"].SetFocus()
				ui.CV = ui.W["Help"].GetView()
				return nil
			})
		case messages.CloseHelpWindow:
			ui.G.Update(func(g *gocui.Gui) error {
				v := ui.Pop("Data")
				g.SetViewOnBottom("Help")
				err := ui.W[v].SetFocus()
				if err != nil {
					logrus.Errorf("CNC: Failed to set view: Bottom")
					return err
				}
				ui.CV = ui.W[v].GetView()
				return nil
			})
		case messages.ShowToastMsg:
			ui.G.Update(func(g *gocui.Gui) error {
				v := ui.CV.Name()
				ui.Push(v)
				err := ui.W["Toast"].SetFocus()
				if err != nil {
					logrus.Errorf("CNC: Failed to set view: Bottom")
					return err
				}
				ui.CV = ui.W["Toast"].GetView()
				return nil
			})
		case messages.HideToastMsg:
			ui.G.Update(func(g *gocui.Gui) error {
				v := ui.Pop("Data")
				g.SetViewOnBottom("Toast")
				err := ui.W[v].SetFocus()
				if err != nil {
					logrus.Errorf("CNC: Failed to set view: Bottom")
					return err
				}
				ui.CV = ui.W[v].GetView()
				return nil
			})
		default:
			logrus.Errorf("CNC: invalid message key: %s", msg.Key)
		}
	}
}

func CreateUI(g *gocui.Gui, filename string) (*UI, error) {
	TheUI := &UI{
		W:    map[string]widgets.Widget{},
		G:    g,
		KS:   keybindings.NewKeyStore(g),
		EB:   eventbus.NewEventBus(),
		F:    filename,
		quit: make(chan int),
	}

	CNCrd, CNCwr := TheUI.EB.RegisterWindow()
	go TheUI.CentralCommand(CNCrd, CNCwr)

	src, err := data.NewCSV(filename)
	if err != nil {
		panic(err)
	}

	cltRd, cltWr := TheUI.EB.RegisterWindow()
	TheUI.AddWidget(widgets.NewToastWindow(g, "Toast", cltRd, cltWr))

	cltRd, cltWr = TheUI.EB.RegisterWindow()
	TheUI.AddWidget(widgets.NewDataWindow(g, "Data", src, TheUI.KS, cltRd, cltWr))

	cltRd, cltWr = TheUI.EB.RegisterWindow()
	TheUI.AddWidget(widgets.NewBottomWindow(g, "Bottom", cltRd, cltWr))

	cltRd, cltWr = TheUI.EB.RegisterWindow()
	TheUI.AddWidget(widgets.NewHelpWindow(g, "Help", cltRd, cltWr))

	TheUI.D = src

	g.SetViewOnBottom("Toast")
	g.SetViewOnTop("Help")

	g.SetManagerFunc(TheUI.layout)

	TheUI.KS.AddKey("", gocui.KeyCtrlC, gocui.ModNone, quit)
	TheUI.KS.AddKey("", gocui.KeyCtrlH, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		g.Update(func(g *gocui.Gui) error {
			msg := &messages.Message{
				Key: messages.ShowHelpWindow,
			}
			CNCwr <- msg
			return nil
		})
		return nil
	})

	for _, w := range TheUI.W {
		w.SetKeys()
	}

	go func() {
		TheUI.animate(g)
	}()

	return TheUI, nil
}

func (ui *UI) AddWidget(w widgets.Widget) {
	ui.W[w.GetName()] = w
}

func (ui *UI) layout(g *gocui.Gui) error {
	for _, w := range ui.W {
		logrus.Debugf("Layout for view %s %p", w.GetName(), g)
		w.Layout()
		if err := w.SetView(); err != nil {
			logrus.Errorf("Failed to setview %+v", err)
			return err
		}
		w.CustomSetup()
	}

	if ui.CV == nil {
		v, err := g.SetCurrentView("Help")
		if err != nil {
			panic(err)
		}
		ui.CV = v
		g.SetViewOnTop("Help")
		g.SetViewOnBottom("Toast")
	}

	return nil
}

func (ui *UI) animate(g *gocui.Gui) {
	// Do it once at the beginning
	for _, w := range ui.W {
		g.Update(w.Animate)
	}

	t := time.NewTicker(time.Millisecond * 100)
	for {
		select {
		case <-t.C:
			for _, w := range ui.W {
				g.Update(w.Animate)
			}
		case <-ui.quit:
			return
		}
	}
}

func (ui *UI) Quit() {
	close(ui.quit)
}

func quit(g *gocui.Gui, v *gocui.View) error {
	logrus.Infof("quit called")
	return gocui.ErrQuit
}

var uiCmd = &cli.Command{
	Name:  "ui",
	Usage: "run the ui to manually edit your csv",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "file",
			Aliases:  []string{"f"},
			Required: true,
			Usage:    "csv file to open",
		},
	},
	Action: uiAction,
}

func uiAction(c *cli.Context) error {
	logrus.Infoln("UI ACTION CALLED")
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		logrus.Panicln(err)
	}

	defer g.Close()

	logrus.Infof("Created gui %p", g)

	filename := c.String("file")

	ui, err := CreateUI(g, filename)
	if err != nil {
		logrus.Errorf("create failed: %v", err)
		return errors.Wrap(err, "failed create ui")
	}

	defer ui.Quit()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		logrus.Panicln(err.Error())
	}

	return nil
}

func init() {
	registerCommand(uiCmd)
}
