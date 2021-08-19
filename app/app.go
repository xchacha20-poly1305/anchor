package app

import (
	_ "embed"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/sagernet/sagerconnect/core/gtkx"
	"os"
)

var mainApp *SagerApplication

type SagerApplication struct {
	*gtk.Application
}

func (a *SagerApplication) onCreate() {
	initMain()
}

func Launch() {
	mainApp = &SagerApplication{}
	app, err := gtk.ApplicationNew("io.nekohasekai.sagerconnect", glib.APPLICATION_FLAGS_NONE)
	gtkx.Must("create main window", err)
	mainApp.Application = app
	mainApp.Connect("activate", mainApp.onCreate)
	mainApp.Run(os.Args)
}
