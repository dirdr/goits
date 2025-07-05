# Goits: Mock Internal Transfers System 🏦

Goits is a simple mocked internal transfers application built with Go.

## Assumptions

- **Single Currency:** All accounts operate under the same currency. No multi-currency support is implemented.
- **No Authentication/Authorization:** The API endpoints are publicly accessible without any authentication or authorization mechanisms. : _This imply that this service is not hosted and that you will need to run it locally_

## Getting Started

### Prerequisites

Go and a PgSQL @17 Instance running (either local or container).

### Run

#### Env

#### Using Docker Compose

For easier setup and management of the application and its dependencies (like PostgreSQL) Docker compose is encouraged

1. Copy `.env.example` to `.env`

2. **Start the services:**

   ```sh
   docker compose up -d --build
   ```

### API Endpoints

API endpoints are available in the swagger instance created with this micro-service. Find it at `http://localhost:<port>/swagger/index.html`
