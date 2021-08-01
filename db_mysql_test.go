package casbinquery

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	gormMySQL "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	testDBURL string
)

func openMySQLForTest() (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	return gorm.Open(gormMySQL.Open(testDBURL), &gorm.Config{
		Logger: newLogger,
	})
}

func initMySQL() {
	testDBHost := os.Getenv("TEST_DB_HOST")
	if testDBHost == "" {
		testDBHost = "127.0.0.1"
	}

	testDBPort := os.Getenv("TEST_DB_PORT")
	if testDBPort == "" {
		testDBPort = "3307"
	}

	testDBURL = fmt.Sprintf("user:password@tcp(%s:%s)/testdb?charset=utf8&parseTime=True&loc=Asia%%2FTokyo", testDBHost, testDBPort)

	setupMySQL()
}

func setupMySQL() {
	db, err := openMySQLForTest()
	if err != nil {
		log.Fatal(err)
	}
	initCasbin(db)
	setupDB(db, "mysql", func(sqlDB *sql.DB) (database.Driver, error) {
		return mysql.WithInstance(sqlDB, &mysql.Config{})
	})
}
