package profileusecase

import (
	"context"
	"github.com/fitzplsr/mgtu-ecg/internal/model"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/profile"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In

	Log  *zap.Logger
	Repo profile.Repo
}

type Profile struct {
	log  *zap.Logger
	repo profile.Repo
}

func New(p Params) *Profile {
	return &Profile{
		log:  p.Log,
		repo: p.Repo,
	}
}

func (u *Profile) Update(ctx context.Context, payload *model.UpdateUserPayload) (*model.User, error) {
	user := &model.User{
		ID:   payload.ID,
		Name: payload.Name,
	}

	updatedUser, err := u.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	return updatedUser, err
}
