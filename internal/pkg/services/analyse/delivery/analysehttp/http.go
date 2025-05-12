package analysehttp

import (
	"errors"
	"strconv"

	"github.com/fitzplsr/mgtu-ecg/internal/model"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/middleware"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/analyse"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/utils"
	"github.com/fitzplsr/mgtu-ecg/pkg/consts"
	"github.com/fitzplsr/mgtu-ecg/pkg/messages"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In

	Usecase   analyse.Usecase
	Logger    *zap.Logger
	Validator *validator.Validate
}

type Analyse struct {
	log       *zap.Logger
	uc        analyse.Usecase
	validator *validator.Validate
}

func New(p Params) (*Analyse, error) {
	return &Analyse{
		log: p.Logger.With(
			zap.String(consts.Service, analyse.AnalyseService),
			zap.String(consts.Layer, consts.DeliveryLayer),
		),
		uc:        p.Usecase,
		validator: p.Validator,
	}, nil
}

// TODO list files, download file

// @Summary Upload a file
// @Description Upload a file associated with the current user
// @Tags analyse
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce application/json
// @Param patient_id formData string true "Patient ID"
// @Param file formData file true "File to upload"
// @Success 201 {object} model.FileInfos
// @Failure 500 {object} model.ErrorResponse
// @Router      /api/v1/analyse/upload [post]
func (a *Analyse) UploadFile(c *fiber.Ctx) error {
	_, ok := c.Locals(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		a.log.Error("can't get user id from locals")
		return utils.Send401(c, messages.Unauthorized)
	}

	patientIDStr := c.FormValue("patient_id")
	patientID, err := strconv.Atoi(patientIDStr)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	a.log.Debug("Upload file for patient",
		zap.Int("patientID", patientID),
	)
	fileInfo, err := a.uc.Upload(c.Context(), file, patientID)
	if err != nil {
		switch {
		case errors.Is(err, analyse.ErrFileNotExist):
			return utils.Send404(c, messages.NotFound)
		case errors.Is(err, analyse.ErrPatientNotExist):
			return utils.Send404(c, messages.NotFound)
		default:
			return utils.Send500(c, messages.InternalServerError)
		}
	}

	return c.Status(fiber.StatusCreated).JSON(model.FileInfos(fileInfo))
}

// @Summary List patient files
// @Description List files associated with the current patient
// @Tags analyse
// @Security BearerAuth
// @Accept application/json
// @Produce application/json
// @Param payload body model.ListPatientFilesRequest true "patient info"
// @Success 200 {object} model.PatientFiles
// @Failure 500 {object} model.ErrorResponse
// @Router      /api/v1/analyse/list_edf [put]
func (a *Analyse) ListPatientFiles(c *fiber.Ctx) error {
	_, ok := c.Locals(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		a.log.Error("can't get user id from locals")
		return utils.Send401(c, messages.Unauthorized)
	}

	var payload model.ListPatientFilesRequest
	err := c.BodyParser(&payload)
	if err != nil {
		a.log.Error("failed to parse request body", zap.Error(err))
		return c.SendStatus(fiber.StatusBadRequest)
	}

	files, err := a.uc.ListPatientFiles(c.Context(), &payload)
	if err != nil {
		return utils.Send500(c, messages.InternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(files)
}

// @Summary Analyse file
// @Description Analyse file by it id
// @Tags analyse
// @Security BearerAuth
// @Accept json
// @Produce application/json
// @Param payload body model.AnalyseRequest true "New role details"
// @Success 200 {object} model.AnalyseTask
// @Failure 500 {object} model.ErrorResponse
// @Router      /api/v1/analyse/run [post]
func (a *Analyse) RunAnalyse(c *fiber.Ctx) error {
	// TODO check fk
	userID, ok := c.Locals(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		a.log.Error("can't get user id from locals")
		return utils.Send401(c, messages.Unauthorized)
	}

	a.log.Info("start run analyse req", zap.String("userID", userID.String()))

	var payload model.AnalyseRequest
	err := c.BodyParser(&payload)
	if err != nil {
		a.log.Error("failed to parse request body", zap.Error(err))
		return c.SendStatus(fiber.StatusBadRequest)
	}

	res, err := a.uc.RunAnalyse(c.Context(), &payload)
	if err != nil {
		a.log.Error("error run analyse", zap.Error(err))
		return utils.Send500(c, messages.InternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// @Summary List patient analyses
// @Description List analyses for patient by id
// @Tags analyse
// @Security BearerAuth
// @Accept json
// @Produce application/json
// @Param payload body model.ListPatientAnalysesRequest true "New role details"
// @Success 200 {object} model.AnalyseTasks
// @Failure 500 {object} model.ErrorResponse
// @Router      /api/v1/analyse/patient/list [put]
func (a *Analyse) ListPatientAnalyses(c *fiber.Ctx) error {
	_, ok := c.Locals(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		a.log.Error("can't get user id from locals")
		return utils.Send401(c, messages.Unauthorized)
	}

	var payload model.ListPatientAnalysesRequest
	err := c.BodyParser(&payload)
	if err != nil {
		a.log.Error("failed to parse request body", zap.Error(err))
		return c.SendStatus(fiber.StatusBadRequest)
	}

	res, err := a.uc.ListPatientAnalyses(c.Context(), &payload)
	if err != nil {
		a.log.Error("error list patient analyses", zap.Error(err))
		return utils.Send500(c, messages.InternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
