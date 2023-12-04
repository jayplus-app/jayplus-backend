package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	Issuer             string        `json:"issuer"`
	Audience           string        `json:"audience"`
	Secret             string        `json:"secret"`
	AccessTokenExpiry  time.Duration `json:"access_token_expiry"`
	RefreshTokenExpiry time.Duration `json:"refresh_token_expiry"`
	CookieDomain       string        `json:"cookie_domain"`
	CookiePath         string        `json:"cookie_path"`
	CookieName         string        `json:"cookie_name"`
}

type AuthUser struct {
	ID         int    `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	BusinessID int    `json:"businessID"`
}

type JWTTokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type JWTClaims struct {
	jwt.RegisteredClaims
	BusinessID int `json:"businessID"`
}
