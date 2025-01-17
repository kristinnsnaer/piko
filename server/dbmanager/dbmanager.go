package dbmanager

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type DBManager struct {
	orm *gorm.DB

	TunnelManager *TunnelManager
}

func NewDBManager() *DBManager {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{TranslateError: true})
	if err != nil {
		panic(err)
	}

	manager := &DBManager{
		orm:           db,
		TunnelManager: NewTunnelManager(db),
	}

	manager.Migrate()

	return manager
}

func (d *DBManager) GetDB() *gorm.DB {
	return d.orm
}

func (d *DBManager) Migrate() {
	d.orm.AutoMigrate(&Tunnel{})
}
