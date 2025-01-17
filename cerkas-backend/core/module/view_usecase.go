package module

import (
	"context"

	"github.com/cerkas/cerkas-backend/config"
	"github.com/cerkas/cerkas-backend/core/repository"
)

type ViewUsecase interface {
	GetViewSchema(ctx context.Context) (err error)
	GetViewLayout(ctx context.Context) (err error)
}

type viewUsecase struct {
	cfg         config.Config
	catalogRepo repository.CatalogRepository
	viewRepo    repository.ViewRepository
}

func NewViewUsecase(cfg config.Config, catalogRepo repository.CatalogRepository, viewRepo repository.ViewRepository) ViewUsecase {
	return &viewUsecase{
		cfg:         cfg,
		catalogRepo: catalogRepo,
		viewRepo:    viewRepo,
	}
}

func (uc *viewUsecase) GetViewSchema(ctx context.Context) (err error) {
	return nil
}

func (uc *viewUsecase) GetViewLayout(ctx context.Context) (err error) {
	return nil
}
