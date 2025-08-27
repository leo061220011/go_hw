package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	// Открываем или создаем базу
	db, err := sql.Open("sqlite", "users.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// Создаем таблицу
	createTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);
	`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

	// Добавим одного тестового пользователя
	email := "test@example.com"
	password := "12345678"
	_, err = db.Exec("INSERT OR IGNORE INTO users (email, password) VALUES (?, ?)", email, password)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Тестовый пользователь создан: test@example.com / 12345678")

	// Авторизация
	var inputEmail, inputPassword string
	fmt.Print("Введите email: ")
	fmt.Scan(&inputEmail)
	fmt.Print("Введите пароль: ")
	fmt.Scan(&inputPassword)

	var storedPassword string
	err = db.QueryRow("SELECT password FROM users WHERE email = ?", inputEmail).Scan(&storedPassword)
	if err == sql.ErrNoRows {
		fmt.Println("Пользователь не найден")
		return
	} else if err != nil {
		fmt.Println("Ошибка базы данных:", err)
		return
	}

	if storedPassword == inputPassword {
		fmt.Println("Пароль верный")
	} else {
		fmt.Println("Пароль неверный")
	}
}
