package service

import (
	"context"
	"tgBank/db/postgres"
	"tgBank/models"
)

type EntryServiceImpl struct {
	repo postgres.EntrySQL
}

func NewEntryServiceImpl(repo postgres.EntrySQL) *EntryServiceImpl {
	return &EntryServiceImpl{
		repo: repo,
	}
}

func (s *EntryServiceImpl) CreateEntry(ctx context.Context, arg models.CreateEntryParams) (models.Entry, error) {
	return s.CreateEntry(ctx, arg)
}
func (s *EntryServiceImpl) GetEntry(ctx context.Context, id int64) (models.Entry, error) {
	return s.GetEntry(ctx, id)
}
func (s *EntryServiceImpl) ListEntries(ctx context.Context, arg models.ListEntriesParams) ([]models.Entry, error) {
	return s.ListEntries(ctx, arg)
}
