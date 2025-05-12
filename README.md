# To-Do App (на Go + PostgreSQL)

Простой To-Do список, реализованный на Go с использованием PostgreSQL.

## 🚀 Стек технологий

- Go (net/http, database/sql, gorilla/sessions)
- PostgreSQL
- HTML/CSS
- Git

## ⚙️ Запуск проекта

1. Создай файл `.env`:

connStr=postgres://postgres:1234@localhost:5432/tododb?sslmode=disable

2. Установи зависимости

go mod tidy

3. Запусти проект

go run main.go