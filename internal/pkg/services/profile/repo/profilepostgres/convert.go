package profilepostgres

import (
	models "github.com/fitzplsr/mgtu-ecg/gen"
	"github.com/fitzplsr/mgtu-ecg/internal/model"
)

func convertUserToModel(userDB *models.User) *model.User {
	return &model.User{
		ID:           userDB.ID.Bytes,
		Role:         model.Role(userDB.Role),
		Name:         userDB.Name,
		Login:        userDB.Login,
		PasswordHash: userDB.PasswordHash,
		CreatedAt:    userDB.CreatedAt.Time,
		UpdatedAt:    userDB.UpdatedAt.Time,
	}
}
