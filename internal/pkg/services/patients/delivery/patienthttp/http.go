package patienthttp

import (
	"github.com/fitzplsr/mgtu-ecg/internal/model"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/patients"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/utils"
	"github.com/fitzplsr/mgtu-ecg/pkg/consts"
	"github.com/fitzplsr/mgtu-ecg/pkg/messages"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In

	Usecase   patients.Usecase
	Logger    *zap.Logger
	Validator *validator.Validate
}

type Patients struct {
	logger    *zap.Logger
	uc        patients.Usecase
	validator *validator.Validate
}

func New(p Params) (*Patients, error) {
	return &Patients{
		logger: p.Logger.With(
			zap.String(consts.Service, patients.PatientsService),
			zap.String(consts.Layer, consts.DeliveryLayer),
		),
		uc:        p.Usecase,
		validator: p.Validator,
	}, nil
}

// @Summary     Create a new patient
// @Description Creates a new patient
// @Security BearerAuth
// @Tags        patients
// @Accept      json
// @Produce     json
// @Param       body  body model.CreatePatient true "Patient creation details"
// @Success     201   {object} model.Patient
// @Failure     400   {object} model.ErrorResponse
// @Failure     500   {object} model.ErrorResponse
// @Router      /api/v1/patients/create [post]
func (p *Patients) Create(c *fiber.Ctx) error {
	log := p.logger.With(
		zap.String(consts.RequestID, c.Get(consts.RequestID)), // or c.Context().ID()
	)

	var payload model.CreatePatient
	err := c.BodyParser(&payload)
	if err != nil {
		log.Error("failed to parse request body", zap.Error(err))
		return c.SendStatus(fiber.StatusBadRequest)
	}

	//if err = h.validator.Struct(payload); err != nil {
	//	return utils.Send400(c, messages.InvalidPayload)
	//}

	patient, err := p.uc.Create(c.Context(), &payload)
	if err != nil {
		return utils.Send500(c, messages.InternalServerError)
	}

	return c.Status(fiber.StatusCreated).JSON(patient)
}

// @Summary     List patients
// @Description List patients
// @Tags        patients
// @Security BearerAuth
// @Accept      json
// @Produce     json
// @Param       body  body model.Filter true "filter list"
// @Success     200   {object} model.Patients
// @Failure     400   {object} model.ErrorResponse
// @Failure     500   {object} model.ErrorResponse
// @Router      /api/v1/patients/list [put]
func (p *Patients) List(c *fiber.Ctx) error {
	log := p.logger.With(
		zap.String(consts.RequestID, c.Get(consts.RequestID)),
	)

	var payload model.Filter
	err := c.BodyParser(&payload)
	if err != nil {
		log.Error("failed to parse request body", zap.Error(err))
		return c.SendStatus(fiber.StatusBadRequest)
	}

	patientsList, err := p.uc.List(c.Context(), &payload)
	if err != nil {
		log.Error("error list patients", zap.Error(err))
		return utils.Send500(c, messages.InternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(patientsList)
}

// @Summary     Get
// @Description Get patient by id
// @Tags        patients
// @Security BearerAuth
// @Accept      json
// @Produce     json
// @Param       body  body model.GetPatient true "patient params"
// @Success     200   {object} model.Patient
// @Failure     400   {object} model.ErrorResponse
// @Failure     500   {object} model.ErrorResponse
// @Router      /api/v1/patients/ [put]
func (p *Patients) Get(c *fiber.Ctx) error {
	log := p.logger.With(
		zap.String(consts.RequestID, c.Get(consts.RequestID)),
	)

	var payload model.GetPatient
	err := c.BodyParser(&payload)
	if err != nil {
		log.Error("failed to parse request body", zap.Error(err))
		return c.SendStatus(fiber.StatusBadRequest)
	}

	patientsList, err := p.uc.Get(c.Context(), &payload)
	if err != nil {
		log.Error("failed to list patients", zap.Error(err))
		return utils.Send500(c, messages.InternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(patientsList)
}
