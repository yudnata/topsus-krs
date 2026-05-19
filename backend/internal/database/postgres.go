package database

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(url string) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	return pool
}

func migrationDir() string {
	candidates := []string{"migrations", filepath.Join("backend", "migrations")}
	for _, dir := range candidates {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			return dir
		}
	}
	log.Fatal("migrations directory not found (run from backend/ or project root)")
	return ""
}

func ensureMigrationTable(ctx context.Context, pool *pgxpool.Pool) {
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			filename   VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`)
	if err != nil {
		log.Fatalf("Create schema_migrations: %v", err)
	}
}

func isMigrationApplied(ctx context.Context, pool *pgxpool.Pool, name string) bool {
	var exists bool
	err := pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE filename = $1)`, name).Scan(&exists)
	return err == nil && exists
}

func markMigrationApplied(ctx context.Context, pool *pgxpool.Pool, name string) {
	_, err := pool.Exec(ctx, `INSERT INTO schema_migrations (filename) VALUES ($1) ON CONFLICT DO NOTHING`, name)
	if err != nil {
		log.Fatalf("Record migration %s: %v", name, err)
	}
}

func Migrate(pool *pgxpool.Pool) {
	dir := migrationDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Read migrations: %v", err)
	}

	names := make([]string, 0)
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	ctx := context.Background()
	ensureMigrationTable(ctx, pool)

	for _, name := range names {
		if isMigrationApplied(ctx, pool, name) {
			log.Printf("⊘ Migration skipped (already applied): %s", name)
			continue
		}

		path := filepath.Join(dir, name)
		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("Read migration %s: %v", name, err)
		}
		if _, err := pool.Exec(ctx, string(sqlBytes)); err != nil {
			log.Fatalf("Migration %s failed: %v", name, err)
		}
		markMigrationApplied(ctx, pool, name)
		log.Printf("✓ Migration applied: %s", name)
	}
	log.Println("✓ Database migrate check complete")
}
