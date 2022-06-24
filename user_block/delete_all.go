package user_block

import (
	"context"
	"github.com/nuntiodev/go-hera/nuntio_authorize"
)

type DeleteAllUserRequest struct {
	// internal required fields
	namespace  string
	userClient go_hera.UserServiceClient
	authorize  nuntio_authorize.Authorize
}

func (r *DeleteAllUserRequest) Execute(ctx context.Context) error {
	accessToken, err := r.authorize.GetAccessToken(ctx)
	if err != nil {
		return err
	}
	_, err = r.userClient.DeleteNamespace(ctx, &go_hera.UserRequest{
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
