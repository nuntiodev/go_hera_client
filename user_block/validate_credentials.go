package user_block

import (
	"context"
	"github.com/nuntiodev/go-hera/nuntio_authorize"
	"github.com/nuntiodev/go-hera/nuntio_options"
	"github.com/nuntiodev/hera-proto/go_hera"
)

type ValidateCredentialsUserRequest struct {
	// external required fields
	password    string
	findOptions *nuntio_options.FindOptions
	// internal required fields
	encryptionKey string
	namespace     string
	client        go_hera.ServiceServer
	authorize     nuntio_authorize.Authorize
}

func (r *ValidateCredentialsUserRequest) Execute(ctx context.Context) (*go_hera.User, error) {
	accessToken, err := r.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	if r.findOptions == nil || r.findOptions.Validate() == false {
		return nil, invalidFindOptionsErr
	}
	validateUser := &go_hera.User{
		Email:    r.findOptions.Email,
		Id:       r.findOptions.Id,
		Username: r.findOptions.Username,
		Password: r.password,
	}
	resp, err := r.client.ValidateCredentials(ctx, &go_hera.UserRequest{
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
		client:        s.userClient,
		authorize:     s.authorize,
	}
}
