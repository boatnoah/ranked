# Ranked

Real-time Leaderboard system for ranking and score in Go

# Features

- User Authentication: Users should be able to register and log in to the system.

- Match Submission: Users submit results that adjust trophies.

- Leaderboard Updates: Display a global leaderboard showing the top users across all games.

- User Rankings: Users should be able to view their rankings on the leaderboard.

- Top Players Report: Generate reports on the top players for a specific period.

# Architecture

- **Postgres** — source of truth. Stores users, sessions, and match history.
- **Redis** — sorted sets for real-time leaderboard rankings.

### Write flow (match submission)

1. Start Postgres transaction
2. Insert match record into Postgres
3. `ZINCRBY` the user's trophy count in Redis
4. If both succeed, commit. If either fails, rollback.

### Read flow (leaderboard)

1. Read from Redis
2. If Redis is empty (data loss), rebuild from Postgres, then read from Redis again
3. If Redis is down, return error — leaderboard is unavailable

# API Spec

## Auth

### Register

```
POST /v1/register
{
    "username": "player1",
    "email": "player1@example.com",
    "password": "secret"
}
```

### Login

```
POST /v1/login
{
    "email": "player1@example.com",
    "password": "secret"
}
```

Response:

```json
{
  "session_token": "abc123..."
}
```

### Logout

```
POST /v1/logout
```

## Leaderboard

### Submit match result

```
POST /v1/ranked/match
{
    "result": "win",
    "crowns": 3
}
```

- `result` — `"win"` or `"loss"`
- `crowns` — `1`, `2`, or `3`

Trophy calculation (server-side, random range based on crowns):

| Crowns | Win range  | Loss range |
| ------ | ---------- | ---------- |
| 1      | +26 to +30 | -18 to -22 |
| 2      | +28 to +32 | -20 to -24 |
| 3      | +30 to +34 | -22 to -26 |

Response:

```json
{
  "trophies_changed": 32,
  "total_trophies": 472,
  "rank": 5
}
```

### Get leaderboard

```
GET /v1/ranked/leaderboard?limit=10
```

Response:

```json
{
  "leaderboard": [
    { "rank": 1, "user_id": 12, "username": "player1", "trophies": 1250 },
    { "rank": 2, "user_id": 7, "username": "player2", "trophies": 1100 }
  ]
}
```

### Get my rank

```
GET /v1/ranked/me
```

Response:

```json
{
  "user_id": 12,
  "username": "player1",
  "trophies": 472,
  "rank": 5
}
```

# Database Schema

## Postgres

### users

| Column     | Type      | Notes            |
| ---------- | --------- | ---------------- |
| id         | BIGSERIAL | Primary key      |
| username   | CITEXT    | Unique, not null |
| email      | CITEXT    | Unique, not null |
| password   | BYTEA     | bcrypt hash      |
| created_at | TIMESTAMP | Default now()    |

### matches

| Column           | Type        | Notes                |
| ---------------- | ----------- | -------------------- |
| id               | BIGSERIAL   | Primary key          |
| user_id          | BIGINT      | FK to users(id)      |
| result           | VARCHAR(10) | "win" or "loss"      |
| crowns           | INT         | 1-3                  |
| trophies_changed | INT         | Positive or negative |
| submitted_at     | TIMESTAMP   | Default now()        |

### user_trophies

| Column     | Type      | Notes                          |
| ---------- | --------- | ------------------------------ |
| user_id    | BIGINT    | PK, FK to users(id)            |
| trophies   | INT       | Current total, non-negative    |
| updated_at | TIMESTAMP | Default now(), last write time |

## Redis

### Sorted set: `leaderboard`

- Member: user ID
- Score: current trophy count
- Updated via `ZINCRBY` on each match submission
