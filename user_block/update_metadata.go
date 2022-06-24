package user_block

import (
	"context"
	"encoding/json"
	"github.com/nuntiodev/go-hera/nuntio_authorize"
	"github.com/nuntiodev/go-hera/nuntio_options"
)

type UpdateMetadataUserRequest struct {
	// external optional fields
	metadata interface{}
	// external required fields
	findOptions *nuntio_options.FindOptions
	// internal required fields
	encryptionKey string
	namespace     string
	userClient    go_hera.UserServiceClient
	authorize     nuntio_authorize.Authorize
}

func (r *UpdateMetadataUserRequest) SetMetadata(metadata interface{}) *UpdateMetadataUserRequest {
	r.metadata = metadata
	return r
}

func (r *UpdateMetadataUserRequest) Execute(ctx context.Context) (*go_hera.User, error) {
	accessToken, err := r.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	if r.findOptions == nil || r.findOptions.Validate() == false {
		return nil, invalidFindOptionsErr
	}
	findUser := &go_hera.User{
		Email:    r.findOptions.Email,
		Id:       r.findOptions.Id,
		Username: r.findOptions.Username,
	}
	updateUser := &go_hera.User{}
	if r.metadata != nil {
		jsonMetadata, err := json.Marshal(r.metadata)
		if err != nil {
			return nil, err
		}
		updateUser.Metadata = string(jsonMetadata)
	} else {
		r.metadata = ""
	}
	userResp, err := r.userClient.UpdateMetadata(ctx, &go_hera.UserRequest{
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

func (s *defaultSocialServiceClient) UpdateMetadata(findOptions *nuntio_options.FindOptions) *UpdateMetadataUserRequest {
	return &UpdateMetadataUserRequest{
		findOptions:   findOptions,
		encryptionKey: s.encryptionKey,
		namespace:     s.namespace,
		userClient:    s.userClient,
		authorize:     s.authorize,
	}
}
