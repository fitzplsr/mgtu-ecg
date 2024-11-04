package server

import (
	_ "github.com/fitzplsr/mgtu-ecg/docs"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/mailru/easyjson"
)

func customJSONMarshal(v any) ([]byte, error) {
	return easyjson.Marshal(v.(easyjson.Marshaler))
}

func customJSONUnmarshal(data []byte, v any) error {
	return easyjson.Unmarshal(data, v.(easyjson.Unmarshaler))
}

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewFiberApp(p AppParams) (*fiber.App, error) {
	app := fiber.New(fiber.Config{
		JSONEncoder: customJSONMarshal,
		JSONDecoder: customJSONUnmarshal,
	})
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: p.Logger,
	}))
	app.Use(recover.New(recover.Config{
		Next:              nil,
		EnableStackTrace:  true,
		StackTraceHandler: nil,
	}))

	app.Use(p.CORSMW.MW)

	api := app.Group("api")
	v1 := api.Group("/v1")
	auth := v1.Group("/auth")
	{
		auth.Post("/signup", p.AuthHandler.SignUp)
		auth.Post("/login", p.AuthHandler.SignIn)
		authUpdate := auth.Group("/update")
		authUpdate.Use(p.ProtectedMW.MW)
		{
			authUpdate.Put("/password", p.AuthHandler.UpdatePassword)
			authUpdate.Put("/role", p.AuthHandler.UpdateRole)
		}
	}

	profile := v1.Group("profile")
	profile.Use(p.ProtectedMW.MW)
	{
		profile.Put("update", p.ProfileHandler.Update)
	}

	v1.Get("/swagger/*", swagger.HandlerDefault) // default
	//v1.Get("/swagger/*", p.ProtectedMW.MW, swagger.HandlerDefault) // default

	p.Logger.Info("start server")

	return app, nil
}
