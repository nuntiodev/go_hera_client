package user_client

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
)

func (s *defaultSocialServiceClient) GetAll(ctx context.Context) ([]*go_block.User, error) {
	accessToken, err := s.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	userResp, err := s.userClient.GetAll(ctx, &go_block.UserRequest{
		CloudToken:    accessToken,
		EncryptionKey: s.encryptionKey,
		Namespace:     s.namespace,
	})
	if err != nil {
		return nil, err
	}
	if userResp == nil || userResp.Users == nil {
		return nil, internalServerError
	}
	return userResp.Users, nil
}
