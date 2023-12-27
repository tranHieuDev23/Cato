package logic

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"time"

	"github.com/golang-jwt/jwt"
	"gitlab.com/pjrpc/pjrpc/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/tranHieuDev23/cato/internal/configs"
	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc"
	"github.com/tranHieuDev23/cato/internal/utils"
)

const (
	rs512KeyPairBitCount = 2048
)

type Token interface {
	GetToken(ctx context.Context, accountID uint64) (string, time.Time, error)
	GetAccountIDAndExpireTime(ctx context.Context, token string) (uint64, time.Time, error)
	GetAccount(ctx context.Context, token string) (*db.Account, error)
	WithDB(db *gorm.DB) Token
}

func generateRSAKeyPair(bits int) (*rsa.PrivateKey, error) {
	privateKeyPair, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}

	return privateKeyPair, nil
}

type token struct {
	accountDataAccessor        db.AccountDataAccessor
	tokenPublicKeyDataAccessor db.TokenPublicKeyDataAccessor
	expiresIn                  time.Duration
	privateKey                 *rsa.PrivateKey
	tokenPublicKeyID           uint64
	tokenConfig                configs.Token
	logger                     *zap.Logger
}

func pemEncodePublicKey(pubKey *rsa.PublicKey) ([]byte, error) {
	pubBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return nil, err
	}
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}
	return pem.EncodeToMemory(block), nil
}

func NewToken(
	accountDataAccessor db.AccountDataAccessor,
	tokenPublicKeyDataAccessor db.TokenPublicKeyDataAccessor,
	tokenConfig configs.Token,
	logger *zap.Logger,
) (Token, error) {
	expiresIn, err := tokenConfig.GetExpiresInDuration()
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to parse expires_in")
		return nil, err
	}

	rsaKeyPair, err := generateRSAKeyPair(rs512KeyPairBitCount)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to generate rsa key pair")
		return nil, err
	}

	publicKeyBytes, err := pemEncodePublicKey(&rsaKeyPair.PublicKey)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to encode public key in pem format")
		return nil, err
	}

	tokenPublicKey := &db.TokenPublicKey{PublicKey: publicKeyBytes}
	err = tokenPublicKeyDataAccessor.CreatePublicKey(context.Background(), tokenPublicKey)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create public key entry in database")
		return nil, err
	}

	return &token{
		accountDataAccessor:        accountDataAccessor,
		tokenPublicKeyDataAccessor: tokenPublicKeyDataAccessor,
		expiresIn:                  expiresIn,
		privateKey:                 rsaKeyPair,
		tokenPublicKeyID:           uint64(tokenPublicKey.ID),
		tokenConfig:                tokenConfig,
		logger:                     logger,
	}, nil
}

func (t token) GetAccountIDAndExpireTime(ctx context.Context, tokenString string) (uint64, time.Time, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)

	parsedToken, err := jwt.Parse(tokenString, func(parsedToken *jwt.Token) (interface{}, error) {
		if _, ok := parsedToken.Method.(*jwt.SigningMethodRSA); !ok {
			logger.Error("unexpected signing method")
			return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeUnauthenticated))
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			logger.Error("cannot get token's claims")
			return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeUnauthenticated))
		}

		tokenPublicKeyID, ok := claims["kid"].(float64)
		if !ok {
			logger.Error("cannot get token's kid claim")
			return nil, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeUnauthenticated))
		}

		tokenPublicKey, err := t.tokenPublicKeyDataAccessor.GetPublicKey(ctx, uint64(tokenPublicKeyID))
		if err != nil {
			logger.Error("cannot get token's public key from database")
			return nil, pjrpc.JRPCErrInternalError()
		}

		return jwt.ParseRSAPublicKeyFromPEM(tokenPublicKey.PublicKey)
	})
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to parse token")
		return 0, time.Time{}, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeUnauthenticated))
	}

	if !parsedToken.Valid {
		logger.Error("invalid token")
		return 0, time.Time{}, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeUnauthenticated))
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		logger.Error("cannot get token's claims")
		return 0, time.Time{}, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeUnauthenticated))
	}

	accountID, ok := claims["sub"].(float64)
	if !ok {
		logger.Error("cannot get token's sub claim")
		return 0, time.Time{}, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeUnauthenticated))
	}

	expireTimeUnix, ok := claims["exp"].(float64)
	if !ok {
		logger.Error("cannot get token's exp claim")
		return 0, time.Time{}, pjrpc.JRPCErrServerError(int(rpc.ErrorCodeUnauthenticated))
	}

	return uint64(accountID), time.Unix(int64(expireTimeUnix), 0), nil
}

func (t token) GetToken(ctx context.Context, accountID uint64) (string, time.Time, error) {
	logger := utils.LoggerWithContext(ctx, t.logger)

	expireTime := time.Now().Add(t.expiresIn)
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"sub": accountID,
		"exp": expireTime.Unix(),
		"kid": t.tokenPublicKeyID,
	})

	tokenString, err := token.SignedString(t.privateKey)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to sign token")
		return "", time.Time{}, pjrpc.JRPCErrInternalError()
	}

	return tokenString, expireTime, nil
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

func (t token) WithDB(db *gorm.DB) Token {
	return &token{
		accountDataAccessor:        t.accountDataAccessor.WithDB(db),
		tokenPublicKeyDataAccessor: t.tokenPublicKeyDataAccessor.WithDB(db),
		expiresIn:                  t.expiresIn,
		privateKey:                 t.privateKey,
		tokenPublicKeyID:           t.tokenPublicKeyID,
		tokenConfig:                t.tokenConfig,
		logger:                     t.logger,
	}
}
