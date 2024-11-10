package analysepostgres

import (
	"context"
	models "github.com/fitzplsr/mgtu-ecg/gen"
	"github.com/fitzplsr/mgtu-ecg/internal/model"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/analyse"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/utils/pghelper"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/utils/txer"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var _ analyse.Repo = (*Analyse)(nil)

type PostgresParams struct {
	fx.In

	Log  *zap.Logger
	Conn *pgxpool.Pool
}

type Analyse struct {
	log  *zap.Logger
	db   *models.Queries
	conn *pgxpool.Pool
}

func New(p PostgresParams) *Analyse {
	return &Analyse{
		log:  p.Log,
		db:   models.New(p.Conn),
		conn: p.Conn,
	}
}

func (r *Analyse) SaveFileMeta(ctx context.Context, meta *model.FileMeta) (*model.FileInfo, error) {
	txConn, err := txer.Tx(ctx, r.conn)
	if err != nil {
		r.log.Error("begin tx", zap.Error(err))
		return nil, err
	}

	tx := r.db.WithTx(txConn)

	res, err := tx.CreateFileMeta(ctx, models.CreateFileMetaParams{
		Format:      int16(meta.Format),
		Size:        meta.Size,
		Filename:    meta.Filename,
		ContentType: meta.ContentType,
		Key:         meta.Key,
		UserID:      pghelper.ToPGUUID(meta.UserID),
	})
	log.Debug("res", zap.Any("res", res))
	if err != nil {
		r.log.Error("save file meta", zap.Error(err))
		return nil, err
	}
	err = txConn.Commit(ctx)
	if err != nil {
		txConn.Rollback(ctx)
		return nil, err
	}
	return convertFileMetaToModel(&res), nil
}
