package user_block

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/go-blocks/nuntio_authorize"
	"github.com/nuntiodev/go-blocks/nuntio_options"
)

type DeleteUserRequest struct {
	// external required fields
	findOptions *nuntio_options.FindOptions
	// internal required fields
	namespace  string
	userClient go_block.UserServiceClient
	authorize  nuntio_authorize.Authorize
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

func (s *defaultSocialServiceClient) Delete(findOptions *nuntio_options.FindOptions) *DeleteUserRequest {
	return &DeleteUserRequest{
		findOptions: findOptions,
		namespace:   s.namespace,
		userClient:  s.userClient,
		authorize:   s.authorize,
	}
}
