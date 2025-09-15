package filters

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// FilterTableService предоставляет функциональность для фильтрации данных
type FilterTableService struct {
	db *sql.DB
}

var Filters = []string{"order_by", "row_limit", "skip", "scopes", "paginate"}

// NewFilterTableService создает новый экземпляр сервиса фильтрации
func NewFilterTableService(db *sql.DB) *FilterTableService {
	return &FilterTableService{db: db}
}

// ApplyFilters применяет фильтры и возвращает SQL запрос и параметры
func (fts *FilterTableService) ApplyFilters(filters map[string]interface{}, tableName string) (string, []interface{}, error) {
	var (
		whereConditions []string
		params          []interface{}
		paramCount      = 1
	)

	// Базовый запрос
	baseQuery := fmt.Sprintf("SELECT * FROM %s", tableName)

	// Обрабатываем фильтры
	for column, value := range filters {
		if contains(Filters, column) {
			continue
		}

		conditions, newParams, newParamCount := fts.buildWhereCondition(column, value, paramCount)
		if conditions != "" {
			whereConditions = append(whereConditions, conditions)
			params = append(params, newParams...)
			paramCount = newParamCount
		}
	}

	// Добавляем WHERE условия
	if len(whereConditions) > 0 {
		baseQuery += " WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Применяем сортировку
	if orderBy, ok := filters["order_by"].(string); ok {
		baseQuery = fts.applySorting(baseQuery, orderBy)
	}

	// Применяем пагинацию и лимиты
	baseQuery = fts.applyPagination(baseQuery, filters)

	return baseQuery, params, nil
}

// buildWhereCondition строит условие WHERE для конкретного столбца
func (fts *FilterTableService) buildWhereCondition(column string, value interface{}, paramCount int) (string, []interface{}, int) {
	var conditions []string
	var params []interface{}

	switch v := value.(type) {
	case []interface{}:
		// Массив значений - используем IN
		placeholders := make([]string, len(v))
		for i, item := range v {
			placeholders[i] = fmt.Sprintf("$%d", paramCount)
			params = append(params, item)
			paramCount++
		}
		conditions = append(conditions, fmt.Sprintf("%s IN (%s)", column, strings.Join(placeholders, ",")))

	case string:
		if strings.Contains(v, "OR") {
			// OR условия в JSON формате
			return fts.buildOrConditions(v, paramCount)
		}
		return fts.buildStringCondition(column, v, paramCount)

	default:
		// Простое равенство
		conditions = append(conditions, fmt.Sprintf("%s = $%d", column, paramCount))
		params = append(params, value)
		paramCount++
	}

	return strings.Join(conditions, " AND "), params, paramCount
}

// buildOrConditions обрабатывает OR условия из JSON
func (fts *FilterTableService) buildOrConditions(value string, paramCount int) (string, []interface{}, int) {
	var orConditions []string
	var params []interface{}

	// Пытаемся распарсить JSON
	var conditions []map[string]interface{}
	if err := json.Unmarshal([]byte(value), &conditions); err != nil {
		log.Printf("Error parsing OR conditions: %v", err)
		return "", nil, paramCount
	}

	for _, condition := range conditions {
		col, ok1 := condition["column"].(string)
		val, ok2 := condition["value"]
		operator, _ := condition["rule"].(string)
		if !ok1 || !ok2 {
			continue
		}

		if operator == "" {
			operator = "="
		}

		switch operator {
		case "IN", "NOT_IN":
			if arr, ok := val.([]interface{}); ok {
				placeholders := make([]string, len(arr))
				for i, item := range arr {
					placeholders[i] = fmt.Sprintf("$%d", paramCount)
					params = append(params, item)
					paramCount++
				}
				sqlOperator := "IN"
				if operator == "NOT_IN" {
					sqlOperator = "NOT IN"
				}
				orConditions = append(orConditions, fmt.Sprintf("%s %s (%s)", col, sqlOperator, strings.Join(placeholders, ",")))
			}

		default:
			orConditions = append(orConditions, fmt.Sprintf("%s %s $%d", col, operator, paramCount))
			params = append(params, val)
			paramCount++
		}
	}

	if len(orConditions) == 0 {
		return "", nil, paramCount
	}

	return "(" + strings.Join(orConditions, " OR ") + ")", params, paramCount
}

// buildStringCondition обрабатывает строковые условия
// Добавьте этот метод в ваш FilterTableService
func (fts *FilterTableService) buildStringCondition(column, value string, paramCount int) (string, []interface{}, int) {
	var conditions []string
	var params []interface{}

	switch {
	case value == "IS NULL":
		conditions = append(conditions, fmt.Sprintf("%s IS NULL", column))

	case value == "IS NOT NULL":
		conditions = append(conditions, fmt.Sprintf("%s IS NOT NULL", column))

	case strings.HasPrefix(value, "LIKE"):
		likeValue := strings.TrimSpace(value[5:])
		conditions = append(conditions, fmt.Sprintf("%s ILIKE $%d", column, paramCount))
		params = append(params, "%"+likeValue+"%")
		paramCount++

	case strings.HasPrefix(value, "in__"):
		values := strings.Split(value[4:], ",")
		placeholders := make([]string, len(values))
		for i, val := range values {
			placeholders[i] = fmt.Sprintf("$%d", paramCount)
			params = append(params, strings.TrimSpace(val))
			paramCount++
		}
		conditions = append(conditions, fmt.Sprintf("%s IN (%s)", column, strings.Join(placeholders, ",")))

	case strings.HasPrefix(value, "not_in__"):
		values := strings.Split(value[8:], ",")
		placeholders := make([]string, len(values))
		for i, val := range values {
			placeholders[i] = fmt.Sprintf("$%d", paramCount)
			params = append(params, strings.TrimSpace(val))
			paramCount++
		}
		conditions = append(conditions, fmt.Sprintf("%s NOT IN (%s)", column, strings.Join(placeholders, ",")))

	case strings.Contains(strings.ToUpper(value), "BETWEEN"):
		// Обработка BETWEEN 'value1' AND 'value2'
		parts := strings.SplitN(value, " ", 4)
		if len(parts) >= 4 && strings.ToUpper(parts[0]) == "BETWEEN" {
			conditions = append(conditions, fmt.Sprintf("%s BETWEEN $%d AND $%d", column, paramCount, paramCount+1))
			params = append(params, parts[1], parts[3])
			paramCount += 2
		}

	default:
		// Проверяем операторы сравнения
		operators := []string{"!=", "<=", ">=", "<", ">"}
		found := false
		for _, op := range operators {
			if strings.Contains(value, op) {
				parts := strings.SplitN(value, " ", 2)
				if len(parts) == 2 {
					conditions = append(conditions, fmt.Sprintf("%s %s $%d", column, parts[0], paramCount))
					params = append(params, strings.TrimSpace(parts[1]))
					paramCount++
					found = true
					break
				}
			}
		}

		// Если не нашли операторов, используем равенство
		if !found {
			conditions = append(conditions, fmt.Sprintf("%s = $%d", column, paramCount))
			params = append(params, value)
			paramCount++
		}
	}

	return strings.Join(conditions, " AND "), params, paramCount
}
// applySorting применяет сортировку
func (fts *FilterTableService) applySorting(query, orderBy string) string {
	if orderBy == "" {
		return query
	}

	parts := strings.Split(orderBy, ":")
	column := parts[0]
	direction := "ASC"
	if len(parts) > 1 && strings.ToLower(parts[1]) == "desc" {
		direction = "DESC"
	}

	return query + fmt.Sprintf(" ORDER BY %s %s", column, direction)
}

// applyPagination применяет пагинацию и лимиты
func (fts *FilterTableService) applyPagination(query string, filters map[string]interface{}) string {
	if limit, ok := getIntFromMap(filters, "row_limit"); ok && limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	if skip, ok := getIntFromMap(filters, "skip"); ok && skip > 0 {
		query += fmt.Sprintf(" OFFSET %d", skip)
	}

	return query
}

// GetFilteredResults выполняет запрос с фильтрацией и возвращает результаты
func (fts *FilterTableService) GetFilteredResults(tableName string, filters map[string]interface{}) (*sql.Rows, error) {
	query, params, err := fts.ApplyFilters(filters, tableName)
	if err != nil {
		return nil, err
	}

	return fts.db.Query(query, params...)
}

// GetFilteredResultsCount возвращает количество записей с учетом фильтров
func (fts *FilterTableService) GetFilteredResultsCount(tableName string, filters map[string]interface{}) (int, error) {
	var count int
	var whereConditions []string
	var params []interface{}
	paramCount := 1

	// Строим условия WHERE
	for column, value := range filters {
		if contains(Filters, column) {
			continue
		}

		conditions, newParams, newParamCount := fts.buildWhereCondition(column, value, paramCount)
		if conditions != "" {
			whereConditions = append(whereConditions, conditions)
			params = append(params, newParams...)
			paramCount = newParamCount
		}
	}

	// Строим запрос COUNT
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	if len(whereConditions) > 0 {
		query += " WHERE " + strings.Join(whereConditions, " AND ")
	}

	err := fts.db.QueryRow(query, params...).Scan(&count)
	return count, err
}

// Вспомогательные функции

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func getIntFromMap(m map[string]interface{}, key string) (int, bool) {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case int:
			return v, true
		case int64:
			return int(v), true
		case float64:
			return int(v), true
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				return i, true
			}
		}
	}
	return 0, false
}

func getIntFromInterface(val interface{}) (int, bool) {
	switch v := val.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i, true
		}
	}
	return 0, false
}