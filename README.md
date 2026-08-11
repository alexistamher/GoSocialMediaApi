# 🌐 Social Media API

RESTful API developed in **Go (Golang)** for managing a social network features including users, posts, comments, reactions, and friendships. The project is designed following **Clean Architecture** principles and **Dependency Injection**, ensuring decoupled, scalable, and testable code.

## 🛠️ Technologies and Tools

- **Language:** Go (Golang)
- **Framework / HTTP Routing:** `gin-gonic/gin`
- **ORM:** `gorm` with `postgres` and `sqlite` drivers
- **Database:** PostgreSQL
- **Containerization:** Docker & Docker Compose
- **Testing:** Unit & Integration tests (`testify`, `httptest`, Mocks, `testcontainers-go`)

## 🐳 Docker

Run both the database and the Go server together with a single command:

```bash
docker compose up --build
```
