package module

import (
	"context"
	"reflect"

	"github.com/cerkas/cerkas-backend/config"
	"github.com/cerkas/cerkas-backend/core/entity"
	"github.com/cerkas/cerkas-backend/core/repository"
	"github.com/cerkas/cerkas-backend/pkg/helper"
)

type ViewUsecase interface {
	GetContentLayoutByKeys(ctx context.Context, request entity.GetViewContentByKeysRequest) (resp entity.ViewContent, err error)
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

func (uc *viewUsecase) GetContentLayoutByKeys(ctx context.Context, request entity.GetViewContentByKeysRequest) (resp entity.ViewContent, err error) {
	viewContentRecord, err := uc.viewRepo.GetViewContentByKeys(ctx, request)
	if err != nil {
		return resp, err
	}

	// Convert map to struct
	if err = mapToStructSnakeCase(viewContentRecord, &resp); err != nil {
		return resp, err
	}

	// tenant record
	if viewContentRecord[entity.TENANT_CODE].Value != nil {
		tenantSt := entity.Tenants{}

		if ok := viewContentRecord[entity.TENANT_CODE].Value.(string); ok != "" {
			tenantRecord, err := uc.catalogRepo.GetObjectDetail(ctx, entity.CatalogQuery{
				Serial:     viewContentRecord[entity.TENANT_CODE].Value.(string),
				ObjectCode: "tenants",
				TenantCode: entity.PUBLIC,
			})
			if err != nil {
				return resp, err
			}

			if err = mapToStructSnakeCase(tenantRecord, &tenantSt); err != nil {
				return resp, err
			}

			resp.Tenant = tenantSt
		}

	}

	// object record
	if viewContentRecord[entity.OBJECT_CODE].Value != nil {
		objectSt := entity.Objects{}

		if ok := viewContentRecord[entity.OBJECT_CODE].Value.(string); ok != "" {
			objectRecord, err := uc.catalogRepo.GetObjectDetail(ctx, entity.CatalogQuery{
				Serial:     viewContentRecord[entity.OBJECT_CODE].Value.(string),
				ObjectCode: "objects",
				TenantCode: entity.PUBLIC,
			})
			if err != nil {
				return resp, err
			}

			if err = mapToStructSnakeCase(objectRecord, &objectSt); err != nil {
				return resp, err
			}

			resp.Object = objectSt
		}
	}

	// product
	if viewContentRecord[entity.PRODUCT_CODE].Value != nil {
		productSt := entity.Products{}

		if ok := viewContentRecord[entity.PRODUCT_CODE].Value.(string); ok != "" {
			objectRecord, err := uc.catalogRepo.GetObjectDetail(ctx, entity.CatalogQuery{
				Serial:     viewContentRecord[entity.PRODUCT_CODE].Value.(string),
				ObjectCode: "products",
				TenantCode: entity.PUBLIC,
			})
			if err != nil {
				return resp, err
			}

			if err = mapToStructSnakeCase(objectRecord, &productSt); err != nil {
				return resp, err
			}

			resp.Product = productSt
		}
	}

	// view schema
	if viewContentRecord[entity.VIEW_SCHEMA_SERIAL].Value != nil {
		viewSchemaSt := entity.ViewSchema{}

		if ok := viewContentRecord[entity.VIEW_SCHEMA_SERIAL].Value.(string); ok != "" {
			objectRecord, err := uc.catalogRepo.GetObjectDetail(ctx, entity.CatalogQuery{
				Serial:     viewContentRecord[entity.VIEW_SCHEMA_SERIAL].Value.(string),
				ObjectCode: "view_schema",
				TenantCode: entity.PUBLIC,
			})
			if err != nil {
				return resp, err
			}

			if err = mapToStructSnakeCase(objectRecord, &viewSchemaSt); err != nil {
				return resp, err
			}

			resp.ViewSchema = viewSchemaSt
		}
	}

	// view layout
	if viewContentRecord[entity.VIEW_LAYOUT_SERIAL].Value != nil {
		viewLayoutSt := entity.ViewLayout{}

		if ok := viewContentRecord[entity.VIEW_LAYOUT_SERIAL].Value.(string); ok != "" {
			objectRecord, err := uc.catalogRepo.GetObjectDetail(ctx, entity.CatalogQuery{
				Serial:     viewContentRecord[entity.VIEW_LAYOUT_SERIAL].Value.(string),
				ObjectCode: "view_layout",
				TenantCode: entity.PUBLIC,
			})
			if err != nil {
				return resp, err
			}

			if err = mapToStructSnakeCase(objectRecord, &viewLayoutSt); err != nil {
				return resp, err
			}

			resp.ViewLayout = viewLayoutSt
		}
	}

	// get original fields
	catalogQuery := entity.CatalogQuery{
		ObjectCode:  request.ObjectCode,
		TenantCode:  request.TenantCode,
		ProductCode: request.ProductCode,
	}

	originalFields, _, err := uc.catalogRepo.GetColumnList(ctx, catalogQuery)
	if err != nil {
		return resp, err
	}

	resp.Fields = originalFields

	return resp, nil
}

// Conversion function
func mapToStructSnakeCase(data map[string]entity.DataItem, target interface{}) error {
	// Get the value of the target struct
	targetVal := reflect.ValueOf(target).Elem()

	// Loop through the fields of the target struct
	for i := 0; i < targetVal.NumField(); i++ {
		// Get the field and its name
		field := targetVal.Type().Field(i)
		fieldName := field.Name

		// Convert the CamelCase field name to snake_case
		snakeKey := helper.CamelToSnake(fieldName)

		// Match the snake_case key with the map
		if item, exists := data[snakeKey]; exists {
			// Get the field in the target struct
			targetField := targetVal.Field(i)

			// Ensure the field is settable
			if targetField.IsValid() && targetField.CanSet() {
				// Set the value from the DataItem.Value
				if item.Value == nil {
					continue
				}

				value := reflect.ValueOf(item.Value)

				// Ensure types match or are assignable
				if value.Type().AssignableTo(targetField.Type()) {
					targetField.Set(value)
				}
			}
		}
	}

	return nil
}
