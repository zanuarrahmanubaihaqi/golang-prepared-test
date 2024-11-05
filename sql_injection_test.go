package main

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

// Perbarui initDBSIT untuk menggunakan VARCHAR di kolom ID
func initDBSIT() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", dbUser, dbPassword, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Update tipe kolom id ke VARCHAR(50)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id VARCHAR(50) PRIMARY KEY, name VARCHAR(50))")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("DELETE FROM users")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("INSERT INTO users (id, name) VALUES (?, ?)", "1", "Alice")
	return db, err
}

func TestSQLInjectionNoPrevention(t *testing.T) {
	db, err := initDBSIT()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Tambahkan data tambahan untuk uji
	_, err = db.Exec("INSERT INTO users (id, name) VALUES (?, ?)", "2", "Bob")
	if err != nil {
		log.Fatal(err)
	}

	// Attempt SQL Injection
	injectionAttempt := "1 OR 1=1" // Modifikasi untuk tipe INT di SQL
	name, err := GetUserByIDNoPrevent(db, injectionAttempt)

	// Verifikasi: tidak ada error dan hasil dari SQL Injection
	assert.NoError(t, err, "Tidak seharusnya ada error dari SQL Injection jika tidak dicegah")
	assert.NotEqual(t, "Alice", name, "Jika SQL Injection berhasil, maka hasil query seharusnya tidak hanya 'Alice'")
	assert.NotEmpty(t, name, "SQL Injection seharusnya menghasilkan hasil yang tidak semestinya")
}

func TestSQLInjectionPrevention(t *testing.T) {
	db, err := initDBSIT()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Prepare the statement
	stmt, err := PrepareGetUserByID(db)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Attempt SQL Injection
	injectionAttempt := "' OR '1'='1"
	name, err := GetUserByID(db, stmt, injectionAttempt)

	// Pastikan tidak ada nama atau error yang berhasil dieksekusi
	assert.Error(t, err, "Expected an error for SQL injection attempt")
	assert.Empty(t, name, "Expected no data to be retrieved with SQL injection attempt")
}
