package authhttp

import (
	"errors"
	"github.com/fitzplsr/mgtu-ecg/internal/model"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/middleware"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/refresh"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/auth"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/profile"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/utils"
	"github.com/fitzplsr/mgtu-ecg/pkg/consts"
	"github.com/fitzplsr/mgtu-ecg/pkg/messages"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mailru/easyjson"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In

	Usecase   auth.Usecase
	Logger    *zap.Logger
	Validator *validator.Validate
}

type Auth struct {
	logger    *zap.Logger
	uc        auth.Usecase
	validator *validator.Validate
}

func New(p Params) (*Auth, error) {
	return &Auth{
		logger: p.Logger.With(
			zap.String(consts.Service, auth.AuthService),
			zap.String(consts.Layer, consts.DeliveryLayer),
		),
		uc:        p.Usecase,
		validator: p.Validator,
	}, nil
}

// TODO таймаут в UserContext

// @Summary     Sign up a new user
// @Description Creates a new user account
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body  body model.SignUpPayload true "User registration details"
// @Success     201   {object} model.AuthResponse
// @Failure     400   {object} model.ErrorResponse
// @Failure     409   {object} model.ErrorResponse
// @Failure     500   {object} model.ErrorResponse
// @Router      /api/v1/auth/signup [post]
func (h *Auth) SignUp(c *fiber.Ctx) error {
	log := h.logger.With(
		zap.String(consts.RequestID, c.Get(consts.RequestID)), // or c.Context().ID()
	)
	log.Debug("start")
	var payload model.SignUpPayload
	err := c.BodyParser(&payload)
	if err != nil {
		log.Error("failed to parse request body", zap.Error(err))
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err = h.validator.Struct(payload); err != nil {
		return utils.Send400(c, messages.InvalidPayload)
	}

	authResp, session, err := h.uc.SignUp(c.UserContext(), &payload)
	if err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			return c.SendStatus(fiber.StatusConflict)
		}
		return utils.Send500(c, messages.InternalServerError)
	}

	refresh.SetRefreshToken(c, session)

	return c.Status(fiber.StatusCreated).JSON(authResp)
}

// @Summary Sign In
// @Description Authenticate user
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body model.SignInPayload true "Login credentials"
// @Success 200 {object} model.AuthResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *Auth) SignIn(c *fiber.Ctx) error {
	log := h.logger.With(
		zap.String(consts.RequestID, c.Get(consts.RequestID)), // or c.Context().ID()
	)
	log.Debug("start")
	var payload model.SignInPayload
	err := c.BodyParser(&payload)
	if err != nil {
		log.Error("failed to parse request body", zap.Error(err))
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err = h.validator.Struct(payload); err != nil {
		return utils.Send400(c, messages.InvalidPayload)
	}

	authResp, session, err := h.uc.SignIn(c.UserContext(), &payload)
	switch {
	case errors.Is(err, profile.ErrProfileNotFound):
		log.Error("sign in", zap.Error(err))
		return utils.Send401(c, messages.WrongLoginOrPassword)
	case errors.Is(err, auth.ErrWrongPassword):
		log.Error("sign in", zap.Error(err))
		return utils.Send401(c, messages.WrongLoginOrPassword)
	case err != nil:
		return utils.Send500(c, messages.InternalServerError)
	}

	_, err = easyjson.Marshal(authResp)
	if err != nil {
		return utils.Send500(c, messages.InternalServerError)
	}

	refresh.SetRefreshToken(c, session)

	return c.Status(fiber.StatusOK).JSON(authResp)
}

// @Summary Update Password
// @Description Change user password
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body model.UpdatePasswordPayload true "New password details"
// @Success 200 {object} model.AuthResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/auth/update/password [put]
func (h *Auth) UpdatePassword(c *fiber.Ctx) error {
	log := h.logger.With(
		zap.String(consts.RequestID, c.Get(consts.RequestID)), // or c.Context().ID()
	)
	log.Debug("start")
	userID, ok := c.Locals(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		log.Error("user id from locals")
		return utils.Send401(c, messages.Unauthorized)
	}

	var payload model.UpdatePasswordPayload
	err := c.BodyParser(&payload)
	if err != nil {
		log.Error("failed to parse request body", zap.Error(err))
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err = h.validator.Struct(payload); err != nil {
		return utils.Send400(c, messages.InvalidPayload)
	}

	payload.ID = userID
	authResp, session, err := h.uc.UpdatePassword(c.UserContext(), &payload)
	switch {
	case errors.Is(err, auth.ErrSamePassword):
		return c.SendStatus(fiber.StatusOK)
	case errors.Is(err, profile.ErrProfileNotFound):
		return utils.Send404(c, messages.NotFound)
	case err != nil:
		return utils.Send500(c, messages.InternalServerError)
	}

	refresh.SetRefreshToken(c, session)

	//return utils.SendBody(c, fiber.StatusOK, res)
	return c.Status(fiber.StatusOK).JSON(authResp)
}

// @Summary Update Role
// @Description Change user role
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body model.UpdateRolePayload true "New role details"
// @Success 200 {object} model.AuthResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/auth/update/role [put]
func (h *Auth) UpdateRole(c *fiber.Ctx) error {
	log := h.logger.With(
		zap.String(consts.RequestID, c.Get(consts.RequestID)), // or c.Context().ID()
	)
	log.Debug("start")

	userID, ok := c.Locals(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		log.Error("user id from locals")
		return utils.Send401(c, messages.Unauthorized)
	}

	var payload model.UpdateRolePayload
	err := c.BodyParser(&payload)
	if err != nil {
		log.Error("failed to parse request body", zap.Error(err))
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err = h.validator.Struct(payload); err != nil {
		return utils.Send400(c, messages.InvalidPayload)
	}

	payload.ID = userID
	authResp, session, err := h.uc.UpdateRole(c.UserContext(), &payload)
	if err != nil {
		if errors.Is(err, auth.ErrWrongRole) {
			return utils.Send400(c, messages.InvalidRole)
		}
		return utils.Send500(c, messages.InternalServerError)
	}

	refresh.SetRefreshToken(c, session)

	return c.Status(fiber.StatusOK).JSON(authResp)
}
