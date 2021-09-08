package transport

import (
	"context"
	"github.com/google/uuid"
	"github.com/kallydev/shadowsocks-gateway/event"
	"github.com/kallydev/shadowsocks-gateway/internal/pool"
	"log"
	"net"
)

type TCPListener struct {
	listener *net.TCPListener
	event    *event.Event
	remote   string
}

func (tcpListener *TCPListener) Run() error {
	for {
		tcpConn, err := tcpListener.listener.AcceptTCP()
		if err != nil {
			continue
		}
		ctx := context.WithValue(context.Background(), event.ID, uuid.New())
		if !tcpListener.event.ConnectHandler(ctx, ctx.Value(event.ID).(uuid.UUID), "tcp", tcpConn.RemoteAddr()) {
			continue
		}
		go func() {
			if err := tcpListener.Handle(ctx, tcpConn); err != nil {
				log.Println(err)
			}
		}()
	}
}

func (tcpListener *TCPListener) Handle(ctx context.Context, src *net.TCPConn) error {
	ips, err := net.LookupIP(tcpListener.remote)
	if err != nil {
		return err
	}
	dst, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   ips[0],
		Port: src.LocalAddr().(*net.TCPAddr).Port,
	})
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(ctx)
	go tcpListener.Forward(ctx, dst, src, func(size int) {
		if !tcpListener.event.ForwardHandler(ctx, ctx.Value(event.ID).(uuid.UUID), event.DirectionInbound, size) {
			cancel()
		}
	})
	tcpListener.Forward(ctx, src, dst, func(size int) {
		if !tcpListener.event.ForwardHandler(ctx, ctx.Value(event.ID).(uuid.UUID), event.DirectionOutbound, size) {
			cancel()
		}
	})
	return nil
}

func (tcpListener *TCPListener) Forward(ctx context.Context, dst, src *net.TCPConn, callback func(size int)) {
	buf := pool.Get()
	defer pool.Put(buf)
	defer dst.Close()
	for {
		select {
		case <-ctx.Done():
			break
		default:
			n, err := src.Read(buf)
			if n > 0 {
				callback(n)
				if _, err := dst.Write(buf[0:n]); err != nil {
					break
				}
			}
			if err != nil {
				break
			}
		}
	}
}

func ListenTCP(event *event.Event, local net.IP, remote string, port int) (tcpListener *TCPListener, err error) {
	tcpListener = &TCPListener{
		event:  event,
		remote: remote,
	}
	tcpListener.listener, err = net.ListenTCP("tcp", &net.TCPAddr{
		IP:   local,
		Port: port,
	})
	if err != nil {
		return nil, err
	}
	return tcpListener, nil
}
