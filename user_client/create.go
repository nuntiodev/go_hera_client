package user_client

import (
	"context"
	"encoding/json"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/go-blocks/softcorp_options"
)

func (s *defaultSocialServiceClient) Create(ctx context.Context, password string, userOptions *softcorp_options.UserOptions, metadataOptions interface{}) (*go_block.User, error) {
	accessToken, err := s.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	createUser := &go_block.User{
		Password: password,
	}
	if userOptions != nil {
		createUser.Id = userOptions.Id
		createUser.OptionalId = userOptions.OptionalId
		createUser.Email = userOptions.Email
		createUser.Image = userOptions.Image
	}
	if metadataOptions != nil {
		jsonMetadata, err := json.Marshal(metadataOptions)
		if err != nil {
			return nil, err
		}
		createUser.Metadata = string(jsonMetadata)
	}
	userResp, err := s.userClient.Create(ctx, &go_block.UserRequest{
		CloudToken:    accessToken,
		EncryptionKey: s.encryptionKey,
		User:          createUser,
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
