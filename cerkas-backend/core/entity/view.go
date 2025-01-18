package entity

type GetViewContentByKeysRequest struct {
	TenantCode      string `json:"tenant_code"`
	ProductCode     string `json:"product_code"`
	ObjectCode      string `json:"object_code"`
	ViewContentCode string `json:"view_content_code"`
	LayoutType      string `json:"layout_type"`
}

type ViewLayout struct {
	Serial       string                 `json:"serial"`
	Code         string                 `json:"code"`
	LayoutConfig map[string]interface{} `json:"layout_config"`
}

type ViewSchema struct {
	Serial        string                 `json:"serial"`
	Code          string                 `json:"code"`
	Name          string                 `json:"name"`
	Query         map[string]interface{} `json:"query"`
	DisplayField  []interface{}          `json:"display_field"`
	StructureType string                 `json:"structure_type"`
	ActionSerial  string                 `json:"action_serial"`
	IsFavorite    bool                   `json:"is_favorite"`
	ObjectSerial  string                 `json:"object_serial"`
	FieldSections map[string]interface{} `json:"field_sections"`
}

type ViewContent struct {
	Serial        string                   `json:"serial"`
	Code          string                   `json:"code"`
	Name          string                   `json:"name"`
	Tenant        Tenants                  `json:"tenant"`
	Product       Products                 `json:"product"`
	Object        Objects                  `json:"object"`
	OwnerSerial   string                   `json:"owner_serial"`
	ViewLayout    ViewLayout               `json:"view_layout"`
	ViewSchema    ViewSchema               `json:"view_schema"`
	LayoutType    string                   `json:"layout_type"`
	IsDefault     bool                     `json:"is_default"`
	IsShownInList bool                     `json:"is_shown_in_list"`
	Fields        []map[string]interface{} `json:"fields"`
}
