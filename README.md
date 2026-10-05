# Go Starter Kit

A production-ready Go REST API starter kit based on **Echo v4**, **PostgreSQL**, **sqlc**, and **JWT**

## 🏗️ Project Structure

```
personal-finance/
├── main.go                         # Entry point
├── go.mod / go.sum
├── .env.example                    # Environment variable template
├── Makefile                        # Common dev commands
├── Dockerfile                      # Multi-stage production build
├── sqlc.yaml                       # sqlc code generation config
│
├── config/
│   ├── type.go                     # All config structs (App, DB, JWT, SMTP…)
│   └── config.go                   # InitConfig, DB init, email helpers
│
├── utils/
│   ├── cors.go                     # CORS middleware config
│   ├── validator.go                # NetworkAPIError, NewError helpers
│   ├── response.go                 # ResponseOK / ResponseError generics
│   ├── token.go                    # JWT generate & verify
│   ├── password.go                 # bcrypt hash & compare
│   ├── regex.go                    # Email regex helper
│   ├── middleware.go               # JWTMiddleware + GetUserID
│   ├── money.go                    # NUMERIC-safe amount parsing & arithmetic
│   ├── date.go                     # Date/month parsing, month boundaries
│   └── errors.go                   # Postgres constraint-violation helpers
│
├── app/
│   ├── app.go                      # Echo bootstrap, middleware, route groups, cron
│   │
│   ├── auth/                       # Authentication module
│   │   ├── router.go               # Route registration
│   │   ├── usecase.go              # Login, Register, ForgotPassword, ResetPassword…
│   │   ├── validator.go            # Request body validation
│   │   └── repository/
│   │       ├── db.go               # DBTX interface, Queries, WithTx
│   │       ├── models.go           # DB model structs
│   │       ├── queries.sql         # Raw SQL queries (sqlc source)
│   │       ├── queries.sql.go      # sqlc-generated query implementations
│   │       └── repository.go       # Repository interface + implementation
│   │
│   ├── user/                       # User management module
│   │   ├── router.go
│   │   ├── usecase.go              # Me, UpdateMe, GetAll, GetByID, Update, Delete
│   │   ├── validator.go
│   │   └── repository/ ...
│   │
│   ├── account/                    # Wallets, bank accounts, cards
│   ├── category/                   # Income/expense/transfer categories
│   ├── transaction/                # Transactions + atomic balance updates
│   ├── budget/                     # Monthly budget per category
│   └── report/                     # Dashboard, summaries, trends
│
└── database/
    ├── tables/
    │   ├── 00_types.sql            # account_type / category_type enums
    │   └── db.sql                  # Tables and indexes
    ├── functions/
    │   └── func.sql                # set_updated_at trigger + triggers
    └── migrations/
        ├── 001_finance_features.sql       # Idempotent upgrade for existing DBs
        └── 002_seed_default_categories.sql
```

Every feature module follows the same shape: `router.go`, `usecase.go`,
`validator.go` and `repository/` (`queries.sql` + sqlc output + `repository.go`).

## 🚀 Quick Start

### 1. Clone & setup environment
```bash
cp .env.example .env
# Edit .env with your DB, JWT, and SMTP settings
```

### 2. Create the database
```bash
psql -U postgres -c "CREATE DATABASE personal_finance;"

# Fresh install: enum types, then tables, then triggers
psql -U postgres -d personal_finance -f database/tables/00_types.sql
psql -U postgres -d personal_finance -f database/tables/db.sql
psql -U postgres -d personal_finance -f database/functions/func.sql

# Global categories every user can pick from
psql -U postgres -d personal_finance -f database/migrations/002_seed_default_categories.sql
```

Already have a database from an earlier version? Run the idempotent upgrade
instead of the files above:

```bash
psql -U postgres -d personal_finance -f database/migrations/001_finance_features.sql
psql -U postgres -d personal_finance -f database/migrations/002_seed_default_categories.sql
```

### 3. Install dependencies & run
```bash
go mod tidy
go run main.go
```

The server starts on `http://localhost:9000`.

---

## 📌 API Endpoints

Everything except `/api/auth/*` and `/api/health` requires
`Authorization: Bearer <token>`. The user is resolved from the token, so no
endpoint takes a user ID — a request can only ever read or write its own data.

### Auth — `/api/auth` (public)
| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/auth/login` | Login, returns JWT token |
| `POST` | `/api/auth/register` | Register a new user |
| `POST` | `/api/auth/forgot-password` | Send password reset email |
| `POST` | `/api/auth/reset-password` | Reset password with token |

### Health — `/api/health` (public)
| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/health` | Liveness probe |

### Users — `/api/users`
| Method | Path | Description |
|--------|------|-------------|
| `GET`    | `/api/users/me` | Profile of the signed-in user |
| `PUT`    | `/api/users/me` | Update own username |
| `GET`    | `/api/users` | List all active users |
| `GET`    | `/api/users/:id` | Get user by ID |
| `PUT`    | `/api/users/:id` | Update user name |
| `DELETE` | `/api/users/:id` | Soft-delete (deactivate) user |

### Accounts — `/api/accounts`
| Method | Path | Description |
|--------|------|-------------|
| `GET`    | `/api/accounts` | List accounts plus total balance. `?includeArchived=true` to include archived ones |
| `POST`   | `/api/accounts` | Create an account with an optional opening balance |
| `GET`    | `/api/accounts/:id` | Get one account |
| `PUT`    | `/api/accounts/:id` | Update name, type and currency |
| `DELETE` | `/api/accounts/:id` | Delete — refused with 409 while transactions exist |
| `PATCH`  | `/api/accounts/:id/archive` | Hide the account, keeping its history |
| `PATCH`  | `/api/accounts/:id/unarchive` | Restore an archived account |

