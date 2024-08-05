package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func DbConnection() (*sql.DB, error) {
	fmt.Println("psql connection establishing ...")

	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	psqlInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DbHost, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbPort)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalln("Error opening database connection:", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln("Error connecting to the database:", err)
		return nil, err
	}

	fmt.Println("psql connected successfully")

	// if err := MigrateTables(); err != nil {
	// 	log.Fatalln("Error performing migration:", err)
	// 	return nil, err
	// }

	return db, nil
}
func MigrateTables() error {

	db, err := DbConnection()

	if err != nil {
		return fmt.Errorf("database connection error: %v", err)
	}
	defer db.Close()

	// Create users table
	CreateUserQuery := `CREATE TABLE IF NOT EXISTS Users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

	_, err = db.Exec(CreateUserQuery)
	if err != nil {
		return err
	}

	// Create feedback table with additional fields
	CreateFeedbackQuery := `CREATE TABLE IF NOT EXISTS Feedback (
    feedback_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES Users(user_id),
    name VARCHAR(100) NOT NULL,
    age INTEGER NOT NULL,
    occupation VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    would_recommend BOOLEAN NOT NULL,
    suggestion TEXT,
    likes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

	_, err = db.Exec(CreateFeedbackQuery)
	if err != nil {
		return err
	}

	return nil
}
