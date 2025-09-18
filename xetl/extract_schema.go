package xetl

// ExtractTableQueryKey 查询参数
type ExtractTableQueryKey struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// ExtractTable 数据表提取规则描述，这里可以直接 select *，数据提取后存储为二进制格式，不会反序列化
type ExtractTable struct {
	Table    string                 `json:"table"`
	Message  string                 `json:"message"`
	Multiple bool                   `json:"multiple"`
	SQL      string                 `json:"sql"`
	Keys     []ExtractTableQueryKey `json:"keys"`
}

// ExtractSchema 数据提取 schema 描述
type ExtractSchema struct {
	Tables []*ExtractTable        `json:"tables"`
	Keys   []ExtractTableQueryKey `json:"keys"`
}
