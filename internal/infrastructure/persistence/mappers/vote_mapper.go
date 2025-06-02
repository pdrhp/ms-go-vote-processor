package mappers

import (
	"github.com/pdrhp/ms-voto-processor-go/internal/core/entity"
	"github.com/pdrhp/ms-voto-processor-go/internal/infrastructure/persistence/models"
)

type VoteMapper struct{}

func NewVoteMapper() *VoteMapper {
	return &VoteMapper{}
}

func (m *VoteMapper) ToModel(vote *entity.Vote) *models.VoteModel {
	model := &models.VoteModel{}
	model.FromEntity(vote)
	return model
}

func (m *VoteMapper) ToEntity(model *models.VoteModel) *entity.Vote {
	return model.ToEntity()
}

func (m *VoteMapper) ToModels(votes []*entity.Vote) []*models.VoteModel {
	models := make([]*models.VoteModel, len(votes))
	for i, vote := range votes {
		models[i] = m.ToModel(vote)
	}
	return models
}

func (m *VoteMapper) ToEntities(models []*models.VoteModel) []*entity.Vote {
	entities := make([]*entity.Vote, len(models))
	for i, model := range models {
		entities[i] = m.ToEntity(model)
	}
	return entities
}
