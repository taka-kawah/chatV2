package ws

type HubManager struct {
	hubs map[uint]*hub
}

func NewHubManager() *HubManager {
	return &HubManager{hubs: make(map[uint]*hub)}
}

func (hm *HubManager) GetOrCreate(roomId uint) *hub {
	if h, exists := hm.hubs[roomId]; exists {
		return h
	}
	nh := newHub(roomId)
	hm.hubs[roomId] = nh
	go nh.run()
	return nh
}
