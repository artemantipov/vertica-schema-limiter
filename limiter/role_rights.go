package limiter

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

// RevokeRW - revoke rw access from role
func RevokeUsage(role, schema string, db *sqlx.DB) (err error) {
	sqlQuery := fmt.Sprintf("REVOKE USAGE ON SCHEMA %v FROM %v", schema, role)
	_, err = db.Exec(sqlQuery)
	return
}

// GrantRW - revoke rw access from role
func GrantUsage(role, schema string, db *sqlx.DB) (err error) {
	sqlQuery := fmt.Sprintf("GRANT USAGE IN SCHEMA %v TO %v", schema, role)
	_, err = db.Exec(sqlQuery)
	return
}
