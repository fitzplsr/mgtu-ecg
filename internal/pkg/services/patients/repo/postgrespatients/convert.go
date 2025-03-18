package postgrespatients

import (
	models "github.com/fitzplsr/mgtu-ecg/gen"
	"github.com/fitzplsr/mgtu-ecg/internal/model"
)

func convertPatientToModel(patientDB *models.Patient) *model.Patient {
	return &model.Patient{
		ID:        int(patientDB.ID),
		Name:      patientDB.Name,
		Surname:   patientDB.Surname,
		Birthday:  patientDB.Bdate.Time,
		CreatedAt: patientDB.CreatedAt.Time,
		UpdatedAt: patientDB.UpdatedAt.Time,
	}
}
