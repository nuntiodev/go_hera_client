package user_block

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/go-blocks/nuntio_authorize"
	"github.com/nuntiodev/go-blocks/nuntio_options"
)

type LoginUserRequest struct {
	// external required fields
	findOptions *nuntio_options.FindOptions
	// external optional fields
	password string
	// internal required fields
	namespace  string
	userClient go_block.UserServiceClient
	authorize  nuntio_authorize.Authorize
}

func (r *LoginUserRequest) SetPassword(password string) *LoginUserRequest {
	r.password = password
	return r
}

func (r *LoginUserRequest) Execute(ctx context.Context) (*go_block.Token, error) {
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
	resp, err := r.userClient.Login(ctx, &go_block.UserRequest{
		CloudToken: accessToken,
		User:       validateUser,
		Namespace:  r.namespace,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.Token == nil {
		return nil, internalServerError
	}
	return resp.Token, nil
}

func (s *defaultSocialServiceClient) Login(findOptions *nuntio_options.FindOptions) *LoginUserRequest {
	return &LoginUserRequest{
		findOptions: findOptions,
		namespace:   s.namespace,
		userClient:  s.userClient,
		authorize:   s.authorize,
	}
}
