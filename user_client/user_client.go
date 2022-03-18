package user_client

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/badoux/checkmail"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/go-blocks/authorize"
	"github.com/softcorp-io/go-blocks/options"
	"google.golang.org/grpc"
)

var (
	// USER_API_URL is the URL of the API user service backend.
	USER_API_URL          = "api.softcorp.io:443"
	internalServerError   = errors.New("internal server error")
	invalidFindOptionsErr = errors.New("at least one find parameter is required")
)

type UserClient interface {
	Create(ctx context.Context, password string, userOptions *options.UserOptions, metadataOptions interface{}) (*go_block.User, error)
	UpdatePassword(ctx context.Context, findOptions *options.FindOptions, password string) (*go_block.User, error)
	UpdateMetadata(ctx context.Context, findOptions *options.FindOptions, metadataOptions interface{}) (*go_block.User, error)
	UpdateEmail(ctx context.Context, findOptions *options.FindOptions, email string) (*go_block.User, error)
	UpdateOptionalId(ctx context.Context, findOptions *options.FindOptions, optionalId string) (*go_block.User, error)
	UpdateImage(ctx context.Context, findOptions *options.FindOptions, imageUrl string) (*go_block.User, error)
	UpdateSecurity(ctx context.Context, findOptions *options.FindOptions, securityOptions *options.SecurityOptions) (*go_block.User, error)
	Get(ctx context.Context, findOptions *options.FindOptions) (*go_block.User, error)
	GetAll(ctx context.Context) ([]*go_block.User, error)
	ValidateCredentials(ctx context.Context, findOptions *options.FindOptions, password string) (*go_block.User, error)
	Delete(ctx context.Context, findOptions *options.FindOptions) error
	DeleteAll(ctx context.Context) error
}

type defaultSocialServiceClient struct {
	userClient    go_block.UserServiceClient
	authorize     authorize.Authorize
	encryptionKey string
}

func New(authorize authorize.Authorize, encryptionKey string, dialOptions grpc.DialOption) (UserClient, error) {
	// setup grpc connection to user service
	userClientConn, err := grpc.Dial(USER_API_URL, dialOptions)
	if err != nil {
		return nil, err
	}
	userClient := go_block.NewUserServiceClient(userClientConn)
	return &defaultSocialServiceClient{
		encryptionKey: encryptionKey,
		userClient:    userClient,
		authorize:     authorize,
	}, nil
}

func (s *defaultSocialServiceClient) Create(ctx context.Context, password string, userOptions *options.UserOptions, metadataOptions interface{}) (*go_block.User, error) {
	accessToken, err := s.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	createUser := &go_block.User{
		Password: password,
	}
	if userOptions != nil {
		createUser.Id = userOptions.Id
		createUser.OptionalId = userOptions.OptionalId
		createUser.Email = userOptions.Email
		createUser.Image = userOptions.Image
	}
	if metadataOptions != nil {
		jsonMetadata, err := json.Marshal(metadataOptions)
		if err != nil {
			return nil, err
		}
		createUser.Metadata = string(jsonMetadata)
	}
	userResp, err := s.userClient.Create(ctx, &go_block.UserRequest{
		AccessToken:   accessToken,
		EncryptionKey: s.encryptionKey,
		User:          createUser,
	})
	if err != nil {
		return nil, err
	}
	if userResp == nil || userResp.User == nil {
		return nil, internalServerError
	}
	return userResp.User, nil
}

func (s *defaultSocialServiceClient) UpdatePassword(ctx context.Context, findOptions *options.FindOptions, newPassword string) (*go_block.User, error) {
	accessToken, err := s.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	if findOptions == nil || findOptions.Validate() == false {
		return nil, invalidFindOptionsErr
	}
	findUser := &go_block.User{
		Email:      findOptions.Email,
		Id:         findOptions.Id,
		OptionalId: findOptions.OptionalId,
	}
	updateUser := &go_block.User{
		Password: newPassword,
	}
	userResp, err := s.userClient.UpdatePassword(ctx, &go_block.UserRequest{
		AccessToken:   accessToken,
		EncryptionKey: s.encryptionKey,
		Update:        updateUser,
		User:          findUser,
	})
	if err != nil {
		return nil, err
	}
	if userResp == nil || userResp.User == nil {
		return nil, internalServerError
	}
	return userResp.User, nil
}

