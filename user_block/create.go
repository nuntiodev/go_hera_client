package user_block

import (
	"context"
	"encoding/json"
	"github.com/nuntiodev/go-hera/nuntio_authorize"
	"github.com/nuntiodev/go-hera/nuntio_options"
)

type CreateUserRequest struct {
	// external optional fields
	userOptions      *nuntio_options.UserOptions
	metadata         interface{}
	password         string
	validatePassword bool
	// internal required fields
	encryptionKey string
	namespace     string
	userClient    go_hera.UserServiceClient
	authorize     nuntio_authorize.Authorize
}

func (r *CreateUserRequest) SetUserOptions(options *nuntio_options.UserOptions) *CreateUserRequest {
	if options != nil {
		r.userOptions = options
	}
	return r
}

func (r *CreateUserRequest) SetMetadata(metadata interface{}) *CreateUserRequest {
	if metadata != nil {
		r.metadata = metadata
	}
	return r
}

func (r *CreateUserRequest) SetPassword(password string) *CreateUserRequest {
	if password != "" {
		r.password = password
	}
	return r
}

func (r *CreateUserRequest) SetValidatePassword(validatePassword bool) *CreateUserRequest {
	if validatePassword {
		r.validatePassword = validatePassword
	}
	return r
}

func (r *CreateUserRequest) Execute(ctx context.Context) (*go_hera.User, error) {
	accessToken, err := r.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	createUser := &go_hera.User{
		Password: r.password,
	}
	if r.userOptions != nil {
		createUser.Id = r.userOptions.Id
		createUser.Username = r.userOptions.Username
		createUser.Email = r.userOptions.Email
		createUser.Image = r.userOptions.Image
	}
	if r.metadata != nil {
		jsonMetadata, err := json.Marshal(r.metadata)
		if err != nil {
			return nil, err
		}
		createUser.Metadata = string(jsonMetadata)
	}
	userResp, err := r.userClient.Create(ctx, &go_hera.UserRequest{
		CloudToken:    accessToken,
		EncryptionKey: r.encryptionKey,
		User:          createUser,
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

func (s *defaultSocialServiceClient) Create() *CreateUserRequest {
	return &CreateUserRequest{
		encryptionKey: s.encryptionKey,
		namespace:     s.namespace,
		userClient:    s.userClient,
		authorize:     s.authorize,
	}
}
