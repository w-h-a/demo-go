# demo-go

This repo provides a Go implementation of testable code architecture. The architecture is built on a simple `inbound -> business logic/domain -> outbound` flow. The goals of this design are not just about clean organization and ease of extension, but testability. The structure is deliberately designed to enable a multi-layered testing strategy and code that is both maintainable and extensible. Most of the insights herein can be found in Vladimir Khorikov's _Unit Testing_.

## Architecture

### The Layers

The main layers are:

1. The Inbound Layer (Handler/Controller): This is a humble object. Its _only responsibility_ is to be a translator. It parses external requests, calls the relevant service at the business logic layer, and maps the result back to the external response.
2. The Business Logic Layer (Service): This is the core orchestrator. This layer makes features like _Create a User_ possible. It coordinates the flow by orchestrating external dependencies via interfaces, although it does not itself perform outbound calls and indeed it neither speaks the language of the inbound layer nor the outbound layer. This layer is usually performing core business logic on domain objects like `User` (`api/user`). 
3. The Outbound Layer (Client/Adapter): These consist primarily of the implementations (`pgUserRepo`) of the interfaces (`UserRepo`) that the business logic layer uses to orchestrate external dependencies to get the job done.

These layers are stacked on top of one another via dependency inversion/injection. The service is willing to accept interfaces and ends up holding implementations of those external client interfaces, and the handler holds a service to which it hands off a DTO (data transfer object).

### Contrasted Anti-Pattern

An alternative model, where the inbound layer orchestrates business logic and calls out to external dependencies is an anti-pattern. This is overcomplicated code because it mixes business logic with external dependencies, making it unnecessarily difficult to test, maintain, and extend.

## Testing

1. Unit Tests (`make unit-test`)
    * Target: Service Layer
    * Goal: Test all business logic in isolation
    * Method: Mock only inter-system dependencies (`UserRepo`), not intra-system collaborators (`User`).
2. Integration Tests (`make integration-test`)
    * Target: Inbound Layer
    * Goal: Test the full slice from request to the database and back.
    * Method: Send real requests to the running app, which connects to a real Postgres database (running in Docker). Mock only unmanaged dependencies (e.g., **Stripe**) with something like Wiremock (running in Docker), for example.

## How to Run

1. `docker-compose up`
2. `cp .env.example .env` (and edit if necessary)
3. `go run main.go demo`