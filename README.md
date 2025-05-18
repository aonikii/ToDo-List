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

**3. Добавь таблицы в базы данных:**

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    done BOOLEAN DEFAULT FALSE,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

**4. Запусти проект**

```console

go run ./cmd

```
