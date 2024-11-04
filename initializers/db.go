package initializers

import (
	"fmt"
	"log"

	"github.com/wpcodevo/go-postgres-crud-rest-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB(env *Env) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", env.DBHost, env.DBUserName, env.DBUserPassword, env.DBName, env.DBPort)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}

	DB.Logger = logger.Default.LogMode(logger.Info)

	DB.AutoMigrate(&models.Feedback{})

	fmt.Println("🚀 Connected Successfully to the Database")
}
