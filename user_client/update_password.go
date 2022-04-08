package user_client

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/go-blocks/softcorp_authorize"
	"github.com/softcorp-io/go-blocks/softcorp_options"
)

type UpdatePasswordUserRequest struct {
	// external optional fields
	validatePassword bool
	// external required fields
	password    string
	findOptions *softcorp_options.FindOptions
	// internal required fields
	encryptionKey string
	namespace     string
	userClient    go_block.UserServiceClient
	authorize     softcorp_authorize.Authorize
}

func (r *UpdatePasswordUserRequest) SetValidatePassword(validatePassword bool) *UpdatePasswordUserRequest {
	r.validatePassword = validatePassword
	return r
}

func (r *UpdatePasswordUserRequest) Execute(ctx context.Context) (*go_block.User, error) {
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
		Password: r.password,
	}
	userResp, err := r.userClient.UpdatePassword(ctx, &go_block.UserRequest{
		CloudToken:       accessToken,
		EncryptionKey:    r.encryptionKey,
		Update:           updateUser,
		User:             findUser,
		Namespace:        r.namespace,
		ValidatePassword: r.validatePassword,
	})
	if err != nil {
		return nil, err
	}
	if userResp == nil || userResp.User == nil {
		return nil, internalServerError
	}
	return userResp.User, nil
}

func (s *defaultSocialServiceClient) UpdatePassword(findOptions *softcorp_options.FindOptions, password string) *UpdatePasswordUserRequest {
	return &UpdatePasswordUserRequest{
		password:      password,
		findOptions:   findOptions,
		encryptionKey: s.encryptionKey,
		namespace:     s.namespace,
		userClient:    s.userClient,
		authorize:     s.authorize,
	}
}
