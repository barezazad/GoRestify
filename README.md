# GoRestify — Go REST API Boilerplate

<p align="center">
  <img src="./GoRestify-logo.jpg" alt="GoRestify logo" width="200" />
</p>

<p align="center">
  <strong>Production-oriented Golang backend starter</strong><br/>
  Gin · GORM · MySQL · Redis · JWT RS256 · Wire · RBAC
</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/BarezAzad/GoRestify"><img src="https://pkg.go.dev/badge/github.com/BarezAzad/GoRestify.svg" alt="Go Reference"></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/license-MIT-brightgreen.svg" alt="MIT license"></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go" alt="Go version"></a>
  <a href="https://gin-gonic.com/"><img src="https://img.shields.io/badge/Gin-HTTP-00ADD8" alt="Gin"></a>
  <a href="https://gorm.io/"><img src="https://img.shields.io/badge/GORM-ORM-00ADD8" alt="GORM"></a>
  <a href="https://redis.io/"><img src="https://img.shields.io/badge/Redis-cache-DC382D?logo=redis" alt="Redis"></a>
  <a href="https://www.mysql.com/"><img src="https://img.shields.io/badge/MySQL-8-4479A1?logo=mysql" alt="MySQL"></a>
</p>

---

## Table of contents

1. [What is GoRestify?](#what-is-gorestify)
2. [What makes it different](#what-makes-it-different)
3. [Architecture](#architecture)
4. [Quick start](#quick-start)
5. [Environment & secrets](#environment--secrets)
6. [Authentication & security](#authentication--security)
7. [RBAC (resources)](#rbac-resources)
8. [Building APIs the GoRestify way](#building-apis-the-gorestify-way)
9. [List / GetAll / filters / pagination](#list--getall--filters--pagination)
10. [Validation (`bind` tags)](#validation-bind-tags)
11. [Errors & error codes](#errors--error-codes)
12. [Responses & activity log](#responses--activity-log)
13. [Redis caching](#redis-caching)
14. [Transactions (`tx`)](#transactions-tx)
15. [Packages reference (`pkg/`)](#packages-reference-pkg)
16. [Wire DI](#wire-di)
17. [i18n, Excel, email, decimals](#i18n-excel-email-decimals)
18. [Docker & CI](#docker--ci)
19. [Security checklist](#security-checklist)
20. [Troubleshooting](#troubleshooting)

---

## What is GoRestify?

**GoRestify** is an open-source **Go REST API boilerplate** for real backend work: layered domain code, dual **admin** / **user** apps, JWT auth, Redis cache, MySQL via GORM, validation, RBAC, activity audit, and a large reusable `pkg/` toolkit.

Clone it, change the domain, ship APIs without rebuilding auth, filters, errors, and wiring from scratch.

**Keywords:** golang rest boilerplate, gin gorm mysql redis jwt, wire di, backend starter kit

---

## What makes it different

Most Go boilerplates stop at “hello Gin + GORM”. GoRestify is opinionated about **backend day-to-day**:

| Area          | Typical starter       | GoRestify                                                  |
| ------------- | --------------------- | ---------------------------------------------------------- |
| Auth          | HS256 secret in env   | **RS256** PEM keys, access + **refresh rotation** in Redis |
| Login abuse   | Often none            | **Progressive rate limit** + admin unblock APIs            |
| List filters  | String-concat SQL     | **Parameterized** DSL + column/operator allow-lists        |
| ORDER BY      | Raw query params      | Regex allow-list (`column` / `table.column`)               |
| Pagination    | Unbounded             | `page_size` **max 100**                                    |
| Get-all       | Huge `LIMIT` hack     | Real repo **`GetAll()`** (no filter/order/limit)           |
| Passwords     | Weak / cost 10        | Policy (10+ + complexity) + bcrypt **cost 12** + salt      |
| Errors        | Plain strings         | Structured `pkg_err` + unique **`E#######`** codes         |
| Permissions   | Role string only      | **Resource RBAC** (`city:read`, `role:write`, …)           |
| Apps          | One binary            | **Admin + user** apps, shared domain                       |
| Cache clear   | `KEYS` (blocks Redis) | **`SCAN`**-based pattern delete                            |
| Docs download | Path join only        | `filepath.Base` + root containment                         |

---

## Architecture

```text
HTTP → Middleware (CORS, JWT, rate limit, logger)
    → API (Gin handlers, bind, response)
      → Service (validation, cache, business rules)
        → Repo (GORM queries)
          → MySQL / Redis
```

```text
cmd/
  admin/     # full CRUD, migrate, seed, settings, activities
  user/      # user-facing API surface
  wire/      # Google Wire providers
domain/
  base/      # accounts, users, roles, cities, regions, documents, auth
  acc/       # currencies, transactions, slots, credits
  service/   # shared services
internal/core/   # Engine, envs, check_access, enums
pkg/             # reusable libraries (see below)
assets/terms/    # i18n TOML
secrets/         # JWT PEMs (gitignored)
```

| App   | Base URL        |
| ----- | --------------- |
| Admin | `/api/admin/v1` |
| User  | `/api/user/v1`  |

---

## Quick start

### Requirements

- Go **1.25+**
- MySQL 8+
- Redis 6+
- Optional: [reflex](https://github.com/cespare/reflex) (`make admin` / `make user`)

### 1) MySQL + Redis

```bash
docker run --rm --name db-mysql -d \
  -v mysql-data:/var/lib/mysql \
  -e MYSQL_ROOT_PASSWORD=change-me \
  -e TZ='Asia/Baghdad' \
  -p 3306:3306 mysql:8 \
  --character-set-server=utf8mb4 \
  --collation-server=utf8mb4_unicode_ci

docker run --rm --name db-redis -d \
  -v redis-data:/data \
  -p 6379:6379 redis:latest
```

```sql
CREATE DATABASE go_restify CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'travis'@'%' IDENTIFIED BY '';
GRANT ALL PRIVILEGES ON *.* TO 'travis'@'%' WITH GRANT OPTION;
FLUSH PRIVILEGES;
```

> Local `sample.env` expects user `travis` with empty password. Use a real password in shared/prod environments.

### 2) JWT keys

```bash
mkdir -p secrets
openssl genrsa -out secrets/jwt_private.pem 2048
openssl rsa -in secrets/jwt_private.pem -pubout -out secrets/jwt_public.pem
```

### 3) Run

```bash
source cmd/admin/sample.env && go run cmd/admin/main.admin.go
# or
make admin
make user
```

With `CORE_AUTO_MIGRATE=true`, admin migrates and seeds.

| Seed admin | Value         |
| ---------- | ------------- |
| Username   | `admin`       |
| Password   | `admin123Aa!` |

Postman: [`GoRestify.postman_collection.json`](./GoRestify.postman_collection.json)

---

## Environment & secrets

From `cmd/admin/sample.env` / `cmd/user/sample.env`:

| Variable                     | Required | Purpose                                          |
| ---------------------------- | -------- | ------------------------------------------------ |
| `CORE_PORT`                  | yes      | HTTP port (default `6969`)                       |
| `GIN_MODE`                   | yes      | `debug` / `release`                              |
| `CORE_DATABASE_DATA_DSN`     | yes      | MySQL DSN                                        |
| `CORE_DATABASE_ACTIVITY_DSN` | admin    | Activity DB                                      |
| `CORE_REDIS_CACHE_API`       | yes      | Redis URL                                        |
| `CORE_JWT_PRIVATE_KEY_PATH`  | yes      | Sign tokens                                      |
| `CORE_JWT_PUBLIC_KEY_PATH`   | yes      | Verify tokens                                    |
| `CORE_JWT_ACCESS_MINUTES`    | yes      | Access TTL                                       |
| `CORE_JWT_REFRESH_HOURS`     | yes      | Refresh TTL                                      |
| `CORE_PASSWORD_SALT`         | yes      | bcrypt pepper — **never reuse defaults in prod** |
| `CORE_AUTO_MIGRATE`          | admin    | Migrate + seed on boot                           |

`secrets/*.pem` is gitignored. Generate keys locally (and in CI/CD secrets), never commit them.

---

## Authentication & security

### Login / refresh

```http
POST /api/admin/v1/login
Content-Type: application/json

{ "username": "admin", "password": "admin123Aa!" }
```

```http
Authorization: Bearer <access_token>
```

```http
POST /api/admin/v1/refresh-token
{ "refresh_token": "<refresh_token>" }
```

**Backend notes**

- Access tokens are **RS256**, short-lived; refresh tokens are stored/rotated in Redis
- `JwtAuthGuard` requires scheme `Bearer` (case-insensitive) and rejects refresh tokens on normal routes
- Login uses **uniform** failure messages (no user enumeration)
- **Rate limit** on login/refresh; admin can inspect/unblock:  
  `GET /login-blocks`, `POST /login-blocks/unblock`
- Password policy (`bind:"password"`): **≥10** chars, lower, upper, digit, **special**
- Hashing: bcrypt **cost 12** + `CORE_PASSWORD_SALT`

### Basic auth (optional)

`BasicAuthGuard` reads `X-Authorization: Basic ...` with **constant-time** compare. Not mounted by default — only attach for third-party integrations.

---

## RBAC (resources)

Protected routes:

```go
rg.GET("/cities", check_access.IsAllow(base.CityRead), cityAPI.List)
rg.POST("/cities", check_access.IsAllow(base.CityWrite), cityAPI.Create)
```

Resources are string constants, e.g. `city:read`, `setting:write` (`domain/base/base_resource.go`, `domain/acc/...`).

Roles hold a list of resources. After JWT auth, `check_access` loads the user’s permissions and returns **403** if missing.

**When adding a feature:** define Read/Write resources, assign them to roles, and wrap every route.

---

## Building APIs the GoRestify way

### Layer checklist

1. **Model** — `domain/<area>/<area>_model/` + GORM tags + `bind` tags
2. **Repo** — `FindByID`, `GetAll`, `List`, `Count`, `Create`, `Save`, `Delete`
3. **Service** — `ValidateModel`, Redis get/set/delete, business rules
4. **API** — `response.NewParam`, `Bind`, `Record`, status/message/JSON
5. **Wire** — add provider in `cmd/wire` and regenerate
6. **Router** — JWT + `IsAllow`
7. **Resource** — new permission constants

### Typical list handler

```go
func (a *CityAPI) List(c *gin.Context) {
    resp, params := response.NewParam(c, base_model.CityTable)
    data := map[string]interface{}{}

    var err error
    if data["list"], data["count"], err = a.Service.List(params); err != nil {
        resp.Error(err).JSON()
        return
    }

    resp.Record(base.ListCity)
    resp.Status(http.StatusOK).
        Message(pkg_terms.ListOfV, base_term.Cities).
        JSON(data)
}
```

### Typical create (with TX)

```go
params.Tx.DB = a.Engine.DB.Begin()
defer func() {
    if err != nil {
        params.Tx.DB.Rollback()
    }
}()

if created, err = a.Service.Create(params.Tx, city); err != nil {
    resp.Error(err).JSON()
    return
}
params.Tx.DB.Commit()
```

### Repo query rules

```go
whereStr, whereArgs, err := params.ParseWhere(r.Cols)
// soft-delete tables:
whereStr, whereArgs, err := params.ParseWhereDelete(r.Cols)

db.Where(whereStr, whereArgs...).Order(params.Order).Limit(params.Limit).Offset(params.Offset)
```

**Never** concatenate user input into SQL. Use `?` + args. `ParseWhereDelete` wraps filters so `OR` cannot bypass `deleted_at IS NULL`.

### Service preconditions

Lock list scope in code (not from the client):

```go
_ = params.SetPreCondition("base_cities.region_id = ?", regionID)
_ = params.AddPreCondition("base_cities.status = ?", status)
```

Must use `?` placeholders; unsafe tokens are rejected.

---

## List / GetAll / filters / pagination

| Endpoint style    | Service/Repo | Behavior                                             |
| ----------------- | ------------ | ---------------------------------------------------- |
| `GET /cities`     | `List`       | Filter + order + page                                |
| `GET /all/cities` | `GetAll`     | All rows, **no** filter/order/limit (cache-friendly) |

### Query params (`List`)

| Param       | Default      | Notes                          |
| ----------- | ------------ | ------------------------------ |
| `page`      | `1`          |                                |
| `page_size` | `10`         | **Capped at 100**              |
| `order_by`  | `{table}.id` | Only `col` or `table.col`      |
| `direction` | `desc`       | Only `asc` / `desc`            |
| `select`    | `*`          | Validated against repo columns |
| `filter`    | empty        | DSL below                      |

### Filter DSL

```text
filter=name[eq]"Erbil"[and]id[gt]"1"
filter=name[like]"%erb%"[or]id[eq]"2"
```

| Token                             | Meaning      |
| --------------------------------- | ------------ |
| `[eq] [ne] [gt] [lt] [gte] [lte]` | Comparisons  |
| `[like]`                          | LIKE         |
| `[and] [or]`                      | Logic        |
| `[date] [date_gte] [date_lte]`    | Date helpers |

Filters are parsed into `WHERE ... ?` with bound args. Invalid columns/operators/values → validation error (not silent SQL).

---

## Validation (`bind` tags)

Call in services:

```go
err = validator.ValidateModel(model, base_term.City, validator.Create)
err = validator.ValidateModel(model, base_term.City, validator.Update)
```

### Rules

| Tag                                    | Meaning                           |
| -------------------------------------- | --------------------------------- |
| `required`                             | Must be present / non-zero        |
| `min` / `max`                          | Length                            |
| `gte` / `lte`                          | Numeric compare                   |
| `email`                                | Email format                      |
| `phone`                                | `964` + 10 digits                 |
| `username`                             | `[a-zA-Z0-9._-]+`                 |
| `password`                             | ≥10, lower, upper, digit, special |
| `one_of=account_status`                | Enum from `MustBeInTypes`         |
| `contain=x`                            | Substring must exist              |
| `birthday`                             | Age window                        |
| `pin`                                  | Digits only                       |
| `if_exist`                             | Validate only if field is set     |
| `create:required\|update:min=7,max=10` | Per-action rules                  |

### Enums

Register lists in `internal/core/action` (`MustBeInTypes`), then:

```go
Status pkg_types.Enum `bind:"one_of=account_status"`
```

### Model tips

```go
Password string `json:"password,omitempty" bind:"if_exist,password"`
Phone    string `json:"phone,omitempty" bind:"create:required|update:min=7,max=10"`
Region   string `gorm:"->;migration:-" json:"region,omitempty" table:"base_regions.name as region"`
```

- `table:"..."` helps joined select aliases
- Omitempty passwords on responses; clear hash fields in `GetAll`/`List` when needed

---

## Errors & error codes

### Creating errors

```go
err = pkg_err.New(pkg_err.SomethingWentWrong, "E1171379").
    Message(pkg_err.SomethingWentWrong).
    Custom(pkg_err.InternalServerErr).
    Build()

err = pkg_err.Take(dbErr, "E1174655").
    Custom(pkg_err.ValidationFailedErr).
    Build()
```

Always pass a unique `"E"` + **7 digits** code. On boot, `utils.GenerateErrCode()` scans `.go` files and **bumps duplicates** (`E1146731` → `E1146732` …). No `err_codes.txt` needed.

### DB errors

In repos, always:

```go
err = db_error.Parse(err, base_term.Cities, validator.Update)
```

Maps duplicates / not found / etc. into API-friendly errors.

### Common customs

`ValidationFailedErr`, `UnauthorizedErr`, `ForbiddenErr`, `NotFoundErr`, `InternalServerErr`, `BadRequestErr`, …

---

## Responses & activity log

```go
resp, params := response.NewParam(c, base_model.CityTable)
resp.Record(base.CreateCity)          // async activity (when enabled)
resp.Status(http.StatusOK).
    Message(pkg_terms.VCreatedSuccessfully, base_term.City).
    JSON(createdCity)

resp.Error(err).Abort().JSON()        // middleware-style abort
```

Activities: admin `GET /activities` (permission-gated).

Messages go through **dictionary** translation based on request language.

---

## Redis caching

```go
key := fmt.Sprintf("%v-%v", base_term.City, id)
if ok := s.Engine.RedisCacheAPI.GetCache(tx, key, &city); ok {
    return
}
city, err = s.Repo.FindByID(tx, id)
_ = s.Engine.RedisCacheAPI.Set(key, city)

// on write
s.Engine.RedisCacheAPI.Delete(key)
s.Engine.RedisCacheAPI.Delete(base_term.Cities) // GetAll list key
```

Admin:

- `PUT /clear-cache/:key`
- `PUT /clear-cache/user/:userID`

Pattern delete uses **SCAN** (not blocking `KEYS`). `FlushDB` is for seed/dev only.

---

## Transactions (`tx`)

```go
// no TX
service.Create(tx.Tx{}, model)

// with TX
params.Tx.DB = engine.DB.Begin()
service.Create(params.Tx, model)
params.Tx.DB.Commit() // or Rollback on error
```

Repos use `tx.GetDB(r.Engine.DB)` so the same code path works with or without an outer transaction.

---

## Packages reference (`pkg/`)

| Package                | What backend devs use it for                                      |
| ---------------------- | ----------------------------------------------------------------- |
| `pkg_config`           | Global config after `pkg.Init`                                    |
| `pkg_sql`              | MySQL connect, filter parser, FK SQL, column extract              |
| `pkg_redis`            | Cache, TTL, pattern reset, refresh-token storage                  |
| `pkg_jwt`              | RS256 issue / parse / refresh                                     |
| `pkg_password`         | Hash / verify                                                     |
| `param`                | Query → `Param`; `ParseWhere` / `ParseWhereDelete`; preconditions |
| `validator`            | `ValidateModel` + `bind` tags                                     |
| `middleware`           | JWT, Basic auth, API logger, login rate limit, domain header      |
| `rate_limiter`         | Progressive lockout store                                         |
| `response`             | JSON envelope, bind helper, activity `Record`                     |
| `pkg_err`              | Typed API errors + codes                                          |
| `db_error`             | GORM/MySQL → API errors                                           |
| `pkg_log`              | Fatal / CheckError; redacted HTTP dumps                           |
| `pkg_http`             | Outbound HTTP (JSON / multipart)                                  |
| `tx`                   | Transaction wrapper                                               |
| `excel`                | Fluent XLSX export                                                |
| `decimal`              | Money-safe arithmetic                                             |
| `dictionary`           | TOML i18n                                                         |
| `activity`             | Audit list/count                                                  |
| `setting`              | DB settings + reload ticker                                       |
| `utils`                | `crypto/rand` strings, AES-GCM, email, `GenerateErrCode`          |
| `pkg_types`            | `Enum`, `Resource`, env helpers                                   |
| `pkg_terms` / `*_term` | Message keys                                                      |

---

## Wire DI

Providers live in `cmd/wire/wire.go` (build tag `wireinject`). Generated code: `wire_gen.go`.

```go
func InitBaseCityAPI(e *core.Engine) base_api.CityAPI {
    wire.Build(
        base_repo.ProvideCityRepo,
        service.ProvideBaseCityService,
        base_api.ProvideCityAPI,
    )
    return base_api.CityAPI{}
}
```

After changing providers:

```bash
cd cmd/wire && wire
```

Routers call `wire.Init...API(engine)` — don’t construct repos/services manually in handlers.

---

## i18n, Excel, email, decimals

### Dictionary

Terms: `assets/terms/terms.toml`

```go
dictionary.Translate(dictionary.Ku, city.Name)
```

Language from request headers / context (`X-LANGUAGE` pattern used by HTTP client too).

### Excel export

```go
ex := excel.New("region").
    AddSheet("Regions").
    Active("Regions").
    WriteHeader("ID", "Name", "Created At").
    SetSheetFields("ID", "Name", "CreatedAt").
    WriteData(regions).
    AddTable()

buf, name, err := ex.Generate()
c.Data(http.StatusOK, "application/octet-stream", buf.Bytes())
```

### Email

```go
engine.EmailConfig = utils.ConfigEmail{ Host, Port, Username, Password }
go engine.EmailConfig.SendEmail(to, cc, fromAlias, subject, body, attachmentPath)
```

### Decimals (money)

Prefer `pkg/decimal` over `float64` for fees, balances, and tax.

---

## Contributing

Stars help others discover this boilerplate. Issues/PRs welcome.

**Topics:** `go` `golang` `gin` `gorm` `redis` `mysql` `jwt` `rest-api` `boilerplate` `wire`

---

## License

MIT — see [`LICENCE`](./LICENCE).

**Author:** [Barez Azad](https://github.com/barezazad)
