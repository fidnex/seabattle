package repo

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"

	"seabattle/internal/domain"
)

const ttl = time.Hour * 24 // Время жизни кеша

type Repo struct {
	cl *redis.Client
}

func New(addr, password string, db int) (*Repo, error) {
	cl := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if _, err := cl.Ping().Result(); err != nil {
		return nil, fmt.Errorf("%w: fail to connect the redis", err)
	}

	return &Repo{cl}, nil
}

func (r *Repo) Close() error {
	return r.cl.Close()
}

func (r *Repo) GetGame(chatID string) (*domain.Game, error) {
	data, err := r.cl.Get(chatID).Result()
	if err != nil {
		return nil, fmt.Errorf("%w: fail to get in redis", err)
	}

	game := &domain.Game{}
	if err := json.Unmarshal([]byte(data), game); err != nil {
		return nil, fmt.Errorf("%w: json unmarshal fail", err)
	}

	return game, nil
}

func (r *Repo) SaveGame(game *domain.Game) error {
	data, err := json.Marshal(game)
	if err != nil {
		return fmt.Errorf("%w: json marshal fail", err)
	}

	if _, err := r.cl.Set(game.ChatID, data, ttl).Result(); err != nil {
		return fmt.Errorf("%w: fail to set in redis", err)
	}

	return nil
}
