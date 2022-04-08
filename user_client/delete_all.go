package user_client

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/go-blocks/softcorp_authorize"
)

type DeleteAllUserRequest struct {
	// internal required fields
	namespace  string
	userClient go_block.UserServiceClient
	authorize  softcorp_authorize.Authorize
}

func (r *DeleteAllUserRequest) Execute(ctx context.Context) error {
	accessToken, err := r.authorize.GetAccessToken(ctx)
	if err != nil {
		return err
	}
	_, err = r.userClient.DeleteNamespace(ctx, &go_block.UserRequest{
		CloudToken: accessToken,
		Namespace:  r.namespace,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *defaultSocialServiceClient) DeleteAll() *DeleteAllUserRequest {
	return &DeleteAllUserRequest{
		namespace:  s.namespace,
		userClient: s.userClient,
		authorize:  s.authorize,
	}
}
