package analysepostgres

import (
	"context"
	models "github.com/fitzplsr/mgtu-ecg/gen"
	"github.com/fitzplsr/mgtu-ecg/internal/model"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/analyse"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/utils/txer"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var _ analyse.Repo = &Analyse{}

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

func (a *Analyse) SaveFileMeta(ctx context.Context, meta *model.FileMeta) (*model.FileInfo, error) {
	txConn, err := txer.Tx(ctx, a.conn)
	if err != nil {
		a.log.Error("begin tx", zap.Error(err))
		return nil, err
	}

	tx := a.db.WithTx(txConn)

	res, err := tx.CreateFileMeta(ctx, models.CreateFileMetaParams{
		Format:      int16(meta.Format),
		Size:        meta.Size,
		Filename:    meta.Filename,
		ContentType: meta.ContentType,
		Key:         meta.Key,
		PatientID: pgtype.Int4{
			Int32: int32(meta.PatientID),
			Valid: true,
		},
	})
	if err != nil {
		a.log.Error("save file meta", zap.Error(err))
		return nil, err
	}
	err = txConn.Commit(ctx)
	if err != nil {
		txConn.Rollback(ctx)
		return nil, err
	}
	return convertFileMetaToModel(&res), nil
}

func (a *Analyse) ListPatientFiles(ctx context.Context, patientID int, filter *model.Filter) (*model.PatientFiles, error) {
	files, err := a.db.GetPatientFileMetas(ctx, models.GetPatientFileMetasParams{
		PatientID: pgtype.Int4{
			Int32: int32(patientID),
			Valid: true,
		},
		Limit:  int32(filter.Limit),
		Offset: int32(filter.Offset),
	})
	if err != nil {
		return nil, err
	}

	res := make([]*model.FileInfo, 0, len(files))
	for _, f := range files {
		res = append(res, convertFileMetaToModel(&f))
	}

	return &model.PatientFiles{PatientID: patientID, Files: res}, nil
}

func (a *Analyse) ListPatientAnalyses(ctx context.Context, patientID int, filter *model.Filter) (*model.AnalyseTasks, error) {
	res, err := a.db.ListAnalyseTasksByPatientID(ctx, models.ListAnalyseTasksByPatientIDParams{
		PatientID: pgtype.Int4{
			Int32: int32(patientID),
			Valid: true,
		},
		Limit:  int32(filter.Limit),
		Offset: int32(filter.Offset),
	})
	if err != nil {
		return nil, err
	}

	converted := make([]*model.AnalyseTask, 0, len(res))
	for _, f := range res {
		converted = append(converted, convertAnalyseTaskToModel(&f))
	}

	return &model.AnalyseTasks{Analyses: converted}, nil
}

func (a *Analyse) GetFileByID(ctx context.Context, fileID int) (*model.FileInfo, error) {
	file, err := a.db.GetFileMetaById(ctx, int32(fileID))
	if err != nil {
		return nil, err
	}

	return convertFileMetaToModel(&file), nil
}

func (a *Analyse) CreateAnalyse(ctx context.Context, name string, fileID int, patientID int, status model.AnalyseStatus) (*model.AnalyseTask, error) {
	res, err := a.db.CreateAnalyseTask(ctx, models.CreateAnalyseTaskParams{
		Name: name,
		PatientID: pgtype.Int4{
			Int32: int32(patientID),
			Valid: true,
		},
		FilemetaID: pgtype.Int4{
			Int32: int32(fileID),
			Valid: true,
		},
		Status: int16(status),
	})
	if err != nil {
		return nil, err
	}

	return convertAnalyseTaskToModel(&res), nil
}

func (a *Analyse) UpdateAnalyseStatus(ctx context.Context, id int, status model.AnalyseStatus) (*model.AnalyseTask, error) {
	res, err := a.db.SetAnalyseTaskStatus(ctx, models.SetAnalyseTaskStatusParams{
		Status: int16(status),
		ID:     int32(id),
	})
	if err != nil {
		return nil, err
	}

	return convertAnalyseTaskToModel(&res), nil
}

func (a *Analyse) SaveAnalyseResult(ctx context.Context, id int, result model.AnalyseResult, predict string, status model.AnalyseStatus) (*model.AnalyseTask, error) {
	res, err := a.db.SaveAnalyseTaskResult(ctx, models.SaveAnalyseTaskResultParams{
		Result: int16(result),
		Predict: pgtype.Text{
			String: predict,
			Valid:  true,
		},
		Status: int16(status),
		ID:     int32(id),
	})
	if err != nil {
		return nil, err
	}

	return convertAnalyseTaskToModel(&res), nil
}
