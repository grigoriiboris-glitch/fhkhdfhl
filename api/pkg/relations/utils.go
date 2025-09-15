package relations

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

func scanRowToMap(rows interface{}) (map[string]interface{}, error) {
	var row *sql.Row
	var sqlRows *sql.Rows
	
	switch r := rows.(type) {
	case *sql.Row:
		row = r
	case *sql.Rows:
		sqlRows = r
	default:
		return nil, fmt.Errorf("unsupported row type")
	}
	
	attrs := make(map[string]interface{})
	
	if row != nil {
		columns, err := getColumnsFromRow(row)
		if err != nil {
			return nil, err
		}
		
		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(interface{})
		}
		
		if err := row.Scan(values...); err != nil {
			return nil, err
		}
		
		for i, col := range columns {
			attrs[col] = *(values[i].(*interface{}))
		}
	}
	
	if sqlRows != nil {
		columns, err := sqlRows.Columns()
		if err != nil {
			return nil, err
		}
		
		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(interface{})
		}
		
		if err := sqlRows.Scan(values...); err != nil {
			return nil, err
		}
		
		for i, col := range columns {
			attrs[col] = *(values[i].(*interface{}))
		}
	}
	
	return attrs, nil
}

func getColumnsFromRow(row *sql.Row) ([]string, error) {
	// Получение колонок из row (упрощенная версия)
	return []string{}, nil
}

func getModelName(model Model) string {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

func getFieldValue(model Model, fieldName string) interface{} {
	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	
	field := val.FieldByName(strings.Title(fieldName))
	if !field.IsValid() {
		return nil
	}
	
	return field.Interface()
}