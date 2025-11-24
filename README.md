# User Service - gRPC Microservice with ScyllaDB

A high-performance user management microservice built with Go, gRPC, REST API gateway, and ScyllaDB. This service provides comprehensive user management capabilities with both gRPC and HTTP/REST interfaces.

## 🏗️ Architecture

### System Architecture
```
┌─────────────────────────────────────────────────────────────────┐
│                         Client Layer                            │
│  ┌──────────────────┐                ┌──────────────────┐       │
│  │   gRPC Client    │                │   HTTP Client    │       │
│  │  (Port 50051)    │                │   (Port 8080)    │       │
│  └────────┬─────────┘                └─────────┬────────┘       │
└───────────┼────────────────────────────────────┼────────────────┘
            │                                    │
            │                                    │
┌───────────┼────────────────────────────────────┼────────────────┐
│           │         Application Layer          │                │
│           │                                    │                │
│  ┌────────▼─────────┐              ┌───────────▼───────┐        │
│  │   gRPC Server    │              │  gRPC Gateway     │        │
│  │  (UserService)   │◄─────────────│  (HTTP→gRPC)      │        │
│  └────────┬─────────┘              └───────────────────┘        │
│           │                                                     │
│  ┌────────▼─────────────────────────────────────────┐           │
│  │            Business Logic Layer                  │           │
│  │  ┌─────────────┐        ┌──────────────┐         │           │
│  │  │   Domain    │        │  Validators  │         │           │
│  │  │   Models    │        │              │         │           │
│  │  └─────────────┘        └──────────────┘         │           │
│  └───────────────────────────┬──────────────────────┘           │
│                              │                                  │
│  ┌───────────────────────────▼──────────────────────┐           │
│  │           Storage Interface Layer                │           │
│  │  ┌────────────────────────────────────────────┐  │           │
│  │  │      ScyllaDB Storage Implementation       │  │           │
│  │  └────────────────────────────────────────────┘  │           │
│  └────────────────────────────┬─────────────────────┘           │
└───────────────────────────────┼─────────────────────────────────┘
                                │
┌───────────────────────────────▼────────────────────────────────┐
│                      Database Layer                            │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    ScyllaDB Cluster                     │   │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │   │
│  │  │    users     │  │users_by_phone│  │users_by_email│   │   │
│  │  │    table     │  │    table     │  │    table     │   │   │
│  │  └──────────────┘  └──────────────┘  └──────────────┘   │   │
│  └─────────────────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────────────────┘
```

### Component Description

#### 1. **Client Layer**
- **gRPC Client**: Direct gRPC communication on port 50051
- **HTTP Client**: RESTful API communication on port 8080

#### 2. **Application Layer**
- **gRPC Server**: Implements `UserService` proto definitions
- **gRPC Gateway**: Translates HTTP/REST requests to gRPC calls
- **Handlers**: Process incoming requests and coordinate business logic

#### 3. **Business Logic Layer**
- **Domain Models**: Core user entity with business rules
- **Validators**: Input validation logic (email, phone, names)
- **Error Handling**: Domain-specific errors

#### 4. **Storage Layer**
- **Storage Interface**: Defines contract for data operations
- **ScyllaDB Implementation**: Concrete implementation using gocql driver

#### 5. **Database Layer**
- **users**: Main user data table (partitioned by user ID)
- **users_by_phone**: Phone number lookup table
- **users_by_email**: Email lookup table

### Data Flow

**Create User Flow:**
```
HTTP POST /api/v1/users
    ↓
gRPC Gateway (translates to gRPC)
    ↓
UserServiceServer.CreateUser()
    ↓
Validate user data
    ↓
Check email/phone uniqueness
    ↓
Insert into users table
    ↓
Insert into users_by_phone table
    ↓
Insert into users_by_email table
    ↓
Return created user
```

## 🚀 Features

- ✅ **Dual Protocol Support**: Both gRPC and HTTP/REST APIs
- ✅ **High Performance**: ScyllaDB for low-latency operations
- ✅ **CRUD Operations**: Complete user management
- ✅ **User Blocking**: Block/unblock users
- ✅ **Contact Updates**: Separate endpoint for email/phone updates
- ✅ **Multiple Lookups**: Query by ID, email, or phone number
- ✅ **Pagination**: List users with page tokens
- ✅ **Validation**: Comprehensive input validation
- ✅ **Uniqueness Constraints**: Email and phone uniqueness
- ✅ **Graceful Shutdown**: Clean server termination

## 📋 Prerequisites

- **Go**: 1.25.3 or higher
- **Docker**: For running ScyllaDB
- **Docker Compose**: For container orchestration
- **Protocol Buffers**: `protoc` compiler
- **Make**: For running build commands

