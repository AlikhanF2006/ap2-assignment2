# Assignment 4 — Performance Optimization & External Integrations

## Student Information

* **Name:** Alikhan Faizrakhman
* **Group:** SE-2408
* **Course:** Advanced Programming 2
* **Assignment:** Assignment 4 — Performance Optimization & External Integrations
## GitHub Repository

Repository link: https://github.com/AlikhanF2006/ap2-assignment2.git

---

## Previous Work: Assignment 3

This project is built on top of Assignment 3.

In Assignment 3, the system used RabbitMQ for event-driven communication between Payment Service and Notification Service. The Payment Service published payment events to the `payment.completed` queue, and the Notification Service consumed those events asynchronously.

Assignment 4 improves this system by adding Redis caching, cache invalidation, reliable background jobs, provider abstraction, retry logic, exponential backoff, and Redis-based idempotency.

---

## Project Overview

Assignment 4 extends the existing microservice system from Assignment 3 by adding performance optimization and reliable background job processing.

The system now includes Redis caching in the Order Service and a reliable Notification Service worker with external provider abstraction, retries, exponential backoff, and Redis-based idempotency.

Main Assignment 4 improvements:

* Redis cache-aside pattern in Order Service
* Cache TTL for order details
* Cache invalidation after order status updates
* Redis added to Docker Compose infrastructure
* Notification Service transformed into a reliable background worker
* External notification provider abstraction using Adapter Pattern
* Simulated provider with latency and random failures
* Retry policy with exponential backoff
* Redis-based idempotency for notification jobs

---

## Updated Services and Infrastructure

The system consists of:

* **Order Service**
* **Payment Service**
* **Notification Service**
* **PostgreSQL**
* **RabbitMQ**
* **Redis**
* **Docker Compose**

Redis is used for two purposes:

1. Caching order details in the Order Service
2. Storing notification idempotency records in the Notification Service

---

## Assignment 4 Architecture Diagram

![Assignment 4 Architecture Diagram](images/assignment4-architecture-diagram.png)
---

## Redis Cache-Aside Strategy

The Order Service implements the cache-aside pattern for order detail requests.

When a client sends:

```http
GET /orders/:id
```

the service first checks Redis using the key:

```text
order:{order_id}
```

If the order exists in Redis, the service returns the cached value immediately. This is called a cache hit.

If the order does not exist in Redis, the service loads the order from PostgreSQL, returns it to the client, and stores it in Redis. This is called a cache miss.

The cache TTL is configured as:

```env
CACHE_TTL_SECONDS=300
```

This means cached order data expires after 5 minutes.

Example logs:

```text
CACHE MISS order_id=e566cd81-1f9f-454e-b7d4-34b1e239f343
CACHE SET order_id=e566cd81-1f9f-454e-b7d4-34b1e239f343 ttl=5m0s
CACHE HIT order_id=e566cd81-1f9f-454e-b7d4-34b1e239f343
```

---

## Cache Invalidation Strategy

To prevent stale data, the Order Service invalidates the Redis cache whenever the order status changes.

For example, after payment processing, the order status is updated in PostgreSQL. Immediately after the database update, the Redis key is deleted:

```text
order:{order_id}
```

Example log:

```text
CACHE INVALIDATED order_id=e566cd81-1f9f-454e-b7d4-34b1e239f343
```

This guarantees that the next `GET /orders/:id` request will load fresh data from PostgreSQL and store the updated order in Redis again.

---

## Notification Background Worker

The Notification Service now works as a background worker.

The Payment Service publishes a `payment.completed` event to RabbitMQ after a payment is created. The Notification Service consumes this event asynchronously and processes the notification in the background.

This design prevents slow external notification providers from blocking the main API request path.

Example worker logs:

```text
[Worker] Notification service is waiting for payment events...
[Worker] Job received event_id=25ae540a-8b1c-43a5-ad1c-10aff004ab20 payment_id=e033519a-d976-4e6b-b1a3-79f8d0949d37 order_id=e566cd81-1f9f-454e-b7d4-34b1e239f343
```

---

## External Provider Adapter Pattern

The Notification Service uses the Adapter Pattern to decouple notification logic from a specific external provider.

The business logic depends on an interface:

