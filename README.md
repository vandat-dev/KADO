# Base Go Backend

A modern Golang backend application built with clean architecture principles, providing a solid foundation for building scalable web services.

## Features

- **Clean Architecture**: Structured using domain-driven design principles
- **RESTful API**: Built with Gin web framework
- **Dependency Injection**: Uses Google Wire for dependency management
- **Database**: MySQL integration with GORM
- **Caching**: Redis for performance optimization
- **Message Broker**: Kafka integration (optional)
- **API Documentation**: Swagger for API documentation
- **Containerization**: Docker and Docker Compose support
- **Logging**: Structured logging with Zap

## Getting Started

### Prerequisites

- Go 1.23+
- Docker and Docker Compose
- Make (optional, for convenience commands)

### Installation

1. Clone the repository:
   ```
   git clone <repository-url>
   cd base_go_be
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Setup environment variables:
   ```
   cp example.env .env
   # Edit .env file with your configuration
   ```

4. Start required services using Docker:
   ```
   docker-compose up -d mysql redis
   ```
   
   Or simply:
   ```
   make mysql
   ```

5. Run database migrations manually:
   ```
   # See migrations/README.md for detailed instructions
   mysql -h 127.0.0.1 -P 33306 -u admin -p123123 go_database < migrations/initial_version_go.sql
   mysql -h 127.0.0.1 -P 33306 -u admin -p123123 go_database < migrations/1_create_table.up.sql
   ```

### Running the Application

Run the application locally:
```
make run
```

Or using Docker:
```
docker-compose up
```

### API Testing

Test a sample endpoint:
```
curl http://localhost:8386/v1/user/test
```

Visit Swagger documentation:
```
http://localhost:8386/docs/index.html
```

## Development

### Dependency Injection

This project uses Google Wire for dependency injection:

1. Install Wire if not already installed:
   ```
   go install github.com/google/wire/cmd/wire@latest
   ```

2. Generate dependency injection code:
   ```
   cd internal/wire
   wire
   ```

### API Documentation
go install github.com/swaggo/swag/cmd/swag@latest
Generate Swagger documentation:
```
make swag
```
or
```
swag init -g ./cmd/server/main.go -o docs
```

## Configuration

The application now uses environment variables for configuration. All settings are loaded from `.env` file:

- Copy `example.env` to `.env`
- Modify the values according to your environment
- Environment variables take precedence over default values

## Project Structure

- `cmd/`: Application entry points
- `config/`: Legacy configuration files (kept for reference)
- `docs/`: API documentation
- `internal/`: Application core code
  - `controller/`: HTTP handlers
  - `service/`: Business logic
  - `repo/`: Data access layer
  - `model/`: Domain models
  - `middlewares/`: HTTP middlewares
  - `routers/`: API routes
- `migrations/`: Database migrations (run manually)
- `pkg/`: Reusable packages
- `tests/`: Test files
- `example.env`: Example environment configuration

## Deployment

The application can be deployed using Docker:

```
docker build -t base_go_be .
docker run -p 8386:8386 base_go_be
```

Or using Docker Compose:

```
docker-compose up -d
```