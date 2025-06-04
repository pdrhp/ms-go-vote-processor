package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pdrhp/ms-voto-processor-go/internal/core/entity"
	"github.com/pdrhp/ms-voto-processor-go/internal/core/util"
)

type VoteMessage struct {
	ID            string    `json:"id"`
	ParticipanteID int      `json:"participanteId"`
	SessionID     string    `json:"sessionId"`
	Timestamp     time.Time `json:"timestamp"`
}

func (v *VoteMessage) ToEntity() *entity.Vote {
	id := v.ID
	if id == "" {
		id = util.GenerateUUID()
	}

	return entity.NewVoteFromData(
		id,
		v.ParticipanteID,
		v.SessionID,
		v.Timestamp,
		entity.VoteStatusReceived,
	)
}

func (v *VoteMessage) FromJSON(data []byte) error {
	return json.Unmarshal(data, v)
}

func (v *VoteMessage) Validate() error {
	if v.ParticipanteID <= 0 {
		return fmt.Errorf("participanteId must be greater than zero")
	}
	if v.SessionID == "" {
		return fmt.Errorf("sessionId is required")
	}
	if v.Timestamp.IsZero() {
		return fmt.Errorf("timestamp is required")
	}
	return nil
}
