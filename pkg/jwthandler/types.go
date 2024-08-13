package jwthandler

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserId string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type CostumClaimsWs struct {
	UserId string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type CostumClaimsPayload struct {
	UserId          string    `json:"user_id"`
	Role            string    `json:"role"`
	TokenExpiration time.Time `json:"token_expiration"`
}

type CostumClaimsPayloadWs struct {
	UserId          string    `json:"user_id"`
	Role            string    `json:"role"`
	TokenExpiration time.Time `json:"token_expiration"`
}
