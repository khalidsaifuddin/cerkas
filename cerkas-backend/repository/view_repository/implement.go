package viewrepository

import (
	"github.com/cerkas/cerkas-backend/config"
	"gorm.io/gorm"

	repository_intf "github.com/cerkas/cerkas-backend/core/repository"
)

type repository struct {
	db  *gorm.DB
	cfg config.Config
}

func New(db *gorm.DB, cfg config.Config) repository_intf.ViewRepository {
	return &repository{
		db:  db,
		cfg: cfg,
	}
}
