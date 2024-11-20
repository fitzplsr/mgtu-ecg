package profilepostgres

import "github.com/fitzplsr/mgtu-ecg/internal/model"

func isChanged(old, new *model.User) bool {
	// doesnt check passHash
	switch {
	case old.Login != new.Login:
	case old.Name != new.Name:
	case old.Role != new.Role:
	case len(new.PasswordHash) > 0:
	default:
		return false
	}
	return true
}
