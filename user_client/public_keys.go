package user_client

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/go-blocks/softcorp_authorize"
)

type PublicKeysUserRequest struct {
	// internal required fields
	namespace  string
	userClient go_block.UserServiceClient
	authorize  softcorp_authorize.Authorize
}

func (r *PublicKeysUserRequest) Execute(ctx context.Context) (*map[string][]byte, error) {
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
