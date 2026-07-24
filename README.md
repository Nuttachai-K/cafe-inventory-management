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

- Create products
- Update product information
- Delete products
- Retrieve product catalog

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

- Store cafe coordinates
- Calculate distance between users and cafes (Haversine Formula)
- Search nearby cafes
- Filter, sort, and paginate search results

## Future Features

- Customer-facing cafe search
- Product availability search
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
| Infrastructure | 	AWS EC2 |
| Version Control | Git & GitHub |

---

# Project Structure

```text
cmd/
└── api/
    └── main.go

internal/
├── config/
├── database/
├── handler/
├── middleware/
├── model/
├── repository/
├── router/
├── service/
└── utils/

migrations/
docs/

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

## Products

| Method | Endpoint |
|---------|----------|
| GET | /api/v1/cafes/{cafeId}/products |
| GET | /api/v1/products/{id} |
| POST | /api/v1/products |
| PATCH | /api/v1/products/{id} |
| DELETE | /api/v1/products/{id} |

---

## Inventory

| Method | Endpoint |
|---------|----------|
| GET | /api/v1/inventory/{productId} |
| PATCH | /api/v1/inventory/{productId} |
| GET | /api/v1/inventory/{productId}/logs |

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

# User Roles

## Admin

- Manage cafes
- Manage products
- Manage users
- Update inventory
- View inventory logs

## Staff

- View cafes
- View products
- View inventory
- Update inventory
- View inventory logs

---

# Search Features

The API supports:

Currently implemented:

 - Search cafes by nearest station
 - Limit result count

Example:

```http
GET /api/v1/cafes?station=Shinjuku%20station&limit=10
```
Distance-based search, category filtering, sorting, and page-based pagination are planned

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

---

# Future Improvements

- Google Maps API integration
- Refresh Token Authentication
- GitHub Actions CI/CD
- Swagger / OpenAPI Documentation
- Product Image Upload
- Redis Caching
- Headquarters Dashboard
- Branch Performance Analytics

---

# License

This project is intended for educational purposes and as a backend engineering portfolio project.