## Tests:
[![Go](https://github.com/Kamillendampf/LiquibaseXgolang/actions/workflows/ci.yml/badge.svg)](https://github.com/Kamillendampf/LiquibaseXgolang/actions/workflows/ci.yml)

## 📦 Liquibase CLI Wrapper for Go

A lightweight Go library to manage **Liquibase database migrations** programmatically using the **Liquibase CLI** – with support for integration testing via **testcontainers-go**.

---

### ✅ Features

- Run Liquibase migrations from Go code
- CLI-based execution (`liquibase` binary or custom path)
- Supports core commands:
- `update`
- `rollback <tag>`
    - `tag`
    - `status`
    - `validate`
    - `clearCheckSums`
    - `releaseLocks`
    - Integration tested with:
    - [`testcontainers-go`](https://github.com/testcontainers/testcontainers-go)
    - PostgreSQL
    - JDBC connection
    - `pgx` to verify state (e.g. applied changesets)

    ---

    ### 🚀 Installation

    ```bash
    go get github.com/Kamillendampf/LiquibaseXgolang
    ```

    ---

    ### 🧰 Requirements

    - Liquibase CLI installed and available in `$PATH` (or provide path manually)
    - Java (required by Liquibase CLI)
    - For testing: Docker + Go ≥ 1.18

    ---

    ### 🧱 Usage

    ```go
    import "github.com/Kamillendampf/LiquibaseXgolang"

    cfg := liquibaseMigrationHelper.Config{
    ChangelogFile: "testdata/changelog.xml",
    URL:           "jdbc:postgresql://localhost:5432/mydb",
    Username:      "postgres",
    Password:      "secret",
    Command:       "liquibase", // Optional – defaults to "liquibase"
    }

    lb := liquibaseMigrationHelper.New(cfg)

    // Run migrations
    if err := lb.Update(); err != nil {
    log.Fatalf("Update failed: %v", err)
    }

    // Tag database state
    _ = lb.Tag("v1.0.0")

    // Validate changelog syntax
    _ = lb.Validate()

    // Rollback to a previous tag
    _ = lb.Rollback("v1.0.0")
    ```

    ---

    ### 🥪 Integration Testing with testcontainers-go

    The project includes integration tests using [testcontainers-go](https://github.com/testcontainers/testcontainers-go) and PostgreSQL. Tests cover:

    - Migration via CLI
    - Rollback
    - Changelog validation
    - Checksum clearing
    - Lock release
    - Changeset verification via direct DB queries (`pgx`)

    To run tests:

    ```bash
    go test -v
    ```

    > Make sure Docker is running and Liquibase is in your `PATH`.

    ---

    ### 📂 Project Structure

    ```plaintext
    liquibasehelper/
    ├ liquibaseMigrationHelper.go           // Config struct, High-level API (Update, Rollback, etc.), CLI runner
    ├ liquibaseMigrationHelper_test.go      // Testcontainers integration tests
    test_resources/
    ├ changelog.xml                         // standard and valid Changelog
    ├ changelog-1.xml                       // test for rollback
    └ invalid-changelog.xml                 // Broken changelog for testing
    ```

    ---

    ### 🧠 Notes

    - You must use a JDBC-compatible database (e.g. PostgreSQL, MySQL, H2)
    - Liquibase runs via CLI (not native Go) – output and error handling is forwarded
    - You can test changesets with `history`, or verify state via SQL queries

    ---

    ### 📜 License

    MIT – use freely in commercial or open source projects.
