# Football Management API

REST API manajemen sepak bola — mengelola tim, pemain, dan pertandingan.

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL 14+
- **Authentication**: JWT RS256
- **Database Library**: sqlx + lib/pq
- **Password Hashing**: bcrypt
- **Swagger**: swaggo/swag
- **Logging**: logrus
- **Configuration**: Viper

## API Endpoints

### Auth (Public)

| Method | Endpoint             | Description   |
| ------ | -------------------- | ------------- |
| POST   | `/v1/auth/register`  | Register user |
| POST   | `/v1/auth/login`     | Login user    |

### Teams (Auth Required)

| Method | Endpoint                  | Description                  |
| ------ | ------------------------- | ---------------------------- |
| GET    | `/v1/teams`               | Get all teams                |
| GET    | `/v1/teams/:id`           | Get team by ID               |
| GET    | `/v1/teams/:id/players`   | Get all players of a team    |
| POST   | `/v1/teams`               | Create team (multipart/form) |
| PUT    | `/v1/teams/:id`           | Update team (multipart/form) |
| DELETE | `/v1/teams/:id`           | Delete team (soft delete)    |

### Players (Auth Required)

| Method | Endpoint           | Description                |
| ------ | ------------------ | -------------------------- |
| GET    | `/v1/players`      | Get all players            |
| GET    | `/v1/players/:id`  | Get player by ID           |
| POST   | `/v1/players`      | Create player              |
| PUT    | `/v1/players/:id`  | Update player              |
| DELETE | `/v1/players/:id`  | Delete player (soft delete)|

### Matches (Auth Required)

| Method | Endpoint                    | Description               |
| ------ | --------------------------- | ------------------------- |
| GET    | `/v1/matches`               | Get all matches           |
| GET    | `/v1/matches/:id`           | Get match by ID           |
| GET    | `/v1/matches/:id/report`    | Get match report          |
| POST   | `/v1/matches`               | Create match schedule     |
| PUT    | `/v1/matches/:id`           | Update match schedule     |
| DELETE | `/v1/matches/:id`           | Delete match (soft delete)|
| POST   | `/v1/matches/:id/result`    | Submit match result       |

## Makefile Commands

```bash
# Install dependencies
make deps

# Run the application
make run
# or manually
go run cmd/main.go

# Database migrations
make migrate.up
# or manually
go run migration/main/main.go up

make migrate.rollback
# or manually
go run migration/main/main.go rollback

# Regenerate Swagger docs (requires swag CLI installed)
make docs-update
# or manually
$(go env GOPATH)/bin/swag init -g cmd/main.go -o swagger/v1 --ot go,json,yaml --pd true
```

## Setup Instructions

### Prerequisites

