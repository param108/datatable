package cmd

import (
	"testing"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/stretchr/testify/assert"
)

func testUIRun(g *gocui.Gui, done chan bool) {
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		done <- false
		return
	}
	done <- true
}

func TestStartUpHappens(t *testing.T) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	assert.Nil(t, err, "Failed to create gocui")

	ui, err := CreateUI(g, "testdata/data.csv")
	assert.Nil(t, err, "Failed to create ui")

	uiDone := make(chan bool)
	go testUIRun(g, uiDone)

	time.Sleep(2 * time.Second)

	t.Run("gui has 4 windows", func(t *testing.T) {
		assert.Equal(t, 4, len(g.Views()), "incorrect number of windows")
		assert.Equal(t, "Help", g.CurrentView().Name(), "incorrect current view")
		for i, v := range g.Views() {
			if v.Name() == "Help" {
				assert.Equal(t, 3, i, "help window not on top")
			}
		}
	})

	g.Update(func(g *gocui.Gui) error {
		return gocui.ErrQuit
	})

	ui.Quit()
	g.Close()
	// should exit within 2 seconds
	timeout := time.NewTimer(2 * time.Second)
	select {
	case ret := <-uiDone:
		assert.True(t, ret, "Error in shutdown")
	case <-timeout.C:
		timeout.Stop()
		assert.True(t, false, "Timedout waiting for exit")
		// Need to wait for proper shutdown
		<-uiDone
	}

}
