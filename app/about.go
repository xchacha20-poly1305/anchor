package app

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/sagernet/sagerconnect/core/gtkx"
	"github.com/sagernet/sagerconnect/res/layout"
)

type AboutDialog struct {
	*gtkx.Layout
	*gtk.AboutDialog
}

func initAbout() {
	about := &AboutDialog{}
	about.onCreate()
}

func (d *AboutDialog) onCreate() {
	d.Layout = gtkx.NewLayout(layout.DialogAbout)
	d.AboutDialog = d.GetAboutDialog("about_dialog")
}
