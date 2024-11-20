package auth

import (
	"context"
	"github.com/fitzplsr/mgtu-ecg/internal/model"
	"github.com/google/uuid"
)

type Usecase interface {
	SignUp(ctx context.Context, payload *model.SignUpPayload) (*model.AuthResponse, *model.Session, error)
	SignIn(ctx context.Context, payload *model.SignInPayload) (*model.AuthResponse, *model.Session, error)
	UpdatePassword(ctx context.Context, payload *model.UpdatePasswordPayload) (*model.AuthResponse, *model.Session, error)
	UpdateRole(ctx context.Context, payload *model.UpdateRolePayload) (*model.AuthResponse, *model.Session, error)
}

type Repo interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	GetByLogin(ctx context.Context, login string) (*model.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
}
