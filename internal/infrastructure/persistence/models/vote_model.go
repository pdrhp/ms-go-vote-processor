package models

import (
	"database/sql/driver"
	"time"

	"github.com/pdrhp/ms-voto-processor-go/internal/core/entity"
)

type VoteModel struct {
	ID               string     `db:"id"`
	ParticipantID    int        `db:"participant_id"`
	SessionID        string     `db:"session_id"`
	Timestamp        time.Time  `db:"timestamp"`
	Status           string     `db:"status"`
	ProcessedAt      *time.Time `db:"processed_at"`
	ProcessingError  *string    `db:"processing_error"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
}

func (v *VoteModel) FromEntity(vote *entity.Vote) {
	v.ID = vote.ID
	v.ParticipantID = vote.ParticipantID
	v.SessionID = vote.SessionID
	v.Timestamp = vote.Timestamp
	v.Status = string(vote.Status)
	v.ProcessedAt = vote.ProcessedAt
	v.ProcessingError = vote.ProcessingError
	v.CreatedAt = time.Now().UTC()
	v.UpdatedAt = time.Now().UTC()
}

func (v *VoteModel) ToEntity() *entity.Vote {
	vote := &entity.Vote{
		ID:              v.ID,
		ParticipantID:   v.ParticipantID,
		SessionID:       v.SessionID,
		Timestamp:       v.Timestamp,
		Status:          entity.VoteStatus(v.Status),
		ProcessedAt:     v.ProcessedAt,
		ProcessingError: v.ProcessingError,
	}
	return vote
}

func (v *VoteModel) Scan(value interface{}) error {
	return nil
}

func (v VoteModel) Value() (driver.Value, error) {
	return v.Status, nil
}
