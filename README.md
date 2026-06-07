# Log-Sense

Multi-tenant observability and incident management platform.

**Stack:** Go 1.21 · Gin · PostgreSQL · Redis · Kafka · OpenSearch · React 18 · Vite · Tailwind CSS

## Architecture

```
                        ┌─────────────────────────────────────────┐
                        │              React Frontend              │
                        │  (Vite dev server or nginx static build) │
                        └────────────────┬────────────────────────┘
                                         │ /api/v1/*
                        ┌────────────────▼────────────────────────┐
                        │           Go / Gin HTTP Server           │
                        │  JWT auth · API-key auth · RBAC          │
                        └──┬──────────┬──────────┬────────────────┘
                           │          │          │
              ┌────────────▼─┐  ┌─────▼────┐  ┌─▼───────────┐
              │  PostgreSQL  │  │  Redis   │  │    Kafka     │
              │  (main store)│  │ (dedup)  │  │ logs, alerts │
              └──────────────┘  └──────────┘  └──┬───────────┘
                                                  │
                              ┌───────────────────┼──────────────────┐
                              │                   │                  │
                    ┌─────────▼──────┐  ┌─────────▼────────┐        │
                    │  LogConsumer   │  │ NotifyConsumer    │        │
                    │  → OpenSearch  │  │ → SMTP (on-call)  │        │
                    └────────────────┘  └──────────────────┘        │
                                                                     │
                    ┌────────────────────────────────────────────────▼┐
                    │  RuleEngine (background, every 60 s)             │
                    │  OpenSearch count → alert → Kafka alerts topic   │
                    └──────────────────────────────────────────────────┘
```

**Log ingestion path:** `POST /api/v1/logs` (API key) → Kafka `logs` → LogConsumer → OpenSearch

**Alert path:** RuleEngine polls OpenSearch → threshold breach → PostgreSQL alert + Redis dedup → Kafka `alerts` → NotificationConsumer → SMTP email to **all** currently on-call developers for that service

## Prerequisites

- Go 1.21+
- Node 18+
- Docker & Docker Compose (for local infrastructure)

## Quick start

### 1. Start infrastructure

```bash
docker compose up -d
```

This starts: PostgreSQL (5433), Redis (6379), Kafka (9092), OpenSearch (9200).

### 2. Apply migrations (in order)

Run each file against the `observability` database:

```
migrations/tenants.up.sql
migrations/users.up.sql
migrations/services.up.sql
migrations/rules.up.sql
migrations/alerts.up.sql
migrations/oncall_rules.up.sql
migrations/api_keys.up.sql
migrations/alerts_v2.up.sql
migrations/api_keys_v2.up.sql
```

Example with psql:

```bash
for f in tenants users services rules alerts oncall_rules api_keys alerts_v2 api_keys_v2; do
  psql -h localhost -p 5433 -U postgres -d observability -f migrations/${f}.up.sql
done
```

### 3. Configure

Copy the example config and fill in your values:

```bash
cp config/config.example.yaml config/config.yaml
```

`config/config.yaml` is gitignored (contains credentials). Key sections:

```yaml
database:
  host: localhost
  port: 5433
  user: postgres
  password: postgres
  name: observability

smtp:
  host: smtp.gmail.com
  port: 587
  username: you@gmail.com
  password: app-password
  from_email: you@gmail.com

opensearch:
  host: localhost
  port: 9200
```

### 4. Seed demo data

```bash
go run ./cmd/seed
```

Creates: demo tenant, admin user (`admin@demo.com` / `admin123`), two services, sample rules, and prints an API key.

### 5. Run the backend

```bash
go run .
# or build first:
go build -o log-sense.exe . && ./log-sense.exe
```

Server starts on `:8081` (configurable in `config.yaml`).

### 6. Run the frontend

```bash
cd frontend
npm install
npm run dev
```

Opens at `http://localhost:5173`. API calls proxy to `:8081` via Vite.

## API reference

### Auth

```bash
# Login — returns JWT
curl -X POST http://localhost:8081/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@demo.com","password":"admin123"}'
```

### Log ingestion (API key)

```bash
curl -X POST http://localhost:8081/api/v1/logs \
  -H "X-API-KEY: ls_..." \
  -H "Content-Type: application/json" \
  -d '{"level":"ERROR","message":"payment timeout","service_id":"<uuid>"}'
```

### Alert rules

```bash
curl -X POST http://localhost:8081/api/v1/rules \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "service_id": "<service-uuid>",
    "name": "Payment errors",
    "severity": "CRITICAL",
    "query": {"level": "ERROR", "message_contains": "payment"},
    "threshold": 5,
    "window_minutes": 5
  }'
```

### On-call schedules (admin only)

```bash
# Create
curl -X POST http://localhost:8081/api/v1/admin/oncall/schedules \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "service_id": "<uuid>",
    "user_id": "<uuid>",
    "start_time": "2026-06-06T00:00:00Z",
    "end_time":   "2026-06-13T00:00:00Z"
  }'
```

## RBAC

| Role | Access |
|------|--------|
| ADMIN | Everything: create services, manage developers, manage on-call schedules |
| DEVELOPER | Rules CRUD, alerts (ack/resolve), dashboard, logs search, service read/update |

## Frontend pages

| Page | Path | Role |
|------|------|------|
| Dashboard | `/dashboard` | All |
| Alerts | `/alerts` | All |
| Rules | `/rules` | All |
| Services | `/services` | All (create: admin only) |
| Logs | `/logs` | All |
| Users | `/users` | Admin only |
| On-Call | `/oncall` | Admin only |

## Alert rule query JSON

```json
{
  "level": "ERROR",
  "message_contains": "timeout",
  "fields": { "env": "prod" }
}
```

All fields are optional and ANDed together.

## Health check

```
GET /health
```

Returns status of PostgreSQL, Redis, and Kafka.

## Deployment

See [PROGRESS.md](PROGRESS.md) for full architecture decisions and deployment guidance.

**Short version:**
- **Development:** `docker compose up -d` for infra, `go run .` for backend, `npm run dev` for frontend
- **Production:** Docker Compose with nginx reverse-proxy (routes `/api/*` to Go, serves `frontend/dist` as static files) — or embed the built frontend into the Go binary with `//go:embed frontend/dist`
