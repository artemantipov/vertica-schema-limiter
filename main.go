package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
	"vertica-schema-limiter/checker"
	"vertica-schema-limiter/config"
	"vertica-schema-limiter/limiter"
	"vertica-schema-limiter/metrics"

	_ "github.com/vertica/vertica-sql-go"
)

var (
	schemaLock    = make(map[string]bool)
	limiterConfig = config.GetConfig()
)

func main() {
	log.Println("Start App")
	go startLimiter()
	metrics.Start()
}

func startLimiter() {
	for {
		interval := limiterConfig.Vertica.CheckInterval
		db := dbConnect()
		for schema := range limiterConfig.Schemas {
			size := checker.SchemaSize(schema, db)
			metricsMap := metrics.GetMap()
			metricsMap[schema] = metrics.Metric{Value: size, CreatedAt: time.Now()}
			limit := limiterConfig.Schemas[schema].Limit
			role := limiterConfig.Schemas[schema].Role
			if size >= limit {
				log.Printf("Limit for %v schema excited: %vGB of %vGB", schema, size, limit)
				if schemaLock[schema] != true {
					err := limiter.RevokeUsage(role, schema, db)
					if err != nil {
						log.Printf("Failed for rewoke with err: %v", err)
					} else {
						log.Printf("REVOKED usage on schema %v from %v ROLE", schema, role)
						schemaLock[schema] = true
					}
				}
			} else {
				if schemaLock[schema] {
					log.Printf("Limit for %v schema is OK: %vGB of %vGB", schema, size, limit)
					err := limiter.GrantUsage(role, schema, db)
					if err != nil {
						log.Printf("Failed for grant with err: %v", err)
					} else {
						log.Printf("GRANTED usage on schema %v to %v ROLE", schema, role)
						schemaLock[schema] = false
					}
				}
			}
		}
		dbClose(db)
		time.Sleep(time.Duration(interval) * time.Minute)
	}
}

func dbConnect() (db *sqlx.DB) {
	conn := limiterConfig.Vertica
	conStr := fmt.Sprintf("vertica://%v:%v@%v:%v/%v", conn.User, conn.Pass, conn.Host, conn.Port, conn.DB)
	db, err := sqlx.Connect("vertica", conStr)
	if err != nil {
		log.Printf("Connection Failed to Open with ERROR: %v", err)
	} else {
		log.Print("Connection Established")
	}
	return
}

func dbClose(db *sqlx.DB) {
	err := db.Close()
	if err != nil {
		log.Printf("Error during close connect: %v", err)
	} else {
		log.Print("Connection closed")
	}
}
