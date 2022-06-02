package nuntio_authorize

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/nuntiodev/cloud-proto/go_cloud"
	"google.golang.org/grpc"
	"sync"
)

type Authorize interface {
	GetAccessToken(ctx context.Context) (string, error)
}

type defaultAuthorize struct {
	apiKeyJwt   *jwt.Token
	accessToken *jwt.Token
	accessApi   go_cloud.CloudServiceClient
	sync.Mutex
}

type NoAuthorization struct{}

func (a *NoAuthorization) GetAccessToken(ctx context.Context) (string, error) {
	return "", nil
}

func New(ctx context.Context, apiUrl string, apiKey string, authorize Authorize, dialOptions grpc.DialOption) (Authorize, error) {
	if authorize != nil {
		return authorize, nil
	}
	// setup grpc connection to access service
	accessClientConn, err := grpc.Dial(apiUrl, dialOptions)
	if err != nil {
		return nil, err
	}
	accessApi := go_cloud.NewCloudServiceClient(accessClientConn)
	// setup access auth token
	accessResp, err := accessApi.GenerateAccessToken(ctx, &go_cloud.CloudRequest{
		PrivateKey: apiKey,
	})
	if err != nil {
		return nil, err
	}
	if accessResp == nil || accessResp.AccessToken == "" {
		return nil, errors.New("could not initialize access token")
	}
	// create jwt
	apiKeyJwt, _ := jwt.Parse(apiKey, nil)
	if apiKeyJwt == nil || apiKeyJwt.Claims == nil {
		return nil, errors.New("invalid secret key")
	}
	if err := apiKeyJwt.Claims.Valid(); err != nil {
		return nil, err
	}
	accessToken, _ := jwt.Parse(accessResp.AccessToken, nil)
	if err := accessToken.Claims.Valid(); err != nil {
		return nil, err
	}
	if accessToken == nil || accessToken.Claims == nil {
		return nil, errors.New("invalid auth token")
	}
	return &defaultAuthorize{
		accessApi:   accessApi,
		apiKeyJwt:   apiKeyJwt,
		accessToken: accessToken,
	}, nil
}

func (sa *defaultAuthorize) GetAccessToken(ctx context.Context) (string, error) {
	sa.Lock()
	defer sa.Unlock()
	if err := sa.accessToken.Claims.Valid(); err != nil {
		accessResp, err := sa.accessApi.GenerateAccessToken(ctx, &go_cloud.CloudRequest{
			PrivateKey: sa.apiKeyJwt.Raw,
		})
		if err != nil {
			return "", err
		}
		authTokenJwt, _ := jwt.Parse(accessResp.AccessToken, nil)
		if authTokenJwt == nil || authTokenJwt.Claims == nil {
			return "", errors.New("invalid token returned")
		}
		if err := authTokenJwt.Claims.Valid(); err != nil {
			return "", err
		}
		sa.accessToken = authTokenJwt
	}
	return sa.accessToken.Raw, nil
}
