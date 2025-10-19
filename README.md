# Gin GORM API Starter Template

## Contents

- [Description](#description)
- [Architecture](#architecture)

  - [Explanation (EN)](#explanation-en)
  - [Explanation (ID)](#explanation-id)

- [Pre-requisites](#pre-requisites)

  - [PostgreSQL Requirements](#postgresql-requirements)
  - [GitHooks Requirements](#githooks-requirements)

- [How to Run?](#how-to-run)
- [API Documentation (Postman)](#api-documentation-postman)

## Description

An API starter template for projects based on Controller-Service-Repository (CSR) Pattern utilizing Gin (Golang) and PostgreSQL, with GORM as the ORM.

## Architecture

```
/backend
│
├── /api
│   └── /v1
│       ├── /controller
│       │   ├── user.go
│       │   └── etc
│       └── /router
│           ├── user.go
│           └── etc
│
├── /config
│   ├── db.go
│   └── etc
│
├── /core
│   ├── /entity
│   │   └── user.go
│   │   └── etc
│   ├── /helper
│   │   ├── /dto
│   │   │   ├── user.go
│   │   │   └── etc
│   │   ├── /errors
│   │   │   ├── other.go
│   │   │   ├── user.go
│   │   │   └── etc
│   │   └── /messages
│   │       ├── file.go
│   │       ├── user.go
│   │       └── etc
│   ├── /interface
│   │   └── /query
│   │       ├── user.go
│   │       └── etc
│   │   └── /repository
│   │       ├── tx.go
│   │       ├── user.go
│   │       └── etc
│   └── /service
│       ├── user.go
│       └── etc
│
├── /database
│   ├── /seeder
│   │   └── user.go
│   │   └── etc
│   └── migrator.go
│
├── /infrastructure
│   ├── /query
│   │   └── user.go
│   │   └── etc
│   └── /repository
│       └── tx.go
│       └── user.go
│       └── etc
│
├── /provider
│   ├── user.go
│   └── etc
│
├── /support
│   ├── /base
│   │   └── model.go
│   │   └── request.go
│   │   └── response.go
│   │   └── etc
│   ├── /constant
│   │   └── default.go
│   │   └── enums.go
│   │   └── etc
│   ├── /middleware
│   │   └── authentication.go
│   │   └── cors.go
│   │   └── authorization.go
│   │   └── etc
│   └── /util
│       └── bcrypt.go
│       └── file.go
│       └── etc
│
├── /tests
│   ├── /testutil
│   └── /integration
│
└── main.go
```

### Explanation (EN)

- `/api/v1` : The directory for things related to API like all available endpoints (routes) and the handlers for each endpoint (controller). Subdirectory `/v1` is used for easy version control in case of several development phases.

  - `/controller` : The directory for things related to the Controller layer, which handles requests and returns responses.
  - `/router` : The directory for things related to routing. Contains all supported routes/endpoints along with request methods and used middleware.

- `/config` : The directory for things related to program configuration like database configuration.

- `/core` : The directory for things related to the core backend logic. Contains business logic, entities, and database interaction.

  - `/entity` : The directory for entities/models that are mapped to the database via migration.
  - `/helper` : The directory to store items that help backend operations, such as DTOs, error variables, and message constants.

    - `/dto` : Stores DTO (Data Transfer Object) used as placeholders to transfer data for requests and responses.
    - `/errors` : Stores error variables for each entity or other needs.
    - `/messages` : Stores message constants for each entity or feature.

  - `/interface` : The directory for all core interfaces used by the service layer, including contracts for repository and query layers.

    - `/repository` : Interfaces for repository layer (entity CRUD).
    - `/query` : Interfaces for query layer (read-only operations, projections).

  - `/service` : The directory for the Service layer, responsible for application flow and business logic.

- `/database` : The directory for things related to database migrations and seeding.

  - `/seeder` : The directory for database seeders for each entity.

- `/infrastructure` : The directory for implementations of interfaces defined in `/core/interface`, it's the only layer capable of interacting directly with the database.

  - `/repository` : Implementations of repository interfaces, handling CRUD and transactional operations.
  - `/query` : Implementations of query interfaces, handling read-only or optimized queries.

- `/provider` : The directory for dependency injection setup, e.g., Samber/Do providers for wiring services, repositories, and queries.

- `/support` : The directory for common supporting things that are frequently used across the architecture.

  - `/base` : The directory for base structures such as variables, constants, and functions used in other directories. Includes response, request, and model base structures.
  - `/middleware` : The directory for Middlewares, mechanisms that intercept HTTP requests/responses before they are handled by controllers.
  - `/util` : The directory for utility/helper functions that can be used in other directories.

- `/tests` : The directory for automated API testing (unit tests and integration tests).

  - `/testutil` : Stores utility/helper functions for testing purposes.
  - `/integration` : Stores integration test functions.

- `main.go` : The entry point of the application.

### Explanation (ID)

- `/api/v1` : Direktori yang berisi berbagai hal yang berkaitan dengan API seperti daftar endpoint yang disediakan (route) serta handler (controller) dari setiap endpoint. Subdirectory `/v1` digunakan untuk version control apabila ada beberapa versi API.

  - `/controller` : Direktori untuk menyimpan hal-hal terkait Controller, yang bertugas menerima request dan memberikan response.
  - `/router` : Direktori untuk menyimpan hal-hal yang terkait dengan routing, berisi semua route/endpoints yang didukung beserta metode request dan middleware yang digunakan.

- `/config` : Direktori yang berisi hal-hal terkait konfigurasi aplikasi, misalnya konfigurasi database.

- `/core` : Direktori yang berisi bagian inti dari backend. Meliputi business logic, entitas, dan interaksi dengan database.

  - `/entity` : Direktori untuk menyimpan entitas atau model yang digunakan di migrasi dan aplikasi.
  - `/helper` : Direktori untuk menyimpan hal-hal yang membantu operasi backend, seperti DTO, variabel error, dan konstanta pesan.

    - `/dto` : Direktori untuk menyimpan DTO (Data Transfer Object), placeholder untuk memindahkan data request dan response.
    - `/errors` : Direktori untuk menyimpan variabel error untuk setiap entitas maupun kebutuhan lain.
    - `/messages` : Direktori untuk menyimpan konstanta pesan untuk response API.

  - `/interface` : Direktori untuk menyimpan semua interface inti yang digunakan service layer, termasuk kontrak untuk repository dan query.

    - `/repository` : Interface untuk repository (CRUD entitas).
    - `/query` : Interface untuk query (read-only, projections).

  - `/service` : Direktori untuk service layer, yang menangani alur aplikasi dan logika bisnis.

- `/database` : Direktori untuk hal-hal terkait migrasi dan seeding database.

  - `/seeder` : Direktori untuk database seeding tiap entitas.

- `/infrastructure` : Direktori untuk implementasi interface yang ada di `/core/interface`. Merupakan satu-satunya layer yang dapat berinteraksi secara langsung dengan basis data.

  - `/repository` : Implementasi repository, menangani operasi CRUD dan transaksi.
  - `/query` : Implementasi query, menangani operasi read-only atau query yang dioptimalkan.

- `/provider` : Direktori untuk setup dependency injection, misalnya provider Samber/Do untuk menghubungkan service, repository, dan query.

- `/support` : Direktori yang berisi hal-hal umum pembantu untuk digunakan di seluruh project.

  - `/base` : Direktori yang berisi struktur dasar seperti variabel, konstanta, dan fungsi yang digunakan di directory lain. Termasuk response, request, dan model base structure.
  - `/middleware` : Direktori untuk Middleware, mekanisme yang menengahi proses HTTP request/response sebelum ditangani controller.
  - `/util` : Direktori untuk fungsi utilitas/pembantu yang dapat digunakan di berbagai directory.

- `/tests` : Direktori untuk automated API testing (unit dan integration tests).

  - `/testutil` : Menyimpan fungsi utilitas/pembantu untuk testing.
  - `/integration` : Menyimpan fungsi integration testing.

- `main.go` : Titik masuk (entry point) aplikasi.

## Pre-requisites

### PostgreSQL Requirements

1. Create the database in PostgreSQL with the name equal to the value of DB_NAME in `.env`

### GitHooks Requirements

> Note : GitHooks is not mandatory for this starter. Only do the steps below if you want to apply & use it.

1. Install golangci-lint as the linters aggregator for pre-commit linting by executing `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`. Alternatively, you can follow the recommended method, which involves installing the binary from the [official source](https://github.com/golangci/golangci-lint/releases)
2. Install commitlint as the conventional commit message checker by executing `go install github.com/conventionalcommit/commitlint@latest`. Alternatively, you can follow the recommended method, which involves installing the binary from the [official source](https://github.com/conventionalcommit/commitlint/releases)
3. Configure your git's hooks path to be linked to the `.githooks` directory on this repository by executing `git config core.hooksPath .githooks`

## How to Run?

1. Use the command `make tidy` (or use `go mod tidy` instead, if `make` is unable to be used) to adjust the dependencies accordingly
2. Use the command `make run` (or use `go run main.go` instead, if `make` is unable to be used) to run the application. You can also use Docker with air to auto-reload by running `make up` (or use `docker-compose up` instead if `make` is unable to be used)
3. Use the command `make test` (or use `go test ./...` instead, if `make` is unable to be used) to run the automated testing

## API Documentation (Postman)

Link : https://documenter.getpostman.com/view/25087235/2s9YXfcizj
