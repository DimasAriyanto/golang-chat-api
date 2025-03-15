// File path: /internal/usecase/user_usecase.go
package usecase

import (
	"github.com/DimasAriyanto/golang-chat-api/internal/domain"
	"github.com/DimasAriyanto/golang-chat-api/internal/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	UserRepo *repository.UserRepository
}

func NewUserUseCase(userRepo *repository.UserRepository) *UserUseCase {
	return &UserUseCase{UserRepo: userRepo}
}

func (uc *UserUseCase) RegisterUser(user domain.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.Password = string(hashedPassword)
	return uc.UserRepo.CreateUser(user)
}

func (uc *UserUseCase) LoginUser(username, password string) (domain.User, error) {
	user, err := uc.UserRepo.GetUserByUsername(username)
	if err != nil {
		return domain.User{}, errors.New("invalid username or password")
	}

	// Verifikasi password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, errors.New("invalid username or password")
	}

	return user, nil
}
