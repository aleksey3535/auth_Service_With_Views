package service

import (
	"auth/internal/pkg"
	"auth/internal/pkg/repository"
	"crypto/sha1"
	"errors"
	"fmt"
	"unicode"

	
)

const PassSalt = "dapjioiASDjopsda"


type Service struct {
	repo *repository.Repository
	ValidateLogin string
	ValidatePassword string
}
func New(repo *repository.Repository) *Service {
	return &Service{
		repo : repo,
		ValidateLogin: "such a login already exist",
		ValidatePassword: "password failed validation",}
}

func(s *Service) CreateUser(login, password string) error {
	if !validatePass(password) {
		return errors.New(s.ValidatePassword)
	}
	if !s.validateLogin(login) {
		return errors.New(s.ValidateLogin)
	}
	hashPass := generatePasswordHash(password)
	return s.repo.CreateUser(login, hashPass)
}


func(s *Service) IsExist(login, password string) (bool, error) {
	passHash := generatePasswordHash(password)
	dbPassHash, err := s.repo.GetUserPass(login)
	if err != nil {
		return false, err
	}
	if passHash != dbPassHash {
		return false, nil
	}
	return true, nil
}

func(s *Service) GetAllUsers() ([]pkg.User, error) {
	return s.repo.GetAllUsers()
}

func(s *Service) DeleteUser(login string) error {
	return s.repo.DeleteUser(login)
}


func generatePasswordHash(password string) string {
	hasher := sha1.New()
	hasher.Write([]byte(password))
	return fmt.Sprintf("%x", hasher.Sum([]byte(PassSalt)))
}

func validatePass(password string) bool {
	var isUpper bool
	var isLetter bool
	var isDigit bool
	var isLen = len(password) > 5
	for _, val := range password {
		switch {
		case unicode.IsUpper(val):
			isUpper = true
		case unicode.IsLetter(val):
			isLetter = true
		case unicode.IsDigit(val):
			isDigit = true
		}
	}
	return isUpper && isLetter && isDigit && isLen
}

func(s *Service) validateLogin(login string) bool {
	return !s.repo.GetUserIfExist(login)


}