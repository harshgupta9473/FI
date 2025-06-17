# FI - Inventory Management Backend

A Go-based inventory management backend with user authentication and product APIs, built using clean architecture principles and Dockerized for easy setup.

---

## Architecture: Clean Layered Design

- **Handlers** receive HTTP requests and delegate to services.
- **Services** contain business logic and call repositories.
- **Repositories** abstract data access (e.g., PostgreSQL).
- **DTOs** define strict request/response shapes.
- **Middleware** adds cross-cutting concerns like JWT authentication.
- **Dependency Injection (DI)** manages wiring of components.

This project follows Clean Architecture, promoting:

- Separation of concerns
- Interface-driven development
- Testability and maintainability

---

## Features

- User registration and login (with JWT)
- Secure password hashing using bcrypt
- Product management (Add, Update, List)
- Pagination support on listing products
- Structured logging using Zap
- Environment-driven configuration
- Fully Dockerized

---

## Getting Started

### Clone the repo

```bash
git clone https://github.com/harshgupta9473/FI.git
cd FI
```

# Run with Docker
```bash
docker-compose up --build
```

# Access the API
```bash
http://localhost:8080
```

# Documeantation Link

[https://documenter.getpostman.com/view/34442065/2sB2x9hpn9](https://documenter.getpostman.com/view/34442065/2sB2x9hpn9)