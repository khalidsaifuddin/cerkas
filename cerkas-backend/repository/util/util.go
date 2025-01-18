package util

import (
	"database/sql"

	"github.com/cerkas/cerkas-backend/core/entity"
)

func HandleSingleRow(columnsList []map[string]interface{}, rows *sql.Rows, request entity.CatalogQuery) (item map[string]entity.DataItem, err error) {
	// Create a slice of interface{} to hold column values
	values := make([]interface{}, len(columnsList))
	valuePointers := make([]interface{}, len(columnsList))
	for i := range values {
		valuePointers[i] = &values[i]
	}

	// Scan the row
	if err := rows.Scan(valuePointers...); err != nil {
		return item, err
	}

	// Create a map for the row
	item = make(map[string]entity.DataItem)
	for i, colName := range columnsList {
		val := values[i]

		// TODO: implement dynamic data type so it will not retun literally from db datatype

		fieldName := colName[entity.FieldColumnCode].(string)
		originalColumnCode := colName[entity.FieldOriginalColumnCode].(string)

		if val, ok := request.Fields[originalColumnCode]; ok && val.FieldName != "" {
			fieldName = val.FieldName
		}

		item[colName[entity.FieldOriginalColumnCode].(string)] = entity.DataItem{
			FieldCode:    colName[entity.FieldColumnCode].(string),
			FieldName:    fieldName,
			DataType:     colName[entity.FieldDataType].(string),
			Value:        val,
			DisplayValue: val,
		}
	}

	return item, nil
}
