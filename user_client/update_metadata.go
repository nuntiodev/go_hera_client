package user_client

import (
	"context"
	"encoding/json"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/go-blocks/softcorp_authorize"
	"github.com/softcorp-io/go-blocks/softcorp_options"
)

type UpdateMetadataUserRequest struct {
	// external optional fields
	metadata interface{}
	// external required fields
	findOptions *softcorp_options.FindOptions
	// internal required fields
	encryptionKey string
	namespace     string
	userClient    go_block.UserServiceClient
	authorize     softcorp_authorize.Authorize
}

func (r *UpdateMetadataUserRequest) SetMetadata(metadata interface{}) *UpdateMetadataUserRequest {
	r.metadata = metadata
	return r
}

func (r *UpdateMetadataUserRequest) Execute(ctx context.Context) (*go_block.User, error) {
	accessToken, err := r.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	if r.findOptions == nil || r.findOptions.Validate() == false {
		return nil, invalidFindOptionsErr
	}
	findUser := &go_block.User{
		Email:      r.findOptions.Email,
		Id:         r.findOptions.Id,
		OptionalId: r.findOptions.OptionalId,
	}
	updateUser := &go_block.User{}
	if r.metadata != nil {
		jsonMetadata, err := json.Marshal(r.metadata)
		if err != nil {
			return nil, err
		}
		updateUser.Metadata = string(jsonMetadata)
	} else {
		r.metadata = ""
	}
	userResp, err := r.userClient.UpdateMetadata(ctx, &go_block.UserRequest{
		CloudToken:    accessToken,
		EncryptionKey: r.encryptionKey,
		Update:        updateUser,
		User:          findUser,
		Namespace:     r.namespace,
	})
	if err != nil {
		return nil, err
	}
	if userResp == nil || userResp.User == nil {
		return nil, internalServerError
	}
	return userResp.User, nil
}

func (s *defaultSocialServiceClient) UpdateMetadata(findOptions *softcorp_options.FindOptions) *UpdateMetadataUserRequest {
	return &UpdateMetadataUserRequest{
		findOptions:   findOptions,
		encryptionKey: s.encryptionKey,
		namespace:     s.namespace,
		userClient:    s.userClient,
		authorize:     s.authorize,
	}
}
