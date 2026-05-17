# Reflect Backend

Modern fashion ecommerce backend built with Golang, Gin, PostgreSQL, and Docker.

Designed with scalable modular architecture, JWT authentication, role-based authorization, secure checkout flow, and ecommerce-ready APIs.

---

# 🚀 Tech Stack

- Golang
- Gin Framework
- PostgreSQL
- GORM
- Docker
- JWT Authentication
- Cloudinary

---

# ✨ Features

## Authentication
- Register
- Login
- JWT Authentication
- Role Authorization

## User
- User Management
- Profile Management
- Avatar Upload

## Ecommerce
- Products
- Categories
- Collections
- Wishlist
- Cart System
- Address Management
- Checkout System
- Order Tracking

## Admin
- Product Management
- Category Management
- Collection Management
- Order Status Management

## Upload System
- Cloudinary Image Upload
- Product Image Upload
- User Avatar Upload

---

# 📁 Project Structure

```txt
cmd/
└── api/

internal/
├── config/
├── database/
├── middleware/
├── routes/
├── utils/
│
├── modules/
│   ├── auth/
│   ├── user/
│   ├── profile/
│   ├── product/
│   ├── category/
│   ├── collection/
│   ├── cart/
│   ├── wishlist/
│   ├── address/
│   ├── order/
│   └── upload/
```

---

# ⚙️ Installation

## Clone Repository

```bash
git clone https://github.com/yourusername/reflect-backend.git
```

```bash
cd reflect-backend
```

---

# 📦 Install Dependencies

```bash
go mod tidy
```

---

# 🔐 Setup Environment

Create `.env`

```env
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=reflect_db
DB_SSLMODE=disable

JWT_SECRET=supersecretjwtkey

CLOUDINARY_CLOUD_NAME=
CLOUDINARY_API_KEY=
CLOUDINARY_API_SECRET=
```

---

---

# ▶️ Run Backend

```bash
go run cmd/api/main.go
```

# 🌐 API Base URL

```txt
http://localhost:8080/api/v1
```

---

# 📌 Main Endpoints

## Auth

```txt
POST /auth/register
POST /auth/login
```

---

## Products

```txt
GET /products
GET /products/featured
GET /products/new-arrivals
GET /products/:id
```

---

## Cart

```txt
GET /cart
POST /cart
PUT /cart/:id
DELETE /cart/:id
```

---

## Orders

```txt
POST /checkout
GET /orders
GET /orders/:id
```

---

# 🏗️ Architecture

This project uses layered architecture:

```txt
Handler
→ Service
→ Repository
→ Database
```

## Benefits

- Scalable
- Maintainable
- Modular
- Easy to test

---

# 🔒 Security

- JWT Authentication
- Protected Routes
- Role-based Authorization
- Stock-safe Checkout
- Race-condition Safe Cart Operations

---

# 📄 License

MIT License

---

# ❤️ Author

Built with passion by RapzDev
