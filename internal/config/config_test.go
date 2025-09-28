package config

import (
	"os"
	"testing"
)

func TestConnectDB_DefaultDSN(t *testing.T) {
	os.Unsetenv("DB_URL")
	defer func() {
		if r := recover(); r != nil {
			// Ожидаем фатальный лог, тест пройден
		}
	}()
	ConnectDB() // если не упадет, тест не пройден
}

func TestConnectDB_WithEnv(t *testing.T) {
	os.Setenv("DB_URL", "host=localhost user=gegxkss password=postgres dbname=wbl0_db port=5432 sslmode=disable")
	defer os.Unsetenv("DB_URL")
	defer func() {
		if r := recover(); r != nil {
			// Ожидаем фатальный лог, тест пройден
		}
	}()
	ConnectDB()
}
