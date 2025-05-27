package analyse

import (
	"context"
	"mime/multipart"

	"github.com/fitzplsr/mgtu-ecg/internal/model"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/filestorage"
)

type Usecase interface {
	Upload(ctx context.Context, file *multipart.FileHeader, patientID int) ([]*model.FileInfo, error)
	ListPatientFiles(ctx context.Context, payload *model.ListPatientFilesRequest) (*model.PatientFiles, error)
	ListPatientAnalyses(ctx context.Context, payload *model.ListPatientAnalysesRequest) (*model.AnalyseTasks, error)
	RunAnalyse(ctx context.Context, req *model.AnalyseRequest) (*model.AnalyseTasks, error)
	GetFileByID(ctx context.Context, payload *model.GetFileByIDRequest) (*model.FileInfo, error)
}

type Repo interface {
	SaveAnalyseResult(ctx context.Context, id int, result model.AnalyseResult, predict string, status model.AnalyseStatus) (*model.AnalyseTask, error)
	CreateAnalyse(ctx context.Context, name string, fileID int, patientID int, status model.AnalyseStatus) (*model.AnalyseTask, error)
	UpdateAnalyseStatus(ctx context.Context, id int, status model.AnalyseStatus) (*model.AnalyseTask, error)
	GetFileByID(ctx context.Context, fileID int) (*model.FileInfo, error)
	SaveFileMeta(ctx context.Context, meta *model.FileMeta) (*model.FileInfo, error)
	ListPatientFiles(ctx context.Context, patientID int, filter *model.Filter) (*model.PatientFiles, error)
	ListPatientAnalyses(ctx context.Context, patientID int, filter *model.Filter) (*model.AnalyseTasks, error)
}

type Analyser interface {
	Run(ctx context.Context, filename string) (*model.InternalAnalyseResult, error)
}

type Converter interface {
	Convert(ctx context.Context, filename string) ([]byte, error)
}

type FileStorage interface {
	Save(ctx context.Context, file *filestorage.File) (string, error)
}
