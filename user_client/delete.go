package user_client

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/go-blocks/options"
)

func (s *defaultSocialServiceClient) Delete(ctx context.Context, findOptions *options.FindOptions) error {
	accessToken, err := s.authorize.GetAccessToken(ctx)
	if err != nil {
		return err
	}
	if findOptions == nil || findOptions.Validate() == false {
		return invalidFindOptionsErr
	}
	deleteUser := &go_block.User{
		Email:      findOptions.Email,
		Id:         findOptions.Id,
		OptionalId: findOptions.OptionalId,
	}
	_, err = s.userClient.Delete(ctx, &go_block.UserRequest{
		CloudToken: accessToken,
		User:       deleteUser,
	})
	if err != nil {
		return err
	}
	return nil
}