func (s *defaultSocialServiceClient) UpdateMetadata(ctx context.Context, findOptions *options.FindOptions, metadataOptions interface{}) (*go_block.User, error) {
	accessToken, err := s.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	if findOptions == nil || findOptions.Validate() == false {
		return nil, invalidFindOptionsErr
	}
	findUser := &go_block.User{
		Email:      findOptions.Email,
		Id:         findOptions.Id,
		OptionalId: findOptions.OptionalId,
	}
	updateUser := &go_block.User{}
	if metadataOptions != nil {
		jsonMetadata, err := json.Marshal(metadataOptions)
		if err != nil {
			return nil, err
		}
		updateUser.Metadata = string(jsonMetadata)
	}
	userResp, err := s.userClient.UpdateMetadata(ctx, &go_block.UserRequest{
		AccessToken:   accessToken,
		EncryptionKey: s.encryptionKey,
		Update:        updateUser,
		User:          findUser,
	})
	if err != nil {
		return nil, err
	}
	if userResp == nil || userResp.User == nil {
		return nil, internalServerError
	}
	return userResp.User, nil
}

func (s *defaultSocialServiceClient) UpdateEmail(ctx context.Context, findOptions *options.FindOptions, email string) (*go_block.User, error) {
	accessToken, err := s.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	if findOptions == nil || findOptions.Validate() == false {
		return nil, invalidFindOptionsErr
	}
	if err := checkmail.ValidateFormat(email); err != nil && email != "" {
		return nil, err
	}
	findUser := &go_block.User{
		Email:      findOptions.Email,
		Id:         findOptions.Id,
		OptionalId: findOptions.OptionalId,
	}
	updateUser := &go_block.User{
		Email: email,
	}
	userResp, err := s.userClient.UpdateEmail(ctx, &go_block.UserRequest{
		AccessToken:   accessToken,
		EncryptionKey: s.encryptionKey,
		Update:        updateUser,
		User:          findUser,
	})
	if err != nil {
		return nil, err
	}
	if userResp == nil || userResp.User == nil {
		return nil, internalServerError
	}
	return userResp.User, nil
}

func (s *defaultSocialServiceClient) UpdateOptionalId(ctx context.Context, findOptions *options.FindOptions, optionalId string) (*go_block.User, error) {
	accessToken, err := s.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	if findOptions == nil || findOptions.Validate() == false {
		return nil, invalidFindOptionsErr
	}
	findUser := &go_block.User{
		Email:      findOptions.Email,
		Id:         findOptions.Id,
		OptionalId: findOptions.OptionalId,
	}
	updateUser := &go_block.User{
		OptionalId: optionalId,
	}
	userResp, err := s.userClient.UpdateOptionalId(ctx, &go_block.UserRequest{
		AccessToken:   accessToken,
		EncryptionKey: s.encryptionKey,
		Update:        updateUser,
		User:          findUser,
	})
	if err != nil {
		return nil, err
	}
	if userResp == nil || userResp.User == nil {
		return nil, internalServerError
	}
	return userResp.User, nil
}

func (s *defaultSocialServiceClient) UpdateImage(ctx context.Context, findOptions *options.FindOptions, imageUrl string) (*go_block.User, error) {
	accessToken, err := s.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	if findOptions == nil || findOptions.Validate() == false {
		return nil, invalidFindOptionsErr
	}
	findUser := &go_block.User{
		Email:      findOptions.Email,
		Id:         findOptions.Id,
		OptionalId: findOptions.OptionalId,
	}
	updateUser := &go_block.User{
		Image: imageUrl,
	}
	userResp, err := s.userClient.UpdateImage(ctx, &go_block.UserRequest{
		AccessToken:   accessToken,
		EncryptionKey: s.encryptionKey,
		Update:        updateUser,
		User:          findUser,
	})
	if err != nil {
		return nil, err
	}
	if userResp == nil || userResp.User == nil {
		return nil, internalServerError
	}
	return userResp.User, nil
}

