# Chirpy - A Twitter-like Microblogging API

[![Go Version](https://img.shields.io/badge/Go-1.24.4-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## 📖 Overview

Chirpy is a robust, production-ready microblogging API built with Go that provides Twitter-like functionality. It features user authentication, content management, real-time metrics, and premium user upgrades through webhook integrations.

### 🎯 What Chirpy Does

Chirpy is a RESTful API that enables users to:
- **Create and manage user accounts** with secure authentication
- **Post short messages (chirps)** with content filtering and validation
- **Authenticate users** using JWT tokens with refresh token support
- **Manage user profiles** with email and password updates
- **Upgrade to premium features** through webhook integrations
- **Track application metrics** with admin dashboard
- **Serve static content** with hit tracking

### 🚀 Why You Should Care

- **Production Ready**: Built with enterprise-grade security and scalability in mind
- **Modern Architecture**: Uses Go 1.24.4 with PostgreSQL and JWT authentication
- **Content Safety**: Built-in profanity filtering and content validation
- **Developer Friendly**: Comprehensive API documentation and easy setup
- **Extensible**: Modular design with clear separation of concerns
- **Monitoring**: Built-in metrics and health checks for production deployment

## 🛠️ Technology Stack

- **Backend**: Go 1.24.4
- **Database**: PostgreSQL
- **Authentication**: JWT with refresh tokens
- **Password Hashing**: bcrypt
- **Database ORM**: sqlc (type-safe SQL)
- **Environment**: godotenv
- **UUID Generation**: Google UUID

## 📋 Prerequisites

Before running Chirpy, ensure you have:

- **Go 1.24.4** or later installed
- **PostgreSQL** database server running
- **Git** for cloning the repository

## 🚀 Installation & Setup

### 1. Clone the Repository

```bash
git clone https://github.com/frozendolphin/Chirpy.git
cd Chirpy
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Database Setup

#### Create Postgres Database
```sql
CREATE DATABASE chirpy_db;
```

#### Run Database Migrations
The project uses Goose for database migrations. Install Goose first:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Then run the migrations:
```bash
goose -dir sql/schema postgres "your_connection_string" up
```

### 4. Environment Configuration

Create a `.env` file in the project root:

```env
DB_URL=postgres://username:password@localhost:5432/chirpy_db?sslmode=disable
PLATFORM=dev
SECRET=your_jwt_secret_key_here
POLKA_KEY=your_polka_webhook_key_here //polka is just like stripe
```

### 5. Generate Database Code

The project uses sqlc for type-safe database operations:

```bash
sqlc generate
```

### 6. Run the Application

```bash
go run .
```

The server will start on `http://localhost:8080`

## 📚 API Documentation

### Authentication Endpoints

#### POST `/api/users`
Create a new user account.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false
}
```

#### POST `/api/login`
Authenticate a user and receive access tokens.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "email": "user@example.com",
  "token": "jwt_access_token",
  "refresh_token": "jwt_refresh_token",
  "is_chirpy_red": false
}
```

#### POST `/api/refresh`
Refresh an access token using a refresh token.

**Headers:**
```
Authorization: Bearer <refresh_token>
```

**Response:**
```json
{
  "token": "new_jwt_access_token"
}
```

#### POST `/api/revoke`
Revoke a refresh token.

**Headers:**
```
Authorization: Bearer <refresh_token>
```

### Chirps (Posts) Endpoints

#### POST `/api/chirps`
Create a new chirp (post).

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "body": "This is my chirp content!"
}
```

**Response:**
```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "body": "This is my chirp content!",
  "user_id": "uuid"
}
```

#### GET `/api/chirps`
Retrieve all chirps.

**Response:**
```json
[
  {
    "id": "uuid",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "body": "This is my chirp content!",
    "user_id": "uuid"
  }
]
```

#### GET `/api/chirps/{chirpID}`
Retrieve a specific chirp by ID.

#### DELETE `/api/chirps/{chirpID}`
Delete a specific chirp (requires authentication).

**Headers:**
```
Authorization: Bearer <access_token>
```

### User Management Endpoints

#### PUT `/api/users`
Update user email and password.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "email": "newemail@example.com",
  "password": "newpassword"
}
```

### Admin Endpoints

#### GET `/admin/metrics`
View application metrics (hit counter).

**Response:**
```html
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited 42 times!</p>
  </body>
</html>
```

#### POST `/admin/reset`
Reset all users (development only).

### Webhook Endpoints

#### POST `/api/polka/webhooks`
Handle Polka webhook events for user upgrades.

**Headers:**
```
Authorization: ApiKey <polka_key>
```

**Request Body:**
```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "uuid"
  }
}
```

### Health & Monitoring

#### GET `/api/healthz`
Health check endpoint.

#### GET `/app/`
Serve static files with hit tracking.

## 🔒 Security Features

- **JWT Authentication**: Secure token-based authentication
- **Password Hashing**: bcrypt for secure password storage
- **Content Filtering**: Automatic profanity filtering
- **Input Validation**: Comprehensive request validation
- **CORS Support**: Cross-origin resource sharing configuration
- **Rate Limiting**: Built-in request limiting (configurable)

## 📊 Content Filtering

Chirpy automatically filters inappropriate content by replacing profane words with asterisks. The current filter list includes:
- kerfuffle
- sharbert
- fornax

## 🏗️ Project Structure

```
Chirpy/
├── main.go                 # Application entry point
├── handlechirps.go        # Chirp-related handlers
├── handleusers.go         # User management handlers
├── handlelogin.go         # Authentication handlers
├── handlerefresh.go       # Token refresh handlers
├── handlewebhook.go       # Webhook handlers
├── metrics.go             # Metrics and monitoring
├── readiness.go           # Health checks
├── json.go                # JSON response utilities
├── go.mod                 # Go module dependencies
├── sqlc.yaml             # SQL code generation config
├── test.http             # API testing examples
├── internal/
│   ├── auth/             # Authentication utilities
│   └── database/         # Generated database code
├── sql/
│   ├── schema/           # Database migrations
│   └── queries/          # SQL queries
└── assets/               # Static assets
```

## 🧪 Testing

Test the API using the provided `test.http` file or with curl:

```bash
# Test health endpoint
curl http://localhost:8080/api/healthz

# Create a user
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

## 🚀 Deployment

### Production Considerations

1. **Environment Variables**: Ensure all required environment variables are set
2. **Database**: Use a production PostgreSQL instance
3. **SSL/TLS**: Configure HTTPS for production
4. **Monitoring**: Set up proper logging and monitoring
5. **Backup**: Implement database backup strategies

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Uses sqlc for type-safe database operations
- JWT implementation for secure authentication
- PostgreSQL for reliable data storage

## 📞 Support

For support, please open an issue on GitHub.

---