package user_block

import (
	"context"
	"github.com/badoux/checkmail"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/go-blocks/nuntio_authorize"
	"github.com/nuntiodev/go-blocks/nuntio_options"
)

type UpdateEmailUserRequest struct {
	// external required fields
	email       string
	findOptions *nuntio_options.FindOptions
	// internal required fields
	encryptionKey string
	namespace     string
	userClient    go_block.UserServiceClient
	authorize     nuntio_authorize.Authorize
}

func (r *UpdateEmailUserRequest) Execute(ctx context.Context) (*go_block.User, error) {
	accessToken, err := r.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	if r.findOptions == nil || r.findOptions.Validate() == false {
		return nil, invalidFindOptionsErr
	}
	if err := checkmail.ValidateFormat(r.email); err != nil && r.email != "" {
		return nil, err
	}
	findUser := &go_block.User{
		Email:      r.findOptions.Email,
		Id:         r.findOptions.Id,
		OptionalId: r.findOptions.OptionalId,
	}
	updateUser := &go_block.User{
		Email: r.email,
	}
	userResp, err := r.userClient.UpdateEmail(ctx, &go_block.UserRequest{
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

func (s *defaultSocialServiceClient) UpdateEmail(findOptions *nuntio_options.FindOptions, email string) *UpdateEmailUserRequest {
	return &UpdateEmailUserRequest{
		email:         email,
		findOptions:   findOptions,
		encryptionKey: s.encryptionKey,
		namespace:     s.namespace,
		userClient:    s.userClient,
		authorize:     s.authorize,
	}
}
