package models

import (
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

const (
	NormalUser = 1
	AdminUser = 3
	ManagerUser = 2
)

type User struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Phone 	 string `json:"phone"`
	Email    string `json:"email"`
	LoginOn  string `json:"loginOn"`
	Role     int    `json:"role"`
}

// CheckAuth checks if authentication information exists
func CheckLogin(username, password string) (bool, error) {
	var auth User
	err := db.Select("id").Where(User{Username: username, Password: password}).First(&auth).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	return auth.ID > 0, nil
}

// ExistUserByPhone check user exist by phone
func ExistUserByPhone(phone string) (bool, error) {
	var auth User
	err := db.Select("id").Where("phone = ? AND deleted_on = ?", phone, 0).Find(&auth).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	return auth.ID > 0, nil
}

func AddUser(username , password, phone, email string) error {
	user := User{
		Username: username,
		Password: password,
		Phone:    phone,
		Email:    email,
		Role:     NormalUser,
	}
	if err := db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func GetUsers(pageNum ,pageSize int, maps interface{}) ([]User, error) {
	var (
		users []User = nil
		err error = nil
	)

	if pageSize > 0 && pageNum > 0 {
		err = db.Where(maps).Find(&users).Offset(pageNum).Limit(pageSize).Error
	} else {
		err = db.Where(maps).Find(&users).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return users, nil
}

func GetUserTotal(maps interface{}) (int, error) {
	var count int
	if err := db.Model(&User{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func DeleteUser(id int) error {
	return db.Where("id = ?", id).Update("deleted_on", 1).Error
	//if err := db.Where("id = ?", id).Delete(&User{}).Error; err != nil {
	//	return err
	//}
	//return nil
}


func UpdateLogin(id int) error {
	now := time.Now().Unix()
	return db.Model(&User{ID: id}).Update("loginOn", strconv.FormatInt(now, 10)).Error
}
