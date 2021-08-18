module github.com/sagernet/sagerconnect

go 1.16

require (
	//github.com/google/uuid v1.2.0
	github.com/miekg/dns v1.1.43
	github.com/pkg/errors v0.9.1
	github.com/ulikunitz/xz v0.5.10
	github.com/xjasonlyu/tun2socks v1.18.4-0.20210813034434-85cf694b8fed
	golang.org/x/sys v0.0.0-20210809222454-d867a43fc93e
)

replace github.com/xjasonlyu/tun2socks v1.18.4-0.20210813034434-85cf694b8fed => github.com/sagernet/tun2socks v1.18.4-0.20210818073943-7767582a6821
