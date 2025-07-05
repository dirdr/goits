![Gitleaks](https://github.com/dirdr/goits/actions/workflows/gitleaks.yaml/badge.svg)
![Tests](https://github.com/dirdr/goits/actions/workflows/tests.yaml/badge.svg)

# Goits: Mock Internal Transfers System 🏦

Goits is a simple, mocked internal transfers service built with Go.

## Assumptions

- **Single Currency:** All accounts operate under the same currency. No multi-currency support is implemented.
- **No Authentication/Authorization:** The API endpoints are publicly accessible without any authentication or authorization mechanisms. : _This imply that this service is not hosted and that you will need to run it locally_

## Technical Requirements

working with money trasnfert imopsent the use of certain paradigms. This project use various mechanism to ensure data consistency and integrity

1. Event Sourcing (as a source of truth for transactions)
2. Double bookkeeping to check at all time that $T_{credit} = T_{debit}$
3. Optimistic locking of the account_balances to prevent lost update while not taking lock for too long

## Getting Started

### Prerequisites

Docker and Compose plugin.

### Run

#### Using Docker Compose

For easier setup and management of the application and its dependencies (like PostgreSQL) Docker compose is encouraged

1. Copy `.env.example` to `.env`
2. Fill env variables
3. **Start the services:**

   ```sh
   docker compose up -d --build
   ```

### Test

You can run business rules unit test with go tests:

```sh
go test ./test/... -v

```

### API Endpoints

A Swagger book is available to see API description and play with endpoints, head down to [Swager](http://localhost:8080/swagger/index.html#/) (Replace the port by the one you set in the `.env` file)
