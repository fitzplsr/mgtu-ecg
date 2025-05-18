package main

import (
	"context"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/converter"
	"os"
	"os/signal"
	"syscall"

	"github.com/fitzplsr/mgtu-ecg/internal/pkg/analyser"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/auther"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/config"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/db"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/filestorage/fsstorage"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/metrics"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/middleware"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/pprof"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/server"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/analyse"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/analyse/analyseusecase"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/analyse/delivery/analysehttp"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/analyse/repo/analysepostgres"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/auth"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/auth/authusecase"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/auth/delivery/authhttp"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/patients"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/patients/delivery/patienthttp"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/patients/patientsusecase"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/patients/repo/postgrespatients"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/profile"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/profile/delivery/profilehttp"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/profile/profileusecase"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/profile/repo/profilepostgres"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/session/storage"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/validate"
	"github.com/fitzplsr/mgtu-ecg/migrations"
	"github.com/fitzplsr/mgtu-ecg/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

// TODO services to module?

// @title Backend API
// @description API server

// @securityDefinitions	AuthKey
// @in					header
// @name				Authorization
func main() {
	app := fx.New(
		fx.Provide(
			logger.Provide,
			server.NewFiberApp,
			config.MustLoad,
			//db.ConnectMongoDB,
			db.NewPostgresConn,
			db.NewPostgresPool,
			db.NewRedisClient,
			//db.NewMinioClient,

			// validator
			validate.New,

			// jwt
			fx.Annotate(auther.New,
				fx.As(new(middleware.JWTer)),
			),

			// redis
			fx.Annotate(storage.NewStorage,
				fx.As(new(middleware.SessionStorage)),
			),

			// filestorage
			// minio
			//fx.Annotate(miniostorage.New,
			//	fx.As(new(analyse.FileStorage)),
			//),
			// fs
			fx.Annotate(fsstorage.New,
				fx.As(new(analyse.FileStorage)),
			),

			// middlewares
			middleware.NewProtectMW,
			middleware.NewCORSMiddleware,

			// auth service
			authhttp.New,
			fx.Annotate(authusecase.New,
				fx.As(new(auth.Usecase)),
			),
			fx.Annotate(profilepostgres.New,
				fx.As(new(auth.Repo)),
			),

			// profile service
			profilehttp.New,
			fx.Annotate(profileusecase.New,
				fx.As(new(profile.Usecase)),
			),
			fx.Annotate(profilepostgres.New,
				fx.As(new(profile.Repo)),
			),

			// analyse service
			analysehttp.New,
			fx.Annotate(analyseusecase.New,
				fx.As(new(analyse.Usecase)),
			),
			fx.Annotate(analysepostgres.New,
				fx.As(new(analyse.Repo)),
			),

			fx.Annotate(analyser.New,
				fx.As(new(analyse.Analyser)),
			),
			fx.Annotate(converter.New,
				fx.As(new(analyse.Converter)),
			),

			// patients service
			patienthttp.New,
			fx.Annotate(patientsusecase.New,
				fx.As(new(patients.Usecase)),
			),
			fx.Annotate(postgrespatients.New,
				fx.As(new(patients.Repo)),
			),
		),
		fx.WithLogger(func(l *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: l}
		}),
		fx.Invoke(
			pprof.StartPprof,
			metrics.InvokeMetricsServer,
			server.RunFiberServer,
			migrations.RunMigrations,
		),
	)

	ctx := context.Background()

	// Signal handling for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start the application
	if err := app.Start(ctx); err != nil {
		panic(err)
	}

	<-stop
	app.Stop(ctx) // Stop the UberFX application

	// Wait for application to finish
	app.Wait()
}
