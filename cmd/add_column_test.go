package cmd

import (
	"testing"
	"time"

	"strings"

	"github.com/jroimartin/gocui"
	"github.com/param108/datatable/widgets"
	"github.com/stretchr/testify/assert"
)

func testAddColumn(t *testing.T, ui *UI) {
	g := ui.G
	g.CurrentView().Editor.Edit(g.CurrentView(), 0, 'a', 0)
	time.Sleep(time.Second)

	assert.Equal(t, "Bottom", g.CurrentView().Name(), "Invalid CurrentView")

	assert.Equal(t, "new_column", strings.TrimSpace(ui.W["Bottom"].GetView().Buffer()),
		"Invalid column input")

	g.CurrentView().Editor.Edit(g.CurrentView(), gocui.KeyEnter, 0, 0)
	time.Sleep(time.Second)

	assert.Equal(t, "Data", g.CurrentView().Name(), "Invalid CurrentView")

	dataWindow := ui.W["Data"].(*widgets.DataWindow)
	d := dataWindow.GetData()
	assert.Equal(t, 6, len(d.GetColumns()))
	assert.Equal(t, "new_column", d.GetColumns()[5].Name)
}