Account `type` is one of `cash`, `bank`, `e_wallet`, `credit_card`,
`investment`. `balance` is maintained by the transaction endpoints and is not
editable through `PUT`.

### Categories — `/api/categories`
| Method | Path | Description |
|--------|------|-------------|
| `GET`    | `/api/categories` | Own categories plus the global ones. `?type=income\|expense\|transfer` |
| `POST`   | `/api/categories` | Create a category |
| `GET`    | `/api/categories/:id` | Get one category |
| `PUT`    | `/api/categories/:id` | Update — 403 for a global category |
| `DELETE` | `/api/categories/:id` | Delete — 403 for a global one, 409 while in use |

Global categories (`userId: null`, flagged `isGlobal: true`) are seeded by
`002_seed_default_categories.sql` and are read-only for every user.

### Transactions — `/api/transactions`
| Method | Path | Description |
|--------|------|-------------|
| `GET`    | `/api/transactions` | Paginated, filtered list |
| `POST`   | `/api/transactions` | Record a transaction and apply it to balances |
| `GET`    | `/api/transactions/:id` | Get one transaction, accounts and category joined |
| `PUT`    | `/api/transactions/:id` | Rewrite a transaction, rebalancing both versions |
| `DELETE` | `/api/transactions/:id` | Delete and reverse its balance effect |

List query params: `page`, `pageSize` (max 100), `startDate`, `endDate`
(`YYYY-MM-DD`), `accountId`, `categoryId`, `type`, `search`.

The transaction type comes from the category, which determines what happens to
the balances:

| Category type | Effect |
|---------------|--------|
| `income` | credits `accountId` |
| `expense` | debits `accountId` |
| `transfer` | debits `accountId`, credits `toAccountId` |

`toAccountId` is required for `transfer` and rejected for the other two. The
write and its balance adjustments share one database transaction, so a balance
can never drift from the transactions behind it.

### Budgets — `/api/budgets`
| Method | Path | Description |
|--------|------|-------------|
| `GET`    | `/api/budgets` | Budgets for a month with amount spent per category. `?period=YYYY-MM` (default: current month) |
| `POST`   | `/api/budgets` | Set a budget. Re-posting the same category and month replaces the amount |
| `PUT`    | `/api/budgets/:id` | Change the amount |
| `DELETE` | `/api/budgets/:id` | Remove a budget |

### Reports — `/api/reports`
| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/reports/dashboard` | Everything the dashboard needs in one request |
| `GET` | `/api/reports/summary` | Income, expense and net for a range |
| `GET` | `/api/reports/by-category` | Per-category breakdown. `?type=` (default `expense`) |
| `GET` | `/api/reports/trend` | Income and expense per month. `?months=` (default 6, max 24) |
| `GET` | `/api/reports/balances` | Balance of every active account |

`summary` and `by-category` take `?startDate=` and `?endDate=` and default to
the current calendar month. Transfers are excluded from income and expense
totals, since moving money between your own accounts is neither.

Amounts are returned as **strings**, not numbers: they come from
`NUMERIC(15,2)` columns and pass through the API without a float64 ever
touching them.

---

## 🔧 Makefile Commands

| Command | Description |
|---------|-------------|
| `make run` | Run the server locally |
| `make build` | Build binary to `bin/main` |
| `make test` | Run all tests with coverage |
| `make tidy` | Run `go mod tidy` |
| `make sqlc` | Regenerate sqlc query files |
| `make init` | Install sqlc if not present |
| `make docker-build` | Build Docker image |
| `make docker-run` | Run Docker container with `.env` |

---

## ➕ Adding a New Module

Follow this pattern for every new feature (e.g., `product`):

```
app/
└── product/
    ├── router.go        # RegisterRoutes(g *echo.Group)
    ├── usecase.go       # Usecase struct + handler methods
    ├── validator.go     # Request structs + validation funcs
    └── repository/
        ├── db.go        # DBTX, Queries, WithTx
        ├── models.go    # DB model structs
        ├── queries.sql  # Raw SQL (sqlc source)
        ├── queries.sql.go  # Generated by sqlc
        └── repository.go   # Interface + implementation
```

Then register in `app/app.go`:
```go
product.RegisterRoutes(API.Group("/products"))
```

And add to `sqlc.yaml`:
```yaml
- name: 'repository'
  path: 'app/product/repository'
  queries: 'app/product/repository/queries.sql'
  schema: 'database/tables/'
  engine: 'postgresql'
  ...
```

---

## 🛡️ Response Format

### Success
```json
{
  "status": 200,
  "data": { ... }
}
```

### Error
```json
{
  "errors": [
    {
      "msg": "Email tidak valid atau tidak diisi",
      "location": "body",
      "path": "email",
      "type": "required"
    }
  ]
}
```

---

## 🐳 Docker

```bash
# Build
make docker-build

# Run
make docker-run
```

---

## 🧪 Testing

```bash
make test
```

Tests live alongside their module:
- `app/app_test.go` — which routes are public and which demand a valid token
- `app/transaction/usecase_test.go` — balance effects per category type
- `utils/middleware_test.go` — JWT middleware and `GetUserID`
- `utils/money_test.go` — amount parsing and decimal arithmetic

These run without a database: they cover the request boundary and the money
arithmetic, which are the parts where a mistake is silent.
