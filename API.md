# Reflect Backend API Documentation

Base URL: `http://localhost:8080/api/v1`

---

## Authentication

### Register

```
POST /auth/register
```

**Request:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (201):**
```json
{
  "status": "success",
  "message": "Register successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": "uuid",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user",
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-01-01T00:00:00Z"
    }
  }
}
```

### Login

```
POST /auth/login
```

**Request:**
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (200):** Same structure as Register (token + user).

**Errors:** `401` Invalid email or password.

---

## Usage

All protected endpoints require:

```
Authorization: Bearer <token>
```

Token is a JWT that expires in **7 days**. Claims: `user_id`, `email`, `role`.

---

## Standard Response Format

### Success
```json
{
  "status": "success",
  "message": "Descriptive message",
  "data": { ... }
}
```

### Error (from service layer)
```json
{
  "status": "error",
  "message": "Error description",
  "error": null
}
```

### Validation Error
```json
{
  "status": "error",
  "message": "Validation error",
  "error": "Field validation error message"
}
```

### Error Status Codes Used

| Code | Meaning |
|------|---------|
| 400 | Bad request / validation error |
| 401 | Unauthorized (missing/invalid token) |
| 403 | Forbidden (wrong role) |
| 404 | Resource not found |
| 409 | Conflict (e.g. duplicate email) |
| 500 | Internal server error |

---

## Public Endpoints (no auth required)

### Health Check
```
GET /api/v1/ping
```
```json
{ "status": "success", "message": "pong" }
```

### Products

#### List all products
```
GET /products
```

#### New arrivals
```
GET /products/new-arrivals
```

#### Featured products
```
GET /products/featured
```

#### By category
```
GET /products/category?category=hoodie
```

Category values: `hoodie`, `t-shirt`, `outerwear`, `accessories`, `shoes`.

#### By slug
```
GET /products/slug/:slug
```

#### By ID
```
GET /products/:id
```

#### Related products (up to 3, same category)
```
GET /products/:id/related
```

**Product response:**
```json
{
  "id": "uuid",
  "name": "Classic Hoodie",
  "slug": "classic-hoodie",
  "description": "A comfortable hoodie",
  "category": "hoodie",
  "collection": "summer-2026",
  "image": "https://res.cloudinary.com/...",
  "price": 299000,
  "stock": 50,
  "stock_status": "in_stock",
  "is_new_arrival": true,
  "is_featured": false,
  "created_at": "2026-01-01T00:00:00Z",
  "updated_at": "2026-01-01T00:00:00Z"
}
```

`stock_status`: `in_stock` | `low_stock` (≤5) | `out_of_stock` (≤0). Price is in **IDR (Rupiah)**.

### Categories

```
GET  /categories
GET  /categories/active
GET  /categories/slug/:slug
GET  /categories/:id
```

**Response:**
```json
{
  "id": "uuid",
  "name": "Hoodie",
  "slug": "hoodie",
  "description": "...",
  "image": "https://...",
  "is_active": true,
  "created_at": "...",
  "updated_at": "..."
}
```

### Collections

```
GET  /collections
GET  /collections/active
GET  /collections/featured
GET  /collections/slug/:slug
GET  /collections/:id
```

**Response:**
```json
{
  "id": "uuid",
  "name": "Summer 2026",
  "slug": "summer-2026",
  "description": "...",
  "image": "https://...",
  "is_featured": false,
  "is_active": true,
  "created_at": "...",
  "updated_at": "..."
}
```

---

## Protected Endpoints (JWT required)

### Profile

#### Get my profile
```
GET /me
```

#### Update my profile
```
PUT /me
```
```json
{
  "name": "New Name",
  "avatar": "https://res.cloudinary.com/..."
}
```
Both fields are optional — only provided fields are updated.

**Profile response:**
```json
{
  "id": "uuid",
  "name": "John Doe",
  "email": "john@example.com",
  "role": "user",
  "avatar": "https://...",
  "created_at": "...",
  "updated_at": "..."
}
```

### Cart

#### Get my cart
```
GET /cart
```

#### Add item
```
POST /cart
```
```json
{
  "product_id": "uuid",
  "quantity": 2
}
```
If the product is already in the cart, quantities are **added together**. Stock is checked against the total.

#### Update item quantity
```
PUT /cart/:id
```
```json
{
  "quantity": 3
}
```

#### Remove item
```
DELETE /cart/:id
```

#### Clear cart
```
DELETE /cart
```

**Cart response:**
```json
{
  "data": {
    "items": [
      {
        "id": "uuid",
        "product_id": "uuid",
        "name": "Classic Hoodie",
        "slug": "classic-hoodie",
        "image": "https://...",
        "price": 299000,
        "quantity": 2,
        "subtotal": 598000,
        "stock": 50,
        "stock_status": "in_stock",
        "created_at": "...",
        "updated_at": "..."
      }
    ],
    "total_items": 2,
    "total_price": 598000
  }
}
```

### Wishlist

#### Get my wishlist
```
GET /wishlist
```

#### Add item
```
POST /wishlist
```
```json
{
  "product_id": "uuid"
}
```

#### Remove by wishlist item ID
```
DELETE /wishlist/:id
```

#### Remove by product ID
```
DELETE /wishlist/product/:productId
```

