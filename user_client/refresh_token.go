package user_client

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
)

func (s *defaultSocialServiceClient) RefreshToken(ctx context.Context, refreshToken string) (*go_block.Token, error) {
	accessToken, err := s.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	if refreshToken == "" {
		return nil, tokenIsEmptyErr
	}
	resp, err := s.userClient.RefreshToken(ctx, &go_block.UserRequest{
		CloudToken: accessToken,
		Token: &go_block.Token{
			RefreshToken: refreshToken,
		},
		Namespace: s.namespace,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.Token == nil {
		return nil, internalServerError
	}
	return resp.Token, nil
}
