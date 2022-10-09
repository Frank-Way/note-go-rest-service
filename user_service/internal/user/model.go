package user

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Login    string `db:"login" json:"login"`
	Password string `db:"password" json:"password"`
}

type Users = []User

func (u *User) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return fmt.Errorf("password does not match")
	}
	return nil
}

func (u *User) GeneratePasswordHash() error {
	pwd, err := generatePasswordHash(u.Password)
	if err != nil {
		return err
	}
	u.Password = pwd
	return nil
}

func generatePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("error occured while hashing password: %w", err)
	}
	return string(hash), nil
}

type CreateUserDTO struct {
	Login          string `json:"login"`
	Password       string `json:"password"`
	RepeatPassword string `json:"repeat_password"`
}

type UpdateUserDTO struct {
	Login             string `json:"login"`
	Password          string `json:"password"`
	OldPassword       string `json:"old_password"`
	NewPassword       string `json:"new_password"`
	RepeatNewPassword string `json:"repeat_new_password"`
}

func NewUser(dto CreateUserDTO) User {
	return User{
		Login:    dto.Login,
		Password: dto.Password,
	}
}

func UpdateUser(dto UpdateUserDTO) User {
	return User{
		Login:    dto.Login,
		Password: dto.Password,
	}
}
