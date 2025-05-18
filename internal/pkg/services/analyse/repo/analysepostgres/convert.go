package analysepostgres

import (
	models "github.com/fitzplsr/mgtu-ecg/gen"
	"github.com/fitzplsr/mgtu-ecg/internal/model"
)

func convertFileMetaToModel(metaDB *models.Filemeta) *model.FileInfo {
	return &model.FileInfo{
		ID:          int64(metaDB.ID),
		PatientID:   int(metaDB.PatientID.Int32),
		Key:         metaDB.Key,
		Filename:    metaDB.Filename,
		Size:        metaDB.Size,
		Format:      model.FileFormat(metaDB.Format).String(),
		ContentType: metaDB.ContentType,
		CreatedAt:   metaDB.CreatedAt.Time,
		UpdatedAt:   metaDB.UpdatedAt.Time,
		Data:        string(metaDB.Data),
	}
}

func convertAnalyseTaskToModel(db *models.AnalyseTask) *model.AnalyseTask {
	return &model.AnalyseTask{
		ID:        int(db.ID),
		Name:      db.Name,
		Result:    model.AnalyseResult(db.Result),
		Predict:   db.Predict.String,
		FileID:    int(db.FilemetaID.Int32),
		PatientID: int(db.PatientID.Int32),
		Status:    model.AnalyseStatus(db.Status),
		CreatedAt: db.CreatedAt.Time,
		UpdatedAt: db.UpdatedAt.Time,
	}
}
