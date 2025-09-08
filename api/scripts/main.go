package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Model struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name string
	Type string
}

func main() {
	modelsDir := "../models"
	
	log.Println("Запуск генератора CRUD...")
	err := filepath.Walk(modelsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			log.Printf("Парсинг файла: %s", path)
			models, err := parseModels(path)
			if err != nil {
				return err
			}
			if len(models) == 0 {
				log.Printf("  В файле %s структур не найдено", path)
			}
			for _, m := range models {
				log.Printf("  Найдена структура: %+v", m)
				generateCRUD(m)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Генерация завершена ✅")
}

// parseModels возвращает все структуры из файла
func parseModels(path string) ([]Model, error) {
	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, path, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	var models []Model
	// Берём только структуры из package "models"
	if node.Name.Name != "models" {
		return models, nil
	}

	for _, decl := range node.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range gen.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}

			model := Model{Name: ts.Name.Name}
			for _, f := range st.Fields.List {
				if len(f.Names) > 0 {
					fieldName := f.Names[0].Name
					fieldType := fmt.Sprint(f.Type)
					model.Fields = append(model.Fields, Field{
						Name: fieldName,
						Type: fieldType,
					})
				}
			}
			models = append(models, model)
		}
	}
	return models, nil
}


func generateCRUD(model Model) {
	files := map[string]string{
		fmt.Sprintf("../internal/http/handlers/%s_handler.go", strings.ToLower(model.Name)): "handler.tmpl",
		fmt.Sprintf("../internal/http/requests/%s_requests.go", strings.ToLower(model.Name)): "requests.tmpl",
		fmt.Sprintf("../internal/services/%s_service.go", strings.ToLower(model.Name)):      "service.tmpl",
		fmt.Sprintf("../repository/%s_repository.go", strings.ToLower(model.Name)):          "repository.tmpl",
	}

	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}

	for path, tmplFile := range files {
		log.Printf("  Генерация файла: %s (шаблон: %s)", path, tmplFile)

		tmpl, err := template.New(tmplFile).Funcs(funcMap).ParseFiles("./templates/" + tmplFile)
		if err != nil {
			log.Fatalf("  Ошибка парсинга шаблона %s: %v", tmplFile, err)
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			log.Fatalf("  Ошибка создания директории для %s: %v", path, err)
		}

		f, err := os.Create(path)
		if err != nil {
			log.Fatalf("  Ошибка создания файла %s: %v", path, err)
		}

		if err := tmpl.Execute(f, model); err != nil {
			log.Fatalf("  Ошибка генерации файла %s: %v", path, err)
		}
		f.Close()

		log.Printf("  ✅ Сгенерирован: %s", path)
	}
}
