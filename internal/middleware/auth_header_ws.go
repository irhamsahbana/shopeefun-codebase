package middleware

import (
	"codebase-app/pkg/jwthandler"
	"context"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
)

type ClaimsKey string

const KeyClaims ClaimsKey = "claims"

func AuthWs(next http.Handler) http.Handler {
	unauthorizedResponse := `
	{
		"message": "Unauthorized",
		"success": false
	}`

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// get the token from the query string
		token := r.URL.Query().Get("token")
		// if the token is empty, return an unauthorized status
		if token == "" {
			log.Error().Msg("middleware::AuthWs - Unauthorized [Token not set]")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(unauthorizedResponse))
			return
		}

		claims, err := jwthandler.ParseEphemeralToken(token)
		if err != nil {
			log.Error().Err(err).Any("payload", token).Msg("middleware::AuthWs - Error while parsing token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(unauthorizedResponse))
			return
		}

		ctx := r.Context()
		claimsMap := map[string]any{
			"user_id": claims.UserId,
			"role":    claims.Role,
		}

		ctx = context.WithValue(ctx, KeyClaims, claimsMap)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetClaims(ctx context.Context) (claims map[string]any, err error) {
	claims, ok := ctx.Value(KeyClaims).(map[string]any)
	if !ok {
		return nil, errors.New("claims not found in context")
	}

	return claims, nil
}
