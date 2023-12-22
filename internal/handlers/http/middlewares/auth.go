package middlewares

import (
	"context"
	"encoding/json"
	"time"

	"gitlab.com/pjrpc/pjrpc/v2"
	"go.uber.org/zap"

	"github.com/tranHieuDev23/cato/internal/configs"
	"github.com/tranHieuDev23/cato/internal/handlers/http/utils"
	"github.com/tranHieuDev23/cato/internal/logic"
)

const (
	AuthContextFieldToken = "Token"
)

type Auth pjrpc.Middleware

func NewAuth(
	tokenLogic logic.Token,
	tokenConfig configs.Token,
	logger *zap.Logger,
) (Auth, error) {
	regenerateTokenBeforeExpiryDuration, err := tokenConfig.GetRegenerateTokenBeforeExpiryDuration()
	if err != nil {
		logger.
			With(zap.String("regenerate_token_before_expiry", tokenConfig.RegenerateTokenBeforeExpiry)).
			With(zap.Error(err)).
			Error("failed to parse regenerate_token_before_expiry")
		return nil, err
	}

	return func(next pjrpc.Handler) pjrpc.Handler {
		return func(ctx context.Context, params json.RawMessage) (any, error) {
			currentTime := time.Now()

			token := utils.GetAuthorizationBearerToken(ctx)
			ctx = context.WithValue(ctx, AuthContextFieldToken, token)
			response, err := next(ctx, params)

			if token != "" {
				accountID, expiresAt, tokenErr := tokenLogic.GetAccountIDAndExpireTime(ctx, token)
				if tokenErr != nil {
					return nil, tokenErr
				}

				if currentTime.Add(regenerateTokenBeforeExpiryDuration).After(expiresAt) {
					newToken, newTokenErr := tokenLogic.GetToken(ctx, accountID)
					if err != nil {
						return nil, newTokenErr
					}

					utils.SetAuthorizationBearerToken(ctx, newToken)
				}
			}

			return response, err
		}
	}, nil
}
