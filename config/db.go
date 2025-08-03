package config

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() {
	dsn := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("не удалось подключиться к БД: %v", err)
	}
	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("не удалось проверить подключение к БД: %v", err)
	}
	DB = pool
	log.Println("✅ База данных подключена")
}
