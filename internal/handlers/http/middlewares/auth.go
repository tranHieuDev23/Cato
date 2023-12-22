package middlewares

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/tranHieuDev23/cato/internal/configs"
	"github.com/tranHieuDev23/cato/internal/logic"
)

const (
	AuthorizationCookie = "CatoAuth"
)

type HTTPAuth func(http.Handler) http.Handler

func getAuthorizationBearerToken(request *http.Request) string {
	authorizationCookie, err := request.Cookie(AuthorizationCookie)
	if err != nil {
		return ""
	}

	return authorizationCookie.Value
}

func setAuthorizationBearerToken(w http.ResponseWriter, token string, expireTime time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     AuthorizationCookie,
		Value:    token,
		HttpOnly: true,
		Expires:  expireTime,
		SameSite: http.SameSiteStrictMode,
	})
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
			accountID, expireTime, err := tokenLogic.GetAccountIDAndExpireTime(ctx, postRequestToken)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(make([]byte, 0))
				return
			}

			if time.Now().Add(regenerateTokenBeforeExpiryDuration).After(expireTime) {
				newToken, newExpireTime, err := tokenLogic.GetToken(ctx, accountID)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(make([]byte, 0))
					return
				}

				postRequestToken = newToken
				expireTime = newExpireTime
			}

			setAuthorizationBearerToken(inMemoryWriter, postRequestToken, expireTime)
			inMemoryWriter.Apply(w)
		})
	}, nil
}
