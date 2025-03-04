package entity

type FilterGroupOperator string
type FilterOperator string

const (
	PUBLIC             = "public"
	PRODUCT_CODE       = "product_code"
	TENANT_CODE        = "tenant_code"
	OBJECT_CODE        = "object_code"
	VIEW_SCHEMA_SERIAL = "view_schema_serial"
	VIEW_LAYOUT_SERIAL = "view_layout_serial"

	FilterOperatorAnd FilterGroupOperator = "AND"
	FilterOperatorOr  FilterGroupOperator = "OR"

	FilterOperatorEqual            FilterOperator = "equal"
	FilterOperatorNotEqual         FilterOperator = "not_equal"
	FilterOperatorContains         FilterOperator = "contains"
	FilterOperatorNotContains      FilterOperator = "not_contains"
	FilterOperatorGreaterThan      FilterOperator = "greater_than"
	FilterOperatorGreaterThanEqual FilterOperator = "greater_than_equal"
	FilterOperatorLessThan         FilterOperator = "less_than"
	FilterOperatorLessThanEqual    FilterOperator = "less_than_equal"

	FieldColumnName            = "column_name"
	FieldDataType              = "data_type"
	FieldColumnCode            = "column_code"
	FieldForeignTableName      = "foreign_table_name"
	FieldForeignColumnName     = "foreign_column_name"
	FieldOriginalColumnCode    = "original_column_code"
	ForeignTable               = "foreign_table"
	ForeignReferenceColumnName = "foreign_reference_column_name"
)

var (
	OperatorQueryMap = map[FilterOperator]string{
		FilterOperatorEqual:            "=",
		FilterOperatorNotEqual:         "!=",
		FilterOperatorContains:         "ILIKE",
		FilterOperatorNotContains:      "NOT ILIKE",
		FilterOperatorGreaterThan:      ">",
		FilterOperatorGreaterThanEqual: ">=",
		FilterOperatorLessThan:         "<",
		FilterOperatorLessThanEqual:    "<=",
	}

	OperatorLIKEList = []FilterOperator{
		FilterOperatorContains,
		FilterOperatorNotContains,
	}
)

type Tenants struct {
	ID     int    `json:"id"`
	Serial string `json:"serial"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	// CreatedBy string    `json:"created_by"`
	// CreatedAt time.Time `json:"created_at"`
	// UpdatedBy string    `json:"updated_by"`
	// UpdatedAt time.Time `json:"updated_at"`
}

type DataSource struct {
	ID          int                    `json:"id"`
	Serial      string                 `json:"serial"`
	Code        string                 `json:"code"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Host        string                 `json:"host"`
	Port        string                 `json:"port"`
	Username    string                 `json:"username"`
	Password    string                 `json:"password"`
	DBName      string                 `json:"db_name"`
	Configs     map[string]interface{} `json:"configs"`
	Tenant      Tenants                `json:"tenant"`
	// CreatedBy   string                 `json:"created_by"`
	// CreatedAt   time.Time              `json:"created_at"`
	// UpdatedBy   string                 `json:"updated_by"`
	// UpdatedAt   time.Time              `json:"updated_at"`
}

type Modules struct {
	ID                 int    `json:"id"`
	Serial             string `json:"serial"`
	Code               string `json:"code"`
	Name               string `json:"name"`
	ParentModuleSerial string `json:"parent_module_serial"`
	Version            string `json:"version"`
	// CreatedBy          string    `json:"created_by"`
	// CreatedAt          time.Time `json:"created_at"`
	// UpdatedBy          string    `json:"updated_by"`
	// UpdatedAt          time.Time `json:"updated_at"`
}

type Products struct {
	ID      int    `json:"id"`
	Serial  string `json:"serial"`
	Code    string `json:"code"`
	Name    string `json:"name"`
	IconURL string `json:"icon_url"`
	// CreatedBy string    `json:"created_by"`
	// CreatedAt time.Time `json:"created_at"`
	// UpdatedBy string    `json:"updated_by"`
	// UpdatedAt time.Time `json:"updated_at"`
}

