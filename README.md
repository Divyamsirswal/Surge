# Surge: High-Throughput Flash Sale Engine

An enterprise-grade, distributed backend architecture designed to handle massive traffic spikes (like Flash Sales or Ticket Bookings) without database locks, race conditions, or server crashes.

## The Core Problem Solved
Traditional synchronous architectures (API -> Database) fail under heavy concurrent load, leading to:
1. **Database Race Conditions:** Resulting in overselling (Negative Stock).
2. **Connection Pool Exhaustion:** Crashing the application server.
3. **High Latency:** Users facing frozen loading screens.

**Surge** solves this by decoupling the components into an **Event-Driven, Asynchronous Pipeline**.

## Architecture Flow

1. **The Gatekeeper (Redis):** All incoming `POST /buy` requests hit Redis first. Using Atomic Decrements (`DECR`), it mathematically guarantees zero overselling without using expensive database Mutex locks.
2. **The Shock Absorber (RabbitMQ & Go Channels):** Validated orders are instantly dropped into an in-memory Go Channel and published to RabbitMQ. The API responds to the user in `< 3ms`.
3. **The Sweeper (PostgreSQL):** A highly reliable background Go Worker consumes the RabbitMQ queue and persists the data into Postgres with `Auto-Ack: false`, ensuring **Zero Data Loss** even if the database temporarily goes down.

## 🛠️ Tech Stack
* **Language:** Golang (Go 1.22+)
* **In-Memory Cache:** Redis
* **Message Broker:** RabbitMQ
* **Database:** PostgreSQL
* **Infrastructure:** Docker & Docker-Compose

##  The "Moment of Truth" (Load Test Results)
I built a custom Worker-Pool based Load Tester in Go to simulate a massive cyber-attack/flash sale. 

**Test Conditions:**
* Initial Stock: 1,000 units
* Concurrent Requests: 10,000

**Results:**
```text
 Time Taken: ~3.1 seconds
 Successful Orders (Got Ticket): 1000 (Exact match, Zero Overselling)
 Sold Out Rejections: 9000
 Server Errors/Drops: 0