package usecase

import (
	"context"

	"github.com/pdrhp/ms-voto-processor-go/internal/core/entity"
	"github.com/pdrhp/ms-voto-processor-go/internal/core/port"
)

type VoteProcessor struct {
    repository port.VoteRepository
    batchSize  int
}

func NewVoteProcessor(repository port.VoteRepository, batchSize int) *VoteProcessor {
    return &VoteProcessor{
		repository: repository,
		batchSize:  batchSize,
	}
}

func (vp *VoteProcessor) ProcessVotesBatch(ctx context.Context, votes []*entity.Vote) error {
	return nil
}
