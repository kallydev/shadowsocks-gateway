package event

import (
	"context"
	"github.com/google/uuid"
	"net"
)

const (
	ID = "id"

	DirectionInbound  = "Inbound"
	DirectionOutbound = "Outbound"
)

type ConnectHandler func(ctx context.Context, id uuid.UUID, network string, remote net.Addr) bool
type ForwardHandler func(ctx context.Context, id uuid.UUID, direction string, size int) bool
type CloseHandler func(ctx context.Context, id uuid.UUID)

type Event struct {
	ConnectHandler ConnectHandler
	ForwardHandler ForwardHandler
	CloseHandler   CloseHandler
}
