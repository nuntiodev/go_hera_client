package user_block

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/go-blocks/nuntio_authorize"
	"github.com/nuntiodev/go-blocks/nuntio_options"
)

type UpdateOptionalIdUserRequest struct {
	// external required fields
	optionalId  string
	findOptions *nuntio_options.FindOptions
	// internal required fields
	encryptionKey string
	namespace     string
	userClient    go_block.UserServiceClient
	authorize     nuntio_authorize.Authorize
}

func (r *UpdateOptionalIdUserRequest) Execute(ctx context.Context) (*go_block.User, error) {
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
	updateUser := &go_block.User{
		Image: r.optionalId,
	}
	userResp, err := r.userClient.UpdateOptionalId(ctx, &go_block.UserRequest{
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

func (s *defaultSocialServiceClient) UpdateOptionalId(findOptions *nuntio_options.FindOptions, optionalId string) *UpdateOptionalIdUserRequest {
	return &UpdateOptionalIdUserRequest{
		optionalId:    optionalId,
		findOptions:   findOptions,
		encryptionKey: s.encryptionKey,
		namespace:     s.namespace,
		userClient:    s.userClient,
		authorize:     s.authorize,
	}
}