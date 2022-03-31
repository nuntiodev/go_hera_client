package user_client

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
)

func (s *defaultSocialServiceClient) DeleteAll(ctx context.Context) error {
	accessToken, err := s.authorize.GetAccessToken(ctx)
	if err != nil {
		return err
	}
	_, err = s.userClient.DeleteNamespace(ctx, &go_block.UserRequest{
		CloudToken: accessToken,
		Namespace:  s.namespace,
	})
	if err != nil {
		return err
	}
	return nil
}
