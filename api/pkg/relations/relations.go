package relations

import (
	"context"
	"fmt"
	//"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"reflect"
	"strings"
	"time"
)

// QueryBuilder построитель запросов с отношениями
type QueryBuilder struct {
	model     Model
	db        *pgxpool.Pool
	relations map[string]Relation
	with      []string
}

// Model interface should include:
type Model interface {
	GetTableName() string
	//SetAttributes(attrs map[string]interface{})
	GetID() interface{}
}

type ModelAdapter struct {
	model interface{}
}

type ModelWithRelations interface {
	Model
	GetRelations() map[string]Relation
}

func NewModelAdapter(model interface{}) *ModelAdapter {
	return &ModelAdapter{model: model}
}

func (a *ModelAdapter) GetID() int64 {
	val := reflect.ValueOf(a.model)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	
	idField := val.FieldByName("ID")
	if !idField.IsValid() {
		return 0
	}
	
	switch idField.Kind() {
	case reflect.Int:
		return int64(idField.Int())
	case reflect.Int32:
		return int64(idField.Int())
	case reflect.Int64:
		return idField.Int()
	default:
		return 0
	}
}

// BaseModel базовая модель
type BaseModel struct {
	ID        int64      `db:"id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (m *BaseModel) GetID() interface{} {
	return m.ID
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
	Load(ctx context.Context, model Model, db *pgxpool.Pool) error
	Query(db *pgxpool.Pool) *QueryBuilder
}

// HasMany отношение "один ко многим"
type HasMany struct {
	relatedModel func() Model
	config       RelationConfig
}

// Query возвращает QueryBuilder для отношения
func (h *HasMany) Query(db *pgxpool.Pool) *QueryBuilder {
	model := h.relatedModel()
	return NewQueryBuilder(model, db)
}

// NewHasMany создает новое отношение "один ко многим"
func NewHasMany(relatedModel func() Model, config RelationConfig) *HasMany {
	return &HasMany{
		relatedModel: relatedModel,
		config:       config,
	}
}

func NewQueryBuilder(model Model, db *pgxpool.Pool) *QueryBuilder {
	builder := &QueryBuilder{
		model:     model,
		db:        db,
		relations: make(map[string]Relation),
	}
	
	// Проверяем, реализует ли модель интерфейс ModelWithRelations
	if modelWithRelations, ok := model.(interface {
		GetRelations() map[string]Relation
	}); ok {
		fmt.Println("Model implements RelationProvider")
		relations := modelWithRelations.GetRelations()
		for name, relation := range relations {
			builder.relations[name] = relation
			fmt.Printf("Registered relation: %s\n", name)
		}
	}
	
	return builder
}

// With eager loading отношений
func (q *QueryBuilder) With(relations ...string) *QueryBuilder {
	for _, relation := range relations {
		// Проверяем, нет ли уже такого отношения в списке
		found := false
		for _, existing := range q.with {
			if existing == relation {
				found = true
				break
			}
		}
		if !found {
			q.with = append(q.with, relation)
		}
	}
	return q
}

// Where условие выборки
func (q *QueryBuilder) Where(column, operator string, value interface{}) *QueryBuilder {
	// Реализация условий
	return q
}

// Get получение модели с отношениями
// Get получение модели с отношениями (упрощенная версия)
func (q *QueryBuilder) Get(ctx context.Context, id interface{}) error {
	if q.db == nil {
		return fmt.Errorf("database connection is nil")
	}

	var idInt64 int64
	switch v := id.(type) {
	case int:
		idInt64 = int64(v)
	case int32:
		idInt64 = int64(v)
	case int64:
		idInt64 = v
	default:
		return fmt.Errorf("unsupported ID type: %T", id)
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", q.model.GetTableName())
	row := q.db.QueryRow(ctx, query, idInt64)

	// Используем рефлексию для создания срезов сканирования
	val := reflect.ValueOf(q.model).Elem()
	t := val.Type()

	// Собираем только поля с тегом db
	var scanFields []interface{}
	fieldMap := make(map[string]int)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag != "" && dbTag != "-" {
			scanFields = append(scanFields, val.Field(i).Addr().Interface())
			fieldMap[dbTag] = i
		}
	}

	// Сканируем напрямую в поля структуры
	if err := row.Scan(scanFields...); err != nil {
		return err
	}

	// Загрузка отношений
	for _, relationName := range q.with {
		if relation, exists := q.relations[relationName]; exists {
			if err := relation.Load(ctx, q.model, q.db); err != nil {
				return err
			}
		}
	}

	return nil
}

// Load загружает связанные модели
func (h *HasMany) Load(ctx context.Context, parent Model, db *pgxpool.Pool) error {
	related := h.relatedModel()

	foreignKey := h.config.ForeignKey
	if foreignKey == "" {
		foreignKey = fmt.Sprintf("%s_id", strings.ToLower(getModelName(parent)))
	}

	localKey := h.config.LocalKey
	if localKey == "" {
		localKey = "id"
	}

	// Используем адаптер для получения ID родительской модели
	adapter := NewModelAdapter(parent)
	parentID := adapter.GetID()

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", 
		related.GetTableName(), foreignKey)

	rows, err := db.Query(ctx, query, parentID)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Получаем описание полей
	fieldDescriptions := rows.FieldDescriptions()
	columns := make([]string, len(fieldDescriptions))
	for i, fd := range fieldDescriptions {
		columns[i] = fd.Name
	}

	// Создаем слайс для хранения связанных моделей
	var relatedModels []Model

	for rows.Next() {
		model := h.relatedModel()

		// Создаем срез для сканирования значений
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return err
		}

		// Создаем map атрибутов
		attrs := make(map[string]interface{})
		for i, col := range columns {
			attrs[col] = values[i]
		}

		//model.SetAttributes(attrs)
		relatedModels = append(relatedModels, model)
	}

	// Устанавливаем связанные модели в родительскую модель
	if err := setRelatedModels(parent, relatedModels, getModelName(related)); err != nil {
		return err
	}

	return nil
}

// setRelatedModels устанавливает связанные модели в родительскую модель используя рефлексию
func setRelatedModels(parent Model, relatedModels []Model, modelName string) error {
	if len(relatedModels) == 0 {
		return nil // Нет моделей для установки
	}
	
	parentVal := reflect.ValueOf(parent).Elem()
	
	// Преобразуем имя модели в имя поля (например: "MindMap" -> "MindMaps")
	// Попробуем разные варианты имен полей
	possibleFieldNames := []string{
		modelName + "s",    // MindMaps
		modelName + "List", // MindMapList
		modelName,          // MindMap
		strings.ToLower(modelName) + "s", // mindmaps
	}
	
	var field reflect.Value
	var fieldName string
	
	// Ищем существующее поле
	for _, name := range possibleFieldNames {
		field = parentVal.FieldByName(name)
		if field.IsValid() {
			fieldName = name
			break
		}
	}
	
	if !field.IsValid() {
		// Если поле не найдено, попробуем поискать по тегам
		t := parentVal.Type()
		for i := 0; i < t.NumField(); i++ {
			structField := t.Field(i)
			if structField.Tag.Get("db") == "-" {
				// Это может быть поле для отношений
				field = parentVal.Field(i)
				fieldName = structField.Name
				break
			}
		}
	}
	
	if !field.IsValid() {
		return fmt.Errorf("relation field for %s not found in parent model", modelName)
	}
	
	if !field.CanSet() {
		return fmt.Errorf("field %s cannot be set", fieldName)
	}
	
	if field.Kind() != reflect.Slice {
		return fmt.Errorf("field %s is not a slice", fieldName)
	}
	
	// Получаем тип элементов слайса
	elemType := field.Type().Elem()
	
	// Создаем новый слайс
	newSlice := reflect.MakeSlice(field.Type(), len(relatedModels), len(relatedModels))
	
	for i, relatedModel := range relatedModels {
		relatedVal := reflect.ValueOf(relatedModel)
		if relatedVal.Type().AssignableTo(elemType) {
			newSlice.Index(i).Set(relatedVal)
		} else {
			// Попробуем разыменовать указатель
			if relatedVal.Kind() == reflect.Ptr && relatedVal.Elem().Type().AssignableTo(elemType) {
				newSlice.Index(i).Set(relatedVal.Elem())
			} else if elemType.Kind() == reflect.Ptr && relatedVal.Type().AssignableTo(elemType.Elem()) {
				// Создаем указатель на значение
				ptr := reflect.New(elemType.Elem())
				ptr.Elem().Set(relatedVal)
				newSlice.Index(i).Set(ptr)
			} else {
				return fmt.Errorf("type mismatch: cannot assign %s to %s", 
					relatedVal.Type(), elemType)
			}
		}
	}
	
	field.Set(newSlice)
	return nil
}

// First получение первой модели
func (q *QueryBuilder) First(ctx context.Context) error {
	if q.db == nil {
		return fmt.Errorf("database connection is nil")
	}

	query := fmt.Sprintf("SELECT * FROM %s LIMIT 1", q.model.GetTableName())
	
	conn, err := q.db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Получаем описание полей
	fieldDescriptions := rows.FieldDescriptions()
	columns := make([]string, len(fieldDescriptions))
	for i, fd := range fieldDescriptions {
		columns[i] = fd.Name
	}

	// Создаем срезы для сканирования
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// Сканируем первую строку
	if rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no records found")
	}

	// Создаем map атрибутов
	attrs := make(map[string]interface{})
	for i, col := range columns {
		attrs[col] = values[i]
	}

	//q.model.SetAttributes(attrs)
	return nil
}

// SetAttributes устанавливает атрибуты только для полей с тегом db
// func (m *BaseModel) SetAttributes(attrs map[string]interface{}) {
// 	val := reflect.ValueOf(m).Elem()
// 	t := val.Type()

// 	for i := 0; i < t.NumField(); i++ {
// 		field := t.Field(i)
// 		dbTag := field.Tag.Get("db")
		
// 		// Устанавливаем значение только для полей с тегом db
// 		if dbTag != "" && dbTag != "-" {
// 			if value, exists := attrs[dbTag]; exists {
// 				fieldValue := val.Field(i)
// 				if fieldValue.IsValid() && fieldValue.CanSet() {
// 					// Конвертируем значение к правильному типу
// 					convertedValue, err := convertValue(value, fieldValue.Type())
// 					if err == nil {
// 						fieldValue.Set(convertedValue)
// 					}
// 				}
// 			}
// 		}
// 	}
// }

// convertValue конвертирует значение к нужному типу
func convertValue(value interface{}, targetType reflect.Type) (reflect.Value, error) {
	val := reflect.ValueOf(value)
	
	if val.Type().ConvertibleTo(targetType) {
		return val.Convert(targetType), nil
	}
	
	// Обработка специальных случаев
	switch targetType.Kind() {
	case reflect.Ptr:
		if val.IsValid() && !val.IsNil() {
			elemValue, err := convertValue(val.Elem().Interface(), targetType.Elem())
			if err != nil {
				return reflect.Value{}, err
			}
			ptr := reflect.New(targetType.Elem())
			ptr.Elem().Set(elemValue)
			return ptr, nil
		}
		return reflect.Zero(targetType), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v := value.(type) {
		case int:
			return reflect.ValueOf(v).Convert(targetType), nil
		case int32:
			return reflect.ValueOf(v).Convert(targetType), nil
		case int64:
			return reflect.ValueOf(v).Convert(targetType), nil
		case float64:
			return reflect.ValueOf(int64(v)).Convert(targetType), nil
		}
	case reflect.String:
		if v, ok := value.(string); ok {
			return reflect.ValueOf(v), nil
		}
		return reflect.ValueOf(fmt.Sprintf("%v", value)), nil
	}
	
	return reflect.Value{}, fmt.Errorf("cannot convert %T to %v", value, targetType)
}

// GetInt64 получение модели с отношениями (специфичная версия для int64)
func (q *QueryBuilder) GetInt64(ctx context.Context, id int64) error {
	return q.Get(ctx, id)
}

// GetInt получение модели с отношениями (специфичная версия для int)
func (q *QueryBuilder) GetInt(ctx context.Context, id int) error {
	return q.Get(ctx, id)
}

// RegisterRelation регистрация отношения
func (q *QueryBuilder) RegisterRelation(name string, relation Relation) {
	q.relations[name] = relation
}

// Вспомогательные функции
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
	
	// Используем правильное преобразование имени поля
	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		// Попробуем с заглавной буквой
		fieldName = strings.Title(fieldName)
		field = val.FieldByName(fieldName)
		if !field.IsValid() {
			return nil
		}
	}
	
	return field.Interface()
}