package user_client

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/go-blocks/nuntio_authorize"
)

type GetAllUserRequest struct {
	// internal required fields
	namespace     string
	encryptionKey string
	userClient    go_block.UserServiceClient
	authorize     nuntio_authorize.Authorize
}

func (r *GetAllUserRequest) Execute(ctx context.Context) ([]*go_block.User, error) {
	accessToken, err := r.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	userResp, err := r.userClient.GetAll(ctx, &go_block.UserRequest{
		CloudToken:    accessToken,
		EncryptionKey: r.encryptionKey,
		Namespace:     r.namespace,
	})
	if err != nil {
		return nil, err
	}
	if userResp == nil || userResp.Users == nil {
		return nil, internalServerError
	}
	return userResp.Users, nil
}

func (s *defaultSocialServiceClient) GetAll() *GetAllUserRequest {
	return &GetAllUserRequest{
		encryptionKey: s.encryptionKey,
		namespace:     s.namespace,
		userClient:    s.userClient,
		authorize:     s.authorize,
	}
}
