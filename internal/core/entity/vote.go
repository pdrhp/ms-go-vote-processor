package entity

import (
	"fmt"
	"time"

	"github.com/pdrhp/ms-voto-processor-go/internal/core/util"
)

type Vote struct {
	ID            string
	ParticipantID int
	SessionID     string
	Timestamp     time.Time
	Status        VoteStatus

	ProcessedAt     *time.Time
	ProcessingError *string
}

func NewVote(participantID int, sessionID string) *Vote {
	return &Vote{
		ID:            util.GenerateUUID(),
		ParticipantID: participantID,
		SessionID:     sessionID,
		Timestamp:     time.Now().UTC(),
		Status:        VoteStatusReceived,
	}
}

func NewVoteFromData(id string, participantID int, sessionID string, timestamp time.Time, status VoteStatus) *Vote {
	return &Vote{
		ID:            id,
		ParticipantID: participantID,
		SessionID:     sessionID,
		Timestamp:     timestamp,
		Status:        status,
	}
}

func (v *Vote) MarkAsProcessing() error {
	if !v.Status.CanTransitionTo(VoteStatusProcessing) {
		return fmt.Errorf("não é possível alterar status de %s para %s", v.Status, VoteStatusProcessing)
	}
	v.Status = VoteStatusProcessing
	return nil
}

func (v *Vote) MarkAsProcessed() error {
	if !v.Status.CanTransitionTo(VoteStatusProcessed) {
		return fmt.Errorf("não é possível alterar status de %s para %s", v.Status, VoteStatusProcessed)
	}
	v.Status = VoteStatusProcessed
	now := time.Now().UTC()
	v.ProcessedAt = &now
	v.ProcessingError = nil
	return nil
}

func (v *Vote) MarkAsFailedWithError(err error) error {
	if !v.Status.CanTransitionTo(VoteStatusFailed) {
		return fmt.Errorf("não é possível alterar status de %s para %s", v.Status, VoteStatusFailed)
	}
	v.Status = VoteStatusFailed
	errStr := err.Error()
	v.ProcessingError = &errStr
	return nil
}

func (v *Vote) SetStatus(status VoteStatus) error {
	if err := status.Validate(); err != nil {
		return fmt.Errorf("falha ao alterar status: %w", err)
	}
	if !v.Status.CanTransitionTo(status) {
		return fmt.Errorf("transição inválida de %s para %s", v.Status, status)
	}
	v.Status = status
	return nil
}

func (v *Vote) MarkAsFailed() error {
	return v.SetStatus(VoteStatusFailed)
}

func (v *Vote) IsProcessed() bool {
	return v.Status == VoteStatusProcessed || v.Status == VoteStatusFailed
}

func (v *Vote) IsValid() bool {
	return v.ID != "" &&
		   v.ParticipantID > 0 &&
		   v.SessionID != "" &&
		   !v.Timestamp.IsZero() &&
		   v.Status.IsValid()
}

func (v *Vote) CanBeProcessed() bool {
	return v.Status == VoteStatusSent || v.Status == VoteStatusReceived
}

func (v *Vote) HasError() bool {
	return v.ProcessingError != nil && *v.ProcessingError != ""
}

func (v *Vote) Validate() error {
	if v.ID == "" {
		return fmt.Errorf("id é obrigatório")
	}
	if v.ParticipantID <= 0 {
		return fmt.Errorf("participantId deve ser maior que zero")
	}
	if v.SessionID == "" {
		return fmt.Errorf("sessionId é obrigatório")
	}
	if v.Timestamp.IsZero() {
		return fmt.Errorf("timestamp é obrigatório")
	}
	return v.Status.Validate()
}

func (v *Vote) String() string {
	return fmt.Sprintf("Vote[ID=%s, ParticipantID=%d, SessionID=%s, Status=%s]",
		v.ID, v.ParticipantID, v.SessionID, v.Status)
}


