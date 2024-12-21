package server

import (
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/middleware"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/analyse/delivery/analysehttp"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/auth/delivery/authhttp"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/patients/delivery/patienthttp"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/profile/delivery/profilehttp"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ServerParams struct {
	fx.In

	App *fiber.App
}

type AppParams struct {
	fx.In

	Logger *zap.Logger
	Config Config

	// handlers
	AuthHandler     *authhttp.Auth
	ProfileHandler  *profilehttp.Profile
	AnalyseHandler  *analysehttp.Analyse
	PatientsHandler *patienthttp.Patients

	// middlewares
	ProtectedMW *middleware.ProtectedMW
	CORSMW      *middleware.CORS
}
