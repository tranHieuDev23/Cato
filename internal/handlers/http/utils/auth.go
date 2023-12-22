package utils

import (
	"context"
	"fmt"
	"strings"

	"gitlab.com/pjrpc/pjrpc/v2"
)

func GetAuthorizationBearerToken(ctx context.Context) string {
	contextData, ok := pjrpc.ContextGetData(ctx)
	if !ok {
		return ""
	}

	authorizationHeader := contextData.HTTPRequest.Header.Get("Authorization")
	authorizationHeaderParts := strings.Split(authorizationHeader, "Bearer ")
	if len(authorizationHeaderParts) != 2 {
		return ""
	}

	return authorizationHeaderParts[1]
}

func SetAuthorizationBearerToken(ctx context.Context, token string) {
	contextData, ok := pjrpc.ContextGetData(ctx)
	if !ok {
		return
	}

	contextData.HTTPRequest.Response.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
}
