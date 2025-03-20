package models

import (
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type OAuthCredential struct {
	gorm.Model

	Token      oauth2.Token `json:"token" gorm:"embedded;embeddedPrefix:token_;not null;"`
	ProxyOwner string       `json:"proxy_owner" gorm:"not null;"`
}
