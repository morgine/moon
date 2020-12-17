package models

import (
	"github.com/morgine/moon/pkg/cache"
	"github.com/morgine/moon/pkg/google_authenticator"
	"github.com/morgine/moon/src/validators"
	"gorm.io/gorm"
)

type Model struct {
	DB                *gorm.DB
	GAC               *google_authenticator.Client
	UserValidator     validators.User
	RecommendersCache *cache.Recommenders
}
