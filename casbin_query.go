package casbinquery

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

const objectSelectSQL = `
SELECT SUBSTRING_INDEX(tp.v1, '_', -1) AS %s 
FROM casbin_rule tg
INNER JOIN casbin_rule tp ON tg.v1 = tp.v0
WHERE tg.v0 = ?
AND tg.ptype = 'g'
AND tp.ptype = 'p'
AND tp.v2 = ?

UNION

SELECT SUBSTRING_INDEX(tp.v1, '_', -1) AS %s 
FROM casbin_rule tp
WHERE tp.v0 = ?
AND tp.ptype = 'p'
AND tp.v2 = ?
`

func QueryObject(db *gorm.DB, columnName, subject, action string) (*gorm.DB, error) {
	if db == nil {
		return nil, errors.New("Invalid argument")
	}

	sql := fmt.Sprintf(objectSelectSQL, columnName, columnName)
	return db.Raw(sql, subject, action, subject, action), nil
}
