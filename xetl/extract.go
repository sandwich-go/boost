package xetl

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
)

type driverExec interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

// QueryBySchema 执行 ExtractSchema SQL 并保存到 ExtractDataFile
func QueryBySchema(ctx context.Context, xd driverExec, schema *ExtractSchema, params map[string]interface{}) (*ExtractDataFile, error) {
	dataFile := &ExtractDataFile{
		Version: "1.0",
		Tables:  make([]*ExtractTableData, 0, len(schema.Tables)),
	}

	for _, tbl := range schema.Tables {
		args := make([]interface{}, 0, len(tbl.Keys))
		for _, k := range tbl.Keys {
			val, ok := params[k.Name]
			if !ok {
				return nil, fmt.Errorf("missing param %s for table %s", k.Name, tbl.Table)
			}
			args = append(args, val)
		}

		rows, err := xd.QueryContext(ctx, tbl.SQL, args...)
		if err != nil {
			return nil, fmt.Errorf("query table %s failed: %w", tbl.Table, err)
		}

		cols, err := rows.Columns()
		if err != nil {
			_ = rows.Close()
			return nil, err
		}

		tableData := &ExtractTableData{
			Table:    tbl.Table,
			Message:  tbl.Message,
			Multiple: tbl.Multiple,
			Columns:  cols,
			Rows:     make([][]string, 0),
		}

		for rows.Next() {
			values := make([]interface{}, len(cols))
			valuePtrList := make([]interface{}, len(cols))
			for i := range values {
				valuePtrList[i] = &values[i]
			}

			if err := rows.Scan(valuePtrList...); err != nil {
				_ = rows.Close()
				return nil, err
			}
			row := stringify(values)
			tableData.Rows = append(tableData.Rows, row)
			if !tbl.Multiple {
				break
			}
		}

		_ = rows.Close()
		dataFile.Tables = append(dataFile.Tables, tableData)
	}

	return dataFile, nil
}

func stringify(values []interface{}) []string {
	row := make([]string, len(values), len(values))

	for i, rawValue := range values {
		if rawValue == nil {
			row[i] = ""
			continue
		}

		byteArray, ok := rawValue.([]byte)
		if ok {
			rawValue = string(byteArray)
		}

		switch castValue := rawValue.(type) {
		case bool:
			row[i] = strconv.FormatBool(castValue)
		case string:
			row[i] = castValue
		case int:
			row[i] = strconv.FormatInt(int64(castValue), 10)
		case int8:
			row[i] = strconv.FormatInt(int64(castValue), 10)
		case int16:
			row[i] = strconv.FormatInt(int64(castValue), 10)
		case int32:
			row[i] = strconv.FormatInt(int64(castValue), 10)
		case int64:
			row[i] = strconv.FormatInt(int64(castValue), 10)
		case uint:
			row[i] = strconv.FormatUint(uint64(castValue), 10)
		case uint8:
			row[i] = strconv.FormatUint(uint64(castValue), 10)
		case uint16:
			row[i] = strconv.FormatUint(uint64(castValue), 10)
		case uint32:
			row[i] = strconv.FormatUint(uint64(castValue), 10)
		case uint64:
			row[i] = strconv.FormatUint(uint64(castValue), 10)
		default:
			row[i] = fmt.Sprintf("%v", castValue)
		}
	}

	return row
}
