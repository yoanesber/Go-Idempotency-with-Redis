# Idempotent Transaction API

This service is responsible for processing **financial transactions** in a **secure** and **idempotent** manner using **Go (Gin framework)**, **Redis**, and **PostgreSQL**. Each transaction request must include a unique **Idempotency-Key** to ensure safe retries and prevent duplicate processing. The system uses **Redis** for fast **idempotency checks** and **PostgreSQL** to persist transaction data reliably. The system ensures that **retries do not result in duplicate operations**.

---


## ✨ Features

This application offers a **robust**, **idempotent transaction processing service** built with **Go (Gin)**, **Redis**, and **PostgreSQL**. It ensures **high reliability**, **safe retries**, and **clear observability** for financial operations or critical APIs.

Each transaction request must include an `Idempotency-Key` (UUID). The system guarantees that **a request with the same key will not be processed more than once**, even if retried due to network failures or client timeouts.

### ♻️ Idempotency Enforcement  

Each transaction request must include an `Idempotency-Key` (UUID). The service ensures the same key cannot be used to create multiple logically different transactions, preventing accidental duplicates on retries.  

✅ Mechanism:
  - The **raw request body** is hashed using **SHA-256**.
  - A Redis key is checked: `idempotency_cache:<Idempotency-Key>`.  

🛡️ Benefits:
  - Prevents **duplicate charges/payments**.  
  - Ensures **safe retries** in unstable network conditions.  
  - Supports **consistent and deterministic** behavior for clients.  

### 🗄️ Logging

Robust logging system for visibility and debugging:  

- Uses `github.com/sirupsen/logrus` for structured, leveled logging.  
- Integrates with `gopkg.in/natefinch/lumberjack.v2` for automatic log rotation based on size and age.  
- Logs are separated by level: **info**, **request**, **warn**, **error**, **fatal**, and **panic**.  


---

## 🧭 Business Process Flow

The following diagram illustrates the full flow of how a **transaction request** is handled, from initial submission to database persistence and idempotency validation. The system ensures **safe retries**, **duplicate prevention**, and **data consistency** using Redis and PostgreSQL.

```pgsql
┌──────────────────────────────────────────────┐
│            [1] Client Sends Request          │
│----------------------------------------------│
│ - POST /transactions                         │
│ - Headers:                                   │
│   - Idempotency-Key: <UUID>                  │
│ - Body: { type, amount, consumerId }         │
└──────────────────────────────────────────────┘
              │
              ▼
┌──────────────────────────────────────────────┐
│   [2] Middleware: Validate Idempotency-Key   │
│----------------------------------------------│
│ - Check format → if invalid → 400            │
│ - Hash raw body (SHA-256)                    │
│ - Query Redis for idempotency_cache:<key>    │
│   - If exists and hash matches → return      │
│     cached response                          │
│   - If exists and hash differs → 409 Conflict│
│   - If not found → proceed                   │
└──────────────────────────────────────────────┘
              │
              ▼
┌──────────────────────────────────────────────┐
│           [3] Context Injection              │
│----------------------------------------------│
│ - Inject Idempotency-Key and body hash       │
│   into request context                       │
└──────────────────────────────────────────────┘
              │
              ▼
┌──────────────────────────────────────────────┐
│     [4] Service Layer: Business Validation   │
│----------------------------------------------│
│ - Validate `consumerId` exists → 404         │
│ - Validate consumer is active → 400          │
└──────────────────────────────────────────────┘
              │
              ▼
┌──────────────────────────────────────────────┐
│ [5] Save Transaction & Idempotency Metadata  │
│----------------------------------------------│
│ - Insert into `transactions` (status=pending)│
│ - Store idempotency key, hash, and response  │
│   into both Redis and PostgreSQL             │
└──────────────────────────────────────────────┘

```
---


## 🤖 Tech Stack

This project uses a clean and efficient stack to deliver reliable and high-performance transaction processing. Below is an overview of the key tools and libraries involved:

| **Component**             | **Description**                                                                             |
|---------------------------|---------------------------------------------------------------------------------------------|
| **Language**              | Go (Golang), a statically typed, compiled language known for concurrency and efficiency     |
| **Web Framework**         | Gin, a fast and minimalist HTTP web framework for Go                                        |
| **ORM**                   | GORM, an ORM library for Go supporting SQL and migrations                                   |
| **Database**              | PostgreSQL — relational storage for transactions and idempotency metadata                   |
| **Cache/Session Store**   | Redis — used for fast idempotency key lookup and temporary response caching                 |
| **Logging**               | Logrus for structured logging, combined with Lumberjack for log rotation                    |
| **Validation**            | `go-playground/validator.v9` for input validation and data integrity enforcement            |

