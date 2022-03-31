package user_client

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/go-blocks/softcorp_options"
)

func (s *defaultSocialServiceClient) ValidateCredentials(ctx context.Context, findOptions *softcorp_options.FindOptions, password string) (*go_block.User, error) {
	accessToken, err := s.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	if findOptions == nil || findOptions.Validate() == false {
		return nil, invalidFindOptionsErr
	}
	validateUser := &go_block.User{
		Email:      findOptions.Email,
		Id:         findOptions.Id,
		OptionalId: findOptions.OptionalId,
		Password:   password,
	}
	resp, err := s.userClient.ValidateCredentials(ctx, &go_block.UserRequest{
		CloudToken: accessToken,
		User:       validateUser,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.User == nil {
		return nil, internalServerError
	}
	return resp.User, nil
}
