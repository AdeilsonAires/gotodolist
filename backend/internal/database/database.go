package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AdeilsonAires/gotodolist/backend/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s timeZone=%s",
		getenv("DB_HOST", "localhost"),
		getenv("DB_USER", "postgres"),
		getrequiredenv("DB_PASSWORD"),
		getenv("DB_NAME", "gotodolisdb"),
		getenv("DB_PORT", "5432"),
		getenv("DB_SSLMODE", "disable"),
		getenv("DB_TIMEZONE", "America/Sao_Paulo"),
	)

	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("erro ao conectar ao banco de dados: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("erro ao database instance: %v", err)
	}

	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("erro ao pingar o banco de dados: %v", err)
	}

	if err := DB.AutoMigrate(&domain.Item{}, &domain.List{}, &domain.User{}); err != nil {
		log.Fatalf("erro ao migrar os dados: %v", err)
	}

	fmt.Println("Banco de dados conectado e migrado com sucesso!")
}

func GetDB() *gorm.DB {
	if DB == nil {
		panic("ERRO: banco de dados nulo")
	}
	return DB
}

func getenv(key, defaulValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaulValue
}

func getrequiredenv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("variavel de ambiente obrigatória não encontrada")
	}
	return value
}

func CloseDB() error {
	if DB == nil {
		return fmt.Errorf("conexão já encerrada")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("erro ao buscar o database instance: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("erro ao fechar a conexão: %v", err)
	}

	fmt.Println("Encerrando a conexão")
	DB = nil
	return nil
}
