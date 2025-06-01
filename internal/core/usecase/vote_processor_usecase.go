package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/pdrhp/ms-voto-processor-go/internal/core/entity"
	"github.com/pdrhp/ms-voto-processor-go/internal/core/port"
)

type VoteProcessorUsecase struct {
    repository port.VoteRepositoryPort
    batchSize  int
}

func NewVoteProcessorUsecase(repository port.VoteRepositoryPort, batchSize int) *VoteProcessorUsecase {
    return &VoteProcessorUsecase{
		repository: repository,
		batchSize:  batchSize,
	}
}

func (vp *VoteProcessorUsecase) ProcessSingleVote(ctx context.Context, vote *entity.Vote) error {
    if err := vote.Validate(); err != nil {
        return fmt.Errorf("voto inválido: %w", err)
    }

    if err := vote.MarkAsProcessing(); err != nil {
        return fmt.Errorf("falha ao marcar como processando: %w", err)
    }

    log.Printf("Processando voto: ID=%s, ParticipantID=%d", vote.ID, vote.ParticipantID)

    if err := vp.repository.Save(ctx, vote); err != nil {
        vote.MarkAsFailedWithError(err)
        return fmt.Errorf("falha ao salvar voto: %w", err)
    }

    vote.MarkAsProcessed()

    return nil
}

func (vp *VoteProcessorUsecase) ProcessVotesBatch(ctx context.Context, votes []*entity.Vote) error {
    if len(votes) == 0 {
        return nil
    }

    log.Printf("Processando batch de %d votos", len(votes))

    validVotes := make([]*entity.Vote, 0, len(votes))
    invalidCount := 0

    for _, vote := range votes {
        if err := vp.validateAndPrepareVote(vote); err != nil {
            log.Printf("Voto inválido ignorado: ID=%s, erro=%v", vote.ID, err)
            invalidCount++
            continue
        }
        validVotes = append(validVotes, vote)
    }

    if len(validVotes) == 0 {
        return fmt.Errorf("nenhum voto válido encontrado no batch de %d votos", len(votes))
    }

    if err := vp.repository.BulkSave(ctx, validVotes); err != nil {
        for _, vote := range validVotes {
            vote.MarkAsFailedWithError(err)
        }
        return fmt.Errorf("falha ao salvar batch: %w", err)
    }

    for _, vote := range validVotes {
        vote.MarkAsProcessed()
    }

    log.Printf("Batch processado com sucesso: %d votos salvos, %d inválidos", len(validVotes), invalidCount)
    return nil
}

func (vp *VoteProcessorUsecase) validateAndPrepareVote(vote *entity.Vote) error {
    if err := vote.Validate(); err != nil {
        return err
    }

    if !vote.CanBeProcessed() {
        return fmt.Errorf("voto não pode ser processado, status atual: %s", vote.Status)
    }

    if err := vote.MarkAsProcessing(); err != nil {
        return fmt.Errorf("falha ao marcar voto como processando: %w", err)
    }

    return nil
}

func (vp *VoteProcessorUsecase) GetBatchSize() int {
    return vp.batchSize
}
