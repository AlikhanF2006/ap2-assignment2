# Order & Payment Microservices (Go, Clean Architecture)

## Overview

This project implements a two-service microservice system using Go and Gin framework.  
The system consists of:

- Order Service – manages customer orders
- Payment Service – processes payments

The services communicate via REST API and follow Clean Architecture principles.

---

## Architecture

Each service follows Clean Architecture:

- domain – business entities
- usecase – business logic
- repository – database interaction
- transport/http – HTTP handlers
- app – dependency injection (composition root)

Services are fully separated:
- no shared code
- no shared database
- communication only via HTTP

---

## Services

### 1. Order Service

Responsible for:
- creating orders
- updating order status
- calling payment service

Runs on port: 8081

Endpoints:

POST /orders  
Creates a new order and triggers payment

GET /orders/{id}  
Returns order details

PATCH /orders/{id}/cancel  
Cancels order (only if status = Pending)

---

### 2. Payment Service

Responsible for:
- processing payments
- validating payment limits

Runs on port: 8082

Endpoints:

POST /payments  
Processes payment

GET /payments/{order_id}  
Returns payment status

---

## Business Rules

- Amount must be greater than 0
- Money is stored as int64 (no float)
- Orders:
    - Pending → Paid / Failed
    - Only Pending orders can be cancelled
- Payments:
    - If amount > 100000 → Declined
    - Otherwise → Authorized

---

## Service Communication

Order Service calls Payment Service via HTTP:

- POST /payments
- Uses custom http.Client with timeout (2 seconds)

Failure handling:
- If Payment Service is unavailable:
    - request times out
    - Order Service returns 503
    - order status becomes "Failed"

---

## Database

Each service has its own database:

- order-service → orders table
- payment-service → payments table

No shared database is used.

---

## How to Run

### 1. Start PostgreSQL

Make sure PostgreSQL is running.

---

### 2. Set environment variables (.env)

Order Service:
DB_URL=postgres://user:password@localhost:5432/order_db?sslmode=disable
PAYMENT_SERVICE_URL=http://localhost:8082
PORT=8081

Payment Service:
DB_URL=postgres://user:password@localhost:5432/payment_db?sslmode=disable
PORT=8082

---

### 3. Run services

Order Service:
go run cmd/main.go

Payment Service:
go run cmd/main.go

---

## API Examples

### Create Order

POST /orders

Body:
{
"customer_id": "123",
"item_name": "Laptop",
"amount": 50000
}

---

### Get Order

GET /orders/{id}

---

### Create Payment (internal)

POST /payments

Body:
{
"order_id": "123",
"amount": 50000
}

---

## Failure Scenario

If Payment Service is down:

- Order Service returns 503 Service Unavailable
- Order status is set to "Failed"

---

## Conclusion

The system demonstrates:

- clean separation of concerns
- microservice architecture
- REST communication
- failure handling with timeouts
- independent data ownership
