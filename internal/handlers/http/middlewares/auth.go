package middlewares

import (
	"context"
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

func unsetAuthorizationBearerToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     AuthorizationCookie,
		Value:    "",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteStrictMode,
	})
}

func regenerateTokenIfDurationToExpireIsWithinThreshold(
	ctx context.Context,
	token string,
	tokenLogic logic.Token,
	regenerateTokenBeforeExpiryDuration time.Duration,
) (string, time.Time, error) {
	accountID, expireTime, err := tokenLogic.GetAccountIDAndExpireTime(ctx, token)
	if err != nil {
		return "", time.Time{}, err
	}

	if !time.Now().Add(regenerateTokenBeforeExpiryDuration).After(expireTime) {
		return token, expireTime, nil
	}

	newToken, newExpireTime, newTokenErr := tokenLogic.GetToken(ctx, accountID)
	if newTokenErr != nil {
		return "", time.Time{}, newTokenErr
	}

	return newToken, newExpireTime, nil
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
			ctx := r.Context()

			inMemoryWriter := newInMemoryResponseWriter()
			baseHandler.ServeHTTP(inMemoryWriter, r)

			postRequestToken := getAuthorizationBearerToken(r)
			if postRequestToken == "" {
				unsetAuthorizationBearerToken(inMemoryWriter)
				applyErr := inMemoryWriter.Apply(w)
				if applyErr != nil {
					logger.With(zap.Error(applyErr)).Error("failed to apply in-memory writer to response writer")
				}

				return
			}

			regenerateToken, regenerateExpireTime, regenerateErr := regenerateTokenIfDurationToExpireIsWithinThreshold(
				ctx, postRequestToken, tokenLogic, regenerateTokenBeforeExpiryDuration)
			if regenerateErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				if _, writeErr := w.Write(make([]byte, 0)); writeErr != nil {
					logger.With(zap.Error(writeErr)).Error("failed to write to response writer")
				}

				return
			}

			setAuthorizationBearerToken(inMemoryWriter, regenerateToken, regenerateExpireTime)
			err = inMemoryWriter.Apply(w)
			if err != nil {
				logger.With(zap.Error(err)).Error("failed to apply in-memory writer to response writer")
			}
		})
	}, nil
}
