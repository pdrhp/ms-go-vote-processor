package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/pdrhp/ms-voto-processor-go/internal/core/entity"
	"github.com/pdrhp/ms-voto-processor-go/internal/core/port"
	"github.com/pdrhp/ms-voto-processor-go/internal/infrastructure/persistence/mappers"
	"github.com/pdrhp/ms-voto-processor-go/internal/infrastructure/persistence/models"
)

type PostgresVoteRepository struct {
	db     *sql.DB
	mapper *mappers.VoteMapper
}

func NewPostgresVoteRepository(database *Database) port.VoteRepositoryPort {
	return &PostgresVoteRepository{
		db:     database.DB,
		mapper: mappers.NewVoteMapper(),
	}
}

func (r *PostgresVoteRepository) Save(ctx context.Context, vote *entity.Vote) error {
	if vote == nil {
		return fmt.Errorf("vote cannot be nil")
	}

	if err := vote.Validate(); err != nil {
		return fmt.Errorf("invalid vote: %w", err)
	}

	model := r.mapper.ToModel(vote)

	query := `
		INSERT INTO votes (
			id, participant_id, session_id, timestamp, status,
			processed_at, processing_error, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		) ON CONFLICT (id) DO UPDATE SET
			participant_id = EXCLUDED.participant_id,
			session_id = EXCLUDED.session_id,
			timestamp = EXCLUDED.timestamp,
			status = EXCLUDED.status,
			processed_at = EXCLUDED.processed_at,
			processing_error = EXCLUDED.processing_error,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(ctx, query,
		model.ID,
		model.ParticipantID,
		model.SessionID,
		model.Timestamp,
		model.Status,
		model.ProcessedAt,
		model.ProcessingError,
		model.CreatedAt,
		model.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save vote: %w", err)
	}

	return nil
}

func (r *PostgresVoteRepository) BulkSave(ctx context.Context, votes []*entity.Vote) error {
	if len(votes) == 0 {
		return nil
	}

	for i, vote := range votes {
		if vote == nil {
			return fmt.Errorf("vote at index %d cannot be nil", i)
		}
		if err := vote.Validate(); err != nil {
			return fmt.Errorf("invalid vote at index %d: %w", i, err)
		}
	}

	models := r.mapper.ToModels(votes)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query, args := r.buildBulkInsertQuery(models)

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to bulk save votes: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *PostgresVoteRepository) buildBulkInsertQuery(models []*models.VoteModel) (string, []interface{}) {
	if len(models) == 0 {
		return "", nil
	}

	query := `
		INSERT INTO votes (
			id, participant_id, session_id, timestamp, status,
			processed_at, processing_error, created_at, updated_at
		) VALUES `

	var placeholders []string
	var args []interface{}
	argIndex := 1

	now := time.Now().UTC()

	for _, model := range models {
		placeholder := fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			argIndex, argIndex+1, argIndex+2, argIndex+3, argIndex+4,
			argIndex+5, argIndex+6, argIndex+7, argIndex+8)

		placeholders = append(placeholders, placeholder)

		args = append(args,
			model.ID,
			model.ParticipantID,
			model.SessionID,
			model.Timestamp,
			model.Status,
			model.ProcessedAt,
			model.ProcessingError,
			now,
			now,
		)

		argIndex += 9
	}

	query += strings.Join(placeholders, ", ")

	query += `
		ON CONFLICT (id) DO UPDATE SET
			participant_id = EXCLUDED.participant_id,
			session_id = EXCLUDED.session_id,
			timestamp = EXCLUDED.timestamp,
			status = EXCLUDED.status,
			processed_at = EXCLUDED.processed_at,
			processing_error = EXCLUDED.processing_error,
			updated_at = EXCLUDED.updated_at
	`

	return query, args
}
