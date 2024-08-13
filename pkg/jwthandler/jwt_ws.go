package jwthandler

import (
	"codebase-app/internal/infrastructure/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

func GenerateEphemeralToken(p CostumClaimsPayloadWs) (string, error) {
	now := time.Now().UTC()
	privateKey := []byte(config.Envs.Guard.JwtPrivateKeyWs)
	exp := time.Now().Add(time.Second * time.Duration(config.Envs.Guard.JwtWsExp))

	claims := CostumClaimsWs{
		UserId: p.UserId,
		Role:   p.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user",
			Issuer:    "codebase-app",
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		log.Warn().Err(err).Msg("jwthandler::GenerateEphemeralToken - Error while signing token")
		return "", err
	}

	return tokenString, nil
}

func ParseEphemeralToken(token string) (*CostumClaimsWs, error) {
	claims := &CostumClaimsWs{}
	privateKey := []byte(config.Envs.Guard.JwtPrivateKeyWs)

	jwtToken, err := jwt.ParseWithClaims(
		token,
		claims,
		func(token *jwt.Token) (any, error) {
			return privateKey, nil
		},
	)
	if err != nil {
		log.Warn().Err(err).Msg("jwthandler::ParseEphemeralToken - Error while parsing token")
		return nil, err
	}

	if !jwtToken.Valid {
		log.Warn().Msg("jwthandler::ParseAphemeralToken - Invalid token")
		return nil, err
	}

	return claims, nil
}
