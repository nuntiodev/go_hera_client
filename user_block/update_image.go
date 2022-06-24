package user_block

import (
	"context"
	"github.com/nuntiodev/go-hera/nuntio_authorize"
	"github.com/nuntiodev/go-hera/nuntio_options"
)

type UpdateImageUserRequest struct {
	// external required fields
	image       string
	findOptions *nuntio_options.FindOptions
	// internal required fields
	encryptionKey string
	namespace     string
	userClient    go_hera.UserServiceClient
	authorize     nuntio_authorize.Authorize
}

func (r *UpdateImageUserRequest) Execute(ctx context.Context) (*go_hera.User, error) {
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
	updateUser := &go_hera.User{
		Image: r.image,
	}
	userResp, err := r.userClient.UpdateImage(ctx, &go_hera.UserRequest{
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

func (s *defaultSocialServiceClient) UpdateImage(findOptions *nuntio_options.FindOptions, image string) *UpdateImageUserRequest {
	return &UpdateImageUserRequest{
		image:         image,
		findOptions:   findOptions,
		encryptionKey: s.encryptionKey,
		namespace:     s.namespace,
		userClient:    s.userClient,
		authorize:     s.authorize,
	}
}
