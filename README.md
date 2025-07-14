# Go Gin Layered Architecture Template

A template/demonstration for building RESTful APIs in Go using the Gin framework. This project showcases a layered architecture pattern designed for scalability, maintainability, and separation of concerns.

## Features

- RESTful API design
- Example JWT-based authentication
- Example CRUD operations
- Request payload validation
- Structured, layered architecture
- Custom error handling
- Formatted request logging

## Project Structure

The project follows a standard layered architecture pattern to separate concerns, making the codebase clean, maintainable, and easy to test.

```
.
├── cmd/api/
│   └── main.go         # Main application entry point
├── internal/
│   ├── handlers/       # HTTP request handlers (controllers)
│   ├── middleware/     # Gin middleware (e.g., logging, auth)
│   ├── models/         # Data structures (request/response models, DB models)
│   ├── repository/     # Data access layer (interacts with the database)
│   ├── routes/         # API route definitions
│   └── services/       # Business logic layer
├── pkg/
│   ├── errors/         # Custom application-wide error types
│   └── utils/          # Shared utility functions (e.g., validator)
├── .env.example        # Example environment variables
├── go.mod              # Go module definitions
├── go.sum              # Go module checksums
└── README.md           # This file
```

### Directory Explanations

*   **`/cmd/api`**: The main entry point for the web application. It is responsible for initializing the server, database connection, configuration, routes, and starting the HTTP server.

*   **`/internal`**: This directory contains all the private application code. According to Go's convention, code within an `internal` directory is not importable by other projects, ensuring encapsulation of your core logic.

    *   **`/handlers`**: This is the presentation or "controller" layer. Handlers are responsible for parsing incoming HTTP requests, calling the appropriate services to handle business logic, and formatting the HTTP response (e.g., sending JSON data and status codes).

    *   **`/services`**: This is the business logic layer. It orchestrates the application's functionality, acting as a bridge between the transport layer (`handlers`) and the data access layer (`repository`). It contains the core logic of what the application does.

    *   **`/repository`**: This is the data access layer. It is responsible for all communication with the database (in this case, MongoDB). It abstracts the data storage details from the rest of the application, providing a clean API for data manipulation (Create, Read, Update, Delete).

    *   **`/models`**: Contains all the Go structs that model the application's data. This includes request/response body structures, database entities, and any other data structures used throughout the application.

    *   **`/middleware`**: Holds custom Gin middleware. Middleware can intercept incoming requests to perform tasks like logging, authentication, authorization, or header manipulation before the request reaches the handler.

    *   **`/routes`**: Defines the API endpoints and maps them to their respective handlers. This helps in organizing all the application's routes in one place.

*   **`/pkg`**: This directory contains shared, public code that could potentially be used by other applications. It's for libraries that are okay to be imported externally.

    *   **`/errors`**: Defines custom, reusable error types for consistent error handling throughout the application (e.g., `AppError`).

    *   **`/utils`**: A collection of helper functions, such as the request validator (`validator.go`), that can be used across different parts of the application.

## API Endpoints

API endpoints are defined in the `/internal/routes` package. This project includes example routes for authentication and resource management to demonstrate how to structure your API routing.

## Getting Started

### Prerequisites

- Go (version 1.18 or newer)
- A running database instance (e.g., MongoDB)

### Installation & Setup

1.  **Clone the repository:**
    ```sh
    git clone <repository-url>
    cd <project-directory-name>
    ```

2.  **Install dependencies:**
    ```sh
    go mod tidy
    ```

3.  **Configure environment:**
    Create a `.env` file by copying the example file.
    ```sh
    cp .env.example .env
    ```
    Update the `.env` file with your database connection string, JWT secret, and other necessary configurations.

### Running the Application

Execute the following command from the project root:

```sh
go run ./cmd/api/main.go
```

The server will start, and by default, it should be listening on `http://localhost:8080`.