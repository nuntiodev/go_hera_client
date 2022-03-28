package user_client

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
)

func (s *defaultSocialServiceClient) PublicKeys(ctx context.Context) (*go_block.Token, error) {
	accessToken, err := s.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := s.userClient.PublicKeys(ctx, &go_block.UserRequest{
		CloudToken: accessToken,
	})
	if err != nil {
		return nil, err
	}
	return resp.Token, nil
}
