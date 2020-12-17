package models

import (
	"fmt"
	"github.com/morgine/moon/pkg/rand"
	"github.com/morgine/moon/src/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID               int
	Username         string `gorm:"index"`
	Password         string
	GoogleAuthSecret string `json:"-"`
	Recommender      int    `gorm:"index"`
	IsBindGoogleAuth bool   `gorm:"index"`
	Avatar           string
}

func (m *Model) RegisterUser(username, password string, recommenderID int) (*User, error) {
	err := m.UserValidator.ValidUsername(username)
	if err != nil {
		return nil, err
	}
	err = m.UserValidator.ValidPassword(password)
	if err != nil {
		return nil, err
	}
	user := &User{}
	err = m.DB.Where("username=?", username).Select("id").First(user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if user.ID > 0 {
		return nil, errors.UsernameAlreadyRegistered
	}
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}
	user.Username = username
	user.Password = string(passwordBytes)
	user.Recommender = recommenderID
	// 生成谷歌验证码 secret
	user.GoogleAuthSecret = rand.Str(16)
	err = m.DB.Create(user).Error
	return user, err
}

func (m *Model) LoginUser(username, password string) (*User, error) {
	user, err := m.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.UsernameOrPasswordIncorrect
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			if err == bcrypt.ErrMismatchedHashAndPassword {
				return nil, errors.UsernameOrPasswordIncorrect
			} else {
				return nil, err
			}
		} else {
			return user, nil
		}
	}
}

func (m *Model) GetUserByUsername(username string) (*User, error) {
	user := &User{}
	err := m.DB.First(user, "username=?", username).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if user.ID > 0 {
		return user, nil
	} else {
		return nil, nil
	}
}

func (m *Model) GetUserByID(id int) (*User, error) {
	user := &User{}
	err := m.DB.First(user, "id=?", id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if user.ID > 0 {
		return user, nil
	} else {
		return nil, nil
	}
}

// 判断用户是否已绑定谷歌验证器
func (m *Model) IsBindGoogleAuth(username string) (bool, error) {
	user, err := m.GetUserByUsername(username)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, fmt.Errorf("用户[username=%s]不存在", username)
	} else {
		return user.IsBindGoogleAuth, nil
	}
}

// 获得谷歌验证器二维码地址
func (m *Model) GetGoogleAuthenticatorQRCodeUrl(loginUserID, with, height int) (string, error) {
	user, err := m.GetUserByID(loginUserID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", fmt.Errorf("用户[id=%d]不存在", loginUserID)
	} else {
		return m.GAC.GetQRCodeURI(user.GoogleAuthSecret, user.Username, with, height), nil
	}
}

// 检测谷歌验证码
func (m *Model) VerifyGoogleAuthCode(username, googleAuthCode string) error {
	user, err := m.GetUserByUsername(username)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("用户名 %s 不存在", username)
	}
	ok, err := m.GAC.Verify(user.GoogleAuthSecret, googleAuthCode)
	if err != nil {
		return err
	}
	if !ok {
		return errors.GoogleAuthCodeIncorrect
	} else {
		return nil
	}
}

// 绑定谷歌验证器，绑定后不可修改
func (m *Model) BindGoogleAuth(userID int, googleAuthCode string) error {
	user, err := m.GetUserByID(userID)
	if err != nil {
		return err
	}
	err := m.VerifyGoogleAuthCode(user.Username, googleAuthCode)
	if err != nil {
		return err
	} else {
		return m.DB.Model(&User{}).Where("id=?", userID).UpdateColumn("is_bind_google_auth", true).Error
	}
}

// 重置密码
func (m *Model) ResetPassword(userID int, googleAuthCode, newPassword string) error {
	err := m.UserValidator.ValidPassword(newPassword)
	if err != nil {
		return err
	}
	user, err := m.GetUserByID(userID)
	if err != nil {
		return err
	}
	err = m.VerifyGoogleAuthCode(user.Username, googleAuthCode)
	if err != nil {
		return err
	}
	password, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		return err
	}
	return m.DB.Where("id=?", userID).Updates(&User{Password: string(password)}).Error
}

// 获得推荐人(自带缓存)
func (m *Model) GetRecommender(userID int) (recommenderID int, err error) {
	recommenderID, err = m.RecommendersCache.Get(userID)
	if err != nil {
		return 0, err
	}
	if recommenderID < 0 {
		user := &User{}
		err = m.DB.Where("id=?", userID).Select("recommender").First(user).Error
		if err != nil {
			return 0, err
		} else {
			err = m.RecommendersCache.Set(userID, user.Recommender)
			if err != nil {
				return 0, err
			} else {
				return user.Recommender, nil
			}
		}
	} else {
		return recommenderID, nil
	}
}

// 设置用户头像
func (m *Model) SetUserAvatar(userID int, avatar string) error {
	return m.DB.Where("id=?", userID).Updates(&User{Avatar: avatar}).Error
}
