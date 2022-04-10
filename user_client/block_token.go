package user_client

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/go-blocks/softcorp_authorize"
)

type BlockTokenUserRequest struct {
	// external required fields
	token string
	// internal required fields
	namespace  string
	userClient go_block.UserServiceClient
	authorize  softcorp_authorize.Authorize
}

func (r *BlockTokenUserRequest) Execute(ctx context.Context) error {
	accessToken, err := r.authorize.GetAccessToken(ctx)
	if err != nil {
		return err
	}
	if r.token == "" {
		return tokenIsEmptyErr
	}
	resp, err := r.userClient.BlockToken(ctx, &go_block.UserRequest{
		CloudToken:   accessToken,
		Namespace:    r.namespace,
		TokenPointer: r.token,
	})
	if err != nil {
		return err
	}
	if resp == nil || resp.User == nil {
		return internalServerError
	}
	return nil
}

func (s *defaultSocialServiceClient) BlockToken(token string) *BlockTokenUserRequest {
	return &BlockTokenUserRequest{
		token:      token,
		namespace:  s.namespace,
		userClient: s.userClient,
		authorize:  s.authorize,
	}
}
