package middleware

import (
	"errors"
	"github.com/fitzplsr/mgtu-ecg/internal/model"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/refresh"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/session"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/utils"
	"github.com/fitzplsr/mgtu-ecg/pkg/consts"
	"github.com/fitzplsr/mgtu-ecg/pkg/messages"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"time"
)

type SessionStorage interface {
	GetByUserID(ctx context.Context, userID string) (*model.Session, error)
	Set(ctx context.Context, session *model.Session) error
}

type JWTer interface {
	ValidateJWT(tokenString string) (*model.UserClaims, error)
	GenerateJWT(session *model.Session) (string, error)
}

type ProtectedMWParmas struct {
	fx.In

	SS    SessionStorage
	Log   *zap.Logger
	JWTer JWTer
}

type ProtectedMW struct {
	MW fiber.Handler
}

// TODO таймаут в UserContext
func NewProtectMW(p ProtectedMWParmas) *ProtectedMW {
	jwter := p.JWTer
	ss := p.SS
	log := p.Log

	return &ProtectedMW{MW: func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			log.Error("empty access token")
			return utils.Send401(c, messages.MissingOrInvalidToken)
		}

		claims, err := jwter.ValidateJWT(tokenString)
		if nil == err {
			setLocalsFromJwt(c, claims)
			return c.Next()
		}

		refreshToken := c.Cookies(refresh.RefreshToken)
		c.ClearCookie(refresh.RefreshToken)

		if !errors.Is(err, jwt.ErrTokenExpired) {
			return utils.Send401(c, messages.Unauthorized)
		}

		log.Debug("ValidateJWT err", zap.Error(err))

		// access token expired or invalid

		if refreshToken == "" {
			log.Debug("empty refresh")
			return utils.Send401(c, messages.Unauthorized)
		}

		userID := claims.ID

		userSession, err := ss.GetByUserID(c.UserContext(), userID.String())
		if err != nil {
			if errors.Is(err, session.ErrSessionNotFound) {
				log.Error("no session in storage", zap.Error(err))
				return utils.Send401(c, messages.RefreshTokenNotFound)
			}
			log.Error("unknown err from storage", zap.Error(err))
			return utils.Send500(c, messages.InternalServerError)
		}

		if refreshToken != userSession.RefreshToken {
			log.Error("cookie token doesnt match storage", zap.Error(err))
			return utils.Send401(c, messages.Unauthorized)
		}

		newAccessToken, err := jwter.GenerateJWT(userSession)
		if err != nil {
			log.Error("generate jwt", zap.Error(err))
			return utils.Send500(c, messages.CouldNotGenerateAccessToken)
		}

		newRefreshToken, err := refresh.GenerateRefreshToken()
		if err != nil {
			log.Error("generate refresh", zap.Error(err))
			return utils.Send500(c, messages.CouldNotGenerateRefreshToken)
		}

		userSession.RefreshToken = newRefreshToken
		userSession.CreatedAt = time.Now()
		err = ss.Set(c.UserContext(), userSession)
		if err != nil {
			log.Error("set session to storage", zap.Error(err))
			return utils.Send500(c, messages.InternalServerError)
		}

		c.Set(consts.Authorization, newAccessToken)
		refresh.SetRefreshToken(c, userSession)
		setLocalsFromSession(c, userSession)

		return c.Next()
	},
	}
}

func setLocalsFromJwt(c *fiber.Ctx, claims *model.UserClaims) {
	c.Locals(UserIDKey, claims.ID)
	c.Locals(UserRoleKey, claims.Role)
}

func setLocalsFromSession(c *fiber.Ctx, s *model.Session) {
	c.Locals(UserIDKey, s.UserId)
	c.Locals(UserRoleKey, s.UserRole)
}
