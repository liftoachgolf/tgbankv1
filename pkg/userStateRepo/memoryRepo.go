package userstate

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type memoryRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewMemoryRepository(client *redis.Client) State {
	return &memoryRepository{
		client: client,
		ctx:    context.Background(),
	}
}

type UserState struct {
	State string
}

const (
	StateStart           = "START"
	StateWaitingForHello = "WAITING_FOR_HELLO"
	StateEnd             = "END"
)

func (r *memoryRepository) GetState(userID int64) (UserState, error) {
	key := r.getKey(userID)
	data, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return UserState{State: "START"}, nil
		}
		return UserState{}, err
	}
	var state UserState
	err = json.Unmarshal([]byte(data), &state)
	if err != nil {
		return UserState{}, err
	}
	return state, nil
}

func (r *memoryRepository) SetState(userID int64, state UserState) error {
	key := r.getKey(userID)
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return r.client.Set(r.ctx, key, data, 0).Err()
}

func (r *memoryRepository) getKey(userID int64) string {
	return fmt.Sprintf("user_state:%d", userID)
}
