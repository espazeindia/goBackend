# Espaze Backend API Documentation

## Base URL
```
http://localhost:8080
```

## Authentication
All protected endpoints require a valid JWT token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

## Login APIs

### 1. Login Operational Guy
**POST** `/login/operational_guy`

Authenticates an operational guy and returns a JWT token.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (Success - 200):**
```json
{
  "success": true,
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response (Error - 401):**
```json
{
  "success": false,
  "error": "Invalid credentials",
  "message": "Email or password is incorrect"
}
```

**Response (Error - 400):**
```json
{
  "success": false,
  "error": "Validation error",
  "message": "Field validation for 'email' failed on the 'email' tag"
}
```

---

### 2. Register Operational Guy
**POST** `/login/operational_guy/register`

Registers a new operational guy account.

**Request Body:**
```json
{
  "email": "newuser@example.com",
  "password": "password123",
  "name": "John Doe",
  "phoneNumber": "1234567890",
  "address": "123 Main Street, City, State 12345",
  "emergencyContactNumber": "0987654321"
}
```

**Response (Success - 201):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "user_id": "507f1f77bcf86cd799439011"
}
```

**Response (Error - 400):**
```json
{
  "success": false,
  "error": "User already exists",
  "message": "An account with this email already exists"
}
```

---

### 3. Refresh Token
**POST** `/login/refresh`

Refreshes an expired access token using a refresh token.

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response (Success - 200):**
```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 3600,
    "token_type": "Bearer"
  }
}
```

**Response (Error - 401):**
```json
{
  "success": false,
  "error": "Invalid refresh token",
  "message": "The provided refresh token is invalid or expired"
}
```

---

### 4. Get Operational Guy Profile
**GET** `/login/operational_guy/profile`

Retrieves the profile of the authenticated operational guy.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response (Success - 200):**
```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "profile": {
    "id": "507f1f77bcf86cd799439011",
    "email": "user@example.com",
    "name": "John Doe",
    "isFirstLogin": false,
    "phoneNumber": "1234567890",
    "address": "123 Main Street, City, State 12345",
    "emergencyContactNumber": "0987654321",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T12:00:00Z",
    "lastLoginAt": "2024-01-01T12:00:00Z"
  }
}
```

**Response (Error - 401):**
```json
{
  "success": false,
  "error": "Authentication required",
  "message": "User not authenticated"
}
```

**Response (Error - 404):**
```json
{
  "success": false,
  "error": "User not found",
  "message": "The requested user profile was not found"
}
```

---

### 5. Update Operational Guy Profile
**PUT** `/login/operational_guy/profile`

Updates the profile of the authenticated operational guy.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Request Body:**
```json
{
  "name": "John Updated Doe",
  "phoneNumber": "1234567890",
  "address": "456 New Street, City, State 12345",
  "emergencyContactNumber": "0987654321"
}
```

**Response (Success - 200):**
```json
{
  "success": true,
  "message": "Profile updated successfully",
  "profile": {
    "id": "507f1f77bcf86cd799439011",
    "email": "user@example.com",
    "name": "John Updated Doe",
    "isFirstLogin": false,
    "phoneNumber": "1234567890",
    "address": "456 New Street, City, State 12345",
    "emergencyContactNumber": "0987654321",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T13:00:00Z",
    "lastLoginAt": "2024-01-01T12:00:00Z"
  }
}
```

**Response (Error - 400):**
```json
{
  "success": false,
  "error": "Update failed",
  "message": "No changes were made to the profile"
}
```

---

### 6. Change Operational Guy Password
**PUT** `/login/operational_guy/password`

Changes the password of the authenticated operational guy.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Request Body:**
```json
{
  "currentPassword": "oldpassword123",
  "newPassword": "newpassword123"
}
```

**Response (Success - 200):**
```json
{
  "success": true,
  "message": "Password changed successfully"
}
```

