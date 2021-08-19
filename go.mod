module github.com/sagernet/sagerconnect

go 1.16

require (
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e
	github.com/miekg/dns v1.1.43
	github.com/pkg/errors v0.9.1
	github.com/ulikunitz/xz v0.5.10
	github.com/xjasonlyu/tun2socks v1.18.4-0.20210813034434-85cf694b8fed
	golang.org/x/sys v0.0.0-20210809222454-d867a43fc93e
)

replace github.com/xjasonlyu/tun2socks v1.18.4-0.20210813034434-85cf694b8fed => github.com/sagernet/tun2socks v1.18.4-0.20210818104726-9b52b624f351
