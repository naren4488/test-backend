# Go Backend: Setup Guide & Implementation Plan

## 1. Overview

You will build a **REST API** in Go with:

- **Users API**: Create, Read by ID, List all, Update, Delete, Reset password  
- **Tasks API**: Create (for a user), List for a user, List all (paginated + search), Get by ID, Update, Delete  

**Database**: SQLite (single file, no server, free, perfect for local learning).  
**Scope**: Local only, learning; no deployment or production hardening.

---

## 2. Tech Stack

| Layer        | Choice                    | Reason |
|-------------|---------------------------|--------|
| Language    | Go 1.21+                  | Your requirement |
| Router      | **chi**                   | Lightweight, idiomatic, middleware support |
| DB          | **SQLite** (pure Go driver) | No install, no CGO, single file |
| DB driver   | `modernc.org/sqlite`      | Pure Go SQLite, no C compiler needed |
| Validation  | `go-playground/validator` | Struct tags, clear errors |
| Config      | Env vars + `.env` (e.g. `godotenv`) | Simple for local |

---

## 3. Project Structure (Layered / Clean)

```
test-backend/
├── cmd/
│   └── api/
│       └── main.go              # Entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Load env, DB path, server port
│   ├── database/
│   │   ├── database.go          # DB connection, init
│   │   └── schema.sql           # Embedded CREATE TABLE statements
│   ├── middleware/
│   │   ├── middleware.go        # Recover, logger, request ID
│   │   └── ...
│   ├── handler/
│   │   ├── handler.go           # Base / shared handler helpers
│   │   ├── user.go              # User HTTP handlers
│   │   └── task.go              # Task HTTP handlers
│   ├── service/
│   │   ├── user.go              # User business logic
│   │   └── task.go              # Task business logic
│   ├── repository/
│   │   ├── user.go              # User DB operations
│   │   └── task.go              # Task DB operations
│   ├── model/
│   │   ├── user.go              # User struct, validation tags
│   │   └── task.go              # Task struct, validation tags
│   └── response/
│       └── response.go          # JSON response helpers, error format
├── pkg/
│   └── errors/
│       └── errors.go            # App error types, codes, messages
├── go.mod
├── go.sum
├── .env.example
├── README.md
└── SETUP_AND_PLAN.md            # This file
```

**Layers:**

- **Handler**: HTTP (parse body, query, call service, write response).
- **Service**: Business rules (e.g. “user must exist before creating task”, “reset password” logic).
- **Repository**: All SQL (CRUD); returns models or app errors.
- **Model**: Structs + validation tags.
- **response**: Consistent JSON shape (data, error, meta for pagination).
- **pkg/errors**: Domain error types and user-facing messages.

---

## 4. Database Choice: SQLite (Easy Local Setup)

- **No separate DB server**: one `.db` file in your project (e.g. `./data/app.db`).
- **Pure Go driver** (`modernc.org/sqlite`): no CGO, no C compiler.
- **Free**, file-based, easy to backup (copy the file) or reset (delete file and re-run schema).

We will create the schema in `internal/database/schema.sql` (embedded) and run it on startup (or via a small init command) so you have a single “run the app” flow.

---

## 5. API Contract (Summary)

### Users

| Method | Path | Description |
|--------|------|-------------|
| POST   | `/api/v1/users`           | Create user |
| GET    | `/api/v1/users`           | List all users |
| GET    | `/api/v1/users/:id`       | Get user by ID |
| PUT    | `/api/v1/users/:id`       | Update user |
| DELETE | `/api/v1/users/:id`       | Delete user |
| POST   | `/api/v1/users/:id/reset-password` | Reset password (body: new password) |

### Tasks

| Method | Path | Description |
|--------|------|-------------|
| POST   | `/api/v1/users/:id/tasks` | Create task for user |
| GET    | `/api/v1/users/:id/tasks` | List all tasks for user |
| GET    | `/api/v1/tasks`           | List all tasks (paginated, optional search) |
| GET    | `/api/v1/tasks/:id`       | Get task by ID |
| PUT    | `/api/v1/tasks/:id`       | Update task |
| DELETE | `/api/v1/tasks/:id`       | Delete task |

Query params for **List all tasks**:

- `page`, `limit` (or `per_page`) for pagination.
- `search` (optional): filter by title/description.
- Optional: `user_id` to filter by user.

---

## 6. Error Handling & Messages

