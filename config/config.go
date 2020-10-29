package config

import "time"

var (
	BaseDirectory            string
	FrontendURL              string
	EmailServicePassword     string
	JWTAccessTokenSecretKey  string
	JWTRefreshTokenSecretKey string
	JWTAccessTokenLifetime   time.Duration
	JWTRefreshTokenLifetime  time.Duration
)
