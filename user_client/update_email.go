package user_client

import (
	"context"
	"github.com/badoux/checkmail"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/go-blocks/softcorp_options"
)

func (s *defaultSocialServiceClient) UpdateEmail(ctx context.Context, findOptions *softcorp_options.FindOptions, email string) (*go_block.User, error) {
	accessToken, err := s.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	if findOptions == nil || findOptions.Validate() == false {
		return nil, invalidFindOptionsErr
	}
	if err := checkmail.ValidateFormat(email); err != nil && email != "" {
		return nil, err
	}
	findUser := &go_block.User{
		Email:      findOptions.Email,
		Id:         findOptions.Id,
		OptionalId: findOptions.OptionalId,
	}
	updateUser := &go_block.User{
		Email: email,
	}
	userResp, err := s.userClient.UpdateEmail(ctx, &go_block.UserRequest{
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