- Go 1.21+
- PostgreSQL 14+
- Make
- [swag CLI](https://github.com/swaggo/swag) (regenerate docs)

### 1. Install dependencies

```bash
make deps
```

### 2. Generate RSA keys

```bash
mkdir -p keys
openssl genpkey -algorithm RSA -out keys/private.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -in keys/private.pem -pubout -out keys/public.pem
```

### 3. Buat database PostgreSQL

```bash
psql -U postgres -c "CREATE DATABASE football_db;"
```

### 4. Konfigurasi environment

Buat file `.env` di root project:

```env
SERVICE_NAME=football-api
ENV=development
BIND_ADDRESS=8080
LOG_LEVEL=5

JWT_PRIVATE_KEY_PATH=keys/private.pem
JWT_PUBLIC_KEY_PATH=keys/public.pem

POSTGRES_CONN_URI=postgres://postgres:yourpassword@localhost:5432/football_db?sslmode=disable
POSTGRES_MAX_POOL_SIZE=100
POSTGRES_MAX_IDLE_CONNECTIONS=10
POSTGRES_MAX_IDLE_TIME=10m
POSTGRES_MAX_LIFE_TIME=30m

TRANSLATION_FILE_PATH=i18n/definitions
TRANSLATION_LANG_PREFERENCES=id-ID
TRANSLATION_DEAULT_LANG=en-ID
```

### 5. Jalankan migrasi database

```bash
make migrate.up
```

### 6. Jalankan aplikasi

```bash
# Menggunakan Makefile
make run

# Manual
go run cmd/main.go
```

Server berjalan di `http://localhost:8080`

### 7. Swagger UI

Buka `http://localhost:8080/swagger/index.html`

Untuk menggunakan endpoint yang membutuhkan auth:
1. Login via `POST /v1/auth/login` untuk mendapatkan token
2. Klik tombol **Authorize** di kanan atas
3. Masukkan JWT token (bisa dengan atau tanpa prefix `Bearer`)
4. Klik **Authorize** → **Close**

---

## API Reference

Base URL: `http://localhost:8080`

Format response selalu:

```json
{
  "data": {},
  "error": null,
  "success": true,
  "metadata": { "request_id": "..." }
}
```

---

### Auth

#### Register

```bash
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Admin",
    "email": "admin@ayo.id",
    "password": "Password123!"
  }'
```

**Success Response (201)**:

```json
{
  "data": {
    "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "error": null,
  "success": true,
  "metadata": { "request_id": "..." }
}
```

#### Login

```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@ayo.id",
    "password": "Password123!"
  }'
```

**Success Response (200)**:

```json
{
  "data": {
    "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "error": null,
  "success": true,
  "metadata": { "request_id": "..." }
}
```

---

### Teams

#### Create Team (with logo upload)

```bash
curl -X POST http://localhost:8080/v1/teams \
  -H "Authorization: Bearer <token>" \
  -F "name=Manchester United" \
  -F "year_founded=1878" \
  -F "city=Manchester" \
  -F "address=Old Trafford" \
  -F "logo=@/path/to/logo.png"
```

**Success Response (201)**:

```json
{
  "data": {
    "id": 1,
    "name": "Manchester United",
    "logo": "/uploads/teams/abc123.png",
    "year_founded": 1878,
    "address": "Old Trafford",
    "city": "Manchester",
    "created_at": "2026-02-22 10:00:00",
    "updated_at": "2026-02-22 10:00:00"
  },
  "error": null,
  "success": true,
  "metadata": { "request_id": "..." }
}
```

#### Get All Teams

```bash
curl http://localhost:8080/v1/teams \
  -H "Authorization: Bearer <token>"
```

#### Get Team by ID

```bash
curl http://localhost:8080/v1/teams/1 \
  -H "Authorization: Bearer <token>"
```

#### Get Players by Team

```bash
curl http://localhost:8080/v1/teams/1/players \
  -H "Authorization: Bearer <token>"
```

#### Update Team

```bash
curl -X PUT http://localhost:8080/v1/teams/1 \
  -H "Authorization: Bearer <token>" \
  -F "city=Manchester Updated"
```

#### Delete Team

```bash
curl -X DELETE http://localhost:8080/v1/teams/1 \
  -H "Authorization: Bearer <token>"
```

---

### Players

#### Create Player

Posisi yang valid: `penyerang`, `gelandang`, `bertahan`, `penjaga_gawang`

```bash
curl -X POST http://localhost:8080/v1/players \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "team_id": 1,
    "name": "Cristiano Ronaldo",
    "height": 187.0,
    "weight": 83.0,
    "position": "penyerang",
    "jersey_number": 7
  }'
```

**Success Response (201)**:

```json
{
  "data": {
    "id": 1,
    "team_id": 1,
    "name": "Cristiano Ronaldo",
    "height": 187,
    "weight": 83,
    "position": "penyerang",
    "jersey_number": 7,
    "created_at": "2026-02-22 10:00:00",
    "updated_at": "2026-02-22 10:00:00"
  },
  "error": null,
  "success": true,
  "metadata": { "request_id": "..." }
}
```

#### Get All Players

```bash
curl http://localhost:8080/v1/players \
  -H "Authorization: Bearer <token>"
```

#### Get Player by ID

```bash
curl http://localhost:8080/v1/players/1 \
  -H "Authorization: Bearer <token>"
```

#### Update Player

```bash
curl -X PUT http://localhost:8080/v1/players/1 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "jersey_number": 9
  }'
```

#### Delete Player

```bash
curl -X DELETE http://localhost:8080/v1/players/1 \
  -H "Authorization: Bearer <token>"
```

---

### Matches

#### Create Match

```bash
curl -X POST http://localhost:8080/v1/matches \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "home_team_id": 1,
    "away_team_id": 2,
    "match_date": "2026-03-01",
    "match_time": "19:00"
  }'
```

**Success Response (201)**:

```json
{
  "data": {
    "id": 1,
    "home_team": { "id": 1, "name": "Manchester United", "logo": "..." },
    "away_team": { "id": 2, "name": "Arsenal", "logo": "..." },
    "match_date": "2026-03-01",
    "match_time": "19:00",
    "home_score": null,
    "away_score": null,
    "status": "scheduled",
    "goals": [],
    "created_at": "2026-02-22 10:00:00",
    "updated_at": "2026-02-22 10:00:00"
  },
  "error": null,
  "success": true,
  "metadata": { "request_id": "..." }
}
```

#### Submit Match Result

```bash
curl -X POST http://localhost:8080/v1/matches/1/result \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "home_score": 2,
    "away_score": 1,
    "goals": [
      { "player_id": 1, "goal_minute": 23 },
      { "player_id": 1, "goal_minute": 67 },
      { "player_id": 5, "goal_minute": 45 }
    ]
  }'
```

#### Get Match Report

```bash
curl http://localhost:8080/v1/matches/1/report \
  -H "Authorization: Bearer <token>"
```

**Success Response (200)**:

```json
{
  "data": {
    "match_id": 1,
    "match_date": "2026-03-01",
    "match_time": "19:00",
    "home_team": { "id": 1, "name": "Manchester United", "logo": "..." },
    "away_team": { "id": 2, "name": "Arsenal", "logo": "..." },
    "home_score": 2,
    "away_score": 1,
    "final_status": "home_win",
    "top_scorer": {
      "player_id": 1,
      "name": "Cristiano Ronaldo",
      "goals": 2
    },
    "home_team_total_wins": 5,
    "away_team_total_wins": 3,
    "goals": [
      { "id": 1, "player_id": 1, "player_name": "Cristiano Ronaldo", "goal_minute": 23 },
      { "id": 2, "player_id": 1, "player_name": "Cristiano Ronaldo", "goal_minute": 67 },
      { "id": 3, "player_id": 5, "player_name": "Bukayo Saka", "goal_minute": 45 }
    ]
  },
  "error": null,
  "success": true,
  "metadata": { "request_id": "..." }
}
```

`final_status` dapat bernilai: `home_win`, `away_win`, atau `draw`

---

## Database Schema

```
users (1) ─────────────────────────────── (auth only)

teams (1) ──────────< (N) players
teams (1) ──────────< (N) matches (as home_team)
teams (1) ──────────< (N) matches (as away_team)
matches (1) ────────< (N) goals
players (1) ────────< (N) goals
```

### Tabel Utama

| Tabel     | Keterangan                                       |
| --------- | ------------------------------------------------ |
| `users`   | Admin credentials (email, password, role=admin)  |
| `teams`   | Data tim sepak bola                              |
| `players` | Data pemain beserta posisi dan nomor jersey      |
| `matches` | Jadwal & hasil pertandingan                      |
| `goals`   | Detail gol per pertandingan                      |

---

## Security

- Password di-hash **bcrypt**
- JWT algoritma **RS256** (RSA 2048-bit asymmetric)
- Private key untuk signing token (AuthService)
- Public key untuk verifikasi token (JWT Middleware)
- Token berlaku **7 hari**
- Semua endpoint selain `/v1/auth/*` membutuhkan token

## Troubleshooting

### PostgreSQL connection refused

```bash
# macOS
brew services start postgresql

# Linux
sudo systemctl start postgresql
```

### Tabel tidak ditemukan

```bash
make migrate.up
```

### Swagger tidak terupdate

```bash
make docs-update
```

## Developers

Developed by [Lutfi M](https://github.com/ramaqqqqq), 2026
