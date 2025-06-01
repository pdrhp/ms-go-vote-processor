package port

import (
	"context"

	"github.com/pdrhp/ms-voto-processor-go/internal/core/entity"
)

type VoteRepository interface {
    BulkSave(ctx context.Context, votes []*entity.Vote) error
}