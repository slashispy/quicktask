package usecases

import (
	"fmt"

	"github.com/slashispy/quicktask/internal/models"
	"github.com/slashispy/quicktask/internal/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserUseCase struct {
	UserRepo  repositories.UserRepository
	Encrypter PasswordEncrypter
}

func (uc *UserUseCase) RegisterUser(input models.UserDTO) error {
	// Check if the user already exists
	existingUser, err := uc.UserRepo.FindByEmail(input.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	if existingUser != nil {
		return fmt.Errorf("user with email %s already exists", input.Email)
	}

	// Encrypt the password
	hashedPassword, err := uc.Encrypter.Encrypt(input.Password)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	// Create a new user model
	newUser := input.ToUser(hashedPassword)

	// Save the new user to the repository
	return uc.UserRepo.Save(newUser)
}
