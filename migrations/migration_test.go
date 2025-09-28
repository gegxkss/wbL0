package migrations

import (
	"testing"
	"github.com/gegxkss/wbL0/internal/config"
)

func TestMigration(t *testing.T) {
	Migration()

	// Проверяем, что таблицы созданы
	tables := []string{"orders", "deliveries", "payments", "items"}
	for _, table := range tables {
		var exists bool
		// Проверяем наличие таблицы через запрос к information_schema.tables
		err := config.DB.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = ?)", table).Scan(&exists).Error
		if err != nil {
			t.Fatalf("Ошибка при проверке таблицы %s: %v", table, err)
		}
		if !exists {
			t.Errorf("Таблица %s не создана", table)
		}
	}
}
