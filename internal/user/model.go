package user

import "golang.org/x/crypto/bcrypt"

type User struct {
	Id       int    `db:"id" json:"id"`
	Login    string `db:"login" json:"login"`
	Password string `db:"password" json:"password"`
	IsActive bool   `db:"is_active" json:"is_active"`
}

type Users = []User

func (u *User) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return err
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
		return "", err
	}
	return string(hash), nil
}

func NewUser(dto CreateUserDTO) User {
	return User{
		Login:    dto.Login,
		Password: dto.Password,
	}
}

func UpdateUser(id int, login string, dto UpdateUserDTO) User {
	return User{
		Id:       id,
		Login:    login,
		Password: dto.NewPassword,
	}
}
