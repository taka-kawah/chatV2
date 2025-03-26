package provider

type ProviderManager struct {
	ap AuthProvider
	cp ChatProvider
	rp RoomProvider
	up UserProvider
	hp HubManagerProvider
}

func NewProviderManager(ap AuthProvider, cp ChatProvider, rp RoomProvider, up UserProvider, hp HubManagerProvider) *ProviderManager {
	return &ProviderManager{
		ap: ap,
		cp: cp,
		rp: rp,
		up: up,
		hp: hp,
	}
}

func (pm *ProviderManager) Ap() AuthProvider {
	return pm.ap
}

func (pm *ProviderManager) Cp() ChatProvider {
	return pm.cp
}

func (pm *ProviderManager) Rp() RoomProvider {
	return pm.rp
}

func (pm *ProviderManager) Up() UserProvider {
	return pm.up
}

func (pm *ProviderManager) Hp() HubManagerProvider {
	return pm.hp
}
