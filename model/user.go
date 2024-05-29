package model

import (
	"errors"
	"strings"
	"wechat-server/common"
)

type User struct {
	Id               int    `json:"id"`
	Username         string `json:"username" gorm:"unique;uniqueIndex" validate:"printascii"`
	Password         string `json:"password" gorm:"not null;" validate:"min=8"`
	DisplayName      string `json:"display_name"`
	Role             int    `json:"role" gorm:"type:int;default:1"`   // admin, common
	Status           int    `json:"status" gorm:"type:int;default:1"` // enabled, disabled
	Token            string `json:"token" gorm:"index"`
	Email            string `json:"email" gorm:"index"`
	VerificationCode string `json:"verification_code" gorm:"-:all"`
}

func GetAllUsers() (users []*User, err error) {
	err = DB.Select([]string{"id", "username", "display_name", "role", "status", "email"}).Find(&users).Error
	return users, err
}

func GetUserById(id int, selectAll bool) (*User, error) {
	user := User{Id: id}
	var err error = nil
	if selectAll {
		err = DB.First(&user, "id = ?", id).Error
	} else {
		err = DB.Select([]string{"id", "username", "display_name", "role", "status", "email"}).First(&user, "id = ?", id).Error
	}
	return &user, err
}

func DeleteUserById(id int) (err error) {
	user := User{Id: id}
	err = DB.Delete(&user).Error
	return err
}

func QueryUsers(query string, startIdx int) (users []*User, err error) {
	query = strings.ToLower(query)
	err = DB.Limit(common.ItemsPerPage).Offset(startIdx).Where("username LIKE ? or display_name LIKE ?", "%"+query+"%", "%"+query+"%").Order("id desc").Find(&users).Error
	return users, err
}

func (user *User) Insert() error {
	var err error
	if user.Password != "" {
		user.Password, err = common.Password2Hash(user.Password)
		if err != nil {
			return err
		}
	}
	err = DB.Create(user).Error
	return err
}

func (user *User) Update(updatePassword bool) error {
	var err error
	if updatePassword {
		user.Password, err = common.Password2Hash(user.Password)
		if err != nil {
			return err
		}
	}
	err = DB.Model(user).Updates(user).Error
	return err
}

func (user *User) Delete() error {
	var err error
	err = DB.Delete(user).Error
	return err
}

// ValidateAndFill check password & user status
func (user *User) ValidateAndFill() (err error) {
	// When querying with struct, GORM will only query with non-zero fields,
	// that means if your field’s value is 0, '', false or other zero values,
	// it won’t be used to build query conditions
	password := user.Password
	DB.Where(User{Username: user.Username}).First(user)
	okay := common.ValidatePasswordAndHash(password, user.Password)
	if !okay || user.Status != common.UserStatusEnabled {
		return errors.New("用户名或密码错误，或者该用户已被封禁")
	}
	return nil
}

func (user *User) FillUserByEmail() {
	DB.Where(User{Email: user.Email}).First(user)
}

func (user *User) FillUserByUsername() {
	DB.Where(User{Username: user.Username}).First(user)
}

func ValidateUserToken(token string) (user *User) {
	if token == "" {
		return nil
	}
	token = strings.Replace(token, "Bearer ", "", 1)
	user = &User{}
	if DB.Where("token = ?", token).First(user).RowsAffected == 1 {
		return user
	}
	return nil
}

func IsEmailAlreadyTaken(email string) bool {
	return DB.Where("email = ?", email).Find(&User{}).RowsAffected == 1
}

func IsUsernameAlreadyTaken(username string) bool {
	return DB.Where("username = ?", username).Find(&User{}).RowsAffected == 1
}

func ResetUserPasswordByEmail(email string, password string) error {
	hashedPassword, err := common.Password2Hash(password)
	if err != nil {
		return err
	}
	err = DB.Model(&User{}).Where("email = ?", email).Update("password", hashedPassword).Error
	return err
}
