
#go run
cd worker
go mod init worker
go get github.com/gofiber/fiber/v2
go get github.com/joho/godotenv
go get -u gorm.io/gorm
go get gorm.io/driver/mysql
go get github.com/golang-jwt/jwt


#데이터베이스 작업
ALTER STEAMPIPE_DB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;
ALTER ad_groups CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
