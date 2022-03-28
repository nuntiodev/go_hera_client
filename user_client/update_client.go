package user_client

import (
	"context"
	"encoding/json"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/go-blocks/options"
)

func (s *defaultSocialServiceClient) UpdateMetadata(ctx context.Context, findOptions *options.FindOptions, metadataOptions interface{}) (*go_block.User, error) {
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
	updateUser := &go_block.User{}
	if metadataOptions != nil {
		jsonMetadata, err := json.Marshal(metadataOptions)
		if err != nil {
			return nil, err
		}
		updateUser.Metadata = string(jsonMetadata)
	}
	userResp, err := s.userClient.UpdateMetadata(ctx, &go_block.UserRequest{
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
