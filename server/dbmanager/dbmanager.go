package dbmanager

import (
	"errors"

	"github.com/andydunstall/piko/pkg/log"
	"github.com/andydunstall/piko/server/config"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type DBRepository struct {
	config *config.DatabaseConfig
	orm    *gorm.DB
}

type DBManager struct {
	config *config.DatabaseConfig
	orm    *gorm.DB
	logger log.Logger

	TunnelManager *TunnelManager
}

var (
	ErrDbDisabled = errors.New("database is not enabled")
)

func NewDBManager(
	conf *config.DatabaseConfig,
	logger log.Logger,
) *DBManager {
	logger = logger.WithSubsystem("database")
	db, err := gorm.Open(GetDialect(conf), &gorm.Config{TranslateError: true})

	if !conf.Enabled {
		logger.Info("database disabled")
	}

	if err != nil {
		panic(err)
	}

	repo := DBRepository{
		orm:    db,
		config: conf,
	}

	manager := &DBManager{
		orm:           db,
		logger:        logger,
		TunnelManager: NewTunnelManager(repo),
		config:        conf,
	}

	manager.Migrate()

	return manager
}

func NewInMemoryDbManager() *DBManager {
	return NewDBManager(&config.DatabaseConfig{
		Enabled: false,
	}, log.NewNopLogger())
}

func (d *DBManager) GetDB() *gorm.DB {
	return d.orm
}

func (d *DBManager) Migrate() {
	err := d.orm.AutoMigrate(&Tunnel{})
	if err != nil {
		panic(err)
	}
}

func (d *DBManager) Enabled() bool {
	return d.config.Enabled
}

func (d *DBRepository) IsEnabled() bool {
	return d.config.Enabled
}

func (d *DBRepository) AssertEnabled() error {
	if !d.IsEnabled() {
		return ErrDbDisabled
	}
	return nil
}

func GetDialect(config *config.DatabaseConfig) gorm.Dialector {
	if !config.Enabled {
		return sqlite.Open(":memory:?cache=shared")
	}
	switch config.DriverName {
	case "sqlite":
		if config.DatasourceConfig.Dsn == "" {
			panic("sqlite: missing dsn")
		}
		return sqlite.Open(config.DatasourceConfig.Dsn)
	default:
		panic("unsupported database driver")
	}
}
