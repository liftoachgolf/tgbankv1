package userstate

import "github.com/go-redis/redis/v8"

type State interface {
	GetState(userID int64) (UserState, error)
	SetState(userID int64, state UserState) error
}

type UserStateRepository struct {
	State
}

func NewUserStateRepository(redis *redis.Client) *UserStateRepository {
	return &UserStateRepository{
		State: NewMemoryRepository(redis),
	}
}
