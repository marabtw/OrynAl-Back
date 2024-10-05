package model

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Register struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type JwtTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserToken struct {
	ID           uint      `gorm:"primary_key;auto_increment" json:"id"`
	UserID       uint      `gorm:"unique;not null"            json:"full_name"`
	Role         string    `gorm:"not null"                   json:"role"`
	Email        string    `gorm:"unique;not null"                  json:"email"`
	AccessToken  string    `gorm:"unique;not null"            json:"access_token"`
	RefreshToken string    `gorm:"unique;not null"            json:"refresh_token"`
	CreatedAt    time.Time `gorm:"default:current_timestamp"  json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime:milli"       json:"updated_at"`
}

type UserClaim struct {
	Email  string `json:"email"`
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}

type JWTClaim struct {
	Email          string `json:"email"`
	UserID         uint   `json:"user_id"`
	Role           string `json:"role"`
	StandardClaims jwt.StandardClaims
}

func (J *JWTClaim) Valid() error {
	return nil
}

type RefreshJWTClaim struct {
	Email          string `json:"email"`
	UserID         uint   `json:"user_id"`
	StandardClaims jwt.StandardClaims
}

func (r *RefreshJWTClaim) Valid() error {
	return nil
}

type ContextUserID struct {
	ID uint `json:"user_id"`
}
type ContextUserRole struct {
	Role string `json:"role"`
}

type contextKey string

var (
	ContextUserIDKey   = contextKey("id")
	ContextUserRoleKey = contextKey("role")
)
