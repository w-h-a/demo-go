# demo-go

## Architecture

* internal/handler (inbound): Handlers responsible for parsing requests, calling the service, and sending responses. No business logic.
* internal/api (domain types): Contains all the shared domain types.
* internal/service (business logic): Contains all the business logic. Any required outbound calls go through the client layer.
* internal/client (outbound): Contains the interfaces (e.g., `UserRepository`) and their concrete implementations (e.g., `PostgresUserRepository`) for all outbound communication (whether that's database or other external dependencies).

## Testing

1. Unit Tests (`make unit-test`)
    * Target: Service Layer
    * Goal: Test all business logic in isolation
    * Method: Mock all outbound dependencies
2. Integration Tests (`make integration-test`)
    * Target: Inbound Layer
    * Goal: Test the full slice from request to the database and back.
    * Method: Send real requests to the running app, which connects to a real Postgres database (running in Docker)

## How to Run

1. `docker-compose up`
2. `cp .env.example .env` (and edit if necessary)
3. `go run main.go demo`