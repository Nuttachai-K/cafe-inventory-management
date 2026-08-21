[English](README.md) | [日本語](README.ja.md)

# Cafe Inventory Management API

A production-style RESTful backend API built with **Go** (`net/http`) and **PostgreSQL** for managing a **multi-branch cafe chain**.

The system is designed for companies that operate multiple cafe branches under a single brand. Headquarters can centrally manage cafes, products, users, and inventory, while branch staff update stock levels in real time. The project also provides location-based search features to help customers find nearby cafes and available products.

---

# Business Scenario

Imagine a company called **ABC Coffee** that operates multiple cafe branches across Tokyo.

As the business grows, several operational challenges arise:

- Each branch has different inventory levels.
- Headquarters needs visibility into inventory across all branches.
- Product catalogs should remain consistent across every branch.
- Managers need an audit trail showing who updated inventory.
- Customers want to find the nearest branch that has a specific product in stock.

This API addresses those challenges by providing centralized management, role-based access control, inventory tracking, and location-based search capabilities.

---

# Features

## Authentication

- JWT Authentication
- Role-based Authorization (Admin / Staff)

## Cafe Management

- Create cafes
- Update cafe information
- Delete cafes
- Retrieve cafe list and details

## Product Management

- Create products (linked to a cafe and a category)
- Update product information
- Delete products
- Retrieve product catalog, including category name

## Category Management

- Create categories
- Update category names
- Delete categories
- Retrieve category list and details

## Inventory Management

- View inventory levels
- Update inventory quantities
- Record inventory movement history
- Track which user performed each inventory operation

## User Management

- Register users
- Update user information
- Delete users
- Manage user roles

## Location Services

- Store cafe coordinates (latitude/longitude)
- Search cafes by nearest station
- Distance-based search from a customer's coordinates (Haversine formula), with optional radius filter
- Limit result count

## Future Features

- Product availability search (find nearby cafes with a specific product in stock)
- Google Maps API integration
- Branch performance dashboard
- Sales analytics

---

# Technology Stack

| Category | Technology |
|----------|------------|
| Language | Go |
| HTTP Server | net/http |
| Database | PostgreSQL |
| Database Driver | pgx |
| Authentication | JWT |
| Password Hashing | bcrypt |
| Container | Docker |
| Architecture | Layered Architecture |
| API Style | RESTful API |
| Infrastructure | AWS ECS |
| Version Control | Git & GitHub |

---

# Project Structure

```text
.github/
└── workflows/
    └── ci.yml

cmd/
└── server/
    └── main.go

internal/
├── database/
├── handler/
├── middleware/
├── model/
├── repository/
├── router/
├── service/
└── utils/

migrations/

Dockerfile
docker-compose.yml
go.mod
go.sum
README.md
```

---

# System Architecture

```text
                Client
                   │
                   ▼
              HTTP Request
                   │
                   ▼
               Router
                   │
                   ▼
             Middleware
        (JWT Authentication,
           Request Logging)
                   │
                   ▼
               Handler
                   │
                   ▼
               Service
                   │
                   ▼
             Repository
                   │
                   ▼
             PostgreSQL
```

---

# Database Design

## Tables

- cafes
- products
- inventory
- inventory_logs
- users

### Entity Relationship

```text
cafes
   │
   │ 1:N
   ▼
products
   │
   │ 1:1
   ▼
inventory
   │
   │ 1:N
   ▼
inventory_logs
        ▲
        │
        │ N:1
      users
```

---

# API Overview

## Server

| Method | Endpoint | Description |
|---------|----------|--------------|
| GET | /health | Liveness check |
| GET | /readyz | Readiness check — pings the database, returns 503 if unreachable |

---

## Authentication

| Method | Endpoint |
|---------|----------|
| POST | /api/v1/auth/login |

---

## Cafes

| Method | Endpoint |
|---------|----------|
| GET | /api/v1/cafes |
| GET | /api/v1/cafes/{id} |
| POST | /api/v1/cafes |
| PATCH | /api/v1/cafes/{id} |
| DELETE | /api/v1/cafes/{id} |

---

## Categories

| Method | Endpoint |
|---------|----------|
| GET | /api/v1/categories |
| GET | /api/v1/categories/{id} |
| POST | /api/v1/categories |
| PATCH | /api/v1/categories/{id} |
| DELETE | /api/v1/categories/{id} |

---


## Products

| Method | Endpoint |
|---------|----------|
| GET | /api/v1/products |
| GET | /api/v1/products/{id} |
| POST | /api/v1/products |
| PATCH | /api/v1/products/{id} |
| DELETE | /api/v1/products/{id} |

---

## Inventory

| Method | Endpoint |
|---------|----------|
| GET | /api/v1/inventory |
| GET | /api/v1/inventory/{id} |
| PATCH | /api/v1/inventory/{id} |
| GET | /api/v1/inventory/logs |

---

## Users

| Method | Endpoint |
|---------|----------|
| GET | /api/v1/users |
| GET | /api/v1/users/{id} |
| POST | /api/v1/users |
| PATCH | /api/v1/users/{id} |
| DELETE | /api/v1/users/{id} |

---

# Authentication

The API uses JWT (JSON Web Token) authentication.

After a successful login, clients must include the access token in the request header.

```http
Authorization: Bearer <JWT Token>
```

---

# How to Use

## Prerequisites

- Docker & Docker Compose
- Go 1.26+ (only needed if running the API outside a container)

## 1. Configure environment

Create a `.env` file in the project root:

