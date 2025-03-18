package profile

import (
	"context"

	"github.com/fitzplsr/mgtu-ecg/internal/model"
)

type Usecase interface {
	Update(ctx context.Context, user *model.UpdateUserPayload) (*model.User, error)
}

type Repo interface {
	Update(ctx context.Context, user *model.User) (*model.User, error)
}
