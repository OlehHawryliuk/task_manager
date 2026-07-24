# Task Manager API

![Go](https://img.shields.io/badge/Go-1.21-blue)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-336791)
![Redis](https://img.shields.io/badge/Redis-7-DC382D)

REST API for task management with JWT authentication and Role-Based Access Control (RBAC).

## Features

- ✅ **User Management** - Registration, login, profile updates
- ✅ **Task Management** - Full CRUD operations for tasks
- ✅ **JWT Authentication** - Secure token-based authentication
- ✅ **Role-Based Access Control** - User and admin roles with permissions
- ✅ **Redis Caching** - In-memory caching for optimized reads
- ✅ **Health Check** - Monitor API, database, and Redis status
- ✅ **Swagger/OpenAPI** - Interactive API documentation
- ✅ **Docker Compose** - Production-ready containerization
- ✅ **Unit Tests** - Repository-level testing
- ✅ **GitHub Actions** - Automated CI/CD pipeline
- ✅ **Error Handling** - Uniform error responses with proper HTTP status codes


## Technologies

- **Go 1.21** - Backend programming language
- **Gin** - HTTP web framework
- **GORM** - Object-Relational Mapping
- **PostgreSQL 15** - Relational database
- **Redis 7** - In-memory data store for caching
- **JWT** - JSON Web Tokens for authentication
- **Docker & Docker Compose** - Container orchestration
- **Swagger** - API documentation

## Installation

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Git

### Quick Start with Docker

```bash
# Clone repository
git clone https://github.com/OlehHawryliuk/task-manager.git
cd task-manager

# Start with Docker Compose
docker-compose up --build
```

The API will be available at: `http://localhost:3000`

### Local Setup (without Docker)

```bash
# Clone repository
git clone https://github.com/OlehHawryliuk/task-manager.git
cd task-manager

# Install dependencies
go mod download
go mod tidy

# Configure environment
cp .env.example .env
# Edit .env with your database credentials

# Run application
go run ./cmd/main.go
```

## Docker Compose

### Available Services

- **API Server** - http://localhost:3000
- **PostgreSQL** - localhost:5432
- **Redis** - localhost:6379
- **Swagger UI** - http://localhost:3000/swagger/index.html

### Commands

```bash
# Start services
docker-compose up --build

# Start in background
docker-compose up -d

# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v

# View logs
docker logs -f task_manager_app
```

## API Documentation

### Swagger UI

Interactive API documentation available at:
http://localhost:3000/swagger/index.html

### Base URL
http://localhost:3000

### Authentication

Protected endpoints require a Bearer token:

```bash
Authorization: Bearer <your_jwt_token>
```

### Core Endpoints

#### Authentication
- `POST /auth/register` - Register new user
- `POST /auth/login` - User login with credentials

#### Tasks
- `POST /tasks` - Create a new task
- `GET /tasks` - Get all tasks (cached)
- `GET /tasks/:id` - Get specific task by ID
- `PUT /tasks/:id` - Update task (owner or admin only)
- `DELETE /tasks/:id` - Delete task (owner or admin only)

#### Users
- `POST /users` - Create user (admin only)
- `GET /users` - Get all users (admin only)
- `GET /users/:id` - Get user by ID (owner or admin)
- `PUT /users/:id` - Update user profile (owner or admin)
- `DELETE /users/:id` - Delete user (owner or admin)
- `GET /users/email/:email` - Get user by email

#### System
- `GET /health` - Health check endpoint

## Testing

### Run Unit Tests

```bash
go test ./tests/unit -v
```

### Test Coverage

```bash
go test ./tests/unit -v -coverprofile=coverage.out
```

### GitHub Actions CI/CD

Tests automatically run on every push:
.github/workflows/test.yml

## Project Structure

```
task-manager/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   ├── database.go         # PostgreSQL configuration
│   │   └── redis.go            # Redis configuration
│   ├── handler/
│   │   ├── auth.go             # Authentication handlers
│   │   ├── task.go             # Task handlers
│   │   ├── user.go             # User handlers
│   │   └── health.go           # Health check handler
│   ├── middleware/
│   │   ├── auth.go             # JWT middleware
│   │   └── error.go            # Error handling middleware
│   ├── model/
│   │   ├── task.go             # Task model
│   │   └── user.go             # User model
│   ├── repository/
│   │   ├── task.go             # Task repository
│   │   └── user.go             # User repository
│   ├── service/
│   │   ├── auth.go             # Authentication service
│   │   └── cache.go            # Caching service
│   └── apierror/
│       └── errors.go           # Custom error definitions
├── tests/
│   └── unit/
│       └── repository_test.go   # Unit tests
├── .github/
│   └── workflows/
│       └── test.yml            # GitHub Actions workflow
├── docker-compose.yml          # Docker Compose configuration
├── Dockerfile                  # Docker image definition
├── .env                        # Environment configuration
├── go.mod                      # Go module file
├── go.sum                      # Go dependencies checksums
└── README.md                   # This file
```


## Environment Variables

```env
# Server Configuration
PORT=3000
GIN_MODE=debug

# Database Configuration
DB_HOST=localhost
DB_USER=gorm
DB_PASSWORD=gorm
DB_NAME=gorm
DB_PORT=5432

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT Configuration
JWT_SECRET=your-super-secret-key-min-32-characters-long!!!
```

## Usage Examples

### Register a New User

```bash
curl -X POST http://localhost:3000/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "john_doe",
    "password": "SecurePassword123"
  }'
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "username": "john_doe",
    "role": "user"
  }
}
```

### Login

```bash
curl -X POST http://localhost:3000/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePassword123"
  }'
```

### Create a Task

```bash
curl -X POST http://localhost:3000/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "title": "Learn Go",
    "description": "Complete Go course",
    "done": false
  }'
```

### Get All Tasks

```bash
curl http://localhost:3000/tasks \
  -H "Authorization: Bearer <token>"
```

### Health Check

```bash
curl http://localhost:3000/health

# Response:
# {
#   "status": "healthy",
#   "database": "healthy",
#   "redis": "healthy",
#   "version": "1.0.0",
#   "time": "2026-07-23T17:30:00Z"
# }
```

## CI/CD Pipeline

GitHub Actions automatically:
1. Runs all unit tests on every push
2. Compiles Go code
3. Tests with PostgreSQL service
4. Generates coverage reports

## RBAC (Role-Based Access Control)

### User Roles

- **user** - Regular user
  - Can create, read, update, delete own tasks
  - Can view own profile
  
- **admin** - Administrator
  - Full access to all tasks and users
  - Access to admin-only endpoints

### Access Control

```
Public Endpoints:
├─ POST /auth/register
├─ POST /auth/login
└─ GET /health

Protected Endpoints (all users):
├─ POST /tasks
├─ GET /tasks
├─ GET /tasks/:id
├─ PUT /tasks/:id (owner or admin)
├─ DELETE /tasks/:id (owner or admin)
├─ GET /users/:id (owner or admin)
├─ PUT /users/:id (owner or admin)
├─ DELETE /users/:id (owner or admin)
└─ GET /users/email/:email

Admin-Only Endpoints:
├─ POST /users
└─ GET /users
```


## Performance Features

### Caching Strategy

- GET /tasks endpoint results cached for 5 minutes
- Cache automatically invalidated on create/update/delete operations
- Redis connection pool: 50 connections
- Minimum idle connections: 10

### Database Optimization

- Connection pooling enabled
- Automatic migrations on startup
- Indexed commonly queried fields

### Request Handling

- Timeout: 2 seconds for health checks
- Read/Write timeout: 3 seconds
- Max retry attempts: 3

## Security Features

- **Password Hashing** - bcrypt with salting
- **JWT Tokens** - 1-hour expiration
- **CORS** - Configurable cross-origin requests
- **Error Handling** - No sensitive info in error messages
- **Environment Variables** - Secrets never committed to git

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## Author

**Oleh Hawryliuk**

- GitHub: [@OlehHawryliuk](https://github.com/OlehHawryliuk)
- Email: edifier373@gmail.com
- LinkedIn: [oleh-havryliuk](https://linkedin.com/in/oleh-havryliuk-2b4233419/)

## Resources & References

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM Documentation](https://gorm.io)
- [PostgreSQL Official Docs](https://www.postgresql.org/docs/)
- [Redis Documentation](https://redis.io/documentation)
- [JWT.io](https://jwt.io)
- [Go Best Practices](https://golang.org/doc/effective_go)

## Support

For issues, questions, or suggestions, please open an issue on GitHub.
