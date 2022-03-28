package user_client

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/go-blocks/options"
)

func (s *defaultSocialServiceClient) Login(ctx context.Context, findOptions *options.FindOptions, password string) (*go_block.Token, error) {
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
	resp, err := s.userClient.Login(ctx, &go_block.UserRequest{
		CloudToken: accessToken,
		User:       validateUser,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.Token == nil {
		return nil, internalServerError
	}
	return resp.Token, nil
}
