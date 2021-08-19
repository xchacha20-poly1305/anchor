package drawable

import (
	_ "embed"
	"github.com/gotk3/gotk3/gdk"
	"github.com/sagernet/sagerconnect/core/gtkx"
)

//go:embed icon.png
var IconPng []byte

func MustPixbuf(file []byte) *gdk.Pixbuf {
	pix, err := gdk.PixbufNewFromBytesOnly(file)
	gtkx.Must("create pixbuf", err)
	return pix
}
