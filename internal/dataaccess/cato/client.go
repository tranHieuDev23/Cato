package cato

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"gitlab.com/pjrpc/pjrpc/v2/client"
	"go.uber.org/zap"

	"github.com/tranHieuDev23/cato/internal/handlers/http/middlewares"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc/rpcclient"
	"github.com/tranHieuDev23/cato/internal/utils"
)

type HTTPClientWithAuthCookie client.HTTPDoer

type httpClientWithAuthCookie struct {
	authCookie *atomic.Pointer[string]
	httpClient *http.Client
	logger     *zap.Logger
}

func NewHTTPClientWithAuthCookie(
	logger *zap.Logger,
) HTTPClientWithAuthCookie {
	return &httpClientWithAuthCookie{
		authCookie: new(atomic.Pointer[string]),
		httpClient: http.DefaultClient,
		logger:     logger,
	}
}

func (h httpClientWithAuthCookie) Do(req *http.Request) (*http.Response, error) {
	logger := utils.LoggerWithContext(req.Context(), h.logger)

	response, err := h.httpClient.Do(req)
	if err != nil {
		logger.With(zap.Error(err)).Error("error occurred when making http request")
		return nil, err
	}

	for _, cookie := range response.Cookies() {
		if cookie.Name != middlewares.AuthorizationCookie {
			continue
		}

		if cookie.Expires.Before(time.Now()) {
			h.authCookie.Swap(nil)
		} else {
			h.authCookie.Swap(&cookie.Value)
		}
	}

	return response, nil
}

func InitializeBaseClient(
	args utils.Arguments,
	httpClientWithAuthCookie HTTPClientWithAuthCookie,
) rpcclient.APIClient {
	hostEndpoint := fmt.Sprintf("http://%s/api", args.HostAddress)
	return rpcclient.NewAPIClient(&client.Client{
		URL:        hostEndpoint,
		HTTPClient: httpClientWithAuthCookie,
	})
}

func InitializeAuthenticatedClient(
	args utils.Arguments,
	httpClientWithAuthCookie HTTPClientWithAuthCookie,
	logger *zap.Logger,
) (rpcclient.APIClient, error) {
	baseClient := InitializeBaseClient(args, httpClientWithAuthCookie)
	_, err := baseClient.CreateSession(context.Background(), &rpc.CreateSessionRequest{
		AccountName: args.WorkerAccountName,
		Password:    args.WorkerAccountPassword,
	})
	if err != nil {
		logger.With(zap.Any("args", args)).With(zap.Error(err)).Error("failed to call create_session")
		return nil, err
	}

	return baseClient, nil
}
