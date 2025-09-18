package xetl

// ExtractTableData 表数据
type ExtractTableData struct {
	Table    string     `json:"table"`
	Message  string     `json:"message"`
	Columns  []string   `json:"column"`
	Multiple bool       `json:"mutiple"`
	Rows     [][]string `json:"rows"`
}

// ExtractDataFile 整个导出的数据文件
type ExtractDataFile struct {
	Version string              `json:"version"`
	Tables  []*ExtractTableData `json:"tables"`
}