```go
type EmailSender interface {
    Send(ctx context.Context, to string, subject string, body string) error
}
```

The current implementation uses a simulated provider.

The provider mode is configured using environment variables:

```env
PROVIDER_MODE=SIMULATED
PROVIDER_LATENCY_MS=1000
PROVIDER_FAILURE_RATE=30
```

The simulated provider imitates real external API behavior by adding artificial network latency and random failures.

Example provider logs:

```text
[Provider] Simulating external email provider latency=1000ms
[Provider] Email sent successfully to=user@example.com subject=Payment Authorized for Order e566cd81-1f9f-454e-b7d4-34b1e239f343
```

---

## Retry Policy and Exponential Backoff

If the simulated external provider fails, the worker retries the job.

Retry configuration:

```env
MAX_RETRIES=3
BACKOFF_BASE_SECONDS=2
```

The retry delay increases after each failed attempt:

```text
Attempt 1 failed -> wait 2 seconds
Attempt 2 failed -> wait 4 seconds
Attempt 3 failed -> final attempt
```

Example logs:

```text
[Worker] Sending notification payment_id=e033519a-d976-4e6b-b1a3-79f8d0949d37 attempt=1/3
[Worker] Provider failed payment_id=e033519a-d976-4e6b-b1a3-79f8d0949d37 error=simulated external provider failure retry_after=2s
[Worker] Sending notification payment_id=e033519a-d976-4e6b-b1a3-79f8d0949d37 attempt=2/3
[Provider] Email sent successfully to=user@example.com
```

This makes the worker resilient to temporary external provider failures.

---

## Redis Idempotency Strategy

The Notification Service uses Redis to prevent duplicate notifications.

Before sending a notification, the worker checks Redis using the key:

```text
notification:payment:{payment_id}
```

If the key already exists, the notification is skipped.

After a successful notification, the worker stores:

```text
notification:payment:{payment_id} = processed
```

with a TTL of 24 hours.

Example Redis check:

```text
GET notification:payment:e033519a-d976-4e6b-b1a3-79f8d0949d37
"processed"

TTL notification:payment:e033519a-d976-4e6b-b1a3-79f8d0949d37
86142
```

This prevents duplicate emails when the same payment event is retried or delivered again.

---

## Environment Configuration

### Order Service

```env
PORT=8081
GRPC_PORT=50052

DB_URL=postgres://postgres:postgres@postgres:5432/ap2db?sslmode=disable
PAYMENT_GRPC_ADDR=payment-service:50051

REDIS_ADDR=redis:6379
REDIS_PASSWORD=
REDIS_DB=0
CACHE_TTL_SECONDS=300
```

### Payment Service

```env
PORT=8082
GRPC_PORT=50051

DB_URL=postgres://postgres:postgres@postgres:5432/ap2db?sslmode=disable
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
```

### Notification Service

```env
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/

REDIS_ADDR=redis:6379
REDIS_PASSWORD=
REDIS_DB=0

PROVIDER_MODE=SIMULATED
MAX_RETRIES=3
BACKOFF_BASE_SECONDS=2
PROVIDER_LATENCY_MS=1000
PROVIDER_FAILURE_RATE=30
```

---

## Docker Compose

The system is started using Docker Compose.

Run all services:

```bash
docker compose up -d --build
```

Check running containers:

```bash
docker compose ps
```

Expected containers:

```text
postgres
rabbitmq
redis
order-service
payment-service
notification-service
```

Stop all services:

```bash
docker compose down
```

---

## Database Tables

If PostgreSQL starts with an empty database, create the required tables:

```sql
CREATE TABLE IF NOT EXISTS orders (
    id TEXT PRIMARY KEY,
    customer_id TEXT NOT NULL,
    item_name TEXT NOT NULL,
    amount BIGINT NOT NULL CHECK (amount > 0),
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS payments (
    id TEXT PRIMARY KEY,
    order_id TEXT NOT NULL,
    transaction_id TEXT NOT NULL UNIQUE,
    amount BIGINT NOT NULL CHECK (amount > 0),
    status TEXT NOT NULL
);
```

---

## How to Test Assignment 4

### 1. Start all services

```bash
docker compose up -d --build
```

### 2. Create an order

Method:

```text
POST
```

