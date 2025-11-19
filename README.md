# Ride-Hail Module — Alem Platform Backend

This module implements the **ride-hailing** functionality of the Alem platform backend.

---

## Table of Contents

1. [Overview](#overview)
2. [Features](#features)
3. [Project Structure](#project-structure)
4. [Configuration](#configuration)
5. [Getting Started](#getting-started)
6. [API](#api)
7. [Data Model](#data-model)
8. [Error Handling](#error-handling)
9. [Deployment](#deployment)
10. [Monitoring & Logging](#monitoring--logging)
11. [ERD](#erd-entity-relationship-diagram)
12. [Contributing](#contributing)
13. [License](#license)

---

## Overview

The **Ride-Hail** module handles all backend logic related to ride-hailing:

* Managing ride requests from users
* Matching riders with available drivers
* Tracking ride status (`requested`, `accepted`, `in_progress`, `completed`)
* Updating driver and rider locations in real-time
* Persisting ride history and related data
* Integrating with payment, notification, and geolocation services

The service is built following **Hexagonal Architecture (Ports & Adapters)** to ensure modularity, testability, and maintainability.

---

## Features

* Ride request lifecycle management (create, cancel, accept, complete)
* Real-time location updates via WebSocket
* Driver-rider matching logic
* Data persistence with PostgreSQL
* Asynchronous processing via RabbitMQ
* Authentication and authorization
* Error handling and structured logging

---

## Project Structure

<details>
<summary>Click to expand project tree</summary>

```text
RIDE-HAIL
│   config.yaml
│   docker-compose.yaml
│   erd-diagram.pdf
│   go.mod
│   go.sum
│   main.go
│   README.md
│
├───.idea
│   └───inspectionProfiles
│           Project_Default.xml
├───cmd
│   └───ride-hail
│           main.go
├───config
│       config.go
│       print.go
├───internal
│   ├───adapters
│   │   ├───http
│   │   │   ├───handle
│   │   │   │       dal-handler.go
│   │   │   │       handle.go
│   │   │   │       ride_handler.go
│   │   │   │
│   │   │   ├───dto
│   │   │   │       dal.go
│   │   │   │       login.go
│   │   │   │       ride.go
│   │   │   │
│   │   │   └───validate
│   │   │           email.go
│   │   │           password.go
│   │   │
│   │   ├───server
│   │   │       middleware.go
│   │   │       router.go
│   │   │       server.go
│   │   │
│   │   └───websocket
│   │           p_ws_handler.go
│   │           p_ws_manager.go
│   │
│   │   ├───postgres
│   │   │       coordinate_repository.go
│   │   │       driver_repository.go
│   │   │       ride_repository.go
│   │   │       user_repository.go
│   │   │
│   │   └───rabbit
│   │           driver_match_consumer.go
│   │           location_consumer.go
│   │           ride_status_consumer.go
│   │           setup.go
│   │
│   ├───app
│   │   │   app.go
│   │   │
│   │   ├───drive
│   │   │       drive.go
│   │   │
│   │   └───ride
│   │           ride.go
│   │
│   └───core
│       ├───domain
│       │   ├───action
│       │   │       action.go
│       │   │
│       │   ├───models
│       │   │       coordinate.go
│       │   │       dal.go
│       │   │       jwt.go
│       │   │       ride.go
│       │   │       user.go
│       │   │
│       │   └───types
│       │           errors.go
│       │           mode.go
│       │           role.go
│       │           status.go
│       │
│       ├───ports
│       │       interface.go
│       │
│       └───service
│               auth_service.go
│               dal_service.go
│               ride_service.go
│
│       ├───calculator
│       │       calculator.go
│       │
│       └───hash
│               hash.go
├───migrations
│       000001_create_tables.down.sql
│       000001_create_tables.up.sql
└───pkg
    ├───executor      executor.go
    ├───logger        logger.go
    ├───postgres      postgres.go
    ├───rabbit
    │       consumer.go
    │       producer.go
    │       rabbit.go
    └───txm           manager.go
```

</details>

---

## Configuration

The service reads configuration from `config.yaml` or environment variables:

```yaml
# Database Configuration
postgres:
  host: ${POSTGRES_HOST:-localhost}
  port: ${POSTGRES_PORT:-5432}
  user: ${POSTGRES_USER:-ridehail_user}
  password: ${POSTGRES_PASSWORD:-ridehail_pass}
  database: ${POSTGRES_DATABASE:-ridehail_db}

# RabbitMQ Configuration
rabbitmq:
  host: ${RABBITMQ_HOST:-localhost}
  port: ${RABBITMQ_PORT:-5672}
  user: ${RABBITMQ_USER:-guest}
  password: ${RABBITMQ_PASSWORD:-guest}

# WebSocket Configuration
websocket:
  port: ${WS_PORT:-8080}

# Service Ports
services:
  ride_service: ${RIDE_SERVICE_PORT:-3000}
  driver_location_service: ${DRIVER_LOCATION_SERVICE_PORT:-3001}
  admin_service: ${ADMIN_SERVICE_PORT:-3004}

# JWT Configuration
jwt:
  secret: ${secret:-V9muwjpb7rRfuAH0fNg+8g80/42v0kT7f7W67cabf3uCpMXATsE0Gzg/3GJtultt}
  expire_hours: 2
```

---

## Getting Started

### Running with Docker

```bash
docker-compose up -d
```

* **Postgres**: Exposed on `5432`
* **RabbitMQ**: Exposed on `5672` (AMQP) and `15672` (management UI)
* **Migrations**: Automatically run on container startup via `migrate` service

After the database and RabbitMQ are up, start your microservice in the desired mode:

```bash
# Ride service
go run cmd/ride-hail/main.go --mode=ride

# Driver service
go run cmd/ride-hail/main.go --mode=driver

# Admin service
go run cmd/ride-hail/main.go --mode=admin
```

> Each mode runs only the components relevant to that service, enabling independent scaling and easier debugging.

---

## API

### HTTP Endpoints

| Service                   | Method | Endpoint                      | Description                 |
| ------------------------- | ------ | ----------------------------- | --------------------------- |
| Ride Service              | POST   | /rides                        | Create a new ride request   |
| Ride Service              | POST   | /rides/{ride_id}/cancel       | Cancel a ride               |
| Driver & Location Service | POST   | /drivers/{driver_id}/online   | Driver goes online          |
| Driver & Location Service | POST   | /drivers/{driver_id}/offline  | Driver goes offline         |
| Driver & Location Service | POST   | /drivers/{driver_id}/location | Update driver location      |
| Driver & Location Service | POST   | /drivers/{driver_id}/start    | Start a ride                |
| Driver & Location Service | POST   | /drivers/{driver_id}/complete | Complete a ride             |
| Admin Service             | GET    | /admin/overview               | Get system metrics overview |
| Admin Service             | GET    | /admin/rides/active           | Get list of active rides    |

### WebSocket Connections

| Service                   | WebSocket URL                            | Purpose                             |
| ------------------------- | ---------------------------------------- | ----------------------------------- |
| Ride Service              | ws://{host}/ws/passengers/{passenger_id} | WebSocket connection for passengers |
| Driver & Location Service | ws://{host}/ws/drivers/{driver_id}       | WebSocket connection for drivers    |

> Each service maintains its own WebSocket connections for real-time updates. Passengers receive ride status updates, and drivers receive ride assignment and location updates.

---

## Data Model

* **Ride**: `id`, `riderId`, `driverId`, `status`, `startLocation`, `endLocation`, `fare`, `timestamps`
* **User**: `id`, `name`, `role`, `email`, `password`
* **Driver**: `userId`, `status`, `location`
* **Coordinate**: `rideId`, `latitude`, `longitude`, `timestamp`

---

## Error Handling

* `400` — Invalid input
* `401` — Unauthorized
* `403` — Forbidden
* `500` — Internal server error
* Async retries via RabbitMQ for critical failures

---

## Deployment

* Build Docker image: `docker build -t ride-hail .`
* Push to registry and deploy via Kubernetes or Docker Compose
* Run DB migrations before deploying

---

## Monitoring & Logging

* Structured logs with `pkg/logger/logger.go`
* Metrics and tracing can be added via Prometheus / OpenTelemetry
* Critical errors alert via notification system

---

## ERD (Entity Relationship Diagram)

The database structure for the Ride-Hail module is represented in the ERD diagram:

[View ERD diagram (PDF)](erd-diagram.pdf)

> This diagram shows the relationships between the main entities:
>
> * **users** — stores user information
> * **drivers** — linked to users, stores driver-specific info like status and location
> * **rides** — stores ride requests and their lifecycle
> * **coordinates** — stores ride tracking coordinates

---

## Contributing

* Follow Git workflow with feature branches
* Write clean, maintainable code
* Follow Go formatting and linting rules

---

## License

This project is licensed under **MIT License**.
