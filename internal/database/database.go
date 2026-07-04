package database

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/medhansh-32/api-gateway/internal/config"
)

func NewMysqlConnection(cfg *config.Config)(*sql.DB,error) {
	
	username := cfg.DBUser
	password := cfg.DBPassword
	host := cfg.DBHost
	port := strconv.Itoa(cfg.DBPort)
	dbName := cfg.DBName
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		username,
		password,
		host,
		port,
		dbName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}