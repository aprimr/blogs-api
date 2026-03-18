# 📝 Blogs API

A RESTful blog API built with **Go**, **Chi Router**, **JWT Authentication**, and **PostgreSQL (Neon)**.

---

## 🚀 Tech Stack

| Tool              | Purpose               |
| ----------------- | --------------------- |
| Chi v5            | HTTP router           |
| PostgreSQL (Neon) | Database              |
| pgx/v5            | PostgreSQL driver     |
| golang-jwt/jwt v5 | JWT authentication    |
| bcrypt            | Password hashing      |
| godotenv          | Environment variables |
| go-chi/cors       | CORS middleware       |

---

## 📁 Project Structure

```
blogs-api/
├── main.go
├── db/
│   └── connect.go
├── handlers/
│   ├── auth.go
│   └── blog.go
├── middlewares/
│   └── auth.go
├── models/
│   ├── user.go
│   ├── blog.go
│   └── response.go
├── repository/
│   ├── user.go
│   └── blog.go
├── utils/
│   ├── jwt.go
│   ├── logger.go
│   └── response.go
└── validation/
    ├── email.go
    ├── password.go
    └── string.go
```

---

## ⚙️ Setup & Installation

**1. Clone the repository**

```bash
git clone https://github.com/aprimr/blogs-api.git
cd blogs-api
```

**2. Install dependencies**

```bash
go mod tidy
```

**3. Create `.env` file**

```env
ENVIRONMENT=development
PORT=8080
DATABASE_URL=your_neon_database_url
JWT_SECRET=your_jwt_secret
```

**4. Set up the database**

Run these SQL statements in your Neon console:

```sql
CREATE TABLE users (
    uid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    is_verified BOOL DEFAULT TRUE,
    last_login TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE blogs (
    blogid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    uid UUID REFERENCES users(uid) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    content TEXT NOT NULL,
    is_deleted BOOL DEFAULT FALSE,
    is_private BOOL DEFAULT FALSE,
    updated_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW()
);
```

**5. Run the server**

```bash
go run main.go
```

---

## 🔐 Authentication

This API uses **JWT Bearer Token** authentication.

After logging in include the token in the `Authorization` header on every protected request:

```
Authorization: Bearer <your_jwt_token>
```

Tokens expire after **24 hours**.

---

## 📡 API Endpoints

### Auth — Public

| Method | Endpoint           | Description                 |
| ------ | ------------------ | --------------------------- |
| POST   | `/api/v1/register` | Register a new user         |
| POST   | `/api/v1/login`    | Login and receive JWT token |

---

### Blogs — Public

| Method | Endpoint                | Description                      |
| ------ | ----------------------- | -------------------------------- |
| GET    | `/api/v1/blogs`         | Get all public blogs (paginated) |
| GET    | `/api/v1/blog/{blogid}` | Get a single public blog by ID   |

---

### Blogs — Protected 🔒

| Method | Endpoint                | Description                |
| ------ | ----------------------- | -------------------------- |
| POST   | `/api/v1/blog`          | Create a new blog post     |
| PUT    | `/api/v1/blog/{blogid}` | Update your blog post      |
| DELETE | `/api/v1/blog/{blogid}` | Soft delete your blog post |

---

## 📋 Request & Response Examples

### POST `/api/v1/register`

**Request:**

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "Password123!"
}
```

**Response `201`:**

```json
{
  "success": true,
  "message": "User registration successful"
}
```

---

### POST `/api/v1/login`

**Request:**

```json
{
  "email": "john@example.com",
  "password": "Password123!"
}
```

**Response `200`:**

```json
{
  "success": true,
  "message": "User login successful",
  "data": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

### GET `/api/v1/blogs?page=1&limit=10`

**Response `200`:**

```json
{
  "success": true,
  "message": "Blogs fetch successful",
  "data": {
    "data": [...],
    "page": 1,
    "limit": 10,
    "total_count": 50,
    "total_pages": 5
  }
}
```

> Defaults to `page=1` and `limit=10` if not provided.

---

### POST `/api/v1/blog` 🔒

**Request:**

```json
{
  "title": "My First Blog Post",
  "description": "A short description of at least 30 characters long",
  "content": "The full content of the blog post goes here and must be at least 60 characters",
  "is_private": false
}
```

**Response `201`:**

```json
{
  "success": true,
  "message": "Blog created successfully",
  "data": {
    "blogid": "uuid",
    "uid": "uuid",
    "title": "My First Blog Post",
    "description": "A short description...",
    "content": "The full content...",
    "is_deleted": false,
    "is_private": false,
    "updated_at": "2024-01-01T00:00:00Z",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## ✅ Validation Rules

### Password

- Minimum 8 characters
- At least one letter
- At least one number
- At least one special character

### Blog Post

| Field       | Rule                  |
| ----------- | --------------------- |
| Title       | Minimum 12 characters |
| Description | Minimum 30 characters |
| Content     | Minimum 60 characters |

---

## 🌍 Environment Variables

| Variable       | Description                       | Example                       |
| -------------- | --------------------------------- | ----------------------------- |
| `ENVIRONMENT`  | App environment                   | `development` or `production` |
| `PORT`         | Server port                       | `8080`                        |
| `DATABASE_URL` | Neon PostgreSQL connection string | `postgresql://...`            |
| `JWT_SECRET`   | Secret key for signing JWT tokens | `your-secret-key`             |

---

## 🛡️ Security Features

- Passwords hashed with **bcrypt**
- JWT tokens signed with **HS256**
- CORS enabled — allows requests from any origin (`*`)
- Soft delete — blogs are never permanently deleted
- Users can only update/delete their **own** blogs
- Sensitive data never exposed in responses (`json:"-"`)
- Environment-based error logging — internal errors hidden in production
