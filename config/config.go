package config

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	BaseDirectory                string
	Port                         string
	FirebaseProjectID            string
	FrontendURL                  string
	BackendURL                   string
	EmailServicePassword         string
	EmailVerficationLifetime     time.Duration
	JWTAccessTokenSecretKey      string
	JWTRefreshTokenSecretKey     string
	JWTAccessTokenLifetime       time.Duration
	JWTRefreshTokenLifetime      time.Duration
	Context                      context.Context
	DefaultAuthenticationBackend gin.HandlerFunc
	OmisePublicKey               string
	OmiseSecretKey               string
)
