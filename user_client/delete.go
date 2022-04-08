package user_client

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/go-blocks/softcorp_authorize"
	"github.com/softcorp-io/go-blocks/softcorp_options"
)

type DeleteUserRequest struct {
	// external required fields
	findOptions *softcorp_options.FindOptions
	// internal required fields
	namespace  string
	userClient go_block.UserServiceClient
	authorize  softcorp_authorize.Authorize
}

func (r *DeleteUserRequest) Execute(ctx context.Context) error {
	accessToken, err := r.authorize.GetAccessToken(ctx)
	if err != nil {
		return err
	}
	if r.findOptions == nil || r.findOptions.Validate() == false {
		return invalidFindOptionsErr
	}
	deleteUser := &go_block.User{
		Email:      r.findOptions.Email,
		Id:         r.findOptions.Id,
		OptionalId: r.findOptions.OptionalId,
	}
	if _, err = r.userClient.Delete(ctx, &go_block.UserRequest{
		CloudToken: accessToken,
		User:       deleteUser,
		Namespace:  r.namespace,
	}); err != nil {
		return err
	}
	return nil
}

func (s *defaultSocialServiceClient) Delete(findOptions *softcorp_options.FindOptions) *DeleteUserRequest {
	return &DeleteUserRequest{
		findOptions: findOptions,
		namespace:   s.namespace,
		userClient:  s.userClient,
		authorize:   s.authorize,
	}
}
