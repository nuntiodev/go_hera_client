package user_client

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/go-blocks/softcorp_options"
)

func (s *defaultSocialServiceClient) UpdateImage(ctx context.Context, findOptions *softcorp_options.FindOptions, imageUrl string) (*go_block.User, error) {
	accessToken, err := s.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	if findOptions == nil || findOptions.Validate() == false {
		return nil, invalidFindOptionsErr
	}
	findUser := &go_block.User{
		Email:      findOptions.Email,
		Id:         findOptions.Id,
		OptionalId: findOptions.OptionalId,
	}
	updateUser := &go_block.User{
		Image: imageUrl,
	}
	userResp, err := s.userClient.UpdateImage(ctx, &go_block.UserRequest{
		CloudToken:    accessToken,
		EncryptionKey: s.encryptionKey,
		Update:        updateUser,
		User:          findUser,
	})
	if err != nil {
		return nil, err
	}
	if userResp == nil || userResp.User == nil {
		return nil, internalServerError
	}
	return userResp.User, nil
}
