package port

import (
	"context"

	"github.com/pdrhp/ms-voto-processor-go/internal/core/entity"
)

type VoteConsumer interface {
    Consume(ctx context.Context) (<-chan *entity.Vote, error)
    Close() error
}

type MessageHandler interface {
    Handle(ctx context.Context, votes []*entity.Vote) error
}