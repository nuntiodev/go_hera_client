package user_client

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
)

func (s *defaultSocialServiceClient) BlockToken(ctx context.Context, token string) error {
	accessToken, err := s.authorize.GetAccessToken(ctx)
	if err != nil {
		return err
	}
	if token == "" {
		return tokenIsEmptyErr
	}
	// when blocking a token the server does not distinguish
	// between whether it is a access or refresh token - it just blocks it,
	// so it does not matter whether we set the access or refresh token - just one of them
	tokenStruct := &go_block.Token{
		AccessToken:  token,
		RefreshToken: token,
	}
	resp, err := s.userClient.BlockToken(ctx, &go_block.UserRequest{
		CloudToken: accessToken,
		Token:      tokenStruct,
		Namespace:  s.namespace,
	})
	if err != nil {
		return err
	}
	if resp == nil || resp.User == nil {
		return internalServerError
	}
	return nil
}