**Response (Error - 400):**
```json
{
  "success": false,
  "error": "Invalid current password",
  "message": "The current password is incorrect"
}
```

---

## Error Codes

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request (Validation errors, business logic errors) |
| 401 | Unauthorized (Authentication required, invalid token) |
| 404 | Not Found |
| 500 | Internal Server Error |

## Validation Rules

### Email
- Required
- Must be a valid email format

### Password
- Required
- Minimum 6 characters

### Name
- Required
- Minimum 2 characters

### Phone Number
- Required
- Minimum 10 characters

### Address
- Required
- Minimum 10 characters

### Emergency Contact Number
- Required
- Minimum 10 characters

## Security Features

1. **Password Hashing**: All passwords are hashed using bcrypt
2. **JWT Tokens**: Secure token-based authentication
3. **Token Expiration**: Access tokens expire after 1 hour
4. **Refresh Tokens**: Long-lived refresh tokens for token renewal
5. **Input Validation**: Comprehensive validation on all inputs
6. **SQL Injection Protection**: Using parameterized queries
7. **CORS**: Cross-origin resource sharing enabled
8. **Security Headers**: XSS protection, content type options, etc.

## Environment Variables

Make sure to set these environment variables:

```env
MONGO_URI=mongodb://localhost:27017/espaze
JWT_SECRET=your-secret-key-here
```

## Testing the APIs

You can test these APIs using tools like:
- Postman
- cURL
- Insomnia
- Thunder Client (VS Code extension)

### Example cURL Commands

**Login:**
```bash
curl -X POST http://localhost:8080/login/operational_guy \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

**Register:**
```bash
curl -X POST http://localhost:8080/login/operational_guy/register \
  -H "Content-Type: application/json" \
  -d '{"email":"newuser@example.com","password":"password123","name":"John Doe","phoneNumber":"1234567890","address":"123 Main St","emergencyContactNumber":"0987654321"}'
```

**Get Profile (with token):**
```bash
curl -X GET http://localhost:8080/login/operational_guy/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Store Management APIs

### 1. Get All Stores
**Endpoint:** `GET /store/`  
**Description:** Retrieve all stores with pagination and search functionality  
**Query Parameters:**
- `warehouse_id` (required): ID of the warehouse
- `limit` (optional): Number of records per page (default: 10)
- `offset` (optional): Number of records to skip (default: 0)
- `search` (optional): Search term for store name

**Example Request:**
```
GET /store/?warehouse_id=warehouse123&limit=10&offset=0&search=store
```

**Example Response:**
```json
{
  "success": true,
  "message": "Stores retrieved successfully",
  "stores": [
    {
      "store_id": "store123",
      "seller_id": "seller456",
      "warehouse_id": "warehouse123",
      "store_name": "My Store",
      "store_address": "123 Main St, City, State",
      "store_contact": "+1234567890",
      "number_of_racks": 50,
      "occupied_racks": 25,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 1,
  "limit": 10,
  "offset": 0
}
```

### 2. Get Store by ID
**Endpoint:** `GET /store/:id`  
**Description:** Retrieve a specific store by its ID  
**Path Parameters:**
- `id`: Store ID

**Example Request:**
```
GET /store/store123
```

