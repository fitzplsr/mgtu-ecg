package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/fitzplsr/mgtu-ecg/internal/model"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/middleware"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/session"
	"github.com/go-redis/redis/v8"
	"github.com/mailru/easyjson"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var _ middleware.SessionStorage = (*Storage)(nil)

type Params struct {
	fx.In

	Client *redis.Client
	Log    *zap.Logger
}

type Storage struct {
	client *redis.Client
	log    *zap.Logger
}

func NewStorage(p Params) *Storage {
	return &Storage{client: p.Client, log: p.Log}
}

func (s *Storage) Set(ctx context.Context, session *model.Session) error {
	b, err := easyjson.Marshal(session)
	if err != nil {
		s.log.Error("error marshal session", zap.Error(err))
		return fmt.Errorf("marshal session: %w", err)
	}

	res, err := s.client.Set(ctx, session.UserId.String(), b, session.ExpiresIn).Result()
	if err != nil {
		s.log.Error("cant set data in redis", zap.Error(err))
		return err
	}
	if res != "OK" {
		err := fmt.Errorf("'set' in redis replies 'not OK'")
		s.log.Error("error set value", zap.Error(err))
		return err
	}
	return nil
}

func (s *Storage) GetByUserID(ctx context.Context, userID string) (*model.Session, error) {
	var userSession model.Session

	result, err := s.client.GetDel(ctx, userID).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {

			return &model.Session{}, session.ErrSessionNotFound
		}
		s.log.Error("cant get data in redis", zap.Error(err))

		return &model.Session{}, fmt.Errorf("cant get data in redis: %w", err)
	}

	err = easyjson.Unmarshal([]byte(result), &userSession)
	if err != nil {
		s.log.Error("error unmarshal session", zap.Error(err))

		return &model.Session{}, fmt.Errorf("unmarshal session: %w", err)
	}
	//userSession.CreatedAt = time.Now()

	return &userSession, nil
}
