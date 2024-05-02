package apisql

type TableField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type DDLCreateTable struct {
	Columns     []TableField `json:"columns"`
	Constraints []TableField `json:"constraints"`
}

type DDLTableOptions struct {
	constraint string
	identity   string
}

type JSONResult struct {
	Message string `json:"message"`
}

type Response struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}