**Wishlist response:**
```json
{
  "data": [
    {
      "id": "uuid",
      "product_id": "uuid",
      "name": "Classic Hoodie",
      "slug": "classic-hoodie",
      "image": "https://...",
      "price": 299000,
      "stock": 50,
      "stock_status": "in_stock",
      "is_new_arrival": true,
      "is_featured": false,
      "created_at": "...",
      "updated_at": "..."
    }
  ]
}
```

### Addresses

#### List my addresses
```
GET /addresses
```

#### Create address
```
POST /addresses
```
```json
{
  "recipient_name": "John Doe",
  "phone_number": "08123456789",
  "province": "Jakarta",
  "city": "Jakarta Selatan",
  "district": "Kebayoran Baru",
  "postal_code": "12120",
  "full_address": "Jl. Contoh No. 123, RT 01 RW 02",
  "label": "Home",
  "is_primary": true
}
```

#### Update address
```
PUT /addresses/:id
```
All fields optional.

#### Delete address
```
DELETE /addresses/:id
```

**Address response:**
```json
{
  "id": "uuid",
  "recipient_name": "John Doe",
  "phone_number": "08123456789",
  "province": "Jakarta",
  "city": "Jakarta Selatan",
  "district": "Kebayoran Baru",
  "postal_code": "12120",
  "full_address": "Jl. Contoh No. 123",
  "label": "Home",
  "is_primary": true,
  "created_at": "...",
  "updated_at": "..."
}
```

### Orders

#### My orders
```
GET /orders
```

#### My order by ID
```
GET /orders/:id
```

#### My order by order number
```
GET /orders/number/:orderNumber
```

### Checkout

```
POST /checkout
```
```json
{
  "address_id": "uuid"
}
```

Checkout converts all items in your cart into an order, deducts stock, and clears the cart. Shipping is fixed at **15,000 IDR**.

**Order response:**
```json
{
  "id": "uuid",
  "order_number": "BLCK1735689600123",
  "subtotal": 598000,
  "shipping": 15000,
  "total": 613000,
  "order_status": "pending",
  "payment_status": "pending",
  "shipping_status": "pending",
  "tracking_number": "TRK17356896001234",
  "items": [
    {
      "id": "uuid",
      "product_id": "uuid",
      "product_name": "Classic Hoodie",
      "product_slug": "classic-hoodie",
      "product_image": "https://...",
      "price": 299000,
      "quantity": 2,
      "subtotal": 598000
    }
  ],
  "tracking": [
    {
      "id": "uuid",
      "status": "Order Created",
      "description": "Your order has been created and waiting for payment.",
      "location": "Reflect Store",
      "created_at": "..."
    }
  ],
  "created_at": "...",
  "updated_at": "..."
}
```

**Order statuses:** `pending` | `paid` | `processing` | `shipped` | `delivered` | `cancelled`

**Payment statuses:** `pending` | `paid` | `failed`

**Shipping statuses:** `pending` | `packed` | `in_transit` | `delivered`

### Image Upload

```
POST /upload
Content-Type: multipart/form-data
```

Field: `image` (file)

**Response:**
```json
{
  "status": "success",
  "message": "Image uploaded successfully",
  "data": {
    "url": "https://res.cloudinary.com/..."
  }
}
```

Uses Cloudinary. Returns the uploaded image's secure URL.

### User Management (Admin only)

```
GET    /users              (admin)
GET    /users/:id          (any authenticated user — gets own or any user)
DELETE /users/:id          (admin)
```

**Users response:**
```json
{
  "id": "uuid",
  "name": "John Doe",
  "email": "john@example.com",
  "role": "user",
  "created_at": "...",
  "updated_at": "..."
}
```

---

## Admin Endpoints (JWT + role `admin`)

### Admin Products

```
POST   /admin/products
PUT    /admin/products/:id
DELETE /admin/products/:id
```

**Create/Update product request:**
```json
{
  "name": "Classic Hoodie",
  "slug": "classic-hoodie",
  "description": "A comfortable hoodie",
  "category": "hoodie",
  "collection": "summer-2026",
  "image": "https://res.cloudinary.com/...",
  "price": 299000,
  "stock": 50,
  "is_new_arrival": true,
  "is_featured": false
}
```

On create: `name`, `slug`, `category`, `price` are required. On update: all fields optional. Stock status is auto-calculated from stock value.

### Admin Categories

```
POST   /admin/categories
PUT    /admin/categories/:id
DELETE /admin/categories/:id
```

**Create category request:**
```json
{
  "name": "Hoodie",
  "slug": "hoodie",
  "description": "All hoodie products",
  "image": "https://...",
  "is_active": true
}
```

`name` and `slug` required on create.

### Admin Collections

```
POST   /admin/collections
PUT    /admin/collections/:id
DELETE /admin/collections/:id
```

**Create collection request:**
```json
{
  "name": "Summer 2026",
  "slug": "summer-2026",
  "description": "Summer collection",
  "image": "https://...",
  "is_featured": false,
  "is_active": true
}
```

`name` and `slug` required on create.

### Admin Orders — Update Status

```
PUT /admin/orders/:id/status
```
```json
{
  "order_status": "processing",
  "payment_status": "paid",
  "shipping_status": "packed",
  "tracking_status": "Order Packed",
  "description": "Your order has been packed and ready to ship.",
  "location": "Reflect Warehouse"
}
```

Only `order_status` is required. If `tracking_status` is provided, a new tracking entry is created. Response is the full updated order object.
