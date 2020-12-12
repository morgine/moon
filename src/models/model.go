package models

import (
	"github.com/morgine/moon/pkg/google_authenticator"
	"gorm.io/gorm"
)

type Model struct {
	db  *gorm.DB
	gac *google_authenticator.Client
}
