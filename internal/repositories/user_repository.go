package repositories

import (
	"github.com/slashispy/quicktask/internal/models"
)

type UserRepository interface {
	Save(user models.User) error
	FindByEmail(email string) (*models.User, error)
}
