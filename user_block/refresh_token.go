package user_block

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/go-blocks/nuntio_authorize"
)

type RefreshTokenUserRequest struct {
	// external required fields
	refreshToken string
	// internal required fields
	namespace  string
	userClient go_block.UserServiceClient
	authorize  nuntio_authorize.Authorize
}

func (r *RefreshTokenUserRequest) Execute(ctx context.Context) (*go_block.Token, error) {
	accessToken, err := r.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	if r.refreshToken == "" {
		return nil, tokenIsEmptyErr
	}
	resp, err := r.userClient.RefreshToken(ctx, &go_block.UserRequest{
		CloudToken: accessToken,
		Token: &go_block.Token{
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
