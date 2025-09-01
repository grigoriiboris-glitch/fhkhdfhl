package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mymindmap/api/auth"
	"github.com/mymindmap/api/models"
	"github.com/mymindmap/api/repository"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

//go:embed templates/*
var templates embed.FS

var postRepo *repository.PostRepository
var userRepo *repository.UserRepository
var mindMapRepo *repository.MindMapRepository
var authService *auth.AuthService

type Config struct {
	PostgresDB       string
	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresURL      string
}

var conf Config

var postFormTmpl = template.Must(template.ParseFS(
	templates,
	filepath.Join(
		"templates",
		"default.html",
	),
	filepath.Join("templates", "create-post.html"),
))

func init() {
	godotenv.Load()

	conf.PostgresDB = os.Getenv("POSTGRES_DB")
	conf.PostgresUser = os.Getenv("POSTGRES_USER")
	conf.PostgresPassword = os.Getenv("POSTGRES_PASSWORD")
	conf.PostgresHost = os.Getenv("POSTGRES_HOST")

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		conf.PostgresUser,
		conf.PostgresPassword,
		conf.PostgresHost,
		conf.PostgresDB,
	)

	conf.PostgresURL = connStr
}

// Database connection using pgxpool.
func main() {
	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, conf.PostgresURL)
	if err != nil {
		log.Fatal("Unable to create database connection pool:", err)
	}

	defer dbpool.Close()

	err = dbpool.Ping(ctx)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	postRepo = repository.NewPostRepository(dbpool)
	userRepo = repository.NewUserRepository(dbpool)
	mindMapRepo = repository.NewMindMapRepository(dbpool)

	// Инициализируем сервис авторизации
	authService, err = auth.NewAuthService(userRepo)
	if err != nil {
		log.Fatal("Unable to create auth service:", err)
	}

	// Создаем общий middleware для авторизации
	authMiddleware := authService.AuthMiddleware

	// Маршруты авторизации (без middleware авторизации)
	http.HandleFunc("/auth/login", loginHandler)
	http.HandleFunc("/auth/register", registerHandler)
	http.HandleFunc("/auth/logout", logoutHandler)

	// Главная страница с постами
	http.HandleFunc("/", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getPosts(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// Создание постов (требует авторизации)
	http.HandleFunc("/posts/new", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			// Применяем middleware для проверки прав на создание постов
			authService.RequirePermission("post", "write")(createPost)(w, r)
		case "GET":
			// Применяем middleware для проверки прав на создание постов
			authService.RequirePermission("post", "write")(func(w http.ResponseWriter, r *http.Request) {
				// Получаем информацию о пользователе из контекста
				user := auth.GetUserFromContext(r.Context())
				
				postFormTmpl.ExecuteTemplate(
					w,
					"default",
					struct{ 
						Post models.Post
						User *auth.Claims
					}{models.Post{}, user},
				)
			})(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// Редактирование постов (требует авторизации)
	http.HandleFunc(
		"/post/{id}/edit",
		authMiddleware(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "GET":
				authService.RequirePermission("post", "write")(func(w http.ResponseWriter, r *http.Request) {
					post, err := getPost(r)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}

					// Получаем информацию о пользователе из контекста
					user := auth.GetUserFromContext(r.Context())

					postFormTmpl.ExecuteTemplate(
						w,
						"default",
						struct{ 
							Post *models.Post
							User *auth.Claims
						}{post, user},
					)
				})(w, r)
			case "POST":
				authService.RequirePermission("post", "write")(updatePost)(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}),
	)

	// Удаление постов (требует авторизации)
	http.HandleFunc(
		"/post/{id}/delete",
		authMiddleware(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "GET":
				authService.RequirePermission("post", "delete")(deletePost)(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}),
	)

	// Просмотр постов (доступно всем, но с информацией о пользователе)
	http.HandleFunc("/post/{id}", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			post, err := getPost(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Получаем информацию о пользователе из контекста
			user := auth.GetUserFromContext(r.Context())

			tmpl := template.Must(template.ParseFS(
				templates,
				filepath.Join(
					"templates",
					"default.html",
				),
				filepath.Join("templates", "post.html"),
			))

			tmpl.ExecuteTemplate(
				w,
				"default",
				struct{ 
					Post *models.Post
					User *auth.Claims
				}{post, user},
			)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// API endpoints для ментальных карт
	http.HandleFunc("/api/mindmaps", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			// Получаем карты пользователя
			user := auth.GetUserFromContext(r.Context())
			mindMaps, err := mindMapRepo.GetByUserID(r.Context(), user.UserID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mindMaps)
		case "POST":
			// Создание новой карты
			user := auth.GetUserFromContext(r.Context())
			
			var req models.CreateMindMapRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
			
			mindMap := &models.MindMap{
				Title:    req.Title,
				Data:     req.Data,
				UserID:   user.UserID,
				IsPublic: req.IsPublic,
			}
			
			if err := mindMapRepo.Create(r.Context(), mindMap); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(mindMap)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// API endpoint для конкретной карты
	http.HandleFunc("/api/mindmaps/", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем ID из URL
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 {
			http.Error(w, "Invalid mindmap ID", http.StatusBadRequest)
			return
		}
		
		idStr := pathParts[3]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid mindmap ID", http.StatusBadRequest)
			return
		}
		
		user := auth.GetUserFromContext(r.Context())
		
		switch r.Method {
		case "GET":
			// Получение карты
			mindMap, err := mindMapRepo.GetByID(r.Context(), id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if mindMap == nil {
				http.Error(w, "Mindmap not found", http.StatusNotFound)
				return
			}
			
			// Проверяем права доступа
			if mindMap.UserID != user.UserID && !mindMap.IsPublic {
				http.Error(w, "Access denied", http.StatusForbidden)
				return
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mindMap)
		case "PUT":
			// Обновление карты
			mindMap, err := mindMapRepo.GetByID(r.Context(), id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if mindMap == nil {
				http.Error(w, "Mindmap not found", http.StatusNotFound)
				return
			}
			
			// Проверяем права доступа
			if mindMap.UserID != user.UserID {
				http.Error(w, "Access denied", http.StatusForbidden)
				return
			}
			
			var req models.UpdateMindMapRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
			
			mindMap.Title = req.Title
			mindMap.Data = req.Data
			mindMap.IsPublic = req.IsPublic
			
			if err := mindMapRepo.Update(r.Context(), mindMap); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mindMap)
		case "DELETE":
			// Удаление карты
			mindMap, err := mindMapRepo.GetByID(r.Context(), id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if mindMap == nil {
				http.Error(w, "Mindmap not found", http.StatusNotFound)
				return
			}
			
			// Проверяем права доступа
			if mindMap.UserID != user.UserID {
				http.Error(w, "Access denied", http.StatusForbidden)
				return
			}
			
			if err := mindMapRepo.Delete(r.Context(), id, user.UserID); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// API endpoint для публичных карт
	http.HandleFunc("/api/mindmaps/public", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			mindMaps, err := mindMapRepo.GetPublic(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mindMaps)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Write([]byte("OK"))
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server started on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

// Auth handlers
func loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tmpl := template.Must(template.ParseFS(
			templates,
			"templates/login.html",
		))
		tmpl.ExecuteTemplate(w, "login", nil)
	case "POST":
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			renderLoginPage(w, "Email и пароль обязательны", "")
			return
		}

		req := &models.LoginRequest{
			Email:    email,
			Password: password,
		}

		token, err := authService.LoginUser(r.Context(), req)
		if err != nil {
			renderLoginPage(w, err.Error(), "")
			return
		}

		// Устанавливаем cookie с токеном
		authService.SetAuthCookie(w, token)

		// Перенаправляем на главную страницу
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tmpl := template.Must(template.ParseFS(
			templates,
			"templates/register.html",
		))
		tmpl.ExecuteTemplate(w, "register", nil)
	case "POST":
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")

		if name == "" || email == "" || password == "" {
			renderRegisterPage(w, "Все поля обязательны", "")
			return
		}

		if password != confirmPassword {
			renderRegisterPage(w, "Пароли не совпадают", "")
			return
		}

		if len(password) < 6 {
			renderRegisterPage(w, "Пароль должен содержать минимум 6 символов", "")
			return
		}

		req := &models.RegisterRequest{
			Name:     name,
			Email:    email,
			Password: password,
		}

		user, err := authService.RegisterUser(r.Context(), req)
		if err != nil {
			renderRegisterPage(w, err.Error(), "")
			return
		}

		renderRegisterPage(w, "", fmt.Sprintf("Пользователь %s успешно зарегистрирован! Теперь вы можете войти.", user.Name))
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Удаляем cookie с токеном
	authService.ClearAuthCookie(w)

	// Перенаправляем на главную страницу
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func renderLoginPage(w http.ResponseWriter, errorMsg, successMsg string) {
	data := struct {
		Error   string
		Success string
	}{
		Error:   errorMsg,
		Success: successMsg,
	}

	tmpl := template.Must(template.ParseFS(
		templates,
		"templates/login.html",
	))
	tmpl.ExecuteTemplate(w, "login", data)
}

func renderRegisterPage(w http.ResponseWriter, errorMsg, successMsg string) {
	data := struct {
		Error   string
		Success string
	}{
		Error:   errorMsg,
		Success: successMsg,
	}

	tmpl := template.Must(template.ParseFS(
		templates,
		"templates/register.html",
	))
	tmpl.ExecuteTemplate(w, "register", data)
}

// Post handlers
func getPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := postRepo.GetAllPosts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем информацию о пользователе из контекста
	user := auth.GetUserFromContext(r.Context())



	data := struct {
		Posts []*models.Post
		User  *auth.Claims
	}{
		Posts: posts,
		User:  user,
	}

	tmpl := template.Must(template.ParseFS(
		templates,
		filepath.Join("templates", "default.html"),
		filepath.Join("templates", "index.html"),
	))

	tmpl.ExecuteTemplate(w, "default", data)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	if title == "" || content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	// Получаем информацию о пользователе из контекста
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	post := &models.Post{
		Title:   title,
		Content: content,
		UserID:  user.UserID, // Добавляем ID пользователя
	}

	if err := postRepo.CreatePost(r.Context(), post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	postID, err := getPostIDFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	if title == "" || content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	// Получаем информацию о пользователе из контекста
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Получаем пост для проверки владельца
	post, err := postRepo.GetPostByID(r.Context(), postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Проверяем, что пользователь является владельцем поста или администратором
	if post.UserID != user.UserID && user.Role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	post.Title = title
	post.Content = content

	if err := postRepo.UpdatePost(r.Context(), post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", post.ID), http.StatusSeeOther)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	postID, err := getPostIDFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем информацию о пользователе из контекста
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Получаем пост для проверки владельца
	post, err := postRepo.GetPostByID(r.Context(), postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Проверяем, что пользователь является владельцем поста или администратором
	if post.UserID != user.UserID && user.Role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := postRepo.DeletePost(r.Context(), postID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getPost(r *http.Request) (*models.Post, error) {
	postID, err := getPostIDFromURL(r.URL.Path)
	if err != nil {
		return nil, err
	}

	post, err := postRepo.GetPostByID(r.Context(), postID)
	if err != nil {
		return nil, err
	}

	return post, nil
}

// getPostIDFromURL извлекает ID поста из URL
func getPostIDFromURL(path string) (int, error) {
	// Извлекаем ID из пути /post/{id}/...
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return 0, fmt.Errorf("invalid URL path")
	}
	
	postID, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, fmt.Errorf("invalid post ID: %s", parts[2])
	}
	
	return postID, nil
}
