package module

import (
	"context"

	"github.com/cerkas/cerkas-backend/config"
	"github.com/cerkas/cerkas-backend/core/entity"
	"github.com/cerkas/cerkas-backend/core/repository"
)

type CatalogUsecase interface {
	GetObjectData(ctx context.Context, request entity.CatalogQuery) (resp entity.CatalogResponse, err error)
	GetObjectDetail(ctx context.Context, request entity.CatalogQuery, serial string) (resp map[string]entity.DataItem, err error)
	GetDataByRawQuery(ctx context.Context, request entity.CatalogQuery) (resp entity.CatalogResponse, err error)
	CreateObjectData(ctx context.Context, request entity.DataMutationRequest) (resp entity.CatalogResponse, err error)
	UpdateObjectData(ctx context.Context, request entity.DataMutationRequest) (resp entity.CatalogResponse, err error)
	DeleteObjectData(ctx context.Context, request entity.DataMutationRequest) (err error)
}

type catalogUsecase struct {
	cfg         config.Config
	catalogRepo repository.CatalogRepository
}

func NewCatalogUsecase(cfg config.Config, catalogRepo repository.CatalogRepository) CatalogUsecase {
	return &catalogUsecase{
		cfg:         cfg,
		catalogRepo: catalogRepo,
	}
}

func (uc *catalogUsecase) GetObjectData(ctx context.Context, request entity.CatalogQuery) (resp entity.CatalogResponse, err error) {
	return uc.catalogRepo.GetObjectData(ctx, request)
}

func (uc *catalogUsecase) GetObjectDetail(ctx context.Context, request entity.CatalogQuery, serial string) (resp map[string]entity.DataItem, err error) {
	request.Serial = serial

	return uc.catalogRepo.GetObjectDetail(ctx, request)
}

func (uc *catalogUsecase) GetDataByRawQuery(ctx context.Context, request entity.CatalogQuery) (resp entity.CatalogResponse, err error) {
	return uc.catalogRepo.GetDataByRawQuery(ctx, request)
}

func (uc *catalogUsecase) CreateObjectData(ctx context.Context, request entity.DataMutationRequest) (resp entity.CatalogResponse, err error) {
	return uc.catalogRepo.CreateObjectData(ctx, request)
}

func (uc *catalogUsecase) UpdateObjectData(ctx context.Context, request entity.DataMutationRequest) (resp entity.CatalogResponse, err error) {
	return uc.catalogRepo.UpdateObjectData(ctx, request)
}

func (uc *catalogUsecase) DeleteObjectData(ctx context.Context, request entity.DataMutationRequest) (err error) {
	return uc.catalogRepo.DeleteObjectData(ctx, request)
}
