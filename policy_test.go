package casbinquery

import (
	"log"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

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

func TestPolicyCheck(t *testing.T) {
	for _, db := range dbList() {
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
}