### Install Protocol Buffers Tools
```bash
# Install protoc compiler
# macOS
brew install protobuf

# Ubuntu/Debian
sudo apt install -y protobuf-compiler

# Install Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

# Add to PATH if not already
export PATH="$PATH:$(go env GOPATH)/bin"
```

## 📦 Installation

### 1. Clone the Repository
```bash
git clone https://github.com/Divyansh031/user-service.git
cd user-service
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Setup Environment
```bash
# Copy example env file
cp .env.example .env

# Edit .env with your settings (optional, defaults work for local development)
```

### 4. Start ScyllaDB
```bash
# Start ScyllaDB container and initialize schema
make docker-up

# This will:
# - Start ScyllaDB container
# - Wait for it to be healthy
# - Apply database schema automatically
```

### 5. Generate Protobuf Code
```bash
make proto
```

### 6. Build the Application
```bash
make build
```

### 7. Run the Service
```bash
make run

# Or directly
./bin/server
```

You should see:
```
time=2025-11-23T... level=INFO msg="Starting user service" env=development
time=2025-11-23T... level=INFO msg="Initializing ScyllaDB" hosts=[localhost] keyspace=userservice
time=2025-11-23T... level=INFO msg="ScyllaDB initialized successfully"
time=2025-11-23T... level=INFO msg="gRPC server listening" port=50051
time=2025-11-23T... level=INFO msg="HTTP REST server listening" port=8080
```

## 🔧 Configuration

Configuration can be set via environment variables or `config/config.yaml`:

### Environment Variables
```bash
# Server Configuration
ENV=development
GRPC_PORT=50051
HTTP_PORT=8080

# ScyllaDB Configuration
SCYLLA_HOSTS=localhost
SCYLLA_PORT=9042
SCYLLA_KEYSPACE=userservice
SCYLLA_CONSISTENCY=QUORUM

# Logging
LOG_LEVEL=info
```

### Configuration File (config/config.yaml)
```yaml
env: "development"

grpc:
  port: 50051

http:
  port: 8080

scylladb:
  hosts:
    - localhost
  port: 9042
  keyspace: userservice
  consistency: QUORUM

log:
  level: info
```

## 📡 API Documentation

### Base URLs

- **HTTP/REST**: `http://localhost:8080/api`
- **gRPC**: `localhost:50051`

### API Endpoints

#### 1. Create User

**HTTP:**
```bash
POST /api/v1/users
Content-Type: application/json

{
  "first_name": "John",
  "last_name": "Doe",
  "gender": "male",
  "date_of_birth": "1990-01-15T00:00:00Z",
  "phone_number": "+919876543210",
  "email": "john.doe@example.com"
}
```

**Response (201 Created):**
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "first_name": "John",
    "last_name": "Doe",
    "gender": "male",
    "date_of_birth": "1990-01-15T00:00:00Z",
    "phone_number": "+919876543210",
    "email": "john.doe@example.com",
    "is_blocked": false,
    "created_at": "2025-11-23T21:47:35Z",
    "updated_at": "2025-11-23T21:47:35Z"
  }
}
```

**gRPC:**
```bash
grpcurl -plaintext -d '{
  "first_name": "John",
  "last_name": "Doe",
  "gender": "male",
  "date_of_birth": {"seconds": 632448000},
  "phone_number": "+919876543210",
  "email": "john.doe@example.com"
}' localhost:50051 user.v1.UserService/CreateUser
```

**Validation Rules:**
- `first_name`: 2-50 characters, required
- `last_name`: 2-50 characters, required
- `gender`: must be "male", "female", or "other"
- `date_of_birth`: must be in the past
- `phone_number`: E.164 format (e.g., +919876543210)
- `email`: valid email format
- Email and phone must be unique

---

#### 2. Get User by ID

**HTTP:**
```bash
GET /api/v1/users/{id}
```

**Example:**
```bash
curl http://localhost:8080/api/v1/users/550e8400-e29b-41d4-a716-446655440000
```

**Response (200 OK):**
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "first_name": "John",
    "last_name": "Doe",
    "gender": "male",
    "date_of_birth": "1990-01-15T00:00:00Z",
    "phone_number": "+919876543210",
    "email": "john.doe@example.com",
    "is_blocked": false,
    "created_at": "2025-11-23T21:47:35Z",
    "updated_at": "2025-11-23T21:47:35Z"
  }
}
```

**gRPC:**
```bash
grpcurl -plaintext -d '{"id": "550e8400-e29b-41d4-a716-446655440000"}' \
  localhost:50051 user.v1.UserService/GetUser
```

---

#### 3. Get User by Email

**HTTP:**
```bash
GET /api/v1/users/email/{email}
```

**Example:**
```bash
curl http://localhost:8080/api/v1/users/email/john.doe@example.com
```