func (s *defaultSocialServiceClient) UpdateSecurity(ctx context.Context, findOptions *options.FindOptions, securityOptions *options.SecurityOptions) (*go_block.User, error) {
	accessToken, err := s.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	if securityOptions == nil {
		return nil, errors.New("security options are not allowed to be nil")
	}
	if findOptions == nil || findOptions.Validate() == false {
		return nil, invalidFindOptionsErr
	}
	findUser := &go_block.User{
		Email:      findOptions.Email,
		Id:         findOptions.Id,
		OptionalId: findOptions.OptionalId,
	}
	updateUser := &go_block.User{}
	userResp, err := s.userClient.UpdateSecurity(ctx, &go_block.UserRequest{
		AccessToken:   accessToken,
		EncryptionKey: s.encryptionKey,
		Update:        updateUser,
		User:          findUser,
	})
	if err != nil {
		return nil, err
	}
	if userResp == nil || userResp.User == nil {
		return nil, internalServerError
	}
	return userResp.User, nil
}

func (s *defaultSocialServiceClient) Get(ctx context.Context, findOptions *options.FindOptions) (*go_block.User, error) {
	accessToken, err := s.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	if findOptions == nil || findOptions.Validate() == false {
		return nil, invalidFindOptionsErr
	}
	getUser := &go_block.User{
		Email:      findOptions.Email,
		Id:         findOptions.Id,
		OptionalId: findOptions.OptionalId,
	}
	userResp, err := s.userClient.Get(ctx, &go_block.UserRequest{
		AccessToken:   accessToken,
		EncryptionKey: s.encryptionKey,
		User:          getUser,
	})
	if err != nil {
		return nil, err
	}
	if userResp == nil || userResp.User == nil {
		return nil, internalServerError
	}
	return userResp.User, nil
}

func (s *defaultSocialServiceClient) GetAll(ctx context.Context) ([]*go_block.User, error) {
	accessToken, err := s.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	userResp, err := s.userClient.Get(ctx, &go_block.UserRequest{
		AccessToken:   accessToken,
		EncryptionKey: s.encryptionKey,
	})
	if err != nil {
		return nil, err
	}
	if userResp == nil || userResp.Users == nil {
		return nil, internalServerError
	}
	return userResp.Users, nil
}

func (s *defaultSocialServiceClient) ValidateCredentials(ctx context.Context, findOptions *options.FindOptions, password string) (*go_block.User, error) {
	accessToken, err := s.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	if findOptions == nil || findOptions.Validate() == false {
		return nil, invalidFindOptionsErr
	}
	validateUser := &go_block.User{
		Email:      findOptions.Email,
		Id:         findOptions.Id,
		OptionalId: findOptions.OptionalId,
		Password:   password,
	}
	resp, err := s.userClient.ValidateCredentials(ctx, &go_block.UserRequest{
		AccessToken: accessToken,
		User:        validateUser,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.User == nil {
		return nil, internalServerError
	}
	return resp.User, nil
}

func (s *defaultSocialServiceClient) Delete(ctx context.Context, findOptions *options.FindOptions) error {
	accessToken, err := s.authorize.GetAccessToken(ctx)
	if err != nil {
		return err
	}
	if findOptions == nil || findOptions.Validate() == false {
		return invalidFindOptionsErr
	}
	deleteUser := &go_block.User{
		Email:      findOptions.Email,
		Id:         findOptions.Id,
		OptionalId: findOptions.OptionalId,
	}
	_, err = s.userClient.Delete(ctx, &go_block.UserRequest{
		AccessToken: accessToken,
		User:        deleteUser,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *defaultSocialServiceClient) DeleteAll(ctx context.Context) error {
	accessToken, err := s.authorize.GetAccessToken(ctx)
	if err != nil {
		return err
	}
	_, err = s.userClient.DeleteNamespace(ctx, &go_block.UserRequest{
		AccessToken: accessToken,
	})
	if err != nil {
		return err
	}
	return nil
}
