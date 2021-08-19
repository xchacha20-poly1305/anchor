package app

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/sagernet/sagerconnect/core/gtkx"
	"github.com/sagernet/sagerconnect/res/drawable"
	"github.com/sagernet/sagerconnect/res/layout"
)

var mainActivity *MainActivity

func initMain() {
	mainActivity = &MainActivity{}
	mainActivity.onCreate()
}

type MainActivity struct {
	*gtkx.Layout
	*gtk.ApplicationWindow
}

func (a *MainActivity) onCreate() {
	a.Layout = gtkx.NewLayout(layout.WindowMain)
	a.ApplicationWindow = mainActivity.GetAppWindow("main_activity")
	mainApp.AddWindow(mainActivity)

	a.SetIcon(drawable.MustPixbuf(drawable.IconPng))
	a.GetMenuItem("exit_item").SetOnClickListener(mainApp.Quit)
	a.GetMenuItem("about_item").SetOnClickListener(initAbout)
}
