package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/DEELAGRA/org-struct-api/internal/config"
	"github.com/DEELAGRA/org-struct-api/internal/repository"
	"github.com/DEELAGRA/org-struct-api/internal/router"
	"github.com/DEELAGRA/org-struct-api/internal/service"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}
	dsn := "host=" + cfg.DBHost +
		" user=" + cfg.DBUser +
		" password=" + cfg.DBPassword +
		" dbname=" + cfg.DBName +
		" port=" + strconv.Itoa(cfg.DBPort) +
		" sslmode=" + cfg.DBSSLMode

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	repo := repository.NewDepartmentRepository(db)
	svc := service.NewDepartmentService(repo)

	mux := router.SetupRouter(svc)

	addr := ":" + strconv.Itoa(cfg.ServerPort)
	log.Printf("Сервер запущен на %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
