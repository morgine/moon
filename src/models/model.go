package models

import (
	"github.com/morgine/moon/pkg/google_authenticator"
	"gorm.io/gorm"
)

type Model struct {
	db  *gorm.DB
	gac *google_authenticator.Client
	//users *usersCache
}

//
//type usersCache struct {
//	client cache.Client
//}
//
//func (uc *usersCache) GetUserByID(id int) (*User, error) {
//	data, err := uc.client.Get(strconv.Itoa(id))
//	if err != nil {
//		return nil, err
//	}
//	if len(data) == 0 {
//		return nil, nil
//	}
//	user := &User{}
//	err = json.Unmarshal(data, user)
//	if err != nil {
//		return nil, err
//	} else {
//		return user, nil
//	}
//}
//
//func (uc *usersCache) GetUserByUsername(username string) (*User, error) {
//
//}