```env
# Used by docker-compose to provision the local Postgres container only —
# the app itself connects using DATABASE_URL below, not these directly.
POSTGRES_USER="cafe"
POSTGRES_PASSWORD="cafe"
POSTGRES_DB="cafe_inventory"
POSTGRES_PORT=5432

DB_MAX_CONNS=10
DATABASE_URL="postgres://cafe:cafe@localhost:5432/cafe_inventory?sslmode=disable"
SERVER_ADDRESS=":8080"
JWT_SECRET=2445f32fe37472df7bd5e626929c471a4b6318e87f1d1e61cd33d7ceba31bb21

# Optional. Only needed when Swagger UI is served behind a reverse proxy or
# load balancer (e.g. AWS ALB) whose public host/scheme differs from the
# container's own address. SWAGGER_SCHEME defaults to "https" if unset.
SWAGGER_HOST=your-domain.example.com
SWAGGER_SCHEME=https
```

## 2. Start everything with Docker Compose

```bash
docker-compose up -d
```

This starts PostgreSQL, runs all migrations (including a seeded admin user, `admin@cafe.local`), builds the API image, and starts the API server.

The API is now available at `http://localhost:8080`.

> **Note:** with this flow, `DATABASE_URL` in your `.env` must point at the Postgres **service name**, not `localhost` — e.g. `postgres://cafe:cafe@postgres:5432/cafe_inventory?sslmode=disable` — since the app runs inside the same Docker network as Postgres.

### Alternative: run the API outside Docker

If you'd rather run the Go binary directly on the host (e.g. for local development), stop the compose-managed app container first so it doesn't also bind port 8080:

```bash
docker-compose stop app
go run cmd/server/main.go
```

In this case, switch `DATABASE_URL`'s host back to `localhost` (Postgres' port is published to the host by docker-compose), e.g. `postgres://cafe:cafe@localhost:5432/cafe_inventory?sslmode=disable`.

## 4. Walkthrough: login → create data → adjust stock → view history

**Login as the seeded admin:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@cafe.local","password":"admin"}'
```

```json
{ "token": "<jwt>" }
```

Save the token — every write below requires `Authorization: Bearer <jwt>`.

**Create a category:**

```bash
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Authorization: Bearer <jwt>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Coffee"}'
```

```json
{ "id": 1, "message": "Category created successfully" }
```

**Create a cafe:**

```bash
curl -X POST http://localhost:8080/api/v1/cafes \
  -H "Authorization: Bearer <jwt>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "ABC Coffee Shinjuku",
    "address": "1-1-1 Shinjuku, Tokyo",
    "latitude": 35.6895,
    "longitude": 139.6917,
    "nearest_station": "Shinjuku",
    "opening_time": "07:00",
    "closing_time": "22:00"
  }'
```

```json
{ "id": 1, "message": "Cafe created successfully" }
```

**Create a product** (this also auto-creates its inventory row with `stock_quantity: 0`):

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer <jwt>" \
  -H "Content-Type: application/json" \
  -d '{
    "cafe_id": 1,
    "category_id": 1,
    "name": "Blend Coffee",
    "description": "House blend, medium roast",
    "price": "350.00"
  }'
```

```json
{ "id": 1, "message": "Product created successfully" }
```

**Add stock** (`operation` is `IN`, `OUT`, or `ADJUST`; the `{id}` here is the product's id):

```bash
curl -X PATCH http://localhost:8080/api/v1/inventory/1 \
  -H "Authorization: Bearer <jwt>" \
  -H "Content-Type: application/json" \
  -d '{"operation":"IN","change_quantity":50}'
```

```json
{ "message": "Inventory updated successfully", "stock_quantity": 50 }
```

**View current stock:**

```bash
curl http://localhost:8080/api/v1/inventory
```

**View the movement history:**

```bash
curl http://localhost:8080/api/v1/inventory/logs \
  -H "Authorization: Bearer <jwt>"
```

## 5. Browse the full API reference (Swagger UI)

```
http://localhost:8080/swagger/index.html
```

Click **Authorize** and paste `Bearer <jwt>` from step 4 to try protected endpoints directly from the browser. Alternatively, the raw spec is available at `docs/swagger.json` / `docs/swagger.yaml`.

---

# User Roles

## Admin

- Manage cafes
- Manage categories
- Manage products
- Manage users

## Staff

- Can log in and view public endpoints (cafes, products, categories)
- No distinct staff-only permissions are enforced yet — all create/update/delete endpoints currently require the Admin role

---

# Search Features

The API supports:

Currently implemented:

 - Search cafes by nearest station (partial, case-insensitive match)
 - Distance-based search from customer coordinates (Haversine formula), sorted nearest-first, with an optional radius filter
 - Limit result count

Example (by station):

```http
GET /api/v1/cafes?station=Shinjuku%20station&limit=10
```

Example (by distance):

```http
GET /api/v1/cafes?lat=35.6895&lng=139.6917&radius=5&limit=10
```

Category filtering, sorting options beyond distance, and page-based pagination are planned.

---

# Learning Objectives

This project demonstrates knowledge of:

- RESTful API Design
- Go Standard Library (`net/http`)
- Layered Architecture
- Repository Pattern
- JWT Authentication
- Role-Based Access Control (RBAC)
- PostgreSQL Database Design
- Database Migration
- Inventory Audit Logging
- Geolocation and Distance Calculation
- Filtering, Sorting, and Pagination
- Production-style Backend Development
- Containerization
- Unit Testing
- Integration Testing
- Swagger / OpenAPI Documentation
- GitHub Actions CI/CD

---

# Future Improvements

- Google Maps API integration
- Refresh Token Authentication
- Product Image Upload
- Redis Caching
- Headquarters Dashboard

---

# License

This project is intended for educational purposes and as a backend engineering portfolio project.