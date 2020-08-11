package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Name         string `json:"name" gorm:"not null"`
	Email        string `json:"email" gorm:"not null;unique_index"`
	Password     string `json:"password,omitempty" gorm:"-"`
	PasswordHash string `json:"passwordhash,omitempty" gorm:"not null"`
}

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	us := UserService{
		db: db,
	}
	if !us.db.HasTable(&User{}) {
		us.db.CreateTable(&User{})
	}
	return &us
}

func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := db.First(&user).Error
	return &user, err
}

func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := db.First(&user).Error
	return &user, err
}

func (us *UserService) Create(user *User) error {
	if user.Password == "" {
		return errors.New("password cannot be empty")
	}
	pwBytes := []byte(user.Password)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return us.db.Create(user).Error
}

func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

func (us *UserService) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}

func (us *UserService) GetAllUsers() ([]User, error) {
	var users []User
	err := us.db.Find(&users).Error
	return users, err
}

func (us *UserService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password))
	if err != nil {
		return nil, err
	}

	return foundUser, nil
}