---

## 🧱 Architecture Overview

This project follows a **modular** and **maintainable** architecture inspired by **Clean Architecture** principles. Each domain feature (e.g., **entity**, **handler**, **repository**, **service**) is organized into self-contained modules with clear separation of concerns.

```bash
📁 go-idempotency-with-redis/
├── 📂cert/                                 # Stores self-signed TLS certificates used for local development
├── 📂cmd/                                  # Contains the application's entry point.
├── 📂config/
│   ├── 📂cache/                            # Config for Redis (host, port, TTL, etc.)
│   └── 📂database/                         # Config for PostgreSQL (DSN, pool settings, migration, etc.)
├── 📂docker/                               # Docker-related configuration for building and running services
│   ├── 📂app/                              # Contains Dockerfile to build the main Go application image
│   ├── 📂postgres/                         # Contains PostgreSQL container configuration
│   └── 📂redis/                            # Contains Redis container configuration
├── 📂internal/                             # Core domain logic and business use cases, organized by module
│   ├── 📂entity/                           # Data models/entities representing business concepts like Transaction, Consumer
│   ├── 📂handler/                          # HTTP handlers (controllers) that parse requests and return responses
│   ├── 📂repository/                       # Data access layer, communicating with DB or cache
│   └── 📂service/                          # Business logic layer orchestrating operations between handlers and repositories
├── 📂logs/                                 # Application log files (error, request, info) written and rotated using Logrus + Lumberjack
├── 📂pkg/                                  # Reusable utility and middleware packages shared across modules
│   ├── 📂contextdata/                      # Stores and retrieves contextual data like Idempotency-Key
│   ├── 📂customtype/                       # Defines custom types, enums, constants used throughout the application
│   ├── 📂diagnostics/                      # Health check endpoints, metrics, and diagnostics handlers for monitoring
│   ├── 📂logger/                           # Centralized log initialization and configuration
│   ├── 📂middleware/                       # Request processing middleware
│   │   ├── 📂headers/                      # Manages request headers like CORS, security
│   │   ├── 📂idempotency/                  # Extracts, validates, and processes Idempotency-Key
│   │   └── 📂logging/                      # Logs incoming requests
│   └── 📂util/                             # General utility functions and helpers
│       ├── 📂hash-util/                    # Functions for hashing request bodies (e.g., SHA-256)
│       ├── 📂http-util/                    # Utilities for common HTTP tasks (e.g., write JSON, status helpers)
│       ├── 📂redis-util/                   # Redis connection and command utilities
│       └── 📂validation-util/              # Common input validators (e.g., UUID, numeric range)
├── 📂routes/                               # Route definitions, groups APIs, and applies middleware per route scope
└── 📂tests/                                # Contains unit or integration tests for business logic
```

---

## 🛠️ Installation & Setup  

Follow the instructions below to get the project up and running in your local development environment. You may run it natively or via Docker depending on your preference.  

### ✅ Prerequisites

Make sure the following tools are installed on your system:

