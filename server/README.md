# Realtime Chat App API

This is the backend API for a realtime chat application built with Go, Gin, and PostgreSQL. It supports user authentication (JWT & refresh tokens), room-based chat via WebSockets, and basic health checks.

## Features

- User registration and login with JWT authentication
- Refresh token support (secure, stored in DB)
- Room creation and joining via WebSockets
- List available rooms and clients in a room
- Health check endpoint

## Project Structure

```
cmd/
server/
  api/
    main.go
  db/
    db.go
    migrations/
  internal/
    auth/
    config/
    handlers/
    model/
    routes/
    services/
    store/
    utils/
    ws/
  pkg/
    env.go
  ...
```

## Getting Started

### Prerequisites

- Go 1.20+
- PostgreSQL
- [migrate](https://github.com/golang-migrate/migrate) (for DB migrations)
- Docker (optional, for running PostgreSQL)

### Environment Variables

Create a `.env` file in the `server/` directory:

```
PORT=8080
ENV=development
DB_ADDR=postgres://admin:adminpassword123@localhost:5434/realtimechatapp?sslmode=disable
JWT_SECRET=your_jwt_secret
JWT_ISS=realtimechatapp
JWT_AUD=realtimechatapp
JWT_EXP=15m
```

### Database Setup

You can start a local PostgreSQL instance using Docker:

```sh
cd server
make dkup
```

Run migrations:

```sh
make mup
```

### Running the API

```sh
cd server
go run ./api/main.go
```

Or use [Air](https://github.com/cosmtrek/air) for live reload:

```sh
air
```

## API Endpoints

### Auth

- `POST /api/v1/signup` — Register a new user
- `POST /api/v1/login` — Login and receive JWT & refresh token
- `POST /api/v1/logout` — Logout (clears cookies)
- `POST /api/v1/refresh` — Refresh JWT using refresh token

### Health

- `GET /api/v1/healthcheck` — API health check

### WebSocket

- `POST /api/v1/ws/createRoom` — Create a new chat room
- `GET /api/v1/ws/joinRoom/:roomId?userId=...&username=...` — Join a chat room via WebSocket
- `GET /api/v1/ws/getRooms` — List all rooms
- `GET /api/v1/ws/getClients/:roomId` — List clients in a room

## License

MIT