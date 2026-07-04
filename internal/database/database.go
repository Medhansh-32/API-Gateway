package database

import (
	"database/sql"
	"fmt"
	"github.com/medhansh-32/api-gateway/internal/config"
)

func NewMysqlConnection(cfg *config.Config)(*sql.DB,error) {
	
	username := "root"
	password := "password"
	host := "localhost"
	port := "3306"
	dbName := "mydb"

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