# Inventory Management System

## Quick Start

```bash
# 1. Clone the repository
cd /path/to/inventory

# 2. Start all services
docker-compose up -d

# 3. Access the application
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
```

## Prerequisites

- Docker Engine 20.10+
- Docker Compose 2.0+
- 4GB RAM minimum

## Services

| Service | Port | Description |
|---------|------|-------------|
| Frontend | 3000 | React UI |
| Backend | 8080 | Go API |
| Database | 5432 | PostgreSQL |

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| DB_HOST | postgres | Database hostname |
| DB_PORT | 5432 | Database port |
| DB_USER | postgres | Database user |
| DB_PASSWORD | postgres | Database password |
| DB_NAME | inventory | Database name |
| JWT_SECRET | (change-me) | JWT signing secret |
| SERVER_PORT | 8080 | Backend port |
| ENV | development | Environment mode |

## Development Commands

### Using Make (recommended)
```bash
make up          # Start services
make down       # Stop services
make restart    # Restart services
make logs       # View logs
make clean      # Remove all data
```

### Using Docker Compose directly
```bash
# Build images
docker-compose build

# Start services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Reset everything (including database)
docker-compose down -v
```

## Development Mode (for local development)

```bash
# Start database only
docker-compose up -d postgres

# Run migrations
# ( You'll need to run these manually for now)
# psql -h localhost -U postgres -d inventory -f db/migrations/001_initial_schema.sql

# Start backend
cd backend
go run ./cmd/server

# Start frontend (new terminal)
cd frontend
npm run dev
```

## Building for Production

```bash
# Build and start
docker-compose up -d --build

# The application will be available at http://localhost:3000
```

## Troubleshooting

### Database connection issues
```bash
# Check if postgres is running
docker-compose ps

# Check postgres logs
docker-compose logs postgres
```

### Reset everything
```bash
make clean
make up
```

### Check API health
```bash
curl http://localhost:8080/api/public/stats
```

## Project Structure

```
├── backend/          # Go API
│   ├── cmd/        # Entry points
│   ├── internal/  # Business logic
│   ├── db/        # Database migrations
│   └── Dockerfile
├── frontend/       # React UI
│   ├── src/       # Source code
│   ├── Dockerfile
│   └── nginx.conf
├── docker-compose.yml
├── Makefile
└── .env.example
```

## Security Notes

- Change `JWT_SECRET` in production
- Use strong database passwords
- Enable HTTPS in production (configure nginx)
- Update CORS settings for production domains