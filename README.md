# Expaze Go Backend

A Go backend API for managing product metadata with MongoDB integration.

## Features

- **Product Management**: CRUD operations for products
- **MongoDB Integration**: Persistent data storage
- **RESTful API**: Standard HTTP methods and responses
- **Pagination**: Built-in pagination support
- **Search**: Full-text search across product names and descriptions
- **Validation**: Request validation and error handling

## API Endpoints

### Products

#### Get All Products
```
GET /api/metadata
```
**Query Parameters:**
- `limit` (optional): Number of products per page (default: 10)
- `offset` (optional): Number of products to skip (default: 0)

**Response:**
```json
{
  "success": true,
  "data": {
    "products": [
      {
        "id": "507f1f77bcf86cd799439011",
        "name": "Product Name",
        "description": "Product Description",
        "price": 99.99,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 1,
    "limit": 10,
    "offset": 0,
    "has_next": false,
    "has_previous": false
  }
}
```

#### Get Product by ID
```
GET /api/metadata/:id
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "Product Name",
    "description": "Product Description",
    "price": 99.99,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Create Product Metadata
```
POST /api/metadata
```

**Request Body:**
```json
{
  "name": "Product Name",
  "description": "Product Description",
  "image": "https://example.com/image.jpg",
  "category_id": "electronics",
  "subcategory_id": "smartphones",
  "mrp": 99.99
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "product_id": "507f1f77bcf86cd799439011",
    "name": "Product Name",
    "description": "Product Description",
    "image": "https://example.com/image.jpg",
    "category_id": "electronics",
    "subcategory_id": "smartphones",
    "mrp": 99.99,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Update Product
```
PUT /api/metadata/:id
```

**Request Body:**
```json
{
  "name": "Updated Product Name",
  "description": "Updated Description",
  "price": 149.99
}
```

**Response:**
```json
{
  "success": true,
  "message": "Product updated successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "Updated Product Name",
    "description": "Updated Description",
    "price": 149.99,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Delete Product
```
DELETE /api/metadata/:id
```

**Response:**
```json
{
  "success": true,
  "message": "Product deleted successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011"
  }
}
```

#### Search Products
```
GET /api/metadata/search?q=search_term
```

**Query Parameters:**
- `q` (required): Search query
- `limit` (optional): Number of results per page (default: 10)
- `offset` (optional): Number of results to skip (default: 0)

**Response:**
```json
{
  "success": true,
  "data": {
    "query": "search_term",
    "results": {
      "products": [
        {
          "id": "507f1f77bcf86cd799439011",
          "name": "Product Name",
          "description": "Product Description",
          "price": 99.99,
          "created_at": "2024-01-01T00:00:00Z",
          "updated_at": "2024-01-01T00:00:00Z"
        }
      ],
      "total": 1,
      "limit": 10,
      "offset": 0,
      "has_next": false,
      "has_previous": false
    }
  }
}
```

## Setup

1. **Install Dependencies:**
   ```bash
   go mod tidy
   ```

2. **Environment Variables:**
   Create a `.env` file with:
   ```
   MONGO_URI=mongodb://localhost:27017
   ```

3. **Run the Application:**
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`

## Project Structure

```
├── config/
│   └── db.go              # Database configuration
├── controllers/
│   └── metadataController.go  # Product controller
├── models/
│   └── product.go         # Product model and database operations
├── routes/
│   ├── routes.go          # Main routes setup
│   └── metadataRoutes.go  # Product routes
└── main.go               # Application entry point
```

## Database Schema

### Metadata Collection
```json
{
  "_id": "ObjectId",
  "metadata_product_id": "ObjectId (auto-generated)",
  "metadata_name": "string (required)",
  "metadata_description": "string",
  "metadata_image": "string",
  "metadata_category_id": "string",
  "metadata_subcategory_id": "string",
  "metadata_mrp": "number (required)",
  "metadata_created_at": "timestamp",
  "metadata_updated_at": "timestamp"
}
```

## Error Handling

All endpoints return consistent error responses:

```json
{
  "error": "Error message description"
}
```

Common HTTP status codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `404` - Not Found
- `500` - Internal Server Error 