URL:

```text
http://localhost:8081/orders
```

Body:

```json
{
  "customer_id": "customer-1",
  "item_name": "Laptop",
  "amount": 1000
}
```

Expected response:

```json
{
  "id": "...",
  "customer_id": "customer-1",
  "item_name": "Laptop",
  "amount": 1000,
  "status": "Paid",
  "created_at": "..."
}
```

### 3. Test Redis cache

Send the same request twice:

```text
GET http://localhost:8081/orders/{order_id}
```

Expected Order Service logs:

```text
CACHE MISS
CACHE SET
CACHE HIT
```

### 4. Test notification worker

Open Notification Service logs:

```bash
docker compose logs -f notification-service
```

Create another order and check the worker logs.

Expected logs:

```text
[Worker] Job received
[Worker] Sending notification payment_id=... attempt=1/3
[Provider] Simulating external email provider latency=1000ms
[Worker] Provider failed ... retry_after=2s
[Worker] Sending notification payment_id=... attempt=2/3
[Provider] Email sent successfully
[Idempotency] Notification marked as processed
[Worker] Notification completed
```

### 5. Check Redis keys

Open Redis CLI:

```bash
docker exec -it redis redis-cli
```

Check all keys:

```redis
KEYS *
```

Expected keys:

```text
order:{order_id}
notification:payment:{payment_id}
```

Check order cache:

```redis
GET order:{order_id}
TTL order:{order_id}
```

Check notification idempotency:

```redis
GET notification:payment:{payment_id}
TTL notification:payment:{payment_id}
```

Expected result:

```text
"processed"
```

---

## Evidence / Screenshots

### 1. Docker Compose Running with Redis

![Docker Compose Assignment 4](images/docker-compose-assignment4.png)

This screenshot shows all required containers running: Order Service, Payment Service, Notification Service, PostgreSQL, RabbitMQ, and Redis.

### 2. Successful Order Creation in Postman

![Postman Create Order](images/postman-create-order.png)

This screenshot shows a successful `POST /orders` request with `201 Created`.

### 3. Successful Get Order Request in Postman

![Postman Get Order](images/postman-get-order.png)

This screenshot shows a successful `GET /orders/{order_id}` request with `200 OK`.

### 4. Redis Cache Miss, Set, and Hit Logs

![Order Cache Logs](images/order-cache-miss-hit.png)

This screenshot proves the Redis cache-aside implementation in the Order Service. The first request produces `CACHE MISS` and `CACHE SET`, while the next request produces `CACHE HIT`.

### 5. Notification Worker Retry and Idempotency Logs

![Notification Worker Retry](images/notification-worker-retry.png)

This screenshot proves that the Notification Service works as a background worker, retries failed provider calls, uses exponential backoff, and marks the job as processed.

### 6. Redis Cache and Idempotency Proof

![Redis Proof](images/redis-cache-idempotency-proof.png)

This screenshot shows both Redis keys: `order:{order_id}` for cache and `notification:payment:{payment_id}` for idempotency. It also shows TTL values for both keys.

### 7. RabbitMQ Management UI

![RabbitMQ UI](images/rabbitmq-ui.png)

This screenshot shows the RabbitMQ Management UI and the message queue used for payment events.

---

## Deliverables Included

* Source code for Order, Payment, and Notification services
* Redis cache-aside implementation in Order Service
* Cache invalidation after order status updates
* Redis-based idempotency in Notification Service
* Notification background worker with RabbitMQ
* Simulated external provider using Adapter Pattern
* Retry policy with exponential backoff
* Docker Compose file with Redis, PostgreSQL, RabbitMQ, and all services
* Architecture diagram
* README documentation
* Evidence screenshots

---

## Conclusion

Assignment 4 improves the previous event-driven microservice system by adding production-ready performance and reliability patterns.

The Order Service now uses Redis cache-aside to reduce database load and improve response time for repeated order lookups.

The Notification Service now works as a reliable background worker. It consumes RabbitMQ events, sends notifications through a provider interface, retries temporary provider failures using exponential backoff, and stores processed payment IDs in Redis to prevent duplicate notifications.

This implementation demonstrates caching, background jobs, external integration abstraction, retries, idempotency, and Docker-based infrastructure orchestration.