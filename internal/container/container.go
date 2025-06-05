package container

import (
	"context"
	"fmt"
	"log"

	"github.com/pdrhp/ms-voto-processor-go/internal/config"
	"github.com/pdrhp/ms-voto-processor-go/internal/core/port"
	"github.com/pdrhp/ms-voto-processor-go/internal/core/usecase"
	"github.com/pdrhp/ms-voto-processor-go/internal/infrastructure/persistence"
)

type Container struct {
	config *config.Config

	database     *persistence.Database
	migrator     *persistence.Migrator

	voteRepository port.VoteRepositoryPort

	voteProcessor *usecase.VoteProcessorUsecase

	isBuilt   bool
	isHealthy bool
}

func NewContainer(cfg *config.Config) *Container {
	return &Container{
		config: cfg,
	}
}

func (c *Container) Build() error {
	if c.isBuilt {
		return fmt.Errorf("container already built")
	}


	if err := c.buildDatabase(); err != nil {
		return fmt.Errorf("failed to build database: %w", err)
	}

	log.Println("Database connection established")

	if err := c.runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err := c.buildRepositories(); err != nil {
		return fmt.Errorf("failed to build repositories: %w", err)
	}

	if err := c.buildUseCases(); err != nil {
		return fmt.Errorf("failed to build use cases: %w", err)
	}

	if err := c.performHealthCheck(); err != nil {
		return fmt.Errorf("failed to perform health check: %w", err)
	}

	c.isBuilt = true
	log.Println("Container built successfully")

	return nil
}

func (c *Container) buildDatabase() error {
	db, err := persistence.NewConnection(&c.config.Database)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	c.database = db

	c.migrator = persistence.NewMigrator(c.database.DB, c.config.Database.MigrationsPath)

	return nil
}

func (c *Container) runMigrations() error {
	if !c.config.Database.RunMigrations {
		log.Println("Migrations disabled, skipping...")
		return nil
	}

	log.Println("Running database migrations...")
	if err := c.migrator.Run(); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("Database migrations completed")
	return nil
}

func (c *Container) buildRepositories() error {
	c.voteRepository = persistence.NewPostgresVoteRepository(c.database)

	return nil
}

func (c *Container) buildUseCases() error {
	c.voteProcessor = usecase.NewVoteProcessorUsecase(
		c.voteRepository,
		c.config.Kafka.BatchSize,
	)

	return nil
}

func (c *Container) performHealthCheck() error {
	log.Println("Performing health checks...")

	if err := c.database.Health(); err != nil {
		return fmt.Errorf("database unhealthy: %w", err)
	}
	log.Println("Database health check passed")

	log.Println("All health checks passed")
	c.isHealthy = true

	return nil
}

func (c *Container) Start(ctx context.Context) error {
	if !c.isBuilt {
		return fmt.Errorf("container not built, call Build() first")
	}

	if !c.isHealthy {
		return fmt.Errorf("container not healthy")
	}

	// TODO implementar inicio do worker

	log.Println("Starting application...")

	c.logConfiguration()

	log.Println("Application started successfully!")
	return nil
}

func (c *Container) Stop() {
	log.Println("Stopping application...")

	// TODO implementar parada do worker

	log.Println("Application stopped successfully")
}

func (c *Container) Close() error {
	log.Println("Closing container resources...")

	if c.database != nil {
		if err := c.database.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
			return err
		}
	}

	log.Println("Container resources closed")
	return nil
}

func (c *Container) logConfiguration() {
	cfg := c.config
	log.Printf("Configuration Summary:")
	log.Printf("   Environment: %s", cfg.App.Environment)
	log.Printf("   Database: %s:%s/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Database)
	log.Printf("   Kafka Topic: %s", cfg.Kafka.Topic)
	log.Printf("   Consumer Group: %s", cfg.Kafka.ConsumerGroup)
	log.Printf("   Batch Size: %d", cfg.Kafka.BatchSize)
	log.Printf("   Workers: %d", cfg.Kafka.Workers)
}

func (c *Container) IsReady() bool {
	return c.isBuilt && c.isHealthy
}