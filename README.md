# Go JWT Authentication Service

A JWT-based authentication service built with Go, implementing refresh tokens and rate limiting for enhanced security.

## Features

- **JWT Authentication**: Access tokens for secure API protection
- **Refresh Token System**: Long-lived refresh tokens for seamless re-authentication
- **Rate Limiting**: Protection against brute-force attacks and spam requests
- **Clean Architecture**: Repository and service pattern for maintainable code
- **SQLite Database**: Lightweight database for storing users and tokens
- **Automatic Cleanup**: Background job to remove expired refresh tokens

## Tech Stack

- [Gin Web Framework](https://github.com/gin-gonic/gin): High-performance HTTP web framework
- [GORM](https://gorm.io): Object-Relational Mapping for database interactions
- [JWT](https://github.com/golang-jwt/jwt): JSON Web Token implementation
- [Rate Limiting](https://github.com/JGLTechnologies/gin-rate-limit): Request rate limiting for security
- [Logrus](https://github.com/sirupsen/logrus): Structured logging
- [GoDotEnv](https://github.com/joho/godotenv): Environment variable management

## Project Structure

```
go-jwt-auth/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── database/
│   │   └── database.go          # Database connection
│   ├── handler/
│   │   └── auth_handler.go      # HTTP request handlers
│   ├── middleware/
│   │   └── auth_middleware.go   # Authentication middleware
│   ├── models/
│   │   ├── refresh_token.go     # Refresh token model
│   │   └── user.go              # User model
│   ├── repository/
│   │   ├── refresh_token_repository.go  # Refresh token data access
│   │   └── user_repository.go   # User data access
│   └── service/
│       └── auth_service.go      # Business logic
├── pkg/
│   ├── jwt/
│   │   └── jwt.go               # JWT utilities
│   └── logger/
│       └── logger.go            # Logging configuration
├── build/
├── .env                         # Environment variables
├── .env.example                 # Example environment variables
├── go.mod                       # Go module dependencies
├── go.sum                       # Dependency checksums
└── Makefile                     # Build automation
```

## Installation

### Prerequisites

- Go 1.20 or newer

### Steps

1. Clone the repository:
    ```bash
    git clone https://github.com/0xirvan/go-jwt-auth.git
    cd go-jwt-auth
    ```

2. Create your environment file:
    ```bash
    cp .env.example .env
    ```

3. Edit the `.env` file with your preferred settings:
    ```
    PORT=8080
    DATABASE_PATH=auth.db
    JWT_SECRET=your_very_secure_jwt_secret_here
    REFRESH_SECRET=another_very_secure_refresh_secret_here
    ```

4. Install dependencies:
    ```bash
    go mod tidy
    ```

5. Build and run the application:
    ```bash
    make run
    ```
    
    Or use Go directly:
    ```bash
    go run cmd/api/main.go
    ```

## API Endpoints

### Public Routes

- **POST `/api/v1/auth/register`**: Register a new user
  ```json
  {
    "name": "John Doe",
    "email": "john@example.com",
    "password": "secure_password"
  }
  ```

- **POST `/api/v1/auth/login`**: Login and get tokens
  ```json
  {
    "email": "john@example.com",
    "password": "secure_password"
  }
  ```

- **POST `/api/v1/auth/refresh`**: Refresh access token using refresh token
  ```json
  {
    "refresh_token": "your_refresh_token_here"
  }
  ```

### Protected Routes

- **POST `/api/v1/auth/logout`**: Logout (invalidate refresh token)
  - Requires `Authorization: Bearer your_access_token` header

## Security Features

### Access Token

- Short-lived (15 minutes)
- Secured with HMAC SHA-256

### Refresh Token

- Long-lived (7 days)
- Stored in database
- One refresh token per user
- Atomically replaced on refresh
- Automatic cleanup of expired tokens

### Rate Limiting

- Login: 5 requests per 30 seconds
- Overall API: Configurable rate limits

## Development

### Run Tests

```bash
make test
```

### Clean Build Files

```bash
make clean
```

## Future Improvements

- Email verification system
- Password reset functionality
- Role-based access control
- Docker containerization

