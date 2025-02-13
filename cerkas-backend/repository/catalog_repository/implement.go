package catalogrepository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cerkas/cerkas-backend/config"
	"github.com/cerkas/cerkas-backend/core/entity"
	repository_intf "github.com/cerkas/cerkas-backend/core/repository"
	"github.com/cerkas/cerkas-backend/pkg/helper"
	"github.com/cerkas/cerkas-backend/repository/util"
	"gorm.io/gorm"
)

type repository struct {
	cfg config.Config
	db  *gorm.DB
}

func New(cfg config.Config, db *gorm.DB) repository_intf.CatalogRepository {
	return &repository{
		cfg: cfg,
		db:  db,
	}
}

func (r *repository) GetColumnList(ctx context.Context, request entity.CatalogQuery) (columns []map[string]interface{}, columnStrings string, err error) {
	// get list of column from request.ObjectCode
	listColumnQuery := fmt.Sprintf(`
	SELECT
    col.column_name as column_code, 
		col.udt_name as data_type,
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name
	FROM
			information_schema.columns AS col
	LEFT JOIN information_schema.key_column_usage AS kcu ON col.table_name = kcu.table_name
	AND col.column_name = kcu.column_name
	LEFT JOIN information_schema.constraint_column_usage AS ccu ON kcu.constraint_name = ccu.constraint_name
	WHERE
		col.table_schema = '%v' 
	AND col.table_name = '%v'
	`, request.TenantCode, request.ObjectCode)

	rows, err := r.db.Raw(listColumnQuery).Rows()
	if err != nil {
		return columns, columnStrings, err
	}
	defer rows.Close()

	// iterate over the result to get value of column_name and data_type
	for rows.Next() {
		column := make(map[string]interface{})

		var columnCode, dataType, foreignTableName, foreignColumnName interface{}
		if err := rows.Scan(&columnCode, &dataType, &foreignTableName, &foreignColumnName); err != nil {
			return columns, columnStrings, err
		}

		column[entity.FieldDataType] = dataType.(string)
		column[entity.FieldColumnCode] = fmt.Sprintf("%v.%v.%v", request.TenantCode, request.ObjectCode, columnCode.(string))
		column[entity.FieldColumnName] = columnCode.(string)
		column[entity.FieldOriginalColumnCode] = columnCode.(string)

		if foreignTableName != nil && foreignTableName.(string) != request.ObjectCode && foreignColumnName != nil && foreignColumnName.(string) != "id" {
			column[entity.FieldForeignTableName] = foreignTableName.(string)
			column[entity.FieldForeignColumnName] = foreignColumnName.(string)
		}

		columns = append(columns, column)
	}

	// filter columns if request.Fields is not empty
	if len(request.Fields) > 0 {
		var filteredColumns []map[string]interface{}
		for fieldNameKey := range request.Fields {
			isFound := false

			if !strings.Contains(fieldNameKey, "__") {
				completeFieldName := fmt.Sprintf("%v.%v.%v", request.TenantCode, request.ObjectCode, fieldNameKey)

				for _, column := range columns {
					if completeFieldName == column[entity.FieldColumnCode] {
						isFound = true
						column[entity.FieldOriginalColumnCode] = fieldNameKey

						filteredColumns = append(filteredColumns, column)
					}
				}
			} else {
				// handle fieldName that has double underscore this indicates that it is a relationship field

				// split fieldName by double underscore
				fieldNameSplit := strings.Split(fieldNameKey, "__")
				foreignColumnName := fmt.Sprintf("%v.%v.%v", request.TenantCode, request.ObjectCode, fieldNameSplit[0])
				referenceColumnName := fieldNameSplit[1]

				// get foreign table name
				foreignTableName := ""
				foreignReferenceColumnName := ""
				for _, column := range columns {
					if column[entity.FieldColumnCode] == foreignColumnName {
						foreignTableName = fmt.Sprintf("%v.%v", request.TenantCode, column[entity.FieldForeignTableName])
						foreignReferenceColumnName = column[entity.FieldForeignColumnName].(string)
						break
					}
				}

				// convert fieldName to tableName.columnName
				fieldCode := foreignTableName + "." + referenceColumnName
				fieldName := fieldCode

				if val := request.Fields[fieldNameKey].FieldName; val != "" {
					fieldName = val
				}

				filteredColumns = append(filteredColumns, map[string]interface{}{
					entity.FieldOriginalColumnCode: fieldNameKey,
					entity.FieldColumnCode:         fieldCode,
					entity.FieldColumnName:         fieldName,
					entity.FieldForeignColumnName:  foreignColumnName,
					entity.FieldDataType:           "text",
					entity.ForeignTable: map[string]string{
						entity.FieldForeignTableName:      foreignTableName,
						entity.FieldForeignColumnName:     referenceColumnName,
						entity.ForeignReferenceColumnName: foreignReferenceColumnName,
					},
				})

				isFound = true
			}

			// after finish iterating columns, if field is not found in columns, return error
			if !isFound {
				return columns, columnStrings, fmt.Errorf("field %v is not found in table %v", fieldNameKey, request.ObjectCode)
			}
		}

		columns = filteredColumns
	}

	// convert columns to string
	for i, col := range columns {
		if i == 0 {
			columnStrings = col[entity.FieldColumnCode].(string)
		} else {
			columnStrings = columnStrings + ", " + col[entity.FieldColumnCode].(string)
		}
	}

	return columns, columnStrings, nil
}

