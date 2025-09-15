package main

import (
	"encoding/json"
	"fmt"
	"log"

	"yourpackage/filters"
)

func main() {
	// Пример 1: Простые фильтры
	simpleFilters := map[string]interface{}{
		"name":        "John Doe",
		"age":         25,
		"email":       "john@example.com",
		"is_active":   true,
		"order_by":    "created_at:desc",
		"row_limit":   10,
		"skip":        0,
	}

	fmt.Println("=== Простые фильтры ===")
	printFilters(simpleFilters)

	// Пример 2: Массивы значений (IN)
	arrayFilters := map[string]interface{}{
		"status":      []interface{}{"active", "pending", "completed"},
		"category_id": []interface{}{1, 2, 3, 5, 8},
		"role":        []interface{}{"admin", "moderator"},
		"order_by":    "name:asc",
		"row_limit":   50,
	}

	fmt.Println("\n=== Массивы значений (IN) ===")
	printFilters(arrayFilters)

	// Пример 3: Специальные строковые операторы
	stringOperators := map[string]interface{}{
		"name":        "LIKE john",           // ILIKE '%john%'
		"description": "LIKE important task", // ILIKE '%important task%'
		"deleted_at":  "IS NULL",             // IS NULL
		"updated_at":  "IS NOT NULL",         // IS NOT NULL
		"tags":        "in__go,backend,api",  // IN ('go', 'backend', 'api')
		"excluded_ids": "not_in__1,5,9",      // NOT IN (1, 5, 9)
		"age":         ">= 18",               // >= 18
		"salary":      "< 5000",              // < 5000
		"rating":      "!= 0",                // != 0
	}

	fmt.Println("\n=== Строковые операторы ===")
	printFilters(stringOperators)

	// Пример 4: Комплексные OR условия через JSON
	orConditions := `[
		{"column": "status", "value": "active", "rule": "="},
		{"column": "status", "value": "pending", "rule": "="},
		{"column": "age", "value": 25, "rule": ">"},
		{"column": "department", "value": "IT", "rule": "=", "method": "OR"}
	]`

	complexFilters := map[string]interface{}{
		"name":        "John",
		"or_conditions": orConditions,
		"order_by":    "created_at:desc",
		"paginate":    20,
	}

	fmt.Println("\n=== OR условия через JSON ===")
	printFilters(complexFilters)

	// Пример 5: Комбинация всех типов фильтров
	combinedFilters := map[string]interface{}{
		// Простые равенства
		"company_id": 123,
		"is_deleted": false,
		
		// Массивы
		"department": []interface{}{"IT", "HR", "Sales"},
		"level":      []interface{}{"junior", "middle"},
		
		// Строковые операторы
		"email":      "LIKE @company.com",
		"start_date": ">= 2024-01-01",
		"end_date":   "<= 2024-12-31",
		"notes":      "IS NOT NULL",
		
		// OR условия
		"complex_or": `[
			{"column": "status", "value": "active", "rule": "="},
			{"column": "status", "value": "on_hold", "rule": "="},
			{"column": "priority", "value": "high", "rule": "="},
			{"column": "category", "value": ["urgent", "important"], "rule": "IN"}
		]`,
		
		// Пагинация и сортировка
		"order_by":   "salary:desc, name:asc",
		"row_limit":  25,
		"skip":       0,
		"paginate":   10,
	}

	fmt.Println("\n=== Комбинация всех фильтров ===")
	printFilters(combinedFilters)

	// Пример 6: Фильтры с датами и временем
	dateFilters := map[string]interface{}{
		"created_at":  ">= 2024-01-01 00:00:00",
		"updated_at":  "< 2024-12-31 23:59:59",
		"deleted_at":  "IS NULL",
		"event_date":  "BETWEEN 2024-06-01 AND 2024-06-30",
		"order_by":    "created_at:desc",
		"row_limit":   100,
	}

	fmt.Println("\n=== Фильтры с датами ===")
	printFilters(dateFilters)

	// Пример 7: Фильтры с числовыми диапазонами
	numericRangeFilters := map[string]interface{}{
		"age":         "BETWEEN 18 AND 65",
		"salary":      ">= 3000",
		"salary":      "<= 10000",
		"rating":      "> 4.0",
		"experience":  ">= 2",
		"order_by":    "salary:desc",
	}

	fmt.Println("\n=== Числовые диапазоны ===")
	printFilters(numericRangeFilters)

	// Пример 8: Фильтры с IN и NOT IN через массивы
	inNotInFilters := map[string]interface{}{
		"status":      []interface{}{"active", "pending", "review"},
		"category":    []interface{}{"electronics", "books", "clothing"},
		"excluded_ids": []interface{}{1, 2, 3, 4, 5},
		"role":        "not_in__guest,blocked,banned",
		"tags":        "in__go,postgres,backend",
	}

	fmt.Println("\n=== IN/NOT IN фильтры ===")
	printFilters(inNotInFilters)
}

// Вспомогательная функция для вывода фильтров в формате JSON
func printFilters(filters map[string]interface{}) {
	jsonData, err := json.MarshalIndent(filters, "", "  ")
	if err != nil {
		log.Printf("Error marshaling filters: %v", err)
		return
	}
	fmt.Println(string(jsonData))
}

// Пример использования с базой данных
func exampleWithDB() {
	// db, err := sql.Open("postgres", "your-connection-string")
	// if err != nil {
	//     log.Fatal(err)
	// }
	// defer db.Close()

	// filterService := filters.NewFilterTableService(db)

	// Пример выполнения запроса с комплексными фильтрами
	filters := map[string]interface{}{
		"name":        "LIKE john",
		"age":         []interface{}{25, 30, 35},
		"status":      "active",
		"department":  "IT",
		"or_conditions": `[
			{"column": "role", "value": "admin", "rule": "="},
			{"column": "role", "value": "moderator", "rule": "="},
			{"column": "permissions", "value": "write", "rule": "="}
		]`,
		"order_by":  "created_at:desc",
		"row_limit": 20,
		"skip":      0,
	}

	fmt.Println("\n=== Пример для выполнения в БД ===")
	printFilters(filters)

	// rows, err := filterService.GetFilteredResults("users", filters)
	// if err != nil {
	//     log.Fatal(err)
	// }
	// defer rows.Close()
	
	// // Обработка результатов...
}