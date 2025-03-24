package ws

type HubManager struct {
	hubs map[uint]*Hub
}

func NewHubManager() *HubManager {
	return &HubManager{hubs: make(map[uint]*Hub)}
}

func (hm *HubManager) GetOrCreate(roomId uint) *Hub {
	if h, exists := hm.hubs[roomId]; exists {
		return h
	}
	nh := newHub(roomId)
	hm.hubs[roomId] = nh
	go nh.run()
	return nh
}