| **Tool**                                                      | **Description**                           |
|---------------------------------------------------------------|-------------------------------------------|
| [Go](https://go.dev/dl/)                                      | Go programming language (v1.20+)          |
| [Make](https://www.gnu.org/software/make/)                    | Build automation tool (`make`)            |
| [Redis](https://redis.io/)                                    | In-memory data store                      |
| [PostgreSQL](https://www.postgresql.org/)                     | Relational database system (v14+)         |
| [Docker](https://www.docker.com/)                             | Containerization platform (optional)      |

### 🔁 Clone the Project  

Clone the repository:  

```bash
git clone https://github.com/yoanesber/Go-Idempotency-with-Redis.git
cd Go-Idempotency-with-Redis
```

### ⚙️ Configure `.env` File  

Set up your **database**, **Redis**, and **JWT configuration** in `.env` file. Create a `.env` file at the project root directory:  

```properties
# Application configuration
ENV=PRODUCTION
API_VERSION=1.0
PORT=1000
IS_SSL=TRUE
SSL_KEYS=./cert/mycert.key
SSL_CERT=./cert/mycert.cer

# Database configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=appuser
DB_PASS=app@123
DB_NAME=golang_demo
DB_SCHEMA=public
DB_SSL_MODE=disable
# Options: disable, require, verify-ca, verify-full
DB_TIMEZONE=Asia/Jakarta
DB_MIGRATE=TRUE
DB_SEED=TRUE
DB_SEED_FILE=import.sql
# Set to INFO for development and staging, SILENT for production
DB_LOG=SILENT

# Redis configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_USER=default
REDIS_PASS=
REDIS_DB=0
REDIS_FLUSH_DB=TRUE

# Idempotency configuration
IDEMPOTENCY_ENABLED=TRUE
IDEMPOTENCY_KEY_HEADER=Idempotency-Key
IDEMPOTENCY_PREFIX=idempotency_cache:
IDEMPOTENCY_TTL_HOURS=24
```

- **🔐 Notes**:  
  - `IS_SSL=TRUE`: Enable this if you want your app to run over `HTTPS`. Make sure to run `generate-certificate.sh` to generate **self-signed certificates** and place them in the `./cert/` directory (e.g., `mycert.key`, `mycert.cer`).
  - Make sure your paths (`./cert/`) exist and are accessible by the application during runtime.
  - `DB_TIMEZONE=Asia/Jakarta`: Adjust this value to your local timezone (e.g., `America/New_York`, etc.).
  - `DB_MIGRATE=TRUE`: Set to `TRUE` to automatically run `GORM` migrations for all entity definitions on app startup.
  - `DB_SEED=TRUE` & `DB_SEED_FILE=import.sql`: Use these settings if you want to insert predefined data into the database using the SQL file provided.
  - `DB_USER=appuser`, `DB_PASS=app@123`: It's strongly recommended to create a dedicated database user instead of using the default postgres superuser.

### 🔐 Generate Certificate for HTTPS (Optional)  

If `IS_SSL=TRUE` in your `.env`, generate the certificate files by running this file:  
```bash
./generate-certificate.sh
```

- **Notes**:  
  - On **Linux/macOS**: Run the script directly
  - On **Windows**: Use **WSL** to execute the `.sh` script

This will generate:
  - `./cert/mycert.key`
  - `./cert/mycert.cer`


Ensure these are correctly referenced in your `.env`:
```properties
IS_SSL=TRUE
SSL_KEYS=./cert/mycert.key
SSL_CERT=./cert/mycert.cer
```

### 👤 Create Dedicated PostgreSQL User (Recommended)

For security reasons, it's recommended to avoid using the default postgres superuser. Use the following SQL script to create a dedicated user (`appuser`) and assign permissions:

```sql
-- Create appuser and database
CREATE USER appuser WITH PASSWORD 'app@123';

-- Allow user to connect to database
GRANT CONNECT, TEMP, CREATE ON DATABASE golang_demo TO appuser;

-- Grant permissions on public schema
GRANT USAGE, CREATE ON SCHEMA public TO appuser;

-- Grant all permissions on existing tables
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO appuser;

-- Grant all permissions on sequences (if using SERIAL/BIGSERIAL ids)
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA public TO appuser;

-- Ensure future tables/sequences will be accessible too
ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO appuser;

-- Ensure future sequences will be accessible too
ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO appuser;
```

Update your `.env` accordingly:
```properties
DB_USER=appuser
DB_PASS=app@123
```

---


## 🚀 Running the Application  

This section provides step-by-step instructions to run the application either **locally** or via **Docker containers**.

- **Notes**:  
  - All commands are defined in the `Makefile`.
  - To run using `make`, ensure that `make` is installed on your system.
  - To run the application in containers, make sure `Docker` is installed and running.
  - Ensure you have `Go` installed on your system

### 📦 Install Dependencies

Make sure all Go modules are properly installed:  

```bash
make tidy
```

### 🧪 Run Unit Tests

```bash
make test
```

### 🔧 Run Locally (Non-containerized)

Ensure Redis and PostgreSQL are running locally, then:

```bash
make run
```

### 🐳 Run Using Docker

To build and run all services (Redis, PostgreSQL, Go app):

```bash
make docker-up
```

To stop and remove all containers:

```bash
make docker-down
```

- **Notes**:  
  - Before running the application inside Docker, make sure to update your environment variables `.env`
    - Change `DB_HOST=localhost` to `DB_HOST=postgres-server`.
    - Change `REDIS_HOST=localhost` to `REDIS_HOST=redis-server`.

### 🟢 Application is Running

Now your application is accessible at:
```bash
http://localhost:1000
```

or 

```bash
https://localhost:1000 (if SSL is enabled)
```

---

## 🧪 Testing Scenarios  

### 👨‍👩‍👧‍👦 Consumer API

#### Scenario 1: Create Consumer

**📌 Endpoint**: 
```http
POST https://localhost:1000/api/v1/consumers
```

**📥 Request Body**:
```json
{
    "fullname": "Austin Libertus",
    "username": "auslibertus",
    "email": "austin.libertus@example.com",
    "phone": "+628997452753",
    "address": "Jl. Anggrek No. 4, Jakarta",
    "birthDate": "1990-03-05"
}
```

**✅ Expected Response**:
```json
{
    "message": "Consumer created successfully",
    "error": null,
    "path": "/api/v1/consumers",
    "status": 201,
    "data": {
        "id": "4c6c42bc-3b82-4f34-9eaf-c4dcfb246ec0",
        "fullname": "Austin Libertus",
        "username": "auslibertus",
        "email": "austin.libertus@example.com",
        "phone": "628997452753",
        "address": "Jl. Anggrek No. 4, Jakarta",
        "birthDate": "1990-03-05",
        "status": "inactive",
        "createdAt": "2025-06-18T11:42:13.165068Z",
        "updatedAt": "2025-06-18T11:42:13.165068Z"
    },
    "timestamp": "2025-06-18T11:42:13.171205664Z"
}
```

#### Scenario 2: Update Consumer Status

**📌 Endpoint**: 
```http
PATCH https://localhost:1000/api/v1/consumers/4c6c42bc-3b82-4f34-9eaf-c4dcfb246ec0?status=active
```

**✅ Expected Response**:
```json
{
    "message": "Consumer status updated successfully",
    "error": null,
    "path": "/api/v1/consumers/4c6c42bc-3b82-4f34-9eaf-c4dcfb246ec0",
    "status": 200,
    "data": {
        "id": "4c6c42bc-3b82-4f34-9eaf-c4dcfb246ec0",
        "fullname": "Austin Libertus",
        "username": "auslibertus",
        "email": "austin.libertus@example.com",
        "phone": "628997452753",
        "address": "Jl. Anggrek No. 4, Jakarta",
        "birthDate": "1990-03-05",
        "status": "active",
        "createdAt": "2025-06-18T11:42:13.165068Z",
        "updatedAt": "2025-06-18T11:44:52.059458364Z"
    },
    "timestamp": "2025-06-18T11:44:52.062880937Z"
}
```

#### Scenario 3: Get All Consumers

**📌 Endpoint**: 
```http
GET https://localhost:1000/api/v1/consumers?page=1&limit=10
```

**✅ Expected Response**:
```json
{
    "message": "All consumers retrieved successfully",
    "error": null,
    "path": "/api/v1/consumers",
    "status": 200,
    "data": [
        {
            "id": "74fe86f3-6324-42c2-97b4-fa3225461299",
            "fullname": "John Doe",
            "username": "johndoe",
            "email": "john.doe@example.com",
            "phone": "6281234567890",
            "address": "Jl. Merdeka No. 123, Jakarta",
            "birthDate": "1990-05-10",
            "status": "active",
            "createdAt": "2025-06-18T11:40:56.66591Z",
            "updatedAt": "2025-06-18T11:40:56.66591Z"
        }
        ...
    ],
    "timestamp": "2025-06-18T13:11:24.539972654Z"
}
```

### 💳 Transaction API

Each `POST` request must also include a unique `Idempotency-Key` header to ensure safe retries:
```http
Idempotency-Key: <UUID>
```

#### Scenario 1: Create a New Transaction with Non-Existent Consumer

**📌 Endpoint**:  
```http
POST https://localhost:1000/api/v1/transactions
```

**📥 Request Body**:
```json
{
  "type": "payment",
  "amount": 150000.00,
  "consumerId": "4c6c42bc-3b82-4f34-9eaf-c4dcfb246ec0"
}
```

**❌ Expected Response**:
```json
{
  "message": "Consumer not found",
  "error": "No consumer found with the given ID",
  "path": "/api/v1/transactions",
  "status": 404,
  "data": null,
  "timestamp": "2025-06-18T16:02:57.380180455Z"
}
```

#### Scenario 2: Create a New Transaction with Inactive Consumer

**📌 Endpoint**:  
```http
POST https://localhost:1000/api/v1/transactions
```

**📥 Request Body**:
```json
{
  "type": "payment",
  "amount": 150000.00,
  "consumerId": "4c6c42bc-3b82-4f34-9eaf-c4dcfb246ec0"
}
```

**❌ Expected Response**:
```json
{
  "message": "Invalid transaction data",
  "error": "Transaction data is invalid, this could be due to missing required fields or incorrect data types",
  "path": "/api/v1/transactions",
  "status": 400,
  "data": null,
  "timestamp": "2025-06-18T16:03:23.349569947Z"
}
```

#### Scenario 3: Create a New Transaction Successfully

**📌 Endpoint**:  
```http
POST https://localhost:1000/api/v1/transactions
```

**📥 Request Body**:
```json
{
  "type": "payment",
  "amount": 150000.00,
  "consumerId": "a1b9d37e-2e7d-42b2-9d3e-7b492162905d"
}
```

**✅ Expected Response**:
```json
{
  "message": "Transaction created successfully",
  "error": null,
  "path": "/api/v1/transactions",
  "status": 201,
  "data": {
    "id": "147735b9-eff7-469d-ac85-3b8108825ce4",
    "idempotencyCacheKey": "06f14f72-dfba-49ca-aa4e-d85b532ca0b7",
    "type": "payment",
    "amount": 150000,
    "status": "pending",
    "consumerId": "a1b9d37e-2e7d-42b2-9d3e-7b492162905d",
    "createdAt": "2025-06-18T16:19:59.952804Z",
    "updatedAt": "2025-06-18T16:19:59.952804Z"
  },
  "timestamp": "2025-06-18T16:20:01.005272013Z"
}
```

#### Scenario 4: Same Idempotency-Key, Different Request Payload

This scenario demonstrates how the system prevents inconsistent processing when the **same `Idempotency-Key`** is used with a **different request body**.  

**📌 Endpoint**:  
```http
POST https://localhost:1000/api/v1/transactions
```

**📥 Request Body**:
```json
{
  "type": "payment",
  "amount": 170000.00,
  "consumerId": "a1b9d37e-2e7d-42b2-9d3e-7b492162905d"
}
```

**❌ Expected Response**:
```json
{
  "message": "Conflict",
  "error": "Request with the same Idempotency-Key but different body has already been processed",
  "path": "/api/v1/transactions",
  "status": 409,
  "data": null,
  "timestamp": "2025-06-18T15:24:50.515722414Z"
}
```

**Explanation**:
- The `Idempotency-Key` matches an existing record.  
- However, the **SHA-256 hash of the current payload** differs from the original request.  
- The system **rejects the request** with an HTTP **409 Conflict** to preserve data integrity and ensure **idempotent guarantees**.  
- The original transaction is **not modified** or **replaced**.  

#### Scenario 5: Reusing the Same Idempotency-Key with Identical Request

This scenario demonstrates how the system responds when a **request is retried** with the same `Idempotency-Key` and an **identical payload**.  

**📌 Endpoint**:  
```http
POST https://localhost:1000/api/v1/transactions
```

**📥 Request Body**:
```json
{
  "type": "payment",
  "amount": 150000.00,
  "consumerId": "a1b9d37e-2e7d-42b2-9d3e-7b492162905d"
}
```

**✅ Expected Response**:
```json
{
  "message": "Request already processed",
  "error": null,
  "path": "/api/v1/transactions",
  "status": 200,
  "data": {
    "amount": 150000,
    "consumerId": "a1b9d37e-2e7d-42b2-9d3e-7b492162905d",
    "createdAt": "2025-06-18T16:19:59.952804Z",
    "id": "147735b9-eff7-469d-ac85-3b8108825ce4",
    "idempotencyCacheKey": "06f14f72-dfba-49ca-aa4e-d85b532ca0b7",
    "status": "failed",
    "type": "payment",
    "updatedAt": "2025-06-18T16:20:08.921759395Z"
  },
  "timestamp": "2025-06-18T16:21:03.791516931Z"
}
```

**Explanation**:
- The system detects that the request has **already been processed** based on:
  - Matching `Idempotency-Key`  
  - Matching SHA-256 hash of the request body  
- The **original response is returned** from Redis cache.  
- The client receives a consistent, successful `200 OK` with the **same transaction data**.  