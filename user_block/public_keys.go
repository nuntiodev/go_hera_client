package user_block

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/go-blocks/nuntio_authorize"
)

type PublicKeysUserRequest struct {
	// internal required fields
	namespace  string
	userClient go_block.UserServiceClient
	authorize  nuntio_authorize.Authorize
}

func (r *PublicKeysUserRequest) Execute(ctx context.Context) (*map[string]string, error) {
	accessToken, err := r.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := r.userClient.PublicKeys(ctx, &go_block.UserRequest{
		CloudToken: accessToken,
		Namespace:  r.namespace,
	})
	if err != nil {
		return nil, err
	}
	return &resp.PublicKeys, nil
}

func (s *defaultSocialServiceClient) PublicKeys() *PublicKeysUserRequest {
	return &PublicKeysUserRequest{
		namespace:  s.namespace,
		userClient: s.userClient,
		authorize:  s.authorize,
	}
}
