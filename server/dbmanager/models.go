package dbmanager

import (
	"time"

	"github.com/dchest/uniuri"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tunnel struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name       string    `gorm:"type:varchar(255);not null"`
	EndpointID string    `gorm:"type:varchar(255);not null;index:idx_endpoint_id,unique"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`

	UpstreamToken string `gorm:"type:varchar(255)"`
	ProxyToken    string `gorm:"type:varchar(255)"`
}

func (t *Tunnel) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}

	t.UpstreamToken = uniuri.NewLen(32)
	t.ProxyToken = uniuri.NewLen(32)

	return
}
