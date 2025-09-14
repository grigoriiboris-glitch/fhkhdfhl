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
	"github.com/joho/godotenv"
)

type Model struct {
	Name        string
	Fields      []Field
	Module      string
	VarName     string
	PackageName string
}

type Field struct {
	Name string
	Type string
	Tag  string
}

type RouteData struct {
	Module string
	Models []RouteModel
}

type RouteModel struct {
	Name    string
	VarName string
}

var allModels []Model
var moduleName string

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}
	
	moduleName = os.Getenv("APP_PACKAGE_NAME")
	if moduleName == "" {
		log.Fatal("APP_PACKAGE_NAME environment variable is required")
	}
	
	modelsDir := "/app/models0"
	err = filepath.Walk(modelsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			model, err := parseModel(path)
			if err != nil {
				return err
			}
			if model != nil {
				allModels = append(allModels, *model)
				generateCRUD(*model)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	generateRoutes(allModels)
}

func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

func parseModel(path string) (*Model, error) {
	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, path, nil, parser.AllErrors)
	if err != nil {
		return nil, err
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

			model := &Model{
				Name:        ts.Name.Name,
				Module:      os.Getenv("APP_PACKAGE_NAME"),
				VarName:     strings.ToLower(ts.Name.Name[:1]) + ts.Name.Name[1:],
				PackageName: strings.ToLower(ts.Name.Name),
			}

			for _, f := range st.Fields.List {
				if len(f.Names) > 0 && f.Names[0].IsExported() {
					fieldName := f.Names[0].Name
					fieldType := exprToString(f.Type)
					
					// Skip ID, CreatedAt, UpdatedAt for create requests
					if fieldName == "ID" || fieldName == "CreatedAt" || fieldName == "UpdatedAt" {
						continue
					}
					
					model.Fields = append(model.Fields, Field{
						Name: fieldName,
						Type: fieldType,
						Tag:  generateJSONTag(fieldName),
					})
				}
			}
			return model, nil
		}
	}
	return nil, nil
}

func exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return exprToString(t.X) + "." + t.Sel.Name
	case *ast.StarExpr:
		return "*" + exprToString(t.X)
	case *ast.ArrayType:
		return "[]" + exprToString(t.Elt)
	default:
		return fmt.Sprintf("%v", expr)
	}
}

func generateJSONTag(fieldName string) string {
	return fmt.Sprintf("`json:\"%s\" validate:\"required\"`", toSnakeCase(fieldName))
}

func generateCRUD(model Model) {
	files := map[string]string{
		fmt.Sprintf("/app/internal/http/handlers/%s_handler.go", strings.ToLower(model.Name)): "handler.tmpl",
		fmt.Sprintf("/app/internal/http/requests/%s/create_%s_request.go", strings.ToLower(model.Name), strings.ToLower(model.Name)): "create_request.tmpl",
		fmt.Sprintf("/app/internal/http/requests/%s/update_%s_request.go", strings.ToLower(model.Name), strings.ToLower(model.Name)): "update_request.tmpl",
		fmt.Sprintf("/app/internal/http/requests/%s/list_%s_request.go", strings.ToLower(model.Name), strings.ToLower(model.Name)): "list_request.tmpl",
		fmt.Sprintf("/app/internal/services/%s/%s_service.go", strings.ToLower(model.Name), strings.ToLower(model.Name)): "service.tmpl",
		fmt.Sprintf("/app/repository/%s_repository.go", strings.ToLower(model.Name)): "repository.tmpl",
	}

	funcMap := template.FuncMap{
		"lower":  strings.ToLower,
		"title":  strings.Title,
		"snake":  toSnakeCase,
	}

	for path, tmplFile := range files {
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("⚠️  file already exists, skipping: %s\n", path)
			//continue
		}

		tmpl, err := template.New(filepath.Base(tmplFile)).Funcs(funcMap).ParseFiles("/app/scripts/templates/" + tmplFile)
		if err != nil {
			log.Printf("Error parsing template %s: %v", tmplFile, err)
			continue
		}

		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Error creating directory %s: %v", dir, err)
			continue
		}

		f, err := os.Create(path)
		if err != nil {
			log.Printf("Error creating file %s: %v", path, err)
			continue
		}

		if err := tmpl.Execute(f, model); err != nil {
			log.Printf("Error executing template for %s: %v", path, err)
			f.Close()
			continue
		}

		f.Close()
		fmt.Println("✅ generated:", path)
	}
}

func generateRoutes(models []Model) {
	routeData := RouteData{
		Module: os.Getenv("APP_PACKAGE_NAME"),
	}
	for _, m := range models {
		routeData.Models = append(routeData.Models, RouteModel{
			Name:    m.Name,
			VarName: strings.ToLower(m.Name[:1]) + m.Name[1:],
		})
	}

	funcMap := template.FuncMap{
		"lower": strings.ToLower,
		"title": strings.Title,
	}

	tmpl, err := template.New("api.tmpl").Funcs(funcMap).ParseFiles("/app/scripts/templates/api.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	os.MkdirAll("/app/internal/http/routes", 0755)
	apiPath := "/app/tmp/api_.go"

	f, err := os.Create(apiPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, routeData); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✅ generated:", apiPath)
}