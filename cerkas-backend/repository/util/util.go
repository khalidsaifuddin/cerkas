package util

import (
	"database/sql"
	"encoding/json"

	"github.com/cerkas/cerkas-backend/core/entity"
)

func HandleSingleRow(columnsList []map[string]interface{}, rows *sql.Rows, request entity.CatalogQuery) (item map[string]entity.DataItem, err error) {
	// Create a slice of interface{} to hold column values
	values := make([]any, len(columnsList))
	valuePointers := make([]any, len(columnsList))
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

		fieldName := colName[entity.FieldColumnCode].(string)
		FieldCode := colName[entity.FieldColumnCode].(string)

		if val, ok := request.Fields[FieldCode]; ok && val.FieldName != "" {
			fieldName = val.FieldName
		}

		// check if val is json
		if IsJSON(val) {
			var jsonData map[string]any

			if err := json.Unmarshal([]byte(val.([]uint8)), &jsonData); err == nil {
				val = jsonData
			}
		}

		item[colName[entity.FieldColumnCode].(string)] = entity.DataItem{
			FieldCode:    colName[entity.FieldColumnCode].(string),
			FieldName:    fieldName,
			DataType:     colName[entity.FieldDataType].(string),
			Value:        val,
			DisplayValue: val,
		}
	}

	return item, nil
}

func IsJSON(input any) bool {
	var temp any

	// Convert input to []byte if it's a string
	var data []byte
	switch v := input.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return false // Not a valid JSON candidate
	}

	// Try to unmarshal into a generic interface
	if err := json.Unmarshal(data, &temp); err != nil {
		return false // Not JSON
	}
	return true // Valid JSON
}
