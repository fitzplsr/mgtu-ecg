package profilehttp

import (
	"github.com/fitzplsr/mgtu-ecg/internal/model"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/middleware"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/profile"
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

	Usecase   profile.Usecase
	Logger    *zap.Logger
	Validator *validator.Validate
}

type Profile struct {
	logger    *zap.Logger
	uc        profile.Usecase
	validator *validator.Validate
}

func New(p Params) (*Profile, error) {
	return &Profile{
		logger: p.Logger.With(
			zap.String(consts.Service, profile.ProfileService),
			zap.String(consts.Layer, consts.DeliveryLayer),
		),
		uc:        p.Usecase,
		validator: p.Validator,
	}, nil
}

// TODO таймаут в UserContext
func (h *Profile) Update(c *fiber.Ctx) error {
	log := h.logger.With(
		zap.String(consts.RequestID, c.Get(consts.RequestID)), // or c.Context().ID()
	)
	log.Debug("start")

	userID, ok := c.Locals(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		log.Error("user id from locals")
		return utils.Send401(c, messages.Unauthorized)
	}

	var payload model.UpdateUserPayload
	err := c.BodyParser(&payload)
	if err != nil {
		log.Error("failed to parse request body", zap.Error(err))
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err = h.validator.Struct(payload); err != nil {
		return utils.Send400(c, messages.InvalidPayload)
	}

	payload.ID = userID
	updatedUser, err := h.uc.Update(c.UserContext(), &payload)
	if err != nil {
		return utils.Send500(c, messages.InternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(updatedUser)
}
