package models

type UserDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

func (dto *UserDTO) ToUser(hashedPassword string) User {
	return User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: hashedPassword,
	}
}
