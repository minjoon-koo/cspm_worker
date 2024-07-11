package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"worker/models"
)

var DB *gorm.DB

func Connection() {
	godotenv.Load(".env")
	dbhost := os.Getenv("MYSQL_HOST")
	dbuser := os.Getenv("MYSQL_USER")
	dbpassword := os.Getenv("MYSQL_PASSWORD")
	dbname := os.Getenv("MYSQL_DBNAME")

	connection := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbuser, dbpassword, dbhost, dbname)
	var db, err = gorm.Open(mysql.Open(connection), &gorm.Config{})

	if err != nil {
		panic("db connect fail")
	}
	DB = db
	fmt.Println("db connect success")

	AutoMigrate(db)
}

func AutoMigrate(connection *gorm.DB) {
	connection.Debug().AutoMigrate(
		&models.AdGroup{},
		&models.DirectoryRole{},
		&models.BackendPool{},
		&models.FrontendResource{},
	)
	connection.Exec("ALTER TABLE ad_groups CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
	connection.Exec("ALTER TABLE directory_roles CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
	connection.Exec("ALTER TABLE backend_pools CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
	connection.Exec("ALTER TABLE frontend_resources CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
}
