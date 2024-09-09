package infra

import (
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dbOnce sync.Once
	db     *gorm.DB
)

func LoadDB() *gorm.DB {
	dbOnce.Do(func() {

		cfg := LoadConfig()

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
			cfg.DB.User,
			cfg.DB.Pass,
			cfg.DB.Host,
			cfg.DB.Port,
			cfg.DB.Name,
			cfg.DB.Param,
		)

		dbConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}

		db = dbConn
	})

	return db
}
