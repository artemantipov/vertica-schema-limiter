package checker

import (
	"github.com/jmoiron/sqlx"
	"log"
)

// SchemaSize - return schema size on GB
func SchemaSize(name string, db *sqlx.DB) (size float64) {
	//noinspection SqlResolve
	checkQuery := `
	SELECT
    	SUM (used_bytes) / (1024^3) AS used_gb
	FROM
    	v_monitor.column_storage
	WHERE anchor_table_schema = ?
	ORDER  BY
	    SUM (used_bytes) DESC;`

	err := db.QueryRow(checkQuery, name).Scan(&size)
	if err != nil {
		log.Println(err)
	}
	return
}