func (r *repository) GetObjectData(ctx context.Context, request entity.CatalogQuery) (resp entity.CatalogResponse, err error) {
	// get list of column from request.ObjectCode
	completeTableName := request.TenantCode + "." + request.ObjectCode

	// Get list of columns
	columnsList, columnsString, err := r.GetColumnList(ctx, request)
	if err != nil {
		return resp, err
	}

	// Get total data count
	countQuery := getTotalCountQuery(completeTableName, request)
	resultCount, err := r.db.Raw(countQuery).Rows()
	if err != nil {
		return resp, err
	}

	for resultCount.Next() {
		resultCount.Scan(&resp.TotalData)
	}

	// Get data with pagination
	dataQuery := getDataWithPagination(columnsList, columnsString, completeTableName, request)
	rows, err := r.db.Raw(dataQuery).Rows()
	if err != nil {
		return resp, err
	}
	defer rows.Close()

	for rows.Next() {
		item, err := util.HandleSingleRow(columnsList, rows, request)
		if err != nil {
			return resp, err
		}

		// Append the item to the catalog
		resp.Items = append(resp.Items, item)
	}

	resp.Page = request.Page
	resp.PageSize = request.PageSize
	resp.TotalPage = int(helper.GenerateTotalPage(int64(resp.TotalData), int64(request.PageSize)))

	return resp, nil
}

func (r *repository) GetObjectDetail(ctx context.Context, request entity.CatalogQuery) (resp map[string]entity.DataItem, err error) {
	// get list of column from request.ObjectCode
	completeTableName := request.TenantCode + "." + request.ObjectCode

	// Get list of columns
	columnsList, columnsString, err := r.GetColumnList(ctx, request)
	if err != nil {
		return resp, err
	}

	// get single data using serial in request
	dataQuery := getSingleData(columnsList, columnsString, completeTableName, request)
	rows, err := r.db.Raw(dataQuery).Rows()
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		item, err := util.HandleSingleRow(columnsList, rows, request)
		if err != nil {
			return resp, err
		}

		resp = item
	}

	return resp, nil
}

func (r *repository) GetDataByRawQuery(ctx context.Context, request entity.CatalogQuery) (resp entity.CatalogResponse, err error) {
	// run raw query from request.RawQuery
	rawQuery := request.RawQuery

	// get total data based on rawQuery
	rawCountQuery := fmt.Sprintf("SELECT SUM(1) as total from (%s) as subquery", rawQuery)

	countRows, err := r.db.Raw(rawCountQuery).Rows()
	if err != nil {
		return resp, err
	}
	defer countRows.Close()

	var total int
	if countRows.Next() {
		err = countRows.Scan(&total)
		if err != nil {
			return resp, err
		}

		resp.TotalData = total
	}

	// add page and page size based on request.Page and request.PageSize
	rawQuery = fmt.Sprintf("%s LIMIT %d OFFSET %d", rawQuery, request.PageSize, (request.Page-1)*request.PageSize)

	rows, err := r.db.Raw(rawQuery).Rows()
	if err != nil {
		return resp, err
	}
	defer rows.Close()

	// get list of column from query result
	columns, err := rows.Columns()
	if err != nil {
		return resp, err
	}

	// iterate over the result to get value of column_name and data_type
	for rows.Next() {
		// Create a slice of interface{} to hold column values
		values := make([]interface{}, len(columns))
		valuePointers := make([]interface{}, len(columns))
		for i := range values {
			valuePointers[i] = &values[i]
		}

		// Scan the row
		if err := rows.Scan(valuePointers...); err != nil {
			return resp, err
		}

		// Create a map for the row
		item := make(map[string]entity.DataItem)
		for i, col := range columns {
			item[col] = entity.DataItem{
				FieldCode:    col,
				FieldName:    helper.CapitalizeWords(helper.ReplaceUnderscoreWithSpace(col)),
				DataType:     "text",
				Value:        values[i],
				DisplayValue: values[i],
			}
		}

		// Append the item to the catalog
		resp.Items = append(resp.Items, item)
	}

	resp.Page = request.Page
	resp.PageSize = request.PageSize
	resp.TotalPage = int(helper.GenerateTotalPage(int64(resp.TotalData), int64(request.PageSize)))

	return resp, nil
}

func (r *repository) CreateObjectData(ctx context.Context, request entity.DataMutationRequest) (resp entity.CatalogResponse, err error) {
	// INSERT INTO table_name (column1, column2, column3, ...)
	// VALUES (value1, value2, value3, ...);

	// get list of column from request.ObjectCode
	completeTableName := request.TenantCode + "." + request.ObjectCode

	// loop through data items and get the values
	var columnCodeString string
	var valueString string
	for _, item := range request.Items {
		columnCodeString = columnCodeString + ", " + item.FieldCode

		if item.Value == nil {
			valueString = valueString + ", NULL"
		} else {
			switch item.DataType {
			case "text":
				valueString = valueString + fmt.Sprintf(", '%v'", item.Value)
			case "integer":
				valueString = valueString + fmt.Sprintf(", %v", item.Value)
			case "boolean":
				valueString = valueString + fmt.Sprintf(", %v", item.Value)
			default:
				valueString = valueString + fmt.Sprintf(", '%v'", item.Value)
			}
		}
	}

	if len(valueString) == 0 {
		return resp, errors.New("no data item found")
	}

	columnCodeString = columnCodeString[2:]
	valueString = valueString[2:]

	// insert into query string
	insertQuery := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)", completeTableName, columnCodeString, valueString)
	log.Printf("insertQuery: %v", insertQuery)

	// execute insert query
	if err := r.db.Exec(insertQuery).Error; err != nil {
		return resp, err
	}

	return resp, nil
}

func (r *repository) UpdateObjectData(ctx context.Context, request entity.DataMutationRequest) (resp entity.CatalogResponse, err error) {
	// UPDATE table_name
	// SET column1 = value1, column2 = value2, ...
	// WHERE condition;

	return resp, nil
}

func (r *repository) DeleteObjectData(ctx context.Context, request entity.DataMutationRequest) (err error) {
	return nil
}

func isOperatorInLIKEList(operator entity.FilterOperator) bool {
	for _, validOperator := range entity.OperatorLIKEList {
		if operator == validOperator {
			return true
		}
	}
	return false
}

