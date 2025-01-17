package dbmanager

import "gorm.io/gorm"

type TunnelManager struct {
	orm *gorm.DB
}

func NewTunnelManager(orm *gorm.DB) *TunnelManager {
	return &TunnelManager{
		orm: orm,
	}
}

func (t *TunnelManager) CreateTunnel(name, endpointID string) (*Tunnel, error) {
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
	var tunnel Tunnel
	err := t.orm.First(&tunnel, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tunnel, nil
}

func (t *TunnelManager) GetTunnelFromEndpointID(endpointID string) (*Tunnel, error) {
	var tunnel Tunnel
	err := t.orm.Where("endpoint_id = ?", endpointID).First(&tunnel).Error
	if err != nil {
		return nil, err
	}
	return &tunnel, nil
}
