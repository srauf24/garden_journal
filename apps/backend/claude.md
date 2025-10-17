#  Garden Journal Project Context

##  Project Overview
**Garden Journal** is a modern web application for plant care and garden management built with **Go (backend)** and **TypeScript/React (frontend)**.  
It follows a **monorepo architecture** using **Turborepo** for efficient builds and development workflow.

---

## ⚙ Technical Stack

### 🐹 Backend (Go)
- **Go 1.24+**
- **Echo** framework for REST API
- **PostgreSQL 16+** with connection pooling
- **Redis 8+** for background jobs
- **Clerk SDK** for authentication
- **New Relic** for APM
- **Resend** for email services

### ⚛ Frontend (TypeScript / React)
- **React 19.1.0**
- **TypeScript 5.8.2**
- **Vite 7.0.4**
- **TanStack Query** for data fetching
- **Clerk** for authentication
- **Tailwind CSS** for styling
- **React Router** for navigation

---

## 🏗 Architecture

###  Backend Structure
```bash
apps/backend/
├── cmd/            # Application entry points
├── internal/       # Private application code
│   ├── config/        # Configuration management
│   ├── database/      # Database connections and migrations
│   ├── handler/       # HTTP request handlers
│   ├── service/       # Business logic layer
│   ├── repository/    # Data access layer
│   ├── model/         # Domain models
│   ├── middleware/    # HTTP middleware
│   ├── validation/    # Request validation
│   └── lib/           # Shared utilities
├── templates/      # Email templates
├── static/         # Static files
└── tests/          # Test suites
````

###  Key Backend Features

1. **Configuration Management**

   * Environment-based configuration using Koanf
   * Structured validation
   * Support for multiple environments

2. **Database Layer**

   * PostgreSQL with connection pooling
   * Migration system using `tern`
   * Configurable connection settings

3. **Authentication & Security**

   * Clerk SDK integration
   * JWT validation
   * CORS configuration
   * Rate limiting
   * Secure headers

4. **Background Processing**

   * Redis-based job queue
   * Async task processing
   * Email notifications

5. **Observability**

   * New Relic APM integration
   * Structured logging (`zerolog`)
   * Health checks
   * Performance monitoring

---

###  Frontend Structure

```bash
apps/frontend/
├── src/
│   ├── components/   # Reusable UI components
│   ├── features/     # Feature-specific code
│   ├── hooks/        # Custom React hooks
│   ├── pages/        # Route pages
│   ├── api/          # API integration
│   ├── utils/        # Utility functions
│   └── styles/       # Global styles
└── tests/            # Frontend tests
```

---

##  Development Workflow

###  Backend Development

1. **Environment Setup**

   ```bash
   cd apps/backend
   go mod download
   cp .env.example .env
   ```

2. **Database Management**

   ```bash
   task migrations:new name=<migration_name>  # Create migration
   task migrations:up                         # Apply migrations
   ```

3. **Running the Server**

   ```bash
   task run    # Start server
   task test   # Run tests
   task tidy   # Format and tidy
   ```

### ⚛ Frontend Development

1. **Setup**

   ```bash
   bun install
   ```

2. **Development**

   ```bash
   bun dev     # Start dev server
   bun build   # Production build
   bun lint    # Run linter
   ```

---

##  API Integration

* REST API with **OpenAPI/Swagger specification**
* Type-safe API client using **ts-rest**
* Automatic type generation from OpenAPI specs
* Request/response validation
* Error handling with retries

---

##  Testing Strategy

### Backend

* Unit tests for business logic
* Integration tests with **testcontainers**
* API endpoint tests
* Performance benchmarks

### Frontend

* Component tests with **React Testing Library**
* Integration tests
* E2E tests with **Cypress**

---

##  Deployment & Operations

1. **Environment Configuration**

   * Environment-specific settings
   * Secret management
   * Feature flags

2. **Monitoring**

   * APM with New Relic
   * Error tracking
   * Performance monitoring
   * Log aggregation

3. **Security**

   * Authentication with Clerk
   * Authorization middleware
   * Input validation
   * Rate limiting
   * CORS policies

---

##  Future Considerations

1. **Scalability**

   * Horizontal scaling of API
   * Caching strategies
   * Database optimization

2. **Feature Enhancements**

   * Real-time updates
   * Mobile responsiveness
   * Offline support
   * Data export/import

3. **Integrations**

   * Weather API integration
   * Plant database
   * Image recognition
   * Social sharing

