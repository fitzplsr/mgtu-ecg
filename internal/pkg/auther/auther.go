package auther

import (
	"errors"
	"github.com/fitzplsr/mgtu-ecg/internal/model"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

type Params struct {
	fx.In

	Cfg Config
	Log *zap.Logger
}

type Auther struct {
	cfg Config
	log *zap.Logger
}

func New(p Params) *Auther {
	return &Auther{
		cfg: p.Cfg,
		log: p.Log,
	}
}

type authClaims struct {
	UserID   uuid.UUID
	UserRole model.Role
	UserIP   string

	jwt.RegisteredClaims
}

func (a *Auther) ValidateJWT(tokenString string) (*model.UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &authClaims{}, func(token *jwt.Token) (interface{}, error) {
		return a.cfg.JwtAccess, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			res, parseErr := parseClaims(token)
			if parseErr != nil {
				return nil, parseErr
			}
			return res, err
		}
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return parseClaims(token)
}

func parseClaims(token *jwt.Token) (*model.UserClaims, error) {
	claims, ok := token.Claims.(*authClaims)
	if !ok {
		return nil, errors.New("could not parse claims")
	}

	return &model.UserClaims{
		ID:   claims.UserID,
		Role: claims.UserRole,
		IP:   claims.UserIP,
	}, nil
}

func (a *Auther) GenerateJWT(session *model.Session) (string, error) {
	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	//	UserID:   session.UserId,
	//	UserRole: session.UserRole,
	//	UserIP:   session.Ip,
	//	Exp:      time.Now().Add(a.cfg.AccessExpirationTime).Unix(),
	//})
	exp := time.Now().Add(a.cfg.AccessExpirationTime)
	a.log.Debug("expiration", zap.Time("exp", exp))
	a.log.Debug("cfg", zap.Duration("cdf", a.cfg.AccessExpirationTime))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &authClaims{
		UserID:   session.UserId,
		UserRole: session.UserRole,
		UserIP:   session.Ip,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	})

	return token.SignedString(a.cfg.JwtAccess)
}
