# To-Do App (на Go + PostgreSQL)

Простой To-Do список, реализованный на Go с использованием PostgreSQL.

## 🚀 Стек технологий

- Go (net/http, database/sql, gorilla/sessions)
- PostgreSQL
- HTML/CSS
- Git

## ⚙️ Запуск проекта

**1. Создай файл** `.env` **и запиши:**

```rb

connStr=postgres://user:password@localhost:5432/database?sslmode=disable

```

**2. Установи зависимости**

```console

go mod tidy

```


**3. Запусти проект**

```console

go run main.go

```