type Objects struct {
	ID          int        `json:"id"`
	Serial      string     `json:"serial"`
	Tenant      Tenants    `json:"tenant"`
	Module      Modules    `json:"module"`
	Code        string     `json:"code"`
	DisplayName string     `json:"display_name"`
	Description string     `json:"description"`
	ObjectType  string     `json:"object_type"`
	DataSource  DataSource `json:"data_source"`
	// CreatedBy   string     `json:"created_by"`
	// CreatedAt   time.Time  `json:"created_at"`
	// UpdatedBy   string     `json:"updated_by"`
	// UpdatedAt   time.Time  `json:"updated_at"`
}

type ObjectFields struct {
	ID                      int                    `json:"id"`
	Serial                  string                 `json:"serial"`
	Object                  Objects                `json:"object"`
	FieldCode               string                 `json:"field_code"`
	IsDisplayName           bool                   `json:"is_display_name"`
	DisplayName             string                 `json:"display_name"`
	FieldReference          string                 `json:"field_reference"`
	Description             string                 `json:"description"`
	DataType                DataType               `json:"data_type"`
	ValidationRules         map[string]interface{} `json:"validation_rules"`
	TargetObject            Objects                `json:"target_object_serial"`
	TargetObjectFieldSerial string                 `json:"target_object_field_serial"`
	Relation                string                 `json:"relation"`
	IsSystem                bool                   `json:"is_system"`
	DefaultValue            string                 `json:"default_value"`
	// CreatedBy               string                 `json:"created_by"`
	// CreatedAt               time.Time              `json:"created_at"`
	// UpdatedBy               string                 `json:"updated_by"`
	// UpdatedAt               time.Time              `json:"updated_at"`
}

type DataType struct {
	ID                int                    `json:"id"`
	Serial            string                 `json:"serial"`
	Code              string                 `json:"code"`
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	PrimitiveDataType string                 `json:"primitive_data_type"`
	ValidationRules   map[string]interface{} `json:"validation_rules"`
	IsActive          bool                   `json:"is_active"`
	DisplayType       string                 `json:"display_type"`
	FieldOptions      map[string]interface{} `json:"field_options"`
	Icon              string                 `json:"icon"`
}

type FilterItem struct {
	FieldName string         `json:"field_name"`
	Operator  FilterOperator `json:"operator"`
	Value     string         `json:"value"`
}

type FilterGroup struct {
	Operator FilterOperator `json:"operator"`
	Filters  []FilterItem   `json:"filter_item"`
}

type Order struct {
	FieldName string `json:"field_name"`
	Direction string `json:"direction"`
}

type Field struct {
	FieldCode string `json:"field_code"`
	FieldName string `json:"field_name"`
}

type CatalogQuery struct {
	Fields      map[string]Field `json:"fields"`
	Filters     []FilterGroup    `json:"filters"`
	Orders      []Order          `json:"orders"`
	Page        int              `json:"page"`
	PageSize    int              `json:"page_size"`
	Serial      string           `json:"serial"`
	ObjectCode  string           `json:"object_code"`
	TenantCode  string           `json:"tenant_code"`
	ProductCode string           `json:"product_code"`
	RawQuery    string           `json:"raw_query"`
}

type DataItem struct {
	FieldCode    string      `json:"field_code"`
	FieldName    string      `json:"field_name"`
	DataType     string      `json:"data_type"`
	Value        interface{} `json:"value"`
	DisplayValue interface{} `json:"display_value"`
}

type CatalogResponse struct {
	Page      int                   `json:"page"`
	PageSize  int                   `json:"page_size"`
	TotalData int                   `json:"total_data"`
	TotalPage int                   `json:"total_page"`
	Items     []map[string]DataItem `json:"items"`
}

type DataMutationRequest struct {
	Serial      string     `json:"serial"`
	Items       []DataItem `json:"items"`
	ObjectCode  string     `json:"object_code"`
	TenantCode  string     `json:"tenant_code"`
	ProductCode string     `json:"product_code"`
	UserSerial  string     `json:"user_serial"`
}
