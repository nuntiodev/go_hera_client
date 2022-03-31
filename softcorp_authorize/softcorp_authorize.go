package softcorp_authorize

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/softcorp-io/cloud-proto/go_cloud"
	"google.golang.org/grpc"
	"sync"
)

var (
	// ACCESS_API_URL  is the URL of the API access service backend.
	ACCESS_API_URL = "api.softcorp.io:443"
	// AUTHORIZE is used to override the default softcorp_authorize interface which is used to validate tokens
	// if you don't want any authorization, set it to softcorp_authorize.AUTHORIZE = softcorp_authorize.NoAuthorization.
	AUTHORIZE Authorize
)

type Authorize interface {
	GetAccessToken(ctx context.Context) (string, error)
}

type defaultSoftcorpAuthorize struct {
	secretKeyJwt *jwt.Token
	authTokenJwt *jwt.Token
	accessApi    go_cloud.ProjectServiceClient
	sync.Mutex
}

type NoAuthorization struct{}

func (a *NoAuthorization) GetAccessToken(ctx context.Context) (string, error) {
	return "", nil
}

func New(ctx context.Context, apiKey string, dialOptions grpc.DialOption) (Authorize, error) {
	if AUTHORIZE != nil {
		return AUTHORIZE, nil
	}
	// setup grpc connection to access service
	accessClientConn, err := grpc.Dial(ACCESS_API_URL, dialOptions)
	if err != nil {
		return nil, err
	}
	accessApi := go_cloud.NewProjectServiceClient(accessClientConn)
	// setup access auth token
	accessResp, err := accessApi.GenerateAccessToken(ctx, &go_cloud.ProjectRequest{
		PrivateKey: apiKey,
	})
	if err != nil {
		return nil, err
	}
	if accessResp == nil || accessResp.AccessToken == "" {
		return nil, errors.New("could not initialize access token")
	}
	// create jwt
	secretKeyJwt, _ := jwt.Parse(apiKey, nil)
	if secretKeyJwt == nil || secretKeyJwt.Claims == nil {
		return nil, errors.New("invalid secret key")
	}
	if err := secretKeyJwt.Claims.Valid(); err != nil {
		return nil, err
	}
	authTokenJwt, _ := jwt.Parse(accessResp.AccessToken, nil)
	if err := authTokenJwt.Claims.Valid(); err != nil {
		return nil, err
	}
	if authTokenJwt == nil || authTokenJwt.Claims == nil {
		return nil, errors.New("invalid auth token")
	}
	return &defaultSoftcorpAuthorize{
		accessApi:    accessApi,
		secretKeyJwt: secretKeyJwt,
		authTokenJwt: authTokenJwt,
	}, nil
}

func (sa *defaultSoftcorpAuthorize) GetAccessToken(ctx context.Context) (string, error) {
	sa.Lock()
	defer sa.Unlock()
	if err := sa.authTokenJwt.Claims.Valid(); err != nil {
		accessResp, err := sa.accessApi.GenerateAccessToken(ctx, &go_cloud.ProjectRequest{
			PrivateKey: sa.secretKeyJwt.Raw,
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
		sa.authTokenJwt = authTokenJwt
	}
	return sa.authTokenJwt.Raw, nil
}
