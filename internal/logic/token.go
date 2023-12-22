package logic

import (
	"context"
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt"
	"gitlab.com/pjrpc/pjrpc/v2"
	"go.uber.org/zap"

	"github.com/tranHieuDev23/cato/internal/configs"
	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc"
	"github.com/tranHieuDev23/cato/internal/utils"
)

type Token interface {
	GetToken(ctx context.Context, accountID uint64) (string, error)
	GetAccountIDAndExpireTime(ctx context.Context, token string) (uint64, time.Time, error)
	GetAccount(ctx context.Context, token string) (*db.Account, error)
}

type token struct {
	accountDataAccessor db.AccountDataAccessor
	expiresIn           time.Duration
	privateKey          *rsa.PrivateKey
	publicKey           *rsa.PublicKey
	tokenConfig         configs.Token
	logger              *zap.Logger
}

func NewToken(
	accountDataAccessor db.AccountDataAccessor,
	tokenConfig configs.Token,
	logger *zap.Logger,
) (Token, error) {
	expiresIn, err := tokenConfig.GetExpiresInDuration()
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to parse expires_in")
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(tokenConfig.PrivateKey))
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to parse rsa private key")
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(tokenConfig.PublicKey))
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to parse rsa public key")
		return nil, err
	}

	return &token{
		accountDataAccessor: accountDataAccessor,
		expiresIn:           expiresIn,
		privateKey:          privateKey,
		publicKey:           publicKey,
		tokenConfig:         tokenConfig,
		logger:              logger,
	}, nil
}

func (t token) GetAccountIDAndExpireTime(ctx context.Context, tokenString string) (uint64, time.Time, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)

	parsedToken, err := jwt.Parse(tokenString, func(parsedToken *jwt.Token) (interface{}, error) {
		if _, ok := parsedToken.Method.(*jwt.SigningMethodRSA); !ok {
			logger.With(zap.Any("alg", parsedToken.Header["alg"])).Error("unexpected signing method")
			return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeUnauthenticated))
		}

		return t.publicKey, nil
	})
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to parse token")
		return 0, time.Time{}, err
	}

	if !parsedToken.Valid {
		logger.Error("invalid token")
		return 0, time.Time{}, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeUnauthenticated))
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, time.Time{}, pjrpc.JRPCErrInternalError()
	}

	accountID, ok := claims["sub"].(float64)
	if !ok {
		return 0, time.Time{}, pjrpc.JRPCErrInternalError()
	}

	expireTimeUnix, ok := claims["exp"].(float64)
	if !ok {
		return 0, time.Time{}, pjrpc.JRPCErrInternalError()
	}

	return uint64(accountID), time.Unix(int64(expireTimeUnix), 0), nil
}

func (t token) GetToken(ctx context.Context, accountID uint64) (string, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"sub": accountID,
		"exp": time.Now().Add(t.expiresIn).Unix(),
	})

	tokenString, err := token.SignedString(t.privateKey)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to sign token")
		return "", pjrpc.JRPCErrInternalError()
	}

	return tokenString, nil
}

func (t token) GetAccount(ctx context.Context, token string) (*db.Account, error) {
	accountID, _, err := t.GetAccountIDAndExpireTime(ctx, token)
	if err != nil {
		return nil, err
	}

	account, err := t.accountDataAccessor.GetAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeUnauthenticated))
	}

	return account, nil
}
