# Assignment 2 — gRPC Migration & Contract-First Development

## Student Information
- **Name:** Alikhan Faizrakhman
- **Group:** SE - 2408
- **Course:** Advanced Programming 2
- **Assignment:** Assignment 2 — gRPC Migration & Contract-First Development

---

## Project Overview
This project is an updated version of Assignment 1.  
The system consists of two microservices:

- **Order Service**
- **Payment Service**

In this assignment, the communication between services was migrated from REST to **gRPC**.  
The external API for the user is still exposed through **REST** by the Order Service, while the internal communication between Order Service and Payment Service is now implemented using **gRPC**.

Additionally, the project demonstrates **server-side streaming** with the method:

`SubscribeToOrderUpdates(OrderRequest) returns (stream OrderStatusUpdate)`

This allows a client to subscribe to order status updates in real time. The assignment requires migrating inter-service communication to gRPC, keeping REST for the external API, and implementing server-side streaming for order tracking.

---

## Repository Links

### 1. Proto Repository

`https://github.com/AlikhanF2006/ap2-protos.git`

### 2. Generated Code Repository

`https://github.com/AlikhanF2006/ap2-protos-gen.git`

### 3. Main Project Repository

`https://github.com/AlikhanF2006/ap2-assignment2.git`
---

## Architecture

### Before Migration
In Assignment 1, the Order Service communicated with the Payment Service using REST.

### After Migration
In Assignment 2:

- **Order Service** still provides REST endpoints for the end user.
- **Order Service** acts as a **gRPC client** when calling the Payment Service.
- **Payment Service** acts as a **gRPC server** and implements `ProcessPayment`.
- **Order Service** also acts as a **gRPC server** for streaming order status updates.
- A client can subscribe to updates using `SubscribeToOrderUpdates`.

This design follows the assignment requirement where the Payment Service must expose a gRPC server interface and the Order Service must provide server-side streaming for order tracking.

---

## Architecture Diagram
```text
+-------------------+          gRPC           +-------------------+
|   Order Service   | ----------------------> |  Payment Service  |
|   (REST + gRPC)   |                         |    (gRPC Server)  |
+-------------------+                         +-------------------+
         |
         | REST
         v
+-------------------+
|   REST Client     |
| (Postman / User)  |
+-------------------+

         ^
         | gRPC Server-Side Streaming
         |
+-----------------------------+
| gRPC Client / Postman       |
| SubscribeToOrderUpdates()   |
+-----------------------------+
```
## Architecture Description

The system consists of two microservices:

- Order Service (REST + gRPC)
- Payment Service (gRPC Server)

The Order Service exposes REST endpoints for external clients (e.g., POST /orders).  
Internally, it communicates with the Payment Service using gRPC for processing payments.

Additionally, the Order Service implements server-side streaming via the method  
`SubscribeToOrderUpdates`, which allows clients to receive real-time order status updates.

This architecture demonstrates the migration from REST to gRPC for internal communication while preserving REST for external access.

## Evidence

### 1. REST API — Create Order
![img_2.png](images/img_2.png)

### 2. gRPC Streaming — SubscribeToOrderUpdates
![img.png](images/img.png)

### 3. Services Running (Logs)
![img_1.png](images/img_1.png)

## How to Run

### 1. Run Payment Service
```bash
cd payment-service
go run cmd/main.go

cd order-service
go run cmd/main.go

2. Run Order Service
cd order-service
go run cmd/main.go
3. Test REST API

POST http://localhost:8081/orders

Body:

{
  "customer_id": "123",
  "item_name": "iphone",
  "amount": 1000
}
4. Test gRPC Streaming

Connect to:
localhost:50052

Method:
SubscribeToOrderUpdates

Message:

{
  "order_id": "YOUR_ORDER_ID"
}

