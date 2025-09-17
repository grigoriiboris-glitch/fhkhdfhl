// pkg/database/seeders/manager.go
package seeders

import (
	"fmt"
	"log"
)

// SeederManager управляет выполнением сидов
type SeederManager struct {
	seeders []Seeder
}

// NewSeederManager создает новый менеджер сидов
func NewSeederManager() *SeederManager {
	return &SeederManager{}
}

// Register регистрирует сидер
func (sm *SeederManager) Register(seeder Seeder) {
	sm.seeders = append(sm.seeders, seeder)
}

// Run запускает все зарегистрированные сиды
func (sm *SeederManager) Run() error {
	log.Println("Запуск всех сидов...")

	for _, seeder := range sm.seeders {
		if err := seeder.Run(); err != nil {
			return fmt.Errorf("ошибка выполнения сидера: %v", err)
		}
	}

	log.Println("Все сиды выполнены успешно!")
	return nil
}

// RunSpecific запускает конкретный сидер
func (sm *SeederManager) RunSpecific(seeder Seeder) error {
	return seeder.Run()
}