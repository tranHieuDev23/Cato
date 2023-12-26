package cato

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"

	"gitlab.com/pjrpc/pjrpc/v2/client"
	"go.uber.org/zap"

	"github.com/tranHieuDev23/cato/internal/handlers/http/middlewares"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc/rpcclient"
	"github.com/tranHieuDev23/cato/internal/utils"
)

func getHostEndpoint(address string) string {
	return fmt.Sprintf("%s/api", address)
}

func InitializeBaseClient(appArguments utils.Arguments) (rpcclient.APIClient, error) {
	hostEndpoint := getHostEndpoint(appArguments.HostAddress)
	pjrpcClient, err := client.New(hostEndpoint)
	if err != nil {
		return nil, err
	}

	return rpcclient.NewAPIClient(pjrpcClient), nil
}

type HTTPClientWithAuthToken client.HTTPDoer

type httpClientWithAuthToken struct {
	appArguments  utils.Arguments
	logger        *zap.Logger
	authCookie    *atomic.Pointer[string]
	baseAPIClient rpcclient.APIClient
	httpClient    *http.Client
}

func NewHTTPClientWithAuthToken(
	appArguments utils.Arguments,
	logger *zap.Logger,
) (HTTPClientWithAuthToken, error) {
	baseClient, err := InitializeBaseClient(appArguments)
	if err != nil {
		return nil, err
	}

	return &httpClientWithAuthToken{
		appArguments:  appArguments,
		logger:        logger,
		authCookie:    new(atomic.Pointer[string]),
		baseAPIClient: baseClient,
		httpClient:    http.DefaultClient,
	}, nil
}

func (h httpClientWithAuthToken) getAuthToken(ctx context.Context) (string, error) {
	response, err := h.baseAPIClient.CreateSession(ctx, &rpc.CreateSessionRequest{
		AccountName: h.appArguments.WorkerAccountName,
		Password:    h.appArguments.WorkerAccountPassword,
	})
	if err != nil {
		return "", err
	}

	return response.Token, nil
}

func (h httpClientWithAuthToken) Do(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	token := h.authCookie.Load()
	if token == nil {
		newToken, err := h.getAuthToken(ctx)
		if err != nil {
			return nil, err
		}

		token = &newToken
		h.authCookie.Store(token)
	}

	req.AddCookie(&http.Cookie{Name: middlewares.AuthorizationCookie, Value: *token})
	return h.httpClient.Do(req)
}

func InitializeAuthenticatedClient(
	appArguments utils.Arguments,
	httpClientWithAuthCookie HTTPClientWithAuthToken,
) rpcclient.APIClient {
	return rpcclient.NewAPIClient(&client.Client{
		URL:        getHostEndpoint(appArguments.HostAddress),
		HTTPClient: httpClientWithAuthCookie,
	})
}
