package analysehttp

import (
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
// @Param file formData file true "File to upload"
// @Success 201 {object} model.FileInfo
// @Failure 500 {object} model.ErrorResponse
// @Router      /api/v1/analyse/upload [post]
func (a *Analyse) UploadFile(c *fiber.Ctx) error {
	userID, ok := c.Locals(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		a.log.Error("can't get user id from locals")
		return utils.Send401(c, messages.Unauthorized)
	}

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	fileInfo, err := a.uc.Upload(c.UserContext(), file, userID)
	if err != nil {
		return utils.Send500(c, messages.InternalServerError)
	}

	return c.Status(fiber.StatusCreated).JSON(fileInfo)
}