**Example Response:**
```json
{
  "success": true,
  "message": "Store retrieved successfully",
  "store": {
    "store_id": "store123",
    "seller_id": "seller456",
    "warehouse_id": "warehouse123",
    "store_name": "My Store",
    "store_address": "123 Main St, City, State",
    "store_contact": "+1234567890",
    "number_of_racks": 50,
    "occupied_racks": 25,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 3. Create Store
**Endpoint:** `POST /store/`  
**Description:** Create a new store  
**Request Body:**
```json
{
  "seller_id": "seller456",
  "warehouse_id": "warehouse123",
  "store_name": "New Store",
  "store_address": "456 Oak St, City, State",
  "store_contact": "+1234567890",
  "number_of_racks": 30
}
```

**Example Response:**
```json
{
  "success": true,
  "message": "Store created successfully",
  "store": {
    "store_id": "store789",
    "seller_id": "seller456",
    "warehouse_id": "warehouse123",
    "store_name": "New Store",
    "store_address": "456 Oak St, City, State",
    "store_contact": "+1234567890",
    "number_of_racks": 30,
    "occupied_racks": 0,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 4. Update Store
**Endpoint:** `PUT /store/:id`  
**Description:** Update an existing store  
**Path Parameters:**
- `id`: Store ID

**Request Body:**
```json
{
  "store_name": "Updated Store Name",
  "store_address": "789 Pine St, City, State",
  "store_contact": "+1234567890",
  "number_of_racks": 40,
  "occupied_racks": 20
}
```

**Example Response:**
```json
{
  "success": true,
  "message": "Store updated successfully",
  "store": {
    "store_id": "store123",
    "seller_id": "seller456",
    "warehouse_id": "warehouse123",
    "store_name": "Updated Store Name",
    "store_address": "789 Pine St, City, State",
    "store_contact": "+1234567890",
    "number_of_racks": 40,
    "occupied_racks": 20,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-02T00:00:00Z"
  }
}
```

### 5. Delete Store
**Endpoint:** `DELETE /store/:id`  
**Description:** Delete a store  
**Path Parameters:**
- `id`: Store ID

**Example Response:**
```json
{
  "success": true,
  "message": "Store deleted successfully"
}
```

### 6. Get Store by Seller ID
**Endpoint:** `GET /store/seller/:seller_id`  
**Description:** Retrieve a store by seller ID  
**Path Parameters:**
- `seller_id`: Seller ID

**Example Request:**
```
GET /store/seller/seller456
```

**Example Response:**
```json
{
  "success": true,
  "message": "Store retrieved successfully",
  "store": {
    "store_id": "store123",
    "seller_id": "seller456",
    "warehouse_id": "warehouse123",
    "store_name": "My Store",
    "store_address": "123 Main St, City, State",
    "store_contact": "+1234567890",
    "number_of_racks": 50,
    "occupied_racks": 25,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 7. Update Store Racks
**Endpoint:** `PATCH /store/:id/racks`  
**Description:** Update the number of occupied racks for a store  
**Path Parameters:**
- `id`: Store ID

**Request Body:**
```json
{
  "occupied_racks": 30
}
```

**Example Response:**
```json
{
  "success": true,
  "message": "Store racks updated successfully"
}
```

## Error Responses

All endpoints return consistent error responses in the following format:

```json
{
  "success": false,
  "error": "Error type",
  "message": "Detailed error message"
}
```

### Common Error Codes:
- `400 Bad Request`: Validation errors, missing required fields
- `404 Not Found`: Store not found
- `500 Internal Server Error`: Server-side errors

## Validation Rules

### Store Creation:
- `seller_id`: Required, must be a valid seller ID
- `warehouse_id`: Required, must be a valid warehouse ID
- `store_name`: Required, non-empty string
- `store_address`: Required, non-empty string
- `store_contact`: Required, non-empty string
- `number_of_racks`: Required, must be greater than 0

### Store Update:
- `store_name`: Optional, non-empty string if provided
- `store_address`: Optional, non-empty string if provided
- `store_contact`: Optional, non-empty string if provided
- `number_of_racks`: Optional, must be greater than 0 if provided
- `occupied_racks`: Optional, must be non-negative and not exceed `number_of_racks`

### Rack Update:
- `occupied_racks`: Required, must be non-negative and not exceed the store's total racks

## Business Rules

1. **One Store per Seller**: Each seller can have only one store
2. **Rack Management**: Occupied racks cannot exceed the total number of racks
3. **Timestamps**: Created and updated timestamps are automatically managed
4. **Store ID**: Automatically generated if not provided during creation
5. **Search**: Store name search is case-insensitive and supports partial matching 