package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

var (
	JWTIssuer          string
	JWTAudience        string
	JWTSecret          string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	JWTCookieDomain    string
	JWTCookiePath      string
	JWTCookieName      string
)

func LoadAuthConfig() {
	JWTIssuer = os.Getenv("JWT_ISSUER")
	JWTAudience = os.Getenv("JWT_AUDIENCE")
	JWTSecret = os.Getenv("JWT_SECRET")

	AccessTokenExpiryInt, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRY"))
	if err != nil {
		log.Fatalf("Error parsing ACCESS_TOKEN_EXPIRY: %v", err)
	}
	AccessTokenExpiry = time.Duration(AccessTokenExpiryInt) * time.Second

	RefreshTokenExpiryInt, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRY"))
	if err != nil {
		log.Fatalf("Error parsing REFRESH_TOKEN_EXPIRY: %v", err)
	}
	RefreshTokenExpiry = time.Duration(RefreshTokenExpiryInt) * time.Second

	JWTCookieDomain = os.Getenv("JWT_COOKIE_DOMAIN")
	JWTCookiePath = os.Getenv("JWT_COOKIE_PATH")
	JWTCookieName = os.Getenv("JWT_COOKIE_NAME")
}
