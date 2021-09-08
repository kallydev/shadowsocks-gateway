package gateway

import (
	"context"
	"github.com/google/uuid"
	"github.com/kallydev/shadowsocks-gateway/config"
	"github.com/kallydev/shadowsocks-gateway/event"
	"github.com/kallydev/shadowsocks-gateway/internal/transport"
	"net"
)

type Gateway struct {
	Event        *event.Event
	conf         *config.Config
	tcpListeners []*transport.TCPListener
}

func New(conf *config.Config) *Gateway {
	return &Gateway{
		Event: &event.Event{
			ConnectHandler: func(ctx context.Context, id uuid.UUID, network string, remote net.Addr) bool {
				return true
			},
			ForwardHandler: func(ctx context.Context, id uuid.UUID, direction string, size int) bool {
				return true
			},
			CloseHandler: func(ctx context.Context, id uuid.UUID) {},
		},
		conf: conf,
	}
}

func (gateway Gateway) Run() error {
	if err := gateway.listen(); err != nil {
		return err
	}
	return gateway.handle()
}

func (gateway *Gateway) listen() error {
	start, end := gateway.conf.Port()
	for port := start; port <= end; port++ {
		tcpListener, err := transport.ListenTCP(gateway.Event, gateway.conf.IP(), gateway.conf.Remote, port)
		if err != nil {
			return err
		}
		gateway.tcpListeners = append(gateway.tcpListeners, tcpListener)
	}
	return nil
}

func (gateway *Gateway) handle() error {
	errCh := make(chan error)
	for _, tcpListener := range gateway.tcpListeners {
		go func(tcpListener *transport.TCPListener) {
			if err := tcpListener.Run(); err != nil {
				errCh <- err
			}
		}(tcpListener)
	}
	select {
	case err := <-errCh:
		return err
	}
}
