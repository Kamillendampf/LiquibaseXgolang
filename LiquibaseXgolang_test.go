package LiquibaseXgolang

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"
)

func TestLiquibase_Update(t *testing.T) {

	//Test Setup
	ctx := context.Background()

	container, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("secret"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithStartupTimeout(30*time.Second),
		),
	)

	if err != nil {
		t.Fatalf("Postgres container faild %v", err)
	}

	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("Could not stop container: %v", err)
		}
	}()

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}

	port, err := container.MappedPort(ctx, nat.Port("5432/tcp"))
	if err != nil {
		t.Fatal(err)
	}

	// test config
	jdbcURL := fmt.Sprintf("jdbc:postgresql://%s:%s/testdb", host, port.Port())

	cfg := Config{
		ChangelogFile: "test_resources/changelog.xml",
		URL:           jdbcURL,
		Username:      "postgres",
		Password:      "secret",
	}

	t.Run("update should pass", func(t *testing.T) {
		lb := New(cfg)
		if err := lb.Update(); err != nil {
			t.Errorf("Update failed: %v", err)
		}
	})

	t.Run("tag without should pass", func(t *testing.T) {
		lb := New(cfg)
		if err := lb.Tag("test"); err != nil {
			t.Errorf("Tag failed: %v", err)
		}
	})

	t.Run("status should pass", func(*testing.T) {
		lb := New(cfg)
		if err := lb.Status(); err != nil {
			t.Errorf("Status failed: %v", err)
		}
	})

	t.Run("Rollback should pass", func(t *testing.T) {
		lb := New(cfg)
		if err := lb.Update(); err != nil {
			t.Errorf("Error while setup Rollback test in 'Update': %v", err)
		}
		if err := lb.Tag("test"); err != nil {
			t.Errorf("Error while setup Rollback test in 'Tag': %v", err)
		}
		lb.cfg.ChangelogFile = "test_resources/changelog-1.xml"

		if err := lb.Update(); err != nil {
			t.Errorf("Error while setup Rollback test in 'Update': %v", err)
		}

		if err := lb.Rollback("test"); err != nil {
			t.Errorf("Rollback failed: %v", err)
		}
	})

	t.Run("Validated should succeed with valid changelog", func(t *testing.T) {
		lb := New(cfg)
		if err := lb.Validate(); err != nil {
			t.Errorf("Expected Validate() to succeed, got error: %v", err)
		}
	})

	t.Run("Validated should fail with invalide changelog", func(t *testing.T) {
		lb := New(cfg)

		lb.cfg.ChangelogFile = "test_resources/invalide-changelog.xml"
		if err := lb.Validate(); err == nil {
			t.Errorf("Expected Validate () to fail due to invalide changelog, but its passed")
		} else {
			t.Logf("Validate failed with expection: %v", err)
		}
	})

	t.Run("ClearChecksum executes successfully", func(t *testing.T) {
		lb := New(cfg)

		// Erstmal Migration durchführen, sonst gibt's nichts zu löschen
		if err := lb.Update(); err != nil {
			t.Fatalf("Initial update failed: %v", err)
		}

		// Dann ClearChecksums ausführen
		if err := lb.ClearChecksums(); err != nil {
			t.Fatalf("ClearChecksums failed: %v", err)
		}
	})

	t.Run("ReleaseLocks clears active DB lock", func(t *testing.T) {

		lb := New(cfg)

		connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			cfg.Username, cfg.Password, host, port.Port(), "testdb",
		)

		// 1. Setup: Migration ausführen (damit Tabelle existiert)
		if err := lb.Update(); err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		// 2. Lock manuell setzen
		if err := setLockManually(connStr); err != nil {
			t.Fatalf("Setup failed: %v", err)
		}

		// 3. Verifizieren, dass Lock gesetzt ist
		locked, err := isLockActive(connStr)
		if err != nil {
			t.Fatalf("Lock status check failed: %v", err)
		}
		if !locked {
			t.Fatalf("Expected lock to be active before release")
		}

		// 4. ReleaseLocks ausführen
		if err := lb.ReleaseLocks(); err != nil {
			t.Fatalf("ReleaseLocks failed: %v", err)
		}

		// 5. Verifizieren, dass Lock aufgehoben wurde
		locked, err = isLockActive(connStr)
		if err != nil {
			t.Fatalf("Lock recheck failed: %v", err)
		}
		if locked {
			t.Errorf("Expected lock to be cleared, but it is still active")
		} else {
			t.Logf("ReleaseLocks successfully cleared the lock ")
		}
	})

	t.Run("Apply update", func(t *testing.T) {
		lb := New(cfg)
		if err := lb.Update(); err != nil {
			t.Fatalf("Update failed before history test: %v", err)
		}
	})

	t.Run("Read history", func(t *testing.T) {
		lb := New(cfg)
		if err := lb.History(); err != nil {
			t.Errorf("Expected History() to succeed, got error: %v", err)
		}
	})

}

func isLockActive(connStr string) (bool, error) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return false, fmt.Errorf("DB connect failed: %w", err)
	}
	defer conn.Close(ctx)

	var locked bool
	err = conn.QueryRow(ctx, `SELECT locked FROM DATABASECHANGELOGLOCK WHERE ID = 1`).Scan(&locked)
	if err != nil {
		return false, fmt.Errorf("query failed: %w", err)
	}
	return locked, nil
}

func setLockManually(connStr string) error {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return fmt.Errorf("DB connect failed: %w", err)
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, `UPDATE DATABASECHANGELOGLOCK SET LOCKED = true, LOCKGRANTED = now(), LOCKEDBY = 'test-runner' WHERE ID = 1`)
	if err != nil {
		return fmt.Errorf("failed to set manual lock: %w", err)
	}
	return nil
}
