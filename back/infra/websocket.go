package infra

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true //一旦開発用に全てのCORSを許可
	},
}

type WsHandler struct {
	clients map[*websocket.Conn]bool
	mux     sync.Mutex
}

func NewWsHandler() *WsHandler {
	return &WsHandler{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (h *WsHandler) Handle(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("failed to establish ws connection")
		return
	}
	h.mux.Lock()
	h.clients[conn] = true
	h.mux.Unlock()
	defer func() {
		h.mux.Lock()
		delete(h.clients, conn)
		h.mux.Unlock()
		conn.Close()
	}()

}