**gRPC:**
```bash
grpcurl -plaintext -d '{"email": "john.doe@example.com"}' \
  localhost:50051 user.v1.UserService/GetUserByEmail
```

---

#### 4. Get User by Phone

**HTTP:**
```bash
GET /api/v1/users/phone/{phone_number}
```

**Example:**
```bash
curl http://localhost:8080/api/v1/users/phone/+919876543210
```

**gRPC:**
```bash
grpcurl -plaintext -d '{"phone_number": "+919876543210"}' \
  localhost:50051 user.v1.UserService/GetUserByPhone
```

---

#### 5. Update User

**HTTP:**
```bash
PUT /api/v1/users/{id}
Content-Type: application/json

{
  "first_name": "Jane",
  "last_name": "Smith"
}
```

**Example:**
```bash
curl -X PUT http://localhost:8080/api/v1/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Jane",
    "last_name": "Smith"
  }'
```

**Response (200 OK):**
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "first_name": "Jane",
    "last_name": "Smith",
    "gender": "male",
    "date_of_birth": "1990-01-15T00:00:00Z",
    "phone_number": "+919876543210",
    "email": "john.doe@example.com",
    "is_blocked": false,
    "created_at": "2025-11-23T21:47:35Z",
    "updated_at": "2025-11-23T22:10:15Z"
  }
}
```

**Note:** Only provided fields are updated. Omitted fields retain their existing values.

**gRPC:**
```bash
grpcurl -plaintext -d '{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "first_name": "Jane",
  "last_name": "Smith"
}' localhost:50051 user.v1.UserService/UpdateUser
```

---

#### 6. Update User Contact

**HTTP:**
```bash
PATCH /api/v1/users/{id}/contact
Content-Type: application/json

{
  "phone_number": "+911234567890",
  "email": "newemail@example.com"
}
```

**Example:**
```bash
curl -X PATCH http://localhost:8080/api/v1/users/550e8400-e29b-41d4-a716-446655440000/contact \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newemail@example.com"
  }'
```

**Note:** You can update phone, email, or both. Omitted fields remain unchanged.

**gRPC:**
```bash
grpcurl -plaintext -d '{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "newemail@example.com"
}' localhost:50051 user.v1.UserService/UpdateUserContact
```

---

#### 7. Block User

**HTTP:**
```bash
POST /api/v1/users/{id}/block
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/users/550e8400-e29b-41d4-a716-446655440000/block
```

**Response (200 OK):**
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "first_name": "Jane",
    "last_name": "Smith",
    "is_blocked": true,
    ...
  }
}
```

**gRPC:**
```bash
grpcurl -plaintext -d '{"id": "550e8400-e29b-41d4-a716-446655440000"}' \
  localhost:50051 user.v1.UserService/BlockUser
```

---

#### 8. Unblock User

**HTTP:**
```bash
POST /api/v1/users/{id}/unblock
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/users/550e8400-e29b-41d4-a716-446655440000/unblock
```

**gRPC:**
```bash
grpcurl -plaintext -d '{"id": "550e8400-e29b-41d4-a716-446655440000"}' \
  localhost:50051 user.v1.UserService/UnblockUser
```

---

#### 9. List Users

**HTTP:**
```bash
GET /api/v1/users?page_size=10&page_token=<token>
```

**Example:**
```bash
curl "http://localhost:8080/api/v1/users?page_size=10"
```

**Response (200 OK):**
```json
{
  "users": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "first_name": "Jane",
      "last_name": "Smith",
      ...
    },
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "first_name": "Bob",
      "last_name": "Johnson",
      ...
    }
  ],
  "next_page_token": "660e8400-e29b-41d4-a716-446655440001",
  "total_count": 2
}
```

**Query Parameters:**
- `page_size`: Number of results (default: 10, max: 100)
- `page_token`: Token for next page (from previous response)

**gRPC:**
```bash
grpcurl -plaintext -d '{"page_size": 10}' \
  localhost:50051 user.v1.UserService/ListUsers
```

---

#### 10. Delete User

**HTTP:**
```bash
DELETE /api/v1/users/{id}
```

**Example:**
```bash
curl -X DELETE http://localhost:8080/api/v1/users/550e8400-e29b-41d4-a716-446655440000
```

**Response (200 OK):**
```json
{}
```

**gRPC:**
```bash
grpcurl -plaintext -d '{"id": "550e8400-e29b-41d4-a716-446655440000"}' \
  localhost:50051 user.v1.UserService/DeleteUser
```

---

### Error Responses

**400 Bad Request - Validation Error:**
```json
{
  "code": 3,
  "message": "invalid email",
  "details": []
}
```

**404 Not Found:**
```json
{
  "code": 5,
  "message": "user not found",
  "details": []
}
```

**409 Conflict - Duplicate:**
```json
{
  "code": 6,
  "message": "email already exists",
  "details": []
}
```

