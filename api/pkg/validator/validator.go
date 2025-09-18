package validator

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

var instance *validator.Validate

// Init инициализирует глобальный валидатор
func Init() {
	instance = validator.New()
	
	// Настраиваем отображение имен полей JSON
	instance.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	
	// Здесь можно добавить кастомные валидации
	// instance.RegisterValidation("custom_validation", customValidationFunc)
}

// Get возвращает экземпляр валидатора
func Get() *validator.Validate {
	if instance == nil {
		panic("validator not initialized. Call validator.Init() first")
	}
	return instance
}

// Validate выполняет валидацию структуры
func Validate(s interface{}) error {
	return Get().Struct(s)
}