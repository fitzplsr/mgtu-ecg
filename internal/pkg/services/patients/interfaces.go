package patients

import (
	"context"
	"github.com/fitzplsr/mgtu-ecg/internal/model"
)

type Usecase interface {
	Create(ctx context.Context, patient *model.CreatePatient) (*model.Patient, error)
	List(ctx context.Context, filter *model.Filter) (*model.Patients, error)
	Get(ctx context.Context, payload *model.GetPatient) (*model.Patient, error)
}

type Repo interface {
	List(ctx context.Context, filter *model.Filter) (*model.Patients, error)
	Create(ctx context.Context, patient *model.CreatePatient) (*model.Patient, error)
	GetByID(ctx context.Context, id int) (*model.Patient, error)
	Search(ctx context.Context, filter *model.Filter) (*model.Patients, error)
}
