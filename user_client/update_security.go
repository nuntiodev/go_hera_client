package user_client

import (
	"context"
	"errors"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/go-blocks/softcorp_options"
)

func (s *defaultSocialServiceClient) UpdateSecurity(ctx context.Context, findOptions *softcorp_options.FindOptions, securityOptions *softcorp_options.SecurityOptions) (*go_block.User, error) {
	accessToken, err := s.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	if securityOptions == nil {
		return nil, errors.New("security softcorp_options are not allowed to be nil")
	}
	if findOptions == nil || findOptions.Validate() == false {
		return nil, invalidFindOptionsErr
	}
	findUser := &go_block.User{
		Email:      findOptions.Email,
		Id:         findOptions.Id,
		OptionalId: findOptions.OptionalId,
	}
	updateUser := &go_block.User{}
	userResp, err := s.userClient.UpdateSecurity(ctx, &go_block.UserRequest{
		CloudToken:    accessToken,
		EncryptionKey: s.encryptionKey,
		Update:        updateUser,
		User:          findUser,
		Namespace:     s.namespace,
	})
	if err != nil {
		return nil, err
	}
	if userResp == nil || userResp.User == nil {
		return nil, internalServerError
	}
	return userResp.User, nil
}
