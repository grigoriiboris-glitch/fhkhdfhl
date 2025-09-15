package repository

import (
    "context"
    "database/sql"
    "fmt"
    "reflect"
    "strings"

    "github.com/mymindmap/api/models"
)

type LogRepository struct {
    db *sql.DB
}

func NewLogRepository(db *sql.DB) *LogRepository {
    return &LogRepository{db: db}
}

func (r *LogRepository) Create(ctx context.Context, entity *models.Log) error {
    var fields []string
    var placeholders []string
    var values []interface{}
    
    fields = append(fields, "title")
    placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)+1))
    values = append(values, entity.Title)
    
    fields = append(fields, "content")
    placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)+1))
    values = append(values, entity.Content)
    
    fields = append(fields, "user_id")
    placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)+1))
    values = append(values, entity.UserId)
    
    if !entity.CreatedAt.IsZero() {
        fields = append(fields, "created_at")
        placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)+1))
        values = append(values, entity.CreatedAt)
    }
    
    query := "INSERT INTO logs (" + strings.Join(fields, ", ") + ") VALUES (" + strings.Join(placeholders, ", ") + ")"
    _, err := r.db.ExecContext(ctx, query, values...)
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

        // Разрешенные поля для обновления
        allowedFields := map[string]bool{
            "title":   true,
            "content": true,
            "user_id": true,
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

            // Пропускаем ID и поля, которые не должны обновляться
            if field.Name == "ID" {
                continue
            }

            // Проверяем, установлено ли поле (не нулевое значение)
            if !fieldValue.IsZero() {
                dbTag := field.Tag.Get("db")
                if dbTag == "" {
                    dbTag = strings.ToLower(field.Name)
                }
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
    query := "UPDATE logs SET " + strings.Join(setClauses, ", ") + " WHERE id = $" + fmt.Sprintf("%d", len(values))
    _, err := r.db.ExecContext(ctx, query, values...)
    return err
}

func (r *LogRepository) Delete(ctx context.Context, id int) error {
    query := "DELETE FROM logs WHERE id = $1"
    _, err := r.db.ExecContext(ctx, query, id)
    return err
}

func (r *LogRepository) List(ctx context.Context, limit, offset int) ([]models.Log, error) {
    query := "SELECT id, title, content, user_id, created_at FROM logs ORDER BY created_at DESC LIMIT $1 OFFSET $2"
    rows, err := r.db.QueryContext(ctx, query, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []models.Log
    for rows.Next() {
        var e models.Log
        if err := rows.Scan(&e.ID, &e.Title, &e.Content, &e.UserId, &e.CreatedAt); err != nil {
            return nil, err
        }
        items = append(items, e)
    }
    
    if err := rows.Err(); err != nil {
        return nil, err
    }
    
    return items, nil
}

func (r *LogRepository) Get(ctx context.Context, id int) (*models.Log, error) {
    query := "SELECT id, title, content, user_id, created_at FROM logs WHERE id = $1"
    row := r.db.QueryRowContext(ctx, query, id)
    
    var log models.Log
    err := row.Scan(&log.ID, &log.Title, &log.Content, &log.UserId, &log.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("log not found")
        }
        return nil, err
    }
    
    return &log, nil
}

// Count возвращает общее количество записей
func (r *LogRepository) Count(ctx context.Context) (int, error) {
    query := "SELECT COUNT(*) FROM logs"
    var count int
    err := r.db.QueryRowContext(ctx, query).Scan(&count)
    return count, err
}

// FindByUserID возвращает логи по ID пользователя
func (r *LogRepository) FindByUserID(ctx context.Context, userID int, limit, offset int) ([]models.Log, error) {
    query := "SELECT id, title, content, user_id, created_at FROM logs WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3"
    rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []models.Log
    for rows.Next() {
        var e models.Log
        if err := rows.Scan(&e.ID, &e.Title, &e.Content, &e.UserId, &e.CreatedAt); err != nil {
            return nil, err
        }
        items = append(items, e)
    }
    
    if err := rows.Err(); err != nil {
        return nil, err
    }
    
    return items, nil
}