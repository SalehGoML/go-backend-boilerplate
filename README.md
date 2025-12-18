# Go Backend Boilerplate

A professional **Golang backend boilerplate** designed for building scalable RESTful CRUD applications. This project follows a **Clean Architecture** pattern to ensure maintainability, testability, and separation of concerns.

It's an excellent starting point for real-world backend services, featuring database integration, migrations, logging, and full Docker support.

## Features
- **Clean Architecture**: Layered structure with `internal/appl` for application logic and `internal/route` for routing.
- **Database Integration**: Connection handling and schema migrations (via `migrations/` folder).
- **Docker Support**: Fully containerized with `docker-compose.yml` (includes app and database services).
- **Environment Configuration**: Loads variables from `.env` file using godotenv.
- **Testing**: Dedicated `test/` folder for unit and integration tests.
- **Production-Ready Server**: HTTP server with proper timeouts and logging.

## How It Works
1. **Initialization** (`main.go`):
   - Loads environment variables from `.env`.
   - Initializes the application with database connection and logger (`app.NewApplication()`).
2. **Routing**:
   - API routes are defined and registered in `internal/route` using `routes.SetupRoutes()`.
3. **Application Layer**:
   - Business logic and repositories are implemented in `internal/appl` and `database/`.
4. **Database**:
   - Establishes connection, applies migrations, and gracefully closes on shutdown.
5. **Server Startup**:
   - Starts an HTTP server on the configured port (default: 8080) and handles incoming CRUD requests.
