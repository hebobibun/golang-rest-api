package model

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name" form:"nama"`
	HP       string `json:"hp" form:"hp"`
	Email    string `json:"email" form:"email"`
	Password string `json:"pwd" form:"password"`
}

type UserModel struct {
	DB *gorm.DB
}

func (um *UserModel) Register(newUser User) (User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Bcrypt error : ", err.Error())
		return User{}, errors.New("Password process error")
	}
	newUser.Password = string(hashed)
	err = um.DB.Create(&newUser).Error
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}

func (um *UserModel) Login(email, password string) (User, error) {
	res := User{}

	if err := um.DB.Where("email = ?", email).First(&res).Error; err != nil {
		log.Println("Login query error : ", err.Error())
		return User{}, errors.New("Data not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(password)); err != nil {
		log.Println("Compare password error : ", err.Error())
		return User{}, errors.New("Email/password doesn't match")
	}

	return res, nil
}

func (um *UserModel) GetAll() ([]User, error) {
	res := []User{}
	if err := um.DB.Find(&res).Error; err != nil {
		log.Println("Get All query error : ", err.Error())
		return nil, err
	}

	return res, nil
}

func (um *UserModel) GetByID(id int) (User, error) {
	res := User{}
	if err := um.DB.Where("id = ?", id).First(&res).Error; err != nil {
		log.Println("Get user by ID query error : ", err.Error())
		return User{}, err
	}

	return res, nil
}

func (um *UserModel) Update(updatedUser User) (User, error) {
	qry := um.DB.Model(&User{}).Where("id = ?", updatedUser.ID).Updates(&updatedUser)
	err := qry.Error

	if err != nil {
		log.Println("Update user query error : ", err.Error())
		return User{}, nil
	}

	return updatedUser, nil
}

func (um *UserModel) Delete(id int) error {
	qry := um.DB.Delete(&User{}, id)

	affRow := qry.RowsAffected

	if affRow <= 0 {
		log.Println("No rows affected")
		return errors.New("Failed to delete, data not found")
	}

	err := qry.Error

	if err != nil {
		log.Println("Delete query error : ", err.Error())
		return errors.New("Unable to delete data")
	}

	return nil
}
