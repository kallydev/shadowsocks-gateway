package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	gateway "github.com/kallydev/shadowsocks-gateway"
	"github.com/kallydev/shadowsocks-gateway/config"
	"log"
	"net"
)

func main() {
	instance := gateway.New(&config.Config{
		Bind: &config.Bind{
			IP:   net.IPv4zero.String(),
			Port: "8000-9000",
		},
		Remote: "127.0.0.1",
	})
	instance.Event.ConnectHandler = func(ctx context.Context, id uuid.UUID, network string, remote net.Addr) bool {
		fmt.Println("connect", ctx, id, network, remote)
		return true
	}
	instance.Event.ForwardHandler = func(ctx context.Context, id uuid.UUID, direction string, size int) bool {
		fmt.Println("forward", ctx, id, direction, size)
		return true
	}
	instance.Event.CloseHandler = func(ctx context.Context, id uuid.UUID) {
		fmt.Println("close", id)
	}
	if err := instance.Run(); err != nil {
		log.Fatalln(err)
	}
}
