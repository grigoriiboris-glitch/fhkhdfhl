// package auth_handler

// import (
// 	"context"
// 	"embed"
// 	"encoding/json"
// 	"fmt"
// 	"html/template")
    

// var postFormTmpl = template.Must(template.ParseFS(
// 	templates,
// 	filepath.Join(
// 		"templates",
// 		"default.html",
// 	),
// 	filepath.Join("templates", "create-post.html"),
// ))

// // Auth handlers
// func loginHandler(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case "GET":
// 		tmpl := template.Must(template.ParseFS(
// 			templates,
// 			"templates/login.html",
// 		))
// 		tmpl.ExecuteTemplate(w, "login", nil)
// 	case "POST":
// 		if err := r.ParseForm(); err != nil {
// 			http.Error(w, "Failed to parse form", http.StatusBadRequest)
// 			return
// 		}

// 		email := r.FormValue("email")
// 		password := r.FormValue("password")

// 		if email == "" || password == "" {
// 			renderLoginPage(w, "Email и пароль обязательны", "")
// 			return
// 		}

// 		req := &models.LoginRequest{
// 			Email:    email,
// 			Password: password,
// 		}

// 		token, err := authService.LoginUser(r.Context(), req)
// 		if err != nil {
// 			renderLoginPage(w, err.Error(), "")
// 			return
// 		}

// 		// Устанавливаем cookie с токеном
// 		authService.SetAuthCookie(w, token)

// 		// Перенаправляем на главную страницу
// 		http.Redirect(w, r, "/", http.StatusSeeOther)
// 	default:
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 	}
// }

// func registerHandler(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case "GET":
// 		tmpl := template.Must(template.ParseFS(
// 			templates,
// 			"templates/register.html",
// 		))
// 		tmpl.ExecuteTemplate(w, "register", nil)
// 	case "POST":
// 		if err := r.ParseForm(); err != nil {
// 			http.Error(w, "Failed to parse form", http.StatusBadRequest)
// 			return
// 		}

// 		name := r.FormValue("name")
// 		email := r.FormValue("email")
// 		password := r.FormValue("password")
// 		confirmPassword := r.FormValue("confirm_password")

// 		if name == "" || email == "" || password == "" {
// 			renderRegisterPage(w, "Все поля обязательны", "")
// 			return
// 		}

// 		if password != confirmPassword {
// 			renderRegisterPage(w, "Пароли не совпадают", "")
// 			return
// 		}

// 		if len(password) < 6 {
// 			renderRegisterPage(w, "Пароль должен содержать минимум 6 символов", "")
// 			return
// 		}

// 		req := &models.RegisterRequest{
// 			Name:     name,
// 			Email:    email,
// 			Password: password,
// 		}

// 		user, err := authService.RegisterUser(r.Context(), req)
// 		if err != nil {
// 			renderRegisterPage(w, err.Error(), "")
// 			return
// 		}

// 		renderRegisterPage(w, "", fmt.Sprintf("Пользователь %s успешно зарегистрирован! Теперь вы можете войти.", user.Name))
// 	default:
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 	}
// }

// func renderLoginPage(w http.ResponseWriter, errorMsg, successMsg string) {
// 	data := struct {
// 		Error   string
// 		Success string
// 	}{
// 		Error:   errorMsg,
// 		Success: successMsg,
// 	}

// 	tmpl := template.Must(template.ParseFS(
// 		templates,
// 		"templates/login.html",
// 	))
// 	tmpl.ExecuteTemplate(w, "login", data)
// }

// func renderRegisterPage(w http.ResponseWriter, errorMsg, successMsg string) {
// 	data := struct {
// 		Error   string
// 		Success string
// 	}{
// 		Error:   errorMsg,
// 		Success: successMsg,
// 	}

// 	tmpl := template.Must(template.ParseFS(
// 		templates,
// 		"templates/register.html",
// 	))
// 	tmpl.ExecuteTemplate(w, "register", data)
// }