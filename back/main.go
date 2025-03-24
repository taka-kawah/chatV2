package main

import (
	"back/db"
	"back/middleware/authentication"
	"back/middleware/ws"
	"back/usecase"
	"fmt"
)

func main() {
	fmt.Println("Hello world 🍣")
	d, err := db.NewDbInstances()
	if err != nil {
		errLog(err)
	}
	defer d.Disconnect()

	pm := &providerManager{
		ap: usecase.NewAuthService(db.NewAuthDriver(d.GormDb), &authentication.AuthMiddleware{}),
		cp: usecase.NewChatService(db.NewChatDriver(d.GormDb), db.NewChatViewDriver(d.GormDb)),
		rp: usecase.NewRoomService(db.NewRoomDriver(d.GormDb)),
		up: usecase.NewUserService(db.NewUserDriver(d.GormDb)),
		hp: ws.NewHubManager(),
	}

	e := setupRouter(*pm)
	e.Run()
}
