package port

import (
	"context"

	"github.com/pdrhp/ms-voto-processor-go/internal/core/entity"
)

type VoteRepositoryPort interface {
    BulkSave(ctx context.Context, votes []*entity.Vote) error
    Save(ctx context.Context, vote *entity.Vote) error
}