// Helper function to build dynamic filters based on CatalogQuery
func buildFilters(filters []entity.FilterGroup) string {
	var filterClauses []string
	for _, filterGroup := range filters {
		var groupClauses []string
		for _, filter := range filterGroup.Filters {
			operator := entity.OperatorQueryMap[filter.Operator]
			value := filter.Value

			// handler value of operator is part of entity.OperatorLIKEList, then we should add %
			if isOperatorInLIKEList(filter.Operator) {
				value = fmt.Sprintf("%%%v%%", value)
			}

			// Create filter conditions based on the field, operator, and value
			groupClauses = append(groupClauses, fmt.Sprintf("%s %s '%s'", filter.FieldName, operator, value))
		}
		// Combine the group clauses with the group operator (AND/OR)
		filterClauses = append(filterClauses, fmt.Sprintf("(%s)", strings.Join(groupClauses, fmt.Sprintf(" %s ", filterGroup.Operator))))
	}
	return strings.Join(filterClauses, " AND ")
}

// Helper function to build dynamic order by clauses
func buildOrderBy(request entity.CatalogQuery) string {
	var orderClauses []string
	for _, order := range request.Orders {
		fieldName := fmt.Sprintf("%v.%v.%v", request.TenantCode, request.ObjectCode, order.FieldName)
		orderClauses = append(orderClauses, fmt.Sprintf("%s %s", fieldName, order.Direction))
	}
	return strings.Join(orderClauses, ", ")
}

func getSingleData(columnList []map[string]interface{}, columnsString, tableName string, request entity.CatalogQuery) string {
	// Start building the base query
	query := fmt.Sprintf(`
		SELECT %v
		FROM %v`, columnsString, tableName)

	// handle join table if any
	for _, column := range columnList {
		if column[entity.ForeignTable] != nil {
			foreignTableName := column[entity.ForeignTable].(map[string]string)[entity.FieldForeignTableName]
			foreignTableReferenceColumnName := column[entity.ForeignTable].(map[string]string)[entity.ForeignReferenceColumnName]

			query = fmt.Sprintf("%s LEFT JOIN %v ON %v = %v.%v", query, foreignTableName, column[entity.FieldForeignColumnName], foreignTableName, foreignTableReferenceColumnName)
		}
	}

	query = query + fmt.Sprintf(" WHERE %v.deleted_at IS NULL", tableName)

	// apply serial to get single data
	identifierColumn := "serial"
	if !helper.IsUUID(request.Serial) {
		identifierColumn = "code"
	}

	query = query + fmt.Sprintf(" AND %v.%v = '%v'", tableName, identifierColumn, request.Serial)

	return query
}

// Main function to get data with pagination, filters, and orders
func getDataWithPagination(columnList []map[string]interface{}, columnsString, tableName string, request entity.CatalogQuery) string {
	// Start building the base query
	query := fmt.Sprintf(`
		SELECT %v
		FROM %v`, columnsString, tableName)

	// handle join table if any
	for _, column := range columnList {
		if column[entity.ForeignTable] != nil {
			foreignTableName := column[entity.ForeignTable].(map[string]string)[entity.FieldForeignTableName]
			foreignTableReferenceColumnName := column[entity.ForeignTable].(map[string]string)[entity.ForeignReferenceColumnName]

			query = fmt.Sprintf("%s LEFT JOIN %v ON %v = %v.%v", query, foreignTableName, column[entity.FieldForeignColumnName], foreignTableName, foreignTableReferenceColumnName)
		}
	}

	query = query + fmt.Sprintf(" WHERE %v.deleted_at IS NULL", tableName)

	// Apply dynamic filters if they exist
	if len(request.Filters) > 0 {
		query = query + " AND " + buildFilters(request.Filters)
	}

	// Apply dynamic order by if they exist
	if len(request.Orders) > 0 {
		query = query + " ORDER BY " + buildOrderBy(request)
	}

	// Apply pagination (LIMIT and OFFSET)
	query = fmt.Sprintf("%s LIMIT %d OFFSET %d", query, request.PageSize, (request.Page-1)*request.PageSize)
	log.Print(query)

	return query
}

func getTotalCountQuery(tableName string, request entity.CatalogQuery) string {
	query := fmt.Sprintf(`
	SELECT COUNT(*)
	FROM %v
	WHERE deleted_at IS NULL`, tableName)

	// Apply dynamic filters if they exist
	if len(request.Filters) > 0 {
		query = query + " AND " + buildFilters(request.Filters)
	}

	log.Print(query)
	return query
}
