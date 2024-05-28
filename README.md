# Paseto Auth Example

This is a simple Golang application demonstrating authentication using PASETO tokens with the Gin framework.

## Requirements

- Go 1.16+
- Gin framework
- Paseto library

## Setup

1. Clone the repository.
2. Install dependencies using `go get`:
   ```bash
   go get -u github.com/gin-gonic/gin
   go get -u github.com/o1egl/paseto
   ```
3. Set the environment variable `PASETO_SECRET_KEY` to a 32-byte secret key:
   ```bash
   export PASETO_SECRET_KEY="your-32-byte-secret-key"
   ```

## Running the Server

1. Build and run the server:
   ```bash
   go run main.go
   ```
2. The server will run on `localhost:2000`.

## Endpoints

### Login

- **URL**: `/login`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "username": "username123",
    "password": "password123"
  }
  ```
- **Response**:
  ```json
  {
    "token": "your-paseto-token"
  }
  ```

### Private Route

- **URL**: `/private`
- **Method**: `GET`
- **Headers**:
  ```json
  {
    "Authorization": "Bearer your-paseto-token"
  }
  ```
- **Response**:
  ```json
  {
    "message": "Hello from a private route"
  }
  ```

## Example Usage

1. **Login**:
   ```bash
   curl -X POST http://localhost:2000/login -H "Content-Type: application/json" -d '{"username": "username123", "password": "password123"}'
   ```

2. **Access Private Route**:
   ```bash
   curl -X GET http://localhost:2000/private -H "Authorization: Bearer <token from login response>"
   ```

## Code Overview

- **`authMiddleware`**: Middleware to verify PASETO tokens.
- **`generateToken`**: Function to generate a PASETO token.
- **`loginUser`**: Handler for the login route.
- **`privateRoute`**: Handler for a protected route.
- **`setupRouter`**: Sets up routes and applies middleware.

## Notes

- Ensure the `PASETO_SECRET_KEY` environment variable is set correctly.
- Modify the login logic to integrate with your user authentication system.
