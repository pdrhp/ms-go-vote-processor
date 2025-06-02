package persistence

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Migration struct {
	Version   string
	Name      string
	SQL       string
	Timestamp time.Time
}

type Migrator struct {
	db            *sql.DB
	migrationsDir string
}

func NewMigrator(db *sql.DB, migrationsDir string) *Migrator {
	return &Migrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

func (m *Migrator) Run() error {
	log.Println("Starting database migrations...")

	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	migrations, err := m.loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	for _, migration := range migrations {
		if err := m.runMigration(migration); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", migration.Version, err)
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}

func (m *Migrator) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			executed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := m.db.Exec(query)
	return err
}

func (m *Migrator) loadMigrations() ([]Migration, error) {
	var migrations []Migration

	err := filepath.WalkDir(m.migrationsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".sql") {
			return nil
		}

		migration, err := m.parseMigrationFile(path)
		if err != nil {
			return fmt.Errorf("failed to parse migration file %s: %w", path, err)
		}

		migrations = append(migrations, migration)
		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func (m *Migrator) parseMigrationFile(filePath string) (Migration, error) {
	fileName := filepath.Base(filePath)
	parts := strings.SplitN(fileName, "_", 2)
	if len(parts) != 2 {
		return Migration{}, fmt.Errorf("invalid migration filename format: %s", fileName)
	}

	version := parts[0]
	name := strings.TrimSuffix(parts[1], ".sql")

	sqlContent, err := os.ReadFile(filePath)
	if err != nil {
		return Migration{}, fmt.Errorf("failed to read migration file %s: %w", filePath, err)
	}

	return Migration{
		Version: version,
		Name:    name,
		SQL:     string(sqlContent),
	}, nil
}

func (m *Migrator) runMigration(migration Migration) error {
	var count int
	err := m.db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = $1", migration.Version).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Printf("Migration %s already executed, skipping", migration.Version)
		return nil
	}

	log.Printf("Executing migration %s: %s", migration.Version, migration.Name)

	if strings.TrimSpace(migration.SQL) != "" {
		if err := m.executeMigrationSQL(migration.SQL); err != nil {
			return fmt.Errorf("failed to execute migration SQL: %w", err)
		}
	}

	_, err = m.db.Exec(
		"INSERT INTO schema_migrations (version, name, executed_at) VALUES ($1, $2, $3)",
		migration.Version,
		migration.Name,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to record migration execution: %w", err)
	}

	log.Printf("Migration %s executed successfully", migration.Version)
	return nil
}

func (m *Migrator) executeMigrationSQL(sqlContent string) error {
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	statements := m.splitSQLStatements(sqlContent)

	for i, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		log.Printf("Executing statement %d: %s", i+1, m.truncateSQL(statement))

		if _, err := tx.Exec(statement); err != nil {
			return fmt.Errorf("failed to execute statement %d: %w\nSQL: %s", i+1, err, statement)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (m *Migrator) splitSQLStatements(sqlContent string) []string {
	statements := []string{}
	current := ""
	inString := false
	var stringChar byte

	for i := 0; i < len(sqlContent); i++ {
		char := sqlContent[i]

		if (char == '\'' || char == '"') && (i == 0 || sqlContent[i-1] != '\\') {
			if !inString {
				inString = true
				stringChar = char
			} else if char == stringChar {
				inString = false
			}
		}

		if char == ';' && !inString {
			statements = append(statements, current)
			current = ""
			continue
		}

		current += string(char)
	}

	if strings.TrimSpace(current) != "" {
		statements = append(statements, current)
	}

	return statements
}

func (m *Migrator) truncateSQL(sql string) string {
	if len(sql) > 100 {
		return sql[:100] + "..."
	}
	return sql
}
