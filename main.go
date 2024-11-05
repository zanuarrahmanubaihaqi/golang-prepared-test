package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func PrepareGetUserByID(db *sql.DB) (*sql.Stmt, error) {
	// Menggunakan prepared statement
	query := "SELECT name FROM users WHERE id = ?"
	return db.Prepare(query)
}

func GetUserByID(db *sql.DB, stmt *sql.Stmt, id string) (string, error) {
	var name string
	err := stmt.QueryRow(id).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func GetUserByIDNoPrevent(db *sql.DB, id string) ([]string, error) {
	var name []string
	query := fmt.Sprintf("SELECT name FROM users WHERE id = '%s'", id)
	err := db.QueryRow(query).Scan(&name)
	if err != nil {
		return []string{}, err
	}
	return name, nil
}
