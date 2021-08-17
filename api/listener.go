package api

import "net"

type EventListener struct {
	listener func(event *Event)
	conn     *net.UDPConn
}

func NewEventListener(listener func(event *Event)) (*EventListener, error) {
	conn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		return nil, err
	}
	return &EventListener{
		listener: listener,
		conn:     conn,
	}, nil
}

func (l EventListener) LocalPort() uint16 {
	return uint16(l.conn.LocalAddr().(*net.UDPAddr).Port)
}

func (l EventListener) Close() {
	_ = l.conn.Close()
}
