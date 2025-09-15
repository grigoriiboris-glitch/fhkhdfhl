package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"your-app/relations"
	
	_ "github.com/lib/pq"
)

// User модель пользователя
type User struct {
	relations.BaseModel
	Name  string `db:"name"`
	Email string `db:"email"`
	Posts []*Post
}

func (u *User) GetTableName() string { return "users" }
func (u *User) GetConnection() *sql.DB { return db }

// Post модель поста
type Post struct {
	relations.BaseModel
	UserID int64  `db:"user_id"`
	Title  string `db:"title"`
	Body   string `db:"body"`
	User   *User
}

func (p *Post) GetTableName() string { return "posts" }
func (p *Post) GetConnection() *sql.DB { return db }

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("postgres", "user=postgres dbname=test sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	ctx := context.Background()
	
	// Создание отношений
	user := &User{}
	post := &Post{}
	
	builder := relations.NewQueryBuilder(user)
	
	// Регистрация отношений
	builder.RegisterRelation("posts", relations.NewHasMany(
		func() relations.Model { return &Post{} },
		relations.RelationConfig{},
	))
	
	// Загрузка пользователя с постами
	if err := builder.With("posts").Get(ctx, 1); err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("User: %s\n", user.Name)
	fmt.Printf("Posts count: %d\n", len(user.Posts))
	
	// Отношение "принадлежит"
	postBuilder := relations.NewQueryBuilder(post)
	postBuilder.RegisterRelation("user", relations.NewBelongsTo(
		&User{},
		relations.RelationConfig{},
	))
	
	if err := postBuilder.With("user").Get(ctx, 1); err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Post: %s, Author: %s\n", post.Title, post.User.Name)
}