package authusecase

import (
	"context"
	"errors"
	"github.com/fitzplsr/mgtu-ecg/internal/model"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/middleware"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/refresh"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/auth"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/profile"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/utils/hasher"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

var _ auth.Usecase = (*Auth)(nil)

type Params struct {
	fx.In

	Log            *zap.Logger
	SessionStorage middleware.SessionStorage
	Repo           auth.Repo
	JWTer          middleware.JWTer
	RefreshConfig  refresh.Config
}

type Auth struct {
	log            *zap.Logger
	sessionStorage middleware.SessionStorage
	repo           auth.Repo
	jwter          middleware.JWTer
	refreshConfig  refresh.Config
}

func New(p Params) *Auth {
	return &Auth{
		log:            p.Log,
		sessionStorage: p.SessionStorage,
		repo:           p.Repo,
		jwter:          p.JWTer,
		refreshConfig:  p.RefreshConfig,
	}
}

func (a *Auth) SignUp(ctx context.Context, payload *model.SignUpPayload) (*model.AuthResponse, *model.Session, error) {
	_, err := a.repo.GetByLogin(ctx, payload.Login)
	if nil == err {
		return &model.AuthResponse{}, &model.Session{}, auth.ErrUserAlreadyExists
	}

	if !errors.Is(err, profile.ErrProfileNotFound) {
		return &model.AuthResponse{}, &model.Session{}, err
	}

	userID, err := uuid.NewV7()
	if err != nil {
		return &model.AuthResponse{}, &model.Session{}, err
	}

	passwordHash, err := hasher.HashPass(payload.Password)
	if err != nil {
		return &model.AuthResponse{}, &model.Session{}, err
	}

	user := model.User{
		ID:           userID,
		Name:         payload.Name,
		Login:        payload.Login,
		PasswordHash: passwordHash,
		Role:         model.RoleAnonymous,
	}

	updatedUser, err := a.repo.Create(ctx, &user)
	if err != nil {
		return nil, nil, err
	}

	return a.makeResponse(ctx, updatedUser)
}

func (a *Auth) SignIn(ctx context.Context, payload *model.SignInPayload) (*model.AuthResponse, *model.Session, error) {
	user, err := a.repo.GetByLogin(ctx, payload.Login)
	if err != nil {
		return nil, nil, err
	}

	if !hasher.CheckPass(user.PasswordHash, payload.Password) {
		return nil, nil, auth.ErrWrongPassword
	}

	return a.makeResponse(ctx, user)
}

func (a *Auth) UpdatePassword(
	ctx context.Context,
	payload *model.UpdatePasswordPayload,
) (*model.AuthResponse, *model.Session, error) {
	if payload.OldPassword == payload.NewPassword {
		return nil, nil, auth.ErrSamePassword
	}

	user, err := a.repo.GetByID(ctx, payload.ID)
	if err != nil {
		return nil, nil, err
	}

	if !hasher.CheckPass(user.PasswordHash, payload.OldPassword) {
		return nil, nil, auth.ErrWrongPassword
	}

	hash, err := hasher.HashPass(payload.NewPassword)
	if err != nil {
		return nil, nil, err
	}
	user.PasswordHash = hash

	updatedUser, err := a.repo.Update(ctx, user)
	if err != nil {
		return nil, nil, err
	}
	return a.makeResponse(ctx, updatedUser)
}

func (a *Auth) UpdateRole(
	ctx context.Context,
	payload *model.UpdateRolePayload,
) (*model.AuthResponse, *model.Session, error) {
	user := &model.User{
		ID:   payload.ID,
		Role: payload.Role,
	}

	if !payload.Role.IsARole() || payload.Role == model.RoleAnonymous {
		return nil, nil, auth.ErrWrongRole
	}

	updatedUser, err := a.repo.Update(ctx, user)
	if err != nil {
		return nil, nil, err
	}
	return a.makeResponse(ctx, updatedUser)
}

func (a *Auth) makeResponse(
	ctx context.Context,
	user *model.User,
) (*model.AuthResponse, *model.Session, error) {
	session, err := a.createSession(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	accessToken, err := a.jwter.GenerateJWT(session)
	if err != nil {
		return nil, nil, err
	}

	return &model.AuthResponse{
		User:        user,
		AccessToken: accessToken,
	}, session, nil
}

func (a *Auth) createSession(ctx context.Context, user *model.User) (*model.Session, error) {
	refreshToken, err := refresh.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	userSession := model.Session{
		UserId:       user.ID,
		UserRole:     0,
		RefreshToken: refreshToken,
		UserAgent:    "",
		Fingerprint:  "",
		Ip:           "",
		ExpiresIn:    a.refreshConfig.RefreshExpirationTime,
		CreatedAt:    time.Now(),
	}

	err = a.sessionStorage.Set(ctx, &userSession)
	if err != nil {
		return nil, err
	}

	return &userSession, nil
}
