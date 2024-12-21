package postgrespatients

import (
	"context"
	"fmt"
	models "github.com/fitzplsr/mgtu-ecg/gen"
	"github.com/fitzplsr/mgtu-ecg/internal/model"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/patients"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var _ patients.Repo = &Patients{}

type PostgresParams struct {
	fx.In

	Log  *zap.Logger
	Conn *pgxpool.Pool
}

type Patients struct {
	log  *zap.Logger
	db   *models.Queries
	conn *pgxpool.Pool
}

func New(p PostgresParams) *Patients {
	return &Patients{
		log:  p.Log,
		db:   models.New(p.Conn),
		conn: p.Conn,
	}
}

func (p *Patients) Create(ctx context.Context, patient *model.CreatePatient) (*model.Patient, error) {
	res, err := p.db.CreatePatient(ctx, models.CreatePatientParams{
		Name:    patient.Name,
		Surname: patient.Surname,
		Bdate: pgtype.Date{
			Time:             patient.Birthday,
			InfinityModifier: 0,
			Valid:            true,
		},
	})
	if err != nil {
		p.log.Error("create patient", zap.Error(err))
		return nil, err
	}

	return convertPatientToModel(&res), nil
}

func (p *Patients) GetByID(ctx context.Context, id int) (*model.Patient, error) {
	res, err := p.db.GetPatientByID(ctx, int32(id))
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}

	return convertPatientToModel(&res), nil
}

func (p *Patients) List(ctx context.Context, filter *model.Filter) (*model.Patients, error) {
	res, err := p.db.ListPatients(ctx, models.ListPatientsParams{
		Limit:  int32(filter.Limit),
		Offset: int32(filter.Offset),
	})
	if err != nil {
		return nil, fmt.Errorf("list patients from postgres: %w", err)
	}

	converted := make([]*model.Patient, 0, len(res))
	for _, p := range res {
		converted = append(converted, convertPatientToModel(&p))
	}

	return &model.Patients{Patients: converted}, nil
}

func (p *Patients) Search(ctx context.Context, filter *model.Filter) (*model.Patients, error) {
	res, err := p.db.SearchPatient(
		ctx,
		models.SearchPatientParams{
			Column1: pgtype.Text{
				String: filter.Search,
				Valid:  true,
			},
			Limit:  int32(filter.Limit),
			Offset: int32(filter.Offset),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("search patients from postgres: %w", err)
	}

	converted := make([]*model.Patient, 0, len(res))
	for _, p := range res {
		converted = append(converted, convertPatientToModel(&p))
	}

	return &model.Patients{Patients: converted}, nil
}
