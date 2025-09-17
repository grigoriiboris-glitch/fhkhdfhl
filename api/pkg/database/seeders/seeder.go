// pkg/database/seeders/seeder.go
package seeders

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/mymindmap/api/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

// Seeder интерфейс для всех сидов
type Seeder interface {
	Run() error
}

// BaseSeeder базовая структура для сидов
type BaseSeeder struct {
	DB *sql.DB
}

// NewBaseSeeder создает новый базовый сидер
func NewBaseSeeder() *BaseSeeder {
	return &BaseSeeder{
		DB: database.GetDB(),
	}
}