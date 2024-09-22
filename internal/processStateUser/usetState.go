package processstateuser

import (
	userstate "tgBank/pkg/userStateRepo"
)

type userStateService struct {
	repo *userstate.UserStateRepository
}

func NewUserStateService(repo *userstate.UserStateRepository) UserStateService {
	return &userStateService{
		repo: repo,
	}
}

func (s *userStateService) GetState(userID int64) (userstate.UserState, error) {
	return s.repo.GetState(userID)
}

func (s *userStateService) SetState(userID int64, state userstate.UserState) error {
	return s.repo.SetState(userID, state)
}
