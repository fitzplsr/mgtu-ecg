package analysepostgres

import (
	models "github.com/fitzplsr/mgtu-ecg/gen"
	"github.com/fitzplsr/mgtu-ecg/internal/model"
)

func convertFileMetaToModel(metaDB *models.Filemeta) *model.FileInfo {
	return &model.FileInfo{
		ID:          int64(metaDB.ID),
		UserID:      metaDB.UserID.Bytes,
		Key:         metaDB.Key,
		Filename:    metaDB.Filename,
		Size:        metaDB.Size,
		Format:      model.FileFormat(metaDB.Format).String(),
		ContentType: metaDB.ContentType,
		CreatedAt:   metaDB.CreatedAt.Time,
		UpdatedAt:   metaDB.UpdatedAt.Time,
	}
}
