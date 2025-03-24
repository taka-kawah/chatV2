package interfaces

import "back/middleware/ws"

type HubManagerProvider interface {
	GetOrCreate(uint) *ws.Hub
}
