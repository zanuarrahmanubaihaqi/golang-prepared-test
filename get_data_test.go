package main

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

const (
	dbUser     = "root"   // Ubah sesuai dengan pengguna MySQL Anda
	dbPassword = ""       // Ubah sesuai dengan password MySQL Anda
	dbName     = "ekocar" // Ubah sesuai dengan nama database Anda
)

func initDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", dbUser, dbPassword, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Inisialisasi tabel untuk keperluan pengujian
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id VARCHAR(4) PRIMARY KEY, name VARCHAR(50))")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("INSERT INTO users (id, name) VALUES (?, ?)", "C001", "Alice")
	return db, err
}

func TestGetUserByID(t *testing.T) {
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Siapkan prepared statement
	stmt, err := PrepareGetUserByID(db)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Test: Memeriksa apakah pengguna dapat diambil dengan ID yang benar
	name, err := GetUserByID(db, stmt, "C001")
	assert.NoError(t, err)
	assert.Equal(t, "Alice", name)

	// Test: Memeriksa ID yang tidak ada
	_, err = GetUserByID(db, stmt, "C999")
	assert.Error(t, err)
}
