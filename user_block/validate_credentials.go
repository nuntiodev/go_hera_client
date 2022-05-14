package user_block

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/go-blocks/nuntio_authorize"
	"github.com/nuntiodev/go-blocks/nuntio_options"
)

type ValidateCredentialsUserRequest struct {
	// external required fields
	password    string
	findOptions *nuntio_options.FindOptions
	// internal required fields
	encryptionKey string
	namespace     string
	userClient    go_block.UserServiceClient
	authorize     nuntio_authorize.Authorize
}

func (r *ValidateCredentialsUserRequest) Execute(ctx context.Context) (*go_block.User, error) {
	accessToken, err := r.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	if r.findOptions == nil || r.findOptions.Validate() == false {
		return nil, invalidFindOptionsErr
	}
	validateUser := &go_block.User{
		Email:    r.findOptions.Email,
		Id:       r.findOptions.Id,
		Username: r.findOptions.Username,
		Password: r.password,
	}
	resp, err := r.userClient.ValidateCredentials(ctx, &go_block.UserRequest{
		CloudToken:    accessToken,
		User:          validateUser,
		Namespace:     r.namespace,
		EncryptionKey: r.encryptionKey,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.User == nil {
		return nil, internalServerError
	}
	return resp.User, nil
}

func (s *defaultSocialServiceClient) ValidateCredentials(findOptions *nuntio_options.FindOptions, password string) *ValidateCredentialsUserRequest {
	return &ValidateCredentialsUserRequest{
		password:      password,
		findOptions:   findOptions,
		encryptionKey: s.encryptionKey,
		namespace:     s.namespace,
		userClient:    s.userClient,
		authorize:     s.authorize,
	}
}