**500 Internal Server Error:**
```json
{
  "code": 13,
  "message": "internal server error",
  "details": []
}
```

## 🧪 Testing

### Run Unit Tests
```bash
make test
```

### Run Integration Tests
```bash
make test-integration
```

### Run Benchmarks
```bash
make benchmark
```

### Manual Testing with cURL

**Create a user:**
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Test",
    "last_name": "User",
    "gender": "male",
    "date_of_birth": "1995-06-20T00:00:00Z",
    "phone_number": "+919123456789",
    "email": "test@example.com"
  }'
```

**Get user:**
```bash
curl http://localhost:8080/api/v1/users/{user_id}
```

**Update user:**
```bash
curl -X PUT http://localhost:8080/api/v1/users/{user_id} \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Updated",
    "last_name": "Name"
  }'
```

**Delete user:**
```bash
curl -X DELETE http://localhost:8080/api/v1/users/{user_id}
```

### Testing with Postman

Import the following collection:

**Collection Variables:**
- `base_url`: `http://localhost:8080/api`
- `user_id`: (set after creating a user)

**Example Requests:**
1. Create User → Save `user.id` to `user_id` variable
2. Get User → Use `{{user_id}}`
3. Update User → Use `{{user_id}}`
4. Delete User → Use `{{user_id}}`

### Testing with grpcurl

**List services:**
```bash
grpcurl -plaintext localhost:50051 list
```

**Describe service:**
```bash
grpcurl -plaintext localhost:50051 describe user.v1.UserService
```

**Call CreateUser:**
```bash
grpcurl -plaintext -d '{
  "first_name": "Alice",
  "last_name": "Wonder",
  "gender": "female",
  "date_of_birth": {"seconds": 757382400},
  "phone_number": "+919999888877",
  "email": "alice@example.com"
}' localhost:50051 user.v1.UserService/CreateUser
```

## 📁 Project Structure
```
user-service/
├── api/
│   └── proto/
│       └── user/
│           └── v1/
│               ├── user.proto      # Protobuf definitions
│               ├── user.pb.go      # Generated Go code
│               ├── user_grpc.pb.go # Generated gRPC code
│               └── user.pb.gw.go   # Generated gateway code
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   ├── config.go               # Configuration loader
│   │   └── config.yaml             # Default config file
│   ├── domain/
│   │   ├── user.go                 # User domain model
│   │   └── errors.go               # Domain errors
│   ├── grpc/
│   │   └── handlers/
│   │       └── user_handler.go     # gRPC service implementation
│   └── storage/
│       ├── storage.go              # Storage interface
│       └── scylla/
│           └── scylla.go           # ScyllaDB implementation
├── pkg/
│   └── validator/
│       └── validator.go            # Input validators
├── scripts/
│   └── schema.cql                  # Database schema
├── tests/
│   └── unit/
│       ├── config_test.go          # Config tests
│       ├── user_test.go            # Domain tests
│       └── validator_test.go       # Validator tests
├── .env.example                    # Environment variables template
├── .gitignore                      # Git ignore rules
├── docker-compose.yml              # Docker services
├── go.mod                          # Go dependencies
├── go.sum                          # Dependency checksums
├── Makefile                        # Build automation
└── README.md                       # This file
```

## 🗄️ Database Schema

### Users Table
```cql
CREATE TABLE users (
    id text PRIMARY KEY,
    first_name text,
    last_name text,
    gender text,
    date_of_birth timestamp,
    phone_number text,
    email text,
    is_blocked boolean,
    created_at timestamp,
    updated_at timestamp
);
```

### Users by Phone (Lookup Table)
```cql
CREATE TABLE users_by_phone (
    phone_number text PRIMARY KEY,
    user_id text,
    created_at timestamp
);
```

### Users by Email (Lookup Table)
```cql
CREATE TABLE users_by_email (
    email text PRIMARY KEY,
    user_id text,
    created_at timestamp
);
```

### Indexes
```cql
CREATE INDEX ON users (email);
CREATE INDEX ON users (phone_number);
CREATE INDEX ON users (is_blocked);
```

## 🛠️ Available Make Commands
```bash
make help          # Show available commands
make proto         # Generate protobuf code
make build         # Build the application
make run           # Run the application
make test          # Run unit tests
make test-integration  # Run integration tests
make benchmark     # Run benchmarks
make docker-up     # Start Docker containers
make docker-down   # Stop Docker containers
make clean         # Clean generated files
```


## 📚 Technology Stack

- **Language**: Go 1.25.3
- **Framework**: gRPC, gRPC-Gateway
- **Database**: ScyllaDB 5.2
- **Protocol**: Protocol Buffers 3
- **Containerization**: Docker, Docker Compose
- **Testing**: testify
- **Configuration**: cleanenv
- **Database Driver**: gocql


---