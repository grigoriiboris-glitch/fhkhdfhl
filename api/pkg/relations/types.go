package relations

import (
	"context"
	//"database/sql"
	//"fmt"
	"reflect"
	"strings"
	"time"
)


// BaseModel базовая модель
type BaseModel struct {
	ID        int64      `db:"id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (m *BaseModel) GetID() int64 {
	return m.ID
}

func (m *BaseModel) SetAttributes(attrs map[string]interface{}) {
	// Реализация установки атрибутов через рефлексию
	val := reflect.ValueOf(m).Elem()
	for key, value := range attrs {
		field := val.FieldByName(strings.Title(key))
		if field.IsValid() && field.CanSet() {
			field.Set(reflect.ValueOf(value))
		}
	}
}

// RelationConfig конфигурация отношения
type RelationConfig struct {
	ForeignKey  string
	LocalKey    string
	PivotTable  string
	PivotLocal  string
	PivotForeign string
}

// Relation интерфейс отношения
type Relation interface {
	Load(ctx context.Context, model Model) error
	Query() *QueryBuilder
}