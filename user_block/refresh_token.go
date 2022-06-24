package user_block

import (
	"context"
	"github.com/nuntiodev/go-hera/nuntio_authorize"
)

type RefreshTokenUserRequest struct {
	// external required fields
	refreshToken string
	// internal required fields
	namespace  string
	userClient go_hera.UserServiceClient
	authorize  nuntio_authorize.Authorize
}

func (r *RefreshTokenUserRequest) Execute(ctx context.Context) (*go_hera.Token, error) {
	accessToken, err := r.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	if r.refreshToken == "" {
		return nil, tokenIsEmptyErr
	}
	resp, err := r.userClient.RefreshToken(ctx, &go_hera.UserRequest{
		CloudToken: accessToken,
		Token: &go_hera.Token{
			RefreshToken: r.refreshToken,
		},
		Namespace: r.namespace,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.Token == nil {
		return nil, internalServerError
	}
	return resp.Token, nil
}

func (s *defaultSocialServiceClient) RefreshToken(refreshToken string) *RefreshTokenUserRequest {
	return &RefreshTokenUserRequest{
		refreshToken: refreshToken,
		namespace:    s.namespace,
		userClient:   s.userClient,
		authorize:    s.authorize,
	}
}
