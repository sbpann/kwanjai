package config

import (
	"context"
	"time"
)

var (
	BaseDirectory            string
	FrontendURL              string
	EmailServicePassword     string
	EmailVerficationLifetime time.Duration
	JWTAccessTokenSecretKey  string
	JWTRefreshTokenSecretKey string
	JWTAccessTokenLifetime   time.Duration
	JWTRefreshTokenLifetime  time.Duration
	Context                  context.Context
)
