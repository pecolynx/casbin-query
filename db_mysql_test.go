package casbinquery

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	gormMySQL "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	testDBURL string
)

func openMySQLForTest() *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	db, err := gorm.Open(gormMySQL.Open(testDBURL), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	return db
}

func init() {
	fmt.Println("init")

	testDBHost := os.Getenv("TEST_DB_HOST")
	if testDBHost == "" {
		testDBHost = "127.0.0.1"
	}

	testDBPort := os.Getenv("TEST_DB_PORT")
	if testDBPort == "" {
		testDBPort = "3307"
	}

	testDBURL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Asia%%2FTokyo", "user", "password", testDBHost, testDBPort, "testdb")

	db := openMySQLForTest()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(wd)
	dir := wd + "/sqls"
	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://"+dir, "mysql", driver)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Fatal(fmt.Errorf("Failed to up. err: %w", err))
		}
	}

	initCasbin(db)
}

const conf = `[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`

func initCasbin(db *gorm.DB) {
	m, err := model.NewModelFromString(conf)
	if err != nil {
		panic(err)
	}
	a, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		panic(err)
	}
	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		panic(err)
	}

	if err := e.LoadPolicy(); err != nil {
		panic(err)
	}

	addNamedPolicy := func(subject, object, action string) {
		if _, err := e.AddNamedPolicy("p", subject, object, action); err != nil {
			panic(err)
		}
	}
	addNamedGroupPolicy := func(user, role string) {
		if _, err := e.AddNamedGroupingPolicy("g", user, role); err != nil {
			panic(err)
		}

	}
	
	addNamedPolicy("owner_A", "pet_ewok", "read")
	addNamedPolicy("owner_A", "pet_fluffy", "read")
	addNamedPolicy("owner_A", "pet_gordo", "update")
	addNamedPolicy("owner_B", "pet_gordo", "read")
	addNamedPolicy("user_david", "pet_ewok", "read")
	addNamedPolicy("user_david", "pet_fluffy", "update")
	addNamedGroupPolicy("user_bob", "owner_A")
	addNamedGroupPolicy("user_charlie", "owner_B")

	if err := e.SavePolicy(); err != nil {
		panic(err)
	}
} 

func TestPolicyCheck(t*testing.T){
	db := openMySQLForTest()
	m, err := model.NewModelFromString(conf)
	if err != nil {
		panic(err)
	}
	a, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		panic(err)
	}
	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		panic(err)
	}


	check := func(subject, object, action string, granted bool) {
		res, err := e.Enforce(subject, object, action)
		if err != nil {
			panic(err)
		}
		if res != granted {
			log.Fatalf("%s, %s, %s. expected: %v, actual: %v", subject, object, action, granted, res)
		}
	}
	check("user_bob", "pet_ewok", "read", true)
	check("user_bob", "pet_fluffy", "read", true)
	check("user_bob", "pet_gordo", "read", false)

	check("user_charlie", "pet_ewok", "read", false)
	check("user_charlie", "pet_fluffy", "read", false)
	check("user_charlie", "pet_gordo", "read", true)

	check("user_david", "pet_ewok", "read", true)
	check("user_david", "pet_fluffy", "read", false)
	check("user_david", "pet_gordo", "read", false)
}
