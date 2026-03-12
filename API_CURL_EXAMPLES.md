# API cURL Examples

Base URL: `http://localhost:8082`  
All IDs are **UUIDs**. Replace the example UUIDs below with real IDs from create/list responses.

**Authentication:** All routes except **Health**, **Login**, and **Create user** require a JWT. Send it in the header:  
`Authorization: Bearer <your_token>`

---

## Health

```bash
curl -X GET http://localhost:8082/health
```

**Example response:** `{"message":"Service is healthy.","status":"ok","database":"connected"}`

---

## Login

Get a JWT by logging in with email and password. Use the returned `token` in the `Authorization: Bearer <token>` header for protected routes.

```bash
curl -X POST http://localhost:8082/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"secret12"}'
```

**Example response:** `{"message":"Login successful.","data":{"token":"eyJ...","user":{...},"expires_in":86400}}`

**Invalid credentials:** Returns `401` with `{"error":{"code":"UNAUTHORIZED","message":"Invalid email or password."}}`

---

## Users

### Create user

*No auth required (public registration).*

```bash
curl -X POST http://localhost:8082/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","name":"Alice","password":"secret12"}'
```

**Example response:** Returns user with `"id": "b77d82df-7131-4a66-b976-4485893a8382"` — use this ID in the requests below.

---

### List all users

*Requires auth. Returns all users.*

```bash
curl -X GET http://localhost:8082/api/v1/users \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### Get user by ID

*Requires auth. You can only get your own user (path id must match token).*

```bash
curl -X GET http://localhost:8082/api/v1/users/b77d82df-7131-4a66-b976-4485893a8382 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### Update user

*Requires auth. You can only update your own user.*

```bash
curl -X PUT http://localhost:8082/api/v1/users/b77d82df-7131-4a66-b976-4485893a8382 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"name":"Alice Updated","email":"alice.new@example.com"}'
```

---

### Reset password

*Requires auth. You can only reset your own password.*

```bash
curl -X POST http://localhost:8082/api/v1/users/b77d82df-7131-4a66-b976-4485893a8382/reset-password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"new_password":"newsecret12"}'
```

---

### Delete user

*Requires auth. You can only delete your own user.*

```bash
curl -X DELETE http://localhost:8082/api/v1/users/b77d82df-7131-4a66-b976-4485893a8382 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Tasks

### Create task for a user

*Requires auth. Path user id must be your own (you can only create tasks for yourself).*

Use a valid **user UUID** from list users or create user.

```bash
curl -X POST http://localhost:8082/api/v1/users/b77d82df-7131-4a66-b976-4485893a8382/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"title":"My first task","description":"Learn Go APIs","status":"pending"}'
```

**Example response:** Returns task with `"id": "2d09bd64-e6c7-4e6c-a506-75c243601edd"` — use this ID for get/update/delete task.

**Status values:** `pending` | `in_progress` | `done` (default: `pending`)

---

### List tasks for a user

*Requires auth. Path user id must be your own.*

```bash
curl -X GET http://localhost:8082/api/v1/users/b77d82df-7131-4a66-b976-4485893a8382/tasks \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### List all tasks (default pagination)

*Requires auth. Returns only the authenticated user's tasks.*

```bash
curl -X GET "http://localhost:8082/api/v1/tasks" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### List all tasks with pagination

```bash
curl -X GET "http://localhost:8082/api/v1/tasks?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### List all tasks with search

```bash
curl -X GET "http://localhost:8082/api/v1/tasks?search=Go" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### List all tasks for a specific user

*When authenticated, list-all returns only your tasks; `user_id` query param is effectively your own id.*

```bash
curl -X GET "http://localhost:8082/api/v1/tasks?user_id=b77d82df-7131-4a66-b976-4485893a8382" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### Get task by ID

*Requires auth. You can only get your own tasks.*

```bash
curl -X GET http://localhost:8082/api/v1/tasks/2d09bd64-e6c7-4e6c-a506-75c243601edd \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### Update task

*Requires auth. You can only update your own tasks.*

```bash
curl -X PUT http://localhost:8082/api/v1/tasks/2d09bd64-e6c7-4e6c-a506-75c243601edd \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"title":"Updated title","description":"Updated description","status":"in_progress"}'
```

---

### Delete task

*Requires auth. You can only delete your own tasks.*

```bash
curl -X DELETE http://localhost:8082/api/v1/tasks/2d09bd64-e6c7-4e6c-a506-75c243601edd \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Quick reference table

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check (no auth) |
| POST | `/api/v1/login` | Login, get JWT (no auth) |
| POST | `/api/v1/users` | Create user (no auth) |
| GET | `/api/v1/users` | List users (auth) |
| GET | `/api/v1/users/{user_id}` | Get user (auth, self only) |
| PUT | `/api/v1/users/{user_id}` | Update user (auth, self only) |
| POST | `/api/v1/users/{user_id}/reset-password` | Reset password (auth, self only) |
| DELETE | `/api/v1/users/{user_id}` | Delete user (auth, self only) |
| POST | `/api/v1/users/{user_id}/tasks` | Create task (auth, self only) |
| GET | `/api/v1/users/{user_id}/tasks` | List user's tasks (auth, self only) |
| GET | `/api/v1/tasks` | List my tasks, paginated (auth) |
| GET | `/api/v1/tasks/{task_id}` | Get task (auth, own only) |
| PUT | `/api/v1/tasks/{task_id}` | Update task (auth, own only) |
| DELETE | `/api/v1/tasks/{task_id}` | Delete task (auth, own only) |

**Note:** `{user_id}` and `{task_id}` must be valid UUIDs. Protected routes require header: `Authorization: Bearer <token>` (from login).
