package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/tranHieuDev23/cato/internal/configs"
	"github.com/tranHieuDev23/cato/internal/logic"
)

type HTTPAuth func(http.Handler) http.Handler

func getAuthorizationBearerToken(request *http.Request) string {
	authorizationHeader := request.Header.Get("Authorization")
	authorizationHeaderParts := strings.Split(authorizationHeader, "Bearer ")
	if len(authorizationHeaderParts) != 2 {
		return ""
	}

	return authorizationHeaderParts[1]
}

func setAuthorizationBearerToken(w http.ResponseWriter, token string) {
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
}

func NewHTTPAuth(
	tokenLogic logic.Token,
	tokenConfig configs.Token,
	logger *zap.Logger,
) (HTTPAuth, error) {
	regenerateTokenBeforeExpiryDuration, err := tokenConfig.GetRegenerateTokenBeforeExpiryDuration()
	if err != nil {
		logger.
			With(zap.String("regenerate_token_before_expiry", tokenConfig.RegenerateTokenBeforeExpiry)).
			With(zap.Error(err)).
			Error("failed to parse regenerate_token_before_expiry")
		return nil, err
	}

	return func(baseHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			inMemoryWriter := newInMemoryResponseWriter()
			baseHandler.ServeHTTP(inMemoryWriter, r)

			postRequestToken := getAuthorizationBearerToken(r)
			if postRequestToken == "" {
				inMemoryWriter.Apply(w)
				return
			}

			ctx := r.Context()
			accountID, expiredAt, err := tokenLogic.GetAccountIDAndExpireTime(ctx, postRequestToken)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(make([]byte, 0))
				return
			}

			if time.Now().Add(regenerateTokenBeforeExpiryDuration).After(expiredAt) {
				newToken, err := tokenLogic.GetToken(ctx, accountID)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(make([]byte, 0))
					return
				}

				postRequestToken = newToken
			}

			setAuthorizationBearerToken(inMemoryWriter, postRequestToken)
			inMemoryWriter.Apply(w)
		})
	}, nil
}
