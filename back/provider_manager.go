package main

import "back/interfaces"

type providerManager struct {
	ap interfaces.AuthProvider
	cp interfaces.ChatProvider
	rp interfaces.RoomProvider
	up interfaces.UserProvider
	hp interfaces.HubManagerProvider
}
