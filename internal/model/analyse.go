package model

import "time"

//easyjson:json
type AnalyseTask struct {
	ID        int           `json:"id"`
	Name      string        `json:"name"`
	Result    AnalyseResult `json:"result"`
	Predict   string        `json:"predict"`
	FileID    int           `json:"file_id"`
	PatientID int           `json:"patient_id"`
	Status    AnalyseStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

//easyjson:json
type AnalyseTasks struct {
	Analyses []*AnalyseTask `json:"analyses"`
}

//easyjson:json
type AnalyseRequest struct {
	Name      string `json:"name"`
	FileID    int    `json:"file_id"`
	PatientID int    `json:"patient_id"`
}

//easyjson:json
type ListPatientFilesRequest struct {
	PatientID int    `json:"patient_id"`
	Filter    Filter `json:"filter"`
}

//easyjson:json
type ListPatientAnalysesRequest struct {
	PatientID int    `json:"patient_id"`
	Filter    Filter `json:"filter"`
}
