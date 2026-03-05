# E-Commerce Backend API

A full-featured e-commerce backend REST API built with Go, providing user authentication, product management, shopping cart, order processing, and integrated payment gateway support.

## Tech Stack

- **Language:** Go
- **Web Framework:** [Gin](https://github.com/gin-gonic/gin)
- **Database:** PostgreSQL with [GORM](https://gorm.io/) ORM (auto-migration on startup)
- **Authentication:** JWT (JSON Web Tokens)
- **Password Hashing:** bcrypt
- **Payment Gateway:** [Midtrans](https://midtrans.com/) Snap
- **Environment Config:** godotenv

## Project Structure

```
.
├── cmd/            # Application entry point (main.go)
├── config/         # Database connection and configuration
├── handler/        # HTTP request handlers
├── middleware/     # JWT authentication & admin authorization middleware
├── models/         # GORM data models (User, Product, Cart, Order, etc.)
├── repository/     # Data access layer
├── routes/         # API route definitions
├── service/        # Business logic layer
└── .env.example    # Example environment variable file
```

## Prerequisites

- Go 1.18+
- PostgreSQL
- A [Midtrans](https://midtrans.com/) account (for payment features)

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/Restartor/E-commerce.git
cd E-commerce
```

### 2. Configure environment variables

Copy the example environment file and fill in your values:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
# Database
DB_HOST=localhost
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
DB_PORT=5432
DB_SSLMODE=disable

# JWT
JWT_SECRET=your_jwt_secret

# Midtrans Payment Gateway
MIDTRANS_SERVER_KEY=your_midtrans_server_key
MIDTRANS_CLIENT_KEY=your_midtrans_client_key
MIDTRANS_ENV=sandbox   # or "production"
```

### 3. Install dependencies

```bash
go mod tidy
```

### 4. Run the application

```bash
go run cmd/main.go
```

The server will start on `http://localhost:3000`. Database tables are created automatically via GORM auto-migration.

## API Endpoints

### Public Routes

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Health check / welcome |
| POST | `/register` | Register a new user |
| POST | `/login` | Login and receive a JWT token |
| GET | `/products` | List all products |

### Authenticated Routes (JWT required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/profile` | Get the current user's profile |
| GET | `/my-orders` | List the current user's orders |

#### Cart

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/cart` | Add an item to the cart |
| GET | `/cart` | View cart contents |
| PUT | `/cart/:id` | Update a cart item quantity |
| DELETE | `/cart/:id` | Remove a cart item |

#### Orders

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/order/checkout` | Checkout and create an order from the cart |

#### Payments

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/payment/notification` | Midtrans payment notification webhook |

### Admin Routes (JWT + admin role required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/admin/products` | Create a new product |
| PUT | `/admin/products/:id` | Update a product |
| DELETE | `/admin/products/:id` | Delete a product |

## Authentication

1. Register via `POST /register` with `username`, `email`, and `password`.
2. Login via `POST /login` to receive a JWT token.
3. Include the token in the `Authorization` header for protected routes:

```
Authorization: Bearer <your_token>
```

Tokens expire after **5 hours**.

## Data Models

- **User** – username, email, hashed password, role (`customer` / `admin`)
- **Product** – name, description, price, stock
- **Cart** – belongs to a user, holds cart items
- **CartItem** – product reference, quantity, subtotal
- **Order** – belongs to a user, order status, total amount
- **OrderItem** – product snapshot per order line

## License

This project is open source. See the repository for more details.
