package module

import (
	"context"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

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
	GetObjectFieldsByObjectCode(ctx context.Context, request entity.CatalogQuery) (resp map[string]any, err error)
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

func (uc *catalogUsecase) GetObjectFieldsByObjectCode(ctx context.Context, request entity.CatalogQuery) (resp map[string]any, err error) {
	result, err := uc.catalogRepo.GetObjectFieldsByObjectCode(ctx, request)
	if err != nil {
		return resp, err
	}

	dataTypeSerials := []string{}
	for _, item := range result {
		data := item.(entity.ObjectFields)

		if data.ID != 0 {
			dataTypeSerials = append(dataTypeSerials, data.DataType.Serial)
		}
	}

	// get data type by serial
	dataType, _ := uc.catalogRepo.GetDataTypeBySerials(ctx, dataTypeSerials)

	// map data type to result
	dataTypeMap := make(map[string]entity.DataType)
	for _, item := range dataType {
		dataTypeMap[item.Serial] = item
	}

	// iterate result and map data type to result
	for i, item := range result {
		data := item.(entity.ObjectFields)
		if data.ID != 0 {
			if dataType, ok := dataTypeMap[data.DataType.Serial]; ok {
				data.DataType = dataType
			}
			result[i] = data
		}
	}

	return result, err
}

func (uc *catalogUsecase) GetObjectData(ctx context.Context, request entity.CatalogQuery) (resp entity.CatalogResponse, err error) {
	results, err := uc.catalogRepo.GetObjectData(ctx, request)
	if err != nil {
		return resp, err
	}

	objects, _ := uc.catalogRepo.GetObjectByCode(ctx, request.ObjectCode, request.TenantCode)
	objectFields := map[string]any{}

	if objects.Serial != "" {
		request.ObjectSerial = objects.Serial
		request.TenantSerial = objects.Tenant.Serial

		// handle custom object fields based on object field table
		objectFields, err = uc.GetObjectFieldsByObjectCode(ctx, request)
		if err != nil {
			return resp, err
		}
	}

	// TODO: inject view schema to get field config and query, and combine it to request fields and filters

	// iterate object fields and map to response
	for i, items := range results.Items {
		for j, item := range items {
			// put field code to complete field code, and assign new value to field code with column name only after splitted
			// example: item.FieldCode = "object.field_code" => item.CompleteFieldCode = "object.field_code" and item.FieldCode = "field_code"
			item.CompleteFieldCode = item.FieldCode

			fieldCode := ""

			// split item.FieldCode by dot
			if len(item.FieldCode) > 0 {
				split := strings.Split(item.FieldCode, ".")

				// find the last element
				fieldCode = split[len(split)-1]
			}

			if fieldCode == "" {
				continue
			}

			item.FieldCode = fieldCode
			// set field name to field code, but remove underscore and make it camel case
			// example: item.FieldCode = "object_field_code" => item.FieldName = "ObjectFieldCode"
			item.FieldName = strings.ReplaceAll(item.FieldCode, "_", " ")
			item.FieldName = cases.Title(language.English).String(item.FieldName)

			if field, ok := objectFields[fieldCode]; ok {
				data, ok := field.(entity.ObjectFields)
				if !ok {
					continue
				}

				// set custom field name
				item.FieldName = data.DisplayName

				// set custom data type
				item.DataType = data.DataType.Name
			}

			if requestDisplayName, ok := request.Fields[item.CompleteFieldCode]; ok {
				// set custom field name
				item.FieldName = requestDisplayName.FieldName
			}

			// set item.DataType to CamelCase
			item.DataType = cases.Title(language.English).String(item.DataType)

			results.Items[i][j] = item
		}
	}

	return results, err
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
