# familyquest-backend
Серверная часть FamilyQuest — семейной системы мотивации детей: еженедельные задания, подтверждение выполнения родителями, баллы, награды и достижения. REST API на Go с PostgreSQL и JWT-авторизацией.

## Локальный запуск

```bash
docker compose up -d postgres
```

Приложение можно запускать из GoLand через пакет `cmd` или командой:

```bash
go run ./cmd
```

Базовые ручки:

- `GET /api/v1/hello` — проверка HTTP API.
- `GET /app/health` — liveness.
- `GET /app/ready` — readiness с проверкой PostgreSQL.
