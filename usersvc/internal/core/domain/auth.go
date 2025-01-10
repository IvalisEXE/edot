package domain

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

type UserTokenData struct {
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
	SessionID string    `json:"session_id"`
}

type UserCustomRegistered struct {
	ID int `json:"id"`
	jwt.RegisteredClaims
}

type UserCustomClaims struct {
	ID        int `json:"id"`
	ExpiredAt int `json:"exp"`
}

type UserSession struct {
	ID       int     `json:"id"`
	ShopID   int     `json:"shop_id"`
	Name     string  `json:"name"`
	Phone    *string `json:"phone"`
	Platform string  `json:"platform"`
}

type UserSessionInfo struct {
	IpAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	ExpiredAt string `json:"expired_at"`
}

type User struct {
	ID        int        `json:"id"`
	ShopID    int        `json:"shop_id"`
	Name      string     `json:"name,omitempty"`
	Phone     string     `json:"phone"`
	Password  *string    `json:"-"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	Platform  string     `json:"platform"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

func (u *User) ToUserSession() UserSession {
	return UserSession{
		ID: u.ID,

		Name:     u.Name,
		Phone:    &u.Phone,
		Platform: u.Platform,
	}
}

func (u *User) GenerateToken() (*UserTokenData, error) {
	tokenExpiration := time.Now().Add(time.Second * time.Duration(viper.GetInt("REDIS_SESSION_TTL")))
	tokenID := base64.StdEncoding.EncodeToString([]byte(time.Now().Format(time.RFC3339)))
	claims := &UserCustomRegistered{
		u.ID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(tokenExpiration),
			ID:        tokenID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	var signed string
	if u.Platform == "customer" {
		signed = viper.GetString("JWT_SECRET_CUSTOMER")
	} else {
		signed = viper.GetString("JWT_SECRET_DASHBOARD")
	}

	tokenString, err := token.SignedString([]byte(signed))
	if err != nil {
		return nil, err
	}

	return &UserTokenData{
		Token:     tokenString,
		ExpiredAt: tokenExpiration,
		SessionID: u.GenerateSessionID("auth", u.ID, tokenID),
	}, nil
}

func (u *User) GenerateSessionID(prefix string, userID int, token string) string {
	return fmt.Sprintf(prefix+":%s&%s", userID, token)
}
