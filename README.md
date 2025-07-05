# Goits: Mock Internal Transfers System 🏦

Goits is a simple, mocked internal transfers service built with Go. It is designed to ensure data consistency and integrity.

![Gitleaks](https://github.com/dirdr/goits/actions/workflows/gitleaks.yaml/badge.svg)
![Tests](https://github.com/dirdr/goits/actions/workflows/tests.yaml/badge.svg)

## Table of Contents

- [Assumptions](#assumptions)
- [Technical Requirements](#technical-requirements)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Run](#run)
  - [Using Docker Compose](#using-docker-compose)
  - [Test](#test)
- [API Endpoints](#api-endpoints)

## Assumptions 🧑‍🔬

- **Single Currency:** All accounts operate under the same currency.
- **No Authentication/Authorization:** The API endpoints are publicly accessible without any authentication or authorization mechanisms.

> [!WARNING]
> Given these assumptions and the fact that this project is a simple test service, do not use it as a real money transfer system!

## Technical Requirements ⚙️

Working with money transfers necessitates the use of certain paradigms. This project uses various mechanisms to ensure data consistency and integrity:

1. Event Sourcing (as a source of truth for transactions)
2. Double-entry bookkeeping to check at all times that $T_{credit} = T_{debit}$
3. Optimistic locking inside `account_balances` projection to prevent lost updates while not holding locks for too long

## Getting Started 🚀

### Prerequisites

Docker and Compose plugin.

### Dowload

Clone the repository and cd into it :

```sh
git clone git@github.com:dirdr/goits.git && cd goits
```

### Run

For easier setup and management of the application and its dependencies (like PostgreSQL), Docker Compose is encouraged:

1. Copy `.env.example` to `.env`
2. Fill in environment variables
3. **Start the services:**

   ```sh
   docker compose up -d --build
   ```

## API Endpoints 🗾

Swagger documentation is available to view API descriptions and interact with endpoints. Navigate to [Swagger](http://localhost:8080/swagger/index.html#/) (Replace the port with the one you set in the `.env` file).

## Test ✅

You can run business rule unit tests with Go tests:

```sh
go test ./test/... -v
```

## Futur improvements 📈

- Add Authn and Authz
- Deploy the service on my VPS
- Load test concurrent transactions
