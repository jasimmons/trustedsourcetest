package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func Main(args map[string]interface{}) map[string]interface{} {
	db, err := dbFromEnv()
	if err != nil {
		return map[string]interface{}{
			"body":   fmt.Sprintf(`{"status":"failed","error":"%w"}`, err),
			"errors": err.Error(),
		}
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return map[string]interface{}{
			"body":   fmt.Sprintf(`{"status":"failed","error":"pinging db: %s"}`, err.Error()),
			"errors": err.Error(),
		}
	}

	return map[string]interface{}{
		"body": fmt.Sprint(`{"status":"connected"}`),
	}
}

func dbFromEnv() (*sql.DB, error) {
	user := os.Getenv("DB_USERNAME")
	if user == "" {
		// default DO mysql user
		user = "doadmin"
	}

	pass := os.Getenv("DB_PASSWORD")
	if pass == "" {
		return nil, errors.New("missing DB_PASSWORD")
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		return nil, errors.New("missing DB_HOST")
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		// default DO mysql port
		port = "25060"
	}

	db := os.Getenv("DB_DATABASE")
	if db == "" {
		db = "defaultdb"
	}

	cfg := &mysql.Config{
		User:   user,
		Passwd: pass,
		Net:    "tcp",
		Addr:   net.JoinHostPort(host, port),
		DBName: db,
	}
	dsn := cfg.FormatDSN()

	return sql.Open("mysql", dsn)
}
