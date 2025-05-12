package analyseusecase

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/fitzplsr/mgtu-ecg/internal/model"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/filestorage"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/analyse"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/patients"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
)

var _ analyse.Usecase = (*Analyse)(nil)

type Params struct {
	fx.In

	Log          *zap.Logger
	FileStorage  analyse.FileStorage
	Repo         analyse.Repo
	PatientsRepo patients.Repo
	Analyser     analyse.Analyser
}

type Analyse struct {
	log *zap.Logger

	repo         analyse.Repo
	patientsRepo patients.Repo
	fileStorage  analyse.FileStorage
	analyser     analyse.Analyser
}

func New(p Params) *Analyse {
	return &Analyse{
		log:          p.Log,
		repo:         p.Repo,
		fileStorage:  p.FileStorage,
		analyser:     p.Analyser,
		patientsRepo: p.PatientsRepo,
	}
}

func (a *Analyse) Upload(ctx context.Context, fileHeader *multipart.FileHeader, patientID int) ([]*model.FileInfo, error) {
	open, err := fileHeader.Open()
	if err != nil {
		a.log.Error("open file", zap.Error(err))
		return nil, err
	}
	defer open.Close()

	data, err := io.ReadAll(open)
	if err != nil {
		a.log.Error("read file", zap.Error(err))
		return nil, err
	}

	contentType := fileHeader.Header.Get("Content-Type")
	filename := fileHeader.Filename

	if strings.HasSuffix(strings.ToLower(filename), ".zip") || contentType == "application/zip" {
		return a.processZip(ctx, data, patientID)
	}

	return a.saveSingleFile(ctx, data, filename, contentType, patientID)
}

func (a *Analyse) processZip(ctx context.Context, data []byte, patientID int) ([]*model.FileInfo, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		a.log.Error("create zip reader", zap.Error(err))
		return nil, err
	}

	var results []*model.FileInfo

	for _, file := range reader.File {
		if !strings.HasSuffix(strings.ToLower(file.Name), ".edf") {
			continue
		}

		if strings.HasPrefix(file.Name, "__MACOSX/") || strings.HasPrefix(filepath.Base(file.Name), "._") {
			continue
		}

		rc, err := file.Open()
		if err != nil {
			a.log.Error("open zip entry", zap.String("file", file.Name), zap.Error(err))
			continue
		}

		fileData, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			a.log.Error("read zip entry", zap.String("file", file.Name), zap.Error(err))
			continue
		}

		contentType := "application/octet-stream"

		info, err := a.saveSingleFile(ctx, fileData, file.Name, contentType, patientID)
		if err != nil {
			a.log.Error("save edf from zip", zap.String("file", file.Name), zap.Error(err))
			continue
		}

		results = append(results, info...)
	}

	return results, nil
}

func (a *Analyse) saveSingleFile(ctx context.Context, data []byte, filename, contentType string, patientID int) ([]*model.FileInfo, error) {
	buffer := bytes.NewBuffer(data)

	file := filestorage.File{
		Data:        buffer,
		Size:        int64(len(data)),
		Filename:    filename,
		ContentType: contentType,
		PatientID:   patientID,
	}

	key, err := a.fileStorage.Save(ctx, &file)
	if err != nil {
		a.log.Error("save file", zap.String("filename", filename), zap.Error(err))
		return nil, err
	}

	meta := model.FileMeta{
		PatientID:   patientID,
		Key:         key,
		Filename:    filename,
		Size:        int32(len(data)),
		Format:      model.EDF,
		ContentType: contentType,
	}

	savedMeta, err := a.repo.SaveFileMeta(ctx, &meta)
	if err != nil {
		a.log.Error("save file meta", zap.String("filename", filename), zap.Error(err))
		return nil, err
	}

	return []*model.FileInfo{savedMeta}, nil
}

func (a *Analyse) RunAnalyse(ctx context.Context, req *model.AnalyseRequest) (*model.AnalyseTask, error) {
	patientInfo, err := a.patientsRepo.GetByID(ctx, req.PatientID)
	if err != nil {
		return nil, errors.Join(err, analyse.ErrPatientNotExist)
	}

	fileInfo, err := a.repo.GetFileByID(ctx, req.FileID)
	if err != nil {
		return nil, errors.Join(err, analyse.ErrFileNotExist)
	}

	res, err := a.repo.CreateAnalyse(ctx, req.Name, int(fileInfo.ID), patientInfo.ID, model.Created)
	if err != nil {
		a.log.Error("failed to create analyse", zap.Error(err))
		return nil, fmt.Errorf("create analyse task: %w", err)
	}

	predictResult, err := a.analyser.Run(ctx, fileInfo.Filename)
	if err != nil {
		a.log.Error("analyse finished with error", zap.Error(err))
		result, err := a.repo.UpdateAnalyseStatus(ctx, res.ID, model.Failed)
		if err != nil {
			a.log.Error("failed to save failed analyse", zap.Error(err))
			return nil, fmt.Errorf("error fail analyse: %w", err)
		}
		return result, nil
	}

	res, err = a.repo.SaveAnalyseResult(
		ctx,
		res.ID,
		model.AnalyseResultFromBool(predictResult.Result),
		predictResult.Predict,
		model.Success,
	)
	if err != nil {
		a.log.Error("failed to save analyse result", zap.Error(err))
		return nil, fmt.Errorf("failed to save analyse result: %w", err)
	}

	return res, nil
}

func (a *Analyse) ListPatientFiles(ctx context.Context, payload *model.ListPatientFilesRequest) (*model.PatientFiles, error) {
	filter := payload.Filter
	filter.AlignLimit()

	files, err := a.repo.ListPatientFiles(ctx, payload.PatientID, &filter)
	if err != nil {
		return nil, err
	}
	return files, err
}

func (a *Analyse) ListPatientAnalyses(ctx context.Context, payload *model.ListPatientAnalysesRequest) (*model.AnalyseTasks, error) {
	filter := payload.Filter
	filter.AlignLimit()
	res, err := a.repo.ListPatientAnalyses(ctx, payload.PatientID, &filter)
	if err != nil {
		return nil, err
	}
	return res, err
}
