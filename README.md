# test-backend

A Go REST API for **users** and **tasks** (CRUD, reset password, pagination, search). Uses SQLite (file-based, no separate server) and UUIDs for all IDs. **Login** returns a JWT; protected routes require `Authorization: Bearer <token>`.

## Prerequisites

- **Go 1.21+** — [Install Go](https://go.dev/dl/)

## Run locally

From the project root:

```bash
go run ./cmd/api
```

The server listens on **http://localhost:8082** by default.

## Environment variables

| Variable         | Default        | Description                    |
|------------------|----------------|--------------------------------|
| `PORT`           | `8082`         | HTTP server port               |
| `DB_PATH`        | `./data/app.db`| Path to SQLite database        |
| `JWT_SECRET`     | *(none)*       | Secret for signing JWTs (required for login/protected routes) |
| `JWT_EXPIRY_HOURS` | `24`         | Token validity in hours        |
| `CORS_ORIGINS`   | *(see below)*  | Comma-separated origins for CORS (e.g. `http://localhost:3000`). Default: localhost:3000, localhost:5173, 127.0.0.1:3000, 127.0.0.1:5173 |

Optional: copy `.env.example` to `.env` and adjust values. The app loads `.env` if present.

## Database

On first run, the app creates the SQLite file and tables (UUID-based IDs). If the schema has been updated to use UUIDs, existing tables are recreated and **existing data will be lost**.

## API

- **[API_CURL_EXAMPLES.md](API_CURL_EXAMPLES.md)** — all endpoints and example cURL commands (including login and `Authorization: Bearer <token>`).
- **[openapi.json](openapi.json)** — OpenAPI 3.0 spec for import into Apidog, Postman, or other API tools.
- **[SETUP_AND_PLAN.md](SETUP_AND_PLAN.md)** — original plan, tech stack, and current implementation summary.

Quick check:

```bash
curl http://localhost:8082/health
```

Expected: `{"message":"Service is healthy.","status":"ok","database":"connected"}`