- **Centralized errors** in `pkg/errors`:
  - Types: e.g. `ErrNotFound`, `ErrBadRequest`, `ErrConflict` (e.g. duplicate email).
  - Each type has an **HTTP status** and a **user-facing message** (and optional code like `USER_NOT_FOUND`).
- **Handlers** never send arbitrary strings; they map service/repository errors to these types and use `response` to send a consistent JSON body, e.g.:

  ```json
  { "error": { "code": "USER_NOT_FOUND", "message": "User not found" } }
  ```

- **Validation errors** (from `validator`): return 400 with a list of field errors.

---

## 7. Setup Guide (Step by Step)

### 7.1 Prerequisites

- **Go 1.21 or later**  
  - Check: `go version`  
  - Install: https://go.dev/dl/

### 7.2 Initialize project (from project root)

```bash
cd /Users/narendrakajla77/Downloads/personal/dev/projects/test-backend
go mod init test-backend
```

### 7.3 Create folder structure

Create the directories as in the structure above (cmd/api, internal/..., pkg/errors, etc.). We will generate files in the next phases.

### 7.4 Dependencies (we will add as we build)

- `github.com/go-chi/chi/v5` – router  
- `modernc.org/sqlite` – SQLite driver  
- `github.com/joho/godotenv` – load `.env` (optional but useful)  
- `github.com/go-playground/validator/v10` – validation  

No database server install: only Go and these packages.

### 7.5 Environment

- `.env.example` with:
  - `PORT=8080`
  - `DB_PATH=./data/app.db`
- Copy to `.env` and adjust if needed. App reads `PORT` and `DB_PATH` from env.

### 7.6 Run

- Create `./data` if needed; app will create `app.db` and run schema on first run (or we use an init step).
- Start server: `go run ./cmd/api` (from project root).

---

## 8. Implementation Phases

We will implement in small steps so you can run and test after each phase.

### Phase 1 – Project skeleton & DB

1. `go.mod` + folder structure.  
2. **Config**: load `PORT`, `DB_PATH` from env; optional `godotenv`.  
3. **DB**: open SQLite in `internal/database`, run `schema.sql` on startup (create `users` and `tasks` tables).  
4. **Schema**: `users` (id, email, name, password_hash, created_at, updated_at), `tasks` (id, user_id, title, description, status, created_at, updated_at).  
5. **Main**: start HTTP server with chi, health route (e.g. `GET /health`).  
6. **Middleware**: recover, logger, request ID.  

**Check**: `GET /health` returns 200; DB file created; tables exist.

---

### Phase 2 – Users API & errors

1. **pkg/errors**: error types and messages (NotFound, BadRequest, Conflict, etc.).  
2. **internal/response**: JSON response helper (success + error shape).  
3. **internal/model/user**: struct + validation tags.  
4. **internal/repository/user**: CRUD + “get by email” (for create/duplicate check).  
5. **internal/service/user**: create, get by ID, list, update, delete, reset password (hash password with bcrypt).  
6. **internal/handler/user**: wire handlers to service, parse path/body, return response/errors.  
7. **Routing**: mount `/api/v1/users` and sub-routes (`/:id`, `/:id/reset-password`).  

**Check**: All user endpoints work with curl/Postman; duplicate email returns proper error.

---

### Phase 3 – Tasks API

1. **internal/model/task**: struct + validation.  
2. **internal/repository/task**: CRUD, list by user_id, list all with pagination and search.  
3. **internal/service/task**: business logic (e.g. ensure user exists when creating task).  
4. **internal/handler/task**: handlers for create, list by user, list all (query params), get, update, delete.  
5. **Routing**: mount `/api/v1/users/:id/tasks` and `/api/v1/tasks`.  

**Check**: Create task for user; list by user; list all with page/limit/search.

---

### Phase 4 – Polish

1. **Pagination**: consistent `page`, `limit`, total count in response (e.g. `meta.total`, `meta.page`).  
2. **Search**: filter tasks by title/description (LIKE or simple full-text).  
3. **Reset password**: validate new password, hash, update in DB; clear error messages.  
4. **README**: how to run, env vars, example requests.  

---

## 9. Next Step

Once you’ve read this and confirmed:

1. You have Go 1.21+ installed.  
2. You’re fine with SQLite and the folder structure above.  
3. You’re fine with the API paths and phases.  

We’ll start with **Phase 1**: create the repo structure, `go.mod`, config, database connection, schema, main with health route and middleware. After that we’ll move slowly through Phase 2 (Users), then Phase 3 (Tasks), then Phase 4 (polish).

If you want to change anything (e.g. different router, add JWT later, or different error format), say so before we start Phase 1.
