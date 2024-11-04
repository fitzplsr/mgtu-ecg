package refresh

import (
	"github.com/fitzplsr/mgtu-ecg/internal/model"
	"github.com/gofiber/fiber/v2"
)

const RefreshToken = "refresh_token"

func SetRefreshToken(c *fiber.Ctx, session *model.Session) {
	cookie := new(fiber.Cookie)
	cookie.Name = RefreshToken
	cookie.Value = session.RefreshToken
	cookie.Expires = session.CreatedAt.Add(session.ExpiresIn)
	cookie.HTTPOnly = true
	//cookie.Secure = true // Отправляем куки только по HTTPS (требуется в продакшене)
	cookie.SameSite = "Lax"

	c.Cookie(cookie)
}
