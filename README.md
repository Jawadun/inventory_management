# Inventory Management System

## Quick Start

```bash
# 1. Clone the repository
cd /path/to/inventory

# 2. Install Task (alternative to Make)
# https://taskfile.dev/installation/

# 3. Start all services
task up

# 4. Access the application
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

### Using Task (recommended)
```bash
# Install task: https://taskfile.dev/installation/

task up          # Start services
task down        # Stop services
task restart     # Restart services
task logs        # View logs
task clean       # Remove all data

# Run locally (requires PostgreSQL running)
task dev:backend    # Start backend
task dev:frontend   # Start frontend with hot reload

# Database
task db:shell       # Open PostgreSQL shell
task db:backup      # Backup database
task db:migrate     # Run migrations
```

### Using Make
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
task up  # or docker compose up -d postgres

# Migrations run automatically from db/migrations/ folder
# Additional migrations can be run manually:
task db:shell

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