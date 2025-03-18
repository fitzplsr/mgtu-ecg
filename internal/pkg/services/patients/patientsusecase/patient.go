package patientsusecase

import (
	"context"

	"github.com/fitzplsr/mgtu-ecg/internal/model"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/patients"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var _ patients.Usecase = &Patients{}

type Params struct {
	fx.In

	Log  *zap.Logger
	Repo patients.Repo
}

type Patients struct {
	log  *zap.Logger
	repo patients.Repo
}

func New(p Params) *Patients {
	return &Patients{
		log:  p.Log,
		repo: p.Repo,
	}
}

func (p *Patients) Create(ctx context.Context, payload *model.CreatePatient) (*model.Patient, error) {
	patient, err := p.repo.Create(ctx, payload)
	if err != nil {
		return nil, err
	}

	return patient, nil
}

func (p *Patients) List(ctx context.Context, filter *model.Filter) (*model.Patients, error) {
	filter.AlignLimit()

	if filter.Search != "" {
		return p.repo.Search(ctx, filter)
	}

	return p.repo.List(ctx, filter)
}

func (p *Patients) Get(ctx context.Context, payload *model.GetPatient) (*model.Patient, error) {
	return p.repo.GetByID(ctx, payload.ID)
}
