package profilepostgres

import (
	"context"
	"errors"
	models "github.com/fitzplsr/mgtu-ecg/gen"
	"github.com/fitzplsr/mgtu-ecg/internal/model"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/profile"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/utils/pghelper"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/utils/txer"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

type PostgresParams struct {
	fx.In

	Log  *zap.Logger
	Conn *pgxpool.Pool
}

type Profile struct {
	log  *zap.Logger
	db   *models.Queries
	conn *pgxpool.Pool
}

func New(p PostgresParams) *Profile {
	return &Profile{
		log:  p.Log,
		db:   models.New(p.Conn),
		conn: p.Conn,
	}
}

// TODO txfunc with runfunc
func (r *Profile) Create(ctx context.Context, user *model.User) (*model.User, error) {
	txConn, err := txer.Tx(ctx, r.conn)
	if err != nil {
		r.log.Error("begin tx", zap.Error(err))
		return nil, err
	}

	tx := r.db.WithTx(txConn)

	res, err := tx.Create(ctx, models.CreateParams{
		ID:           pghelper.ToPGUUID(user.ID),
		Role:         int32(user.Role),
		Name:         user.Name,
		Login:        user.Login,
		PasswordHash: user.PasswordHash,
	})
	log.Debug("res", zap.Any("res", res))
	if err != nil {
		r.log.Error("create user", zap.Error(err))
		return nil, err
	}
	err = txConn.Commit(ctx)
	if err != nil {
		txConn.Rollback(ctx)
		return nil, err
	}
	return convertUserToModel(&res), nil
}

func (r *Profile) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return r.getByID(ctx, r.db, id)
}

func (r *Profile) getByID(ctx context.Context, db *models.Queries, id uuid.UUID) (*model.User, error) {
	user, err := db.GetByID(ctx, pghelper.ToPGUUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, profile.ErrProfileNotFound
		}
		r.log.Error("get user by id", zap.Error(err))
		return nil, err
	}
	return convertUserToModel(&user), err
}

func (r *Profile) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	user, err := r.db.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, profile.ErrProfileNotFound
		}
		r.log.Error("get user by login", zap.Error(err))
		return nil, err
	}
	return convertUserToModel(&user), err
}

func (r *Profile) Update(ctx context.Context, user *model.User) (*model.User, error) {
	txConn, err := txer.Tx(ctx, r.conn)
	if err != nil {
		r.log.Error("begin tx", zap.Error(err))
		return nil, err
	}

	tx := r.db.WithTx(txConn)

	dbUser, err := r.getByID(ctx, tx, user.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.log.Error("user not exist", zap.Error(err))
			return nil, err
		}
		return nil, err
	}

	if !isChanged(user, dbUser) {
		return dbUser, nil
	}

	newUser := model.User{
		Role:         lo.Ternary(user.Role == 0, dbUser.Role, user.Role),
		Name:         lo.Ternary(user.Name == "", dbUser.Name, user.Name),
		Login:        lo.Ternary(user.Login == "", dbUser.Login, user.Login),
		PasswordHash: lo.Ternary(len(user.PasswordHash) == 0, dbUser.PasswordHash, user.PasswordHash),
	}

	res, err := tx.Update(ctx, models.UpdateParams{
		ID:           pghelper.ToPGUUID(user.ID),
		Role:         int32(newUser.Role),
		Name:         newUser.Name,
		Login:        newUser.Login,
		PasswordHash: newUser.PasswordHash,
		UpdatedAt:    pghelper.ToPGTimestamp(time.Now()),
	})

	if err != nil {
		r.log.Error("get user by login", zap.Error(err))
		return nil, err
	}

	err = txConn.Commit(ctx)
	if err != nil {
		r.log.Error("commit tx", zap.Error(err))
		err := txConn.Rollback(ctx)
		if err != nil {
			r.log.Error("rollback tx", zap.Error(err))
			return nil, err
		}
		return nil, err
	}

	return convertUserToModel(&res), err
}
