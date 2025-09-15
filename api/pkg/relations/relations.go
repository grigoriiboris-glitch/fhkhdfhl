package relations

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

// QueryBuilder построитель запросов с отношениями
type QueryBuilder struct {
	model     Model
	relations map[string]Relation
	with      []string
}

// Model interface should include:
type Model interface {
	GetTableName() string
	GetConnection() *pgxpool.Pool
	SetAttributes(attrs map[string]interface{}) // Add this method
}

// HasMany отношение "один ко многим"
type HasMany struct {
	relatedModel func() Model
	config       RelationConfig
}

// Query возвращает QueryBuilder для отношения
func (h *HasMany) Query() *QueryBuilder {
	model := h.relatedModel()
	return NewQueryBuilder(model)
}

// NewHasMany создает новое отношение "один ко многим"
func NewHasMany(relatedModel func() Model, config RelationConfig) *HasMany {
	return &HasMany{
		relatedModel: relatedModel,
		config:       config,
	}
}

// Load загружает связанные модели
func (h *HasMany) Load(ctx context.Context, parent Model) error {
	related := h.relatedModel()
	
	// Определяем foreign key и local key
	foreignKey := h.config.ForeignKey
	if foreignKey == "" {
		foreignKey = fmt.Sprintf("%s_id", parent.GetTableName())
	}
	
	localKey := h.config.LocalKey
	if localKey == "" {
		localKey = "id"
	}
	
	// Получаем значение local key родителя
	// Для простоты предположим, что у родителя есть метод GetID()
	// В реальности нужно использовать рефлексию или другие методы
	
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", 
		related.GetTableName(), foreignKey)
	
	rows, err := parent.GetConnection().Query(ctx, query, 1) // Здесь нужно реальное значение ID
	if err != nil {
		return err
	}
	defer rows.Close()
	
	// Обрабатываем результаты
	for rows.Next() {
		model := h.relatedModel()
		attrs, err := scanRowToMap(rows)
		if err != nil {
			return err
		}
		model.SetAttributes(attrs)
		
		// Здесь нужно добавить модель в коллекцию родителя
		// Это требует рефлексии или конкретной реализации для каждого типа
	}
	
	return nil
}

func NewQueryBuilder(model Model) *QueryBuilder {
	return &QueryBuilder{
		model:     model,
		relations: make(map[string]Relation),
	}
}

// With eager loading отношений
func (q *QueryBuilder) With(relations ...string) *QueryBuilder {
	q.with = relations
	return q
}

// Where условие выборки
func (q *QueryBuilder) Where(column, operator string, value interface{}) *QueryBuilder {
	// Реализация условий
	return q
}

// Get получение модели с отношениями
func (q *QueryBuilder) Get(ctx context.Context, id int64) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", q.model.GetTableName())
	
	row := q.model.GetConnection().QueryRow(ctx, query, id) // Changed to QueryRow
	attrs, err := scanRowToMap(row)
	if err != nil {
		return err
	}
	
	q.model.SetAttributes(attrs)
	
	// Загрузка отношений
	for _, relationName := range q.with {
		if relation, exists := q.relations[relationName]; exists {
			if err := relation.Load(ctx, q.model); err != nil {
				return err
			}
		}
	}
	
	return nil
}

// First получение первой модели
func (q *QueryBuilder) First(ctx context.Context) error {
	query := fmt.Sprintf("SELECT * FROM %s LIMIT 1", q.model.GetTableName())
	
	row := q.model.GetConnection().QueryRow(ctx, query) // Changed to QueryRow
	attrs, err := scanRowToMap(row)
	if err != nil {
		return err
	}
	
	q.model.SetAttributes(attrs)
	return nil
}

// RegisterRelation регистрация отношения
func (q *QueryBuilder) RegisterRelation(name string, relation Relation) {
	q.relations[name] = relation
}