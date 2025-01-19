package dbmanager

type TunnelManager struct {
	DBRepository
}

func NewTunnelManager(repo DBRepository) *TunnelManager {
	return &TunnelManager{
		DBRepository: repo,
	}
}

func (t *TunnelManager) CreateTunnel(name, endpointID string) (*Tunnel, error) {
	if err := t.AssertEnabled(); err != nil {
		return nil, err
	}
	tunnel := Tunnel{
		Name:       name,
		EndpointID: endpointID,
	}
	err := t.orm.Create(&tunnel).Error

	if err != nil {
		return nil, err
	}

	return &tunnel, nil
}

func (t *TunnelManager) GetTunnel(id string) (*Tunnel, error) {
	if err := t.AssertEnabled(); err != nil {
		return nil, err
	}
	var tunnel Tunnel
	err := t.orm.First(&tunnel, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tunnel, nil
}

func (t *TunnelManager) GetTunnelFromEndpointID(endpointID string) (*Tunnel, error) {
	if err := t.AssertEnabled(); err != nil {
		return &Tunnel{
			EndpointID: endpointID,
		}, nil
	}
	var tunnel Tunnel
	err := t.orm.Where("endpoint_id = ?", endpointID).First(&tunnel).Error
	if err != nil {
		return nil, err
	}
	return &tunnel, nil
}
