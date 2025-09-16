package repository

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mymindmap/api/models"
)

type LogRepository struct {
	pool *pgxpool.Pool
	fields []string
	tableName string
}

func NewLogRepository(pool *pgxpool.Pool) *LogRepository {
	// Use a simple default table name based on the model name
	tableName := "logs"
	
	return &LogRepository{
		pool: pool,
		fields: []string{}, // Will be populated dynamically when needed
		tableName: tableName,
	}
}

// Add this method to populate fields on demand
func (r *LogRepository) getFields() []string {
	if len(r.fields) == 0 {
		// Dynamically generate fields from the model structure
		modelType := reflect.TypeOf(models.Log{})
		for i := 0; i < modelType.NumField(); i++ {
			field := modelType.Field(i)
			dbTag := r.getDBTag(field)
			r.fields = append(r.fields, dbTag)
		}
	}
	return r.fields
}

func (r *LogRepository) getDBTag(field reflect.StructField) string {
	dbTag := field.Tag.Get("db")
	if dbTag == "" {
		return strings.ToLower(field.Name)
	}
	return dbTag
}

func (r *LogRepository) getInsertFields(entity *models.Log) ([]string, []interface{}) {
	var fields []string
	var values []interface{}
	
	val := reflect.ValueOf(entity).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)
		
		// Пропускаем ID и автоматически генерируемые поля
		if field.Name == "ID" || field.Name == "CreatedAt" || field.Name == "UpdatedAt" {
			continue
		}
		
		// Пропускаем нулевые значения, если не требуется
		if !fieldValue.IsZero() {
			dbTag := r.getDBTag(field)
			fields = append(fields, dbTag)
			values = append(values, fieldValue.Interface())
		}
	}
	
	return fields, values
}

func (r *LogRepository) Create(ctx context.Context, entity *models.Log) error {
	fields, values := r.getInsertFields(entity)
	
	if len(fields) == 0 {
		return fmt.Errorf("no fields to insert")
	}
	
	// Создаем плейсхолдеры
	placeholders := make([]string, len(fields))
	for i := range fields {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", 
		r.tableName, 
		strings.Join(fields, ", "), 
		strings.Join(placeholders, ", "))
	
	_, err := r.pool.Exec(ctx, query, values...)
	return err
}

func (r *LogRepository) Update(ctx context.Context, id int, updates interface{}) error {
	var setClauses []string
	var values []interface{}

	switch u := updates.(type) {
	case map[string]interface{}:
		// Обработка карты полей
		if len(u) == 0 {
			return fmt.Errorf("no fields to update")
		}

		// Автоматическая проверка допустимых полей на основе модели
		val := reflect.ValueOf(&models.Log{}).Elem()
		typ := val.Type()
		allowedFields := make(map[string]bool)
		
		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			if field.Name != "ID" && field.Name != "CreatedAt" {
				dbTag := r.getDBTag(field)
				allowedFields[dbTag] = true
			}
		}

		for field, value := range u {
			if allowedFields[field] {
				setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, len(values)+1))
				values = append(values, value)
			}
		}

	case *models.Log:
		// Обработка частичной структуры
		val := reflect.ValueOf(u).Elem()
		typ := val.Type()

		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			fieldValue := val.Field(i)

			// Пропускаем ID и автоматически генерируемые поля
			if field.Name == "ID" || field.Name == "CreatedAt" {
				continue
			}

			// Проверяем, установлено ли поле (не нулевое значение)
			if !fieldValue.IsZero() {
				dbTag := r.getDBTag(field)
				setClauses = append(setClauses, fmt.Sprintf("%s = $%d", dbTag, len(values)+1))
				values = append(values, fieldValue.Interface())
			}
		}

	default:
		return fmt.Errorf("unsupported update type: %T", updates)
	}

	if len(setClauses) == 0 {
		return fmt.Errorf("no valid fields to update")
	}

	values = append(values, id)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", 
		r.tableName, 
		strings.Join(setClauses, ", "), 
		len(values))
	
	_, err := r.pool.Exec(ctx, query, values...)
	return err
}

func (r *LogRepository) Delete(ctx context.Context, id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", r.tableName)
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *LogRepository) List(ctx context.Context, limit, offset int) ([]models.Log, error) {
	query := fmt.Sprintf("SELECT %s FROM %s ORDER BY created_at DESC LIMIT $1 OFFSET $2", 
		strings.Join(r.getFields(), ", "), 
		r.tableName)
	
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Log
	for rows.Next() {
		var item models.Log
		
		// Автоматическое сканирование на основе полей модели
		scanValues := make([]interface{}, len(r.fields))
		val := reflect.ValueOf(&item).Elem()
		
		for i, fieldName := range r.fields {
			for j := 0; j < val.NumField(); j++ {
				field := val.Type().Field(j)
				dbTag := r.getDBTag(field)
				if dbTag == fieldName {
					scanValues[i] = val.Field(j).Addr().Interface()
					break
				}
			}
		}
		
		if err := rows.Scan(scanValues...); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return items, nil
}

func (r *LogRepository) Get(ctx context.Context, id int) (*models.Log, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", 
		strings.Join(r.fields, ", "), 
		r.tableName)
	
	row := r.pool.QueryRow(ctx, query, id)
	
	var item models.Log
	
	// Автоматическое сканирование на основе полей модели
	scanValues := make([]interface{}, len(r.fields))
	val := reflect.ValueOf(&item).Elem()
	
	for i, fieldName := range r.fields {
		for j := 0; j < val.NumField(); j++ {
			field := val.Type().Field(j)
			dbTag := r.getDBTag(field)
			if dbTag == fieldName {
				scanValues[i] = val.Field(j).Addr().Interface()
				break
			}
		}
	}
	
	err := row.Scan(scanValues...)
	if err != nil {
		return nil, fmt.Errorf("%s not found: %w", strings.ToLower("Log"), err)
	}
	
	return &item, nil
}

func (r *LogRepository) Count(ctx context.Context) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", r.tableName)
	var count int
	err := r.pool.QueryRow(ctx, query).Scan(&count)
	return count, err
}

// FindByUserID возвращает записи по ID пользователя
func (r *LogRepository) FindByUserID(ctx context.Context, userID int, limit, offset int) ([]models.Log, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3", 
		strings.Join(r.fields, ", "), 
		r.tableName)
	
	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Log
	for rows.Next() {
		var item models.Log
		
		scanValues := make([]interface{}, len(r.fields))
		val := reflect.ValueOf(&item).Elem()
		
		for i, fieldName := range r.fields {
			for j := 0; j < val.NumField(); j++ {
				field := val.Type().Field(j)
				dbTag := r.getDBTag(field)
				if dbTag == fieldName {
					scanValues[i] = val.Field(j).Addr().Interface()
					break
				}
			}
		}
		
		if err := rows.Scan(scanValues...); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return items, nil
}