# Onion Architecture Implementation

This project has been refactored to follow the Onion Architecture pattern, providing better separation of concerns, testability, and maintainability.

## Architecture Layers

### 1. Domain Layer (Core)
**Location**: `domain/`

#### Entities (`domain/entities/`)
- **`product.go`**: Core Product entity with business logic and validation
- **`errors.go`**: Domain-specific error definitions

The domain layer contains:
- Core business entities
- Business rules and validation
- Domain-specific errors
- No external dependencies

#### Repositories (`domain/repositories/`)
- **`product_repository.go`**: Interface defining data access operations

The repository layer contains:
- Data access interfaces
- No implementation details
- Defines contracts for data operations

### 2. Infrastructure Layer
**Location**: `infrastructure/mongodb/`

- **`product_repository.go`**: MongoDB implementation of the ProductRepository interface

The infrastructure layer contains:
- External framework implementations
- Database connections and operations
- External service integrations
- Implements repository interfaces

### 3. Use Case Layer
**Location**: `usecase/`

- **`product_usecase.go`**: Business logic orchestration

The use case layer contains:
- Application business logic
- Orchestrates operations between domain and repositories
- Handles data transformation
- Implements application-specific rules

### 4. Handler Layer
**Location**: `handlers/`

- **`product_handler.go`**: HTTP request/response handling

The handler layer contains:
- HTTP request/response handling
- Input validation and parsing
- Delegates business logic to use cases
- Handles HTTP-specific concerns

## Project Structure

```
Expaze_Go_Backend/
├── domain/
│   ├── entities/
│   │   ├── product.go
│   │   └── errors.go
│   └── repositories/
│       └── product_repository.go
├── infrastructure/
│   └── mongodb/
│       └── product_repository.go
├── usecase/
│   └── product_usecase.go
├── handlers/
│   └── product_handler.go
├── routes/
│   ├── routes.go
│   └── product_routes.go
├── config/
│   └── db.go
├── main.go
└── ARCHITECTURE.md
```

## Dependency Flow

```
Handler → UseCase → Repository Interface ← Repository Implementation
```

- Dependencies point inward
- Inner layers don't know about outer layers
- Outer layers depend on inner layer interfaces

## API Endpoints

The product endpoints are available at:

- `GET /api/products` - Get all products with pagination
- `GET /api/products/:id` - Get specific product by ID
- `POST /api/products` - Create new product
- `PUT /api/products/:id` - Update existing product
- `DELETE /api/products/:id` - Delete product
- `GET /api/products/search?q=query` - Search products

## Benefits of Onion Architecture

1. **Separation of Concerns**: Each layer has a specific responsibility
2. **Testability**: Easy to mock dependencies and test in isolation
3. **Maintainability**: Changes in one layer don't affect others
4. **Flexibility**: Easy to swap implementations (e.g., different databases)
5. **Scalability**: Clear structure for adding new features

## Migration Completed

✅ **Old Structure Removed**: 
- Removed `models/product.go` (mixed concerns)
- Removed `controllers/metadataController.go` (business logic in controllers)
- Removed `routes/metadataRoutes.go` (old routing)

✅ **New Structure Implemented**:
- Clean separation of concerns
- Dependency inversion
- Better testability
- Maintainable codebase

## Running the Application

The application now uses a clean onion architecture:

```bash
go run main.go
```

All functionality is now handled through the new onion architecture layers, providing a clean and maintainable codebase. 