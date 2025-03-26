package main

import (
	"back/db"
	"back/http"
	"back/middleware/authentication"
	"back/middleware/ws"
	"back/provider"
	"back/usecase"
	"log"
)

func main() {
	d, err := db.NewDbInstances()
	if err != nil {
		log.Fatalf("%v, %v", err.Error(), err.Unwrap().Error())
	}
	defer d.Disconnect()

	pm := provider.NewProviderManager(
		usecase.NewAuthService(db.NewAuthDriver(d.GormDb), &authentication.AuthMiddleware{}),
		usecase.NewChatService(db.NewChatDriver(d.GormDb), db.NewChatViewDriver(d.GormDb)),
		usecase.NewRoomService(db.NewRoomDriver(d.GormDb)),
		usecase.NewUserService(db.NewUserDriver(d.GormDb)),
		ws.NewHubManager(),
	)

	e := http.SetupRouter(*pm)
	e.Run()
}
