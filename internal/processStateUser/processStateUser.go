package processstateuser

import userstate "tgBank/pkg/userStateRepo"

type UserStateService interface {
	GetState(userID int64) (userstate.UserState, error)
	SetState(userID int64, state userstate.UserState) error
}

type ProcessUserStateService struct {
	UserStateService
}

func NewProcessUserStateService(repo *userstate.UserStateRepository) *ProcessUserStateService {
	return &ProcessUserStateService{
		UserStateService: NewUserStateService(repo),
	}
}
