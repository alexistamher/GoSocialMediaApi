# 🌐 Social Media API

RESTful API developed in **Go (Golang)** for managing a social network features including users, posts, comments, reactions, and friendships. The project is designed following **Clean Architecture** principles and **Dependency Injection**, ensuring decoupled, scalable, and testable code.

## 🛠️ Technologies and Tools

- **Language:** Go (Golang)
- **Framework / HTTP Routing:** `gin-gonic/gin`
- **ORM:** `gorm` with `postgres` and `sqlite` drivers
- **Database:** PostgreSQL
- **Containerization:** Docker & Docker Compose
- **Testing:** Unit & Integration tests (`testify`, `httptest`, Mocks, `testcontainers-go`)

## ⚙️ Environment Variables

To run this project locally or with Docker, you must create a `.env` file in the root directory and configure the required environment variables:

```env
PORT=3000
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/social_media?sslmode=disable
JWT_SECRET=this is the most secret in the world

POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=social_media
```

| Variable | Description | Example / Default |
| :--- | :--- | :--- |
| `PORT` | Port on which the API server will listen | `3000` |
| `DATABASE_URL` | PostgreSQL connection string | `postgresql://postgres:postgres@localhost:5432/social_media?sslmode=disable` |
| `JWT_SECRET` | Secret key used for signing JWT tokens | `this is the most secret in the world` |
| `POSTGRES_USER` | PostgreSQL database user | `postgres` |
| `POSTGRES_PASSWORD` | PostgreSQL database password | `postgres` |
| `POSTGRES_DB` | PostgreSQL database name | `social_media` |

## 🐳 Docker

Run both the database and the Go server together with a single command:

```bash
docker compose up --build
```
