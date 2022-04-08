package user_client

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/go-blocks/softcorp_authorize"
	"github.com/softcorp-io/go-blocks/softcorp_options"
	"google.golang.org/grpc"
)

var (
	internalServerError   = errors.New("internal server error")
	invalidFindOptionsErr = errors.New("at least one find parameter is required")
	tokenIsEmptyErr       = errors.New("token is empty")
)

type PublicKey struct {
	publicKey []byte
	fetchedAt time.Time
	sync.Mutex
}

type UserClient interface {
	Create() *CreateUserRequest
	UpdatePassword(findOptions *softcorp_options.FindOptions, password string) *UpdatePasswordUserRequest
	UpdateMetadata(findOptions *softcorp_options.FindOptions) *UpdateMetadataUserRequest
	UpdateEmail(findOptions *softcorp_options.FindOptions, email string) *UpdateEmailUserRequest
	UpdateOptionalId(findOptions *softcorp_options.FindOptions, optionalId string) *UpdateOptionalIdUserRequest
	UpdateImage(findOptions *softcorp_options.FindOptions, image string) *UpdateImageUserRequest
	UpdateSecurity(findOptions *softcorp_options.FindOptions) *UpdateSecurityUserRequest
	Get(findOptions *softcorp_options.FindOptions) *GetUserRequest
	GetAll() *GetAllUserRequest
	ValidateCredentials(findOptions *softcorp_options.FindOptions, password string) *ValidateCredentialsUserRequest
	Login(findOptions *softcorp_options.FindOptions) *LoginUserRequest
	PublicKeys() *PublicKeysUserRequest
	RefreshToken(refreshToken string) *RefreshTokenUserRequest
	ValidateToken(jwtToken string) (*go_block.User, error)
	BlockToken(token string) *BlockTokenUserRequest
	Delete(findOptions *softcorp_options.FindOptions) *DeleteUserRequest
	DeleteAll() *DeleteAllUserRequest
}

type defaultSocialServiceClient struct {
	userClient    go_block.UserServiceClient
	authorize     softcorp_authorize.Authorize
	publicKey     *PublicKey
	namespace     string
	encryptionKey string
}

func (c *defaultSocialServiceClient) getPublicKey() ([]byte, error) {
	// get public key
	c.publicKey.Lock()
	defer c.publicKey.Unlock()
	if time.Now().Sub(c.publicKey.fetchedAt) < time.Hour*6 {
		return c.publicKey.publicKey, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	accessToken, err := c.authorize.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	publicKeyResp, err := c.userClient.PublicKeys(context.Background(), &go_block.UserRequest{
		CloudToken: accessToken,
	})
	if err != nil {
		return nil, err
	}
	publicKey, ok := publicKeyResp.PublicKeys["public-jwt-key"]
	if !ok || len(publicKey) <= 10 {
		return nil, errors.New("could not fetch public jwt key")
	}
	c.publicKey.publicKey = publicKey
	c.publicKey.fetchedAt = time.Now()
	return publicKey, nil
}

func New(apiUrl string, authorize softcorp_authorize.Authorize, encryptionKey, namespace string, dialOptions grpc.DialOption) (UserClient, error) {
	// setup grpc connection to user service
	userClientConn, err := grpc.Dial(apiUrl, dialOptions)
	if err != nil {
		return nil, err
	}
	userClient := go_block.NewUserServiceClient(userClientConn)
	return &defaultSocialServiceClient{
		encryptionKey: encryptionKey,
		userClient:    userClient,
		publicKey:     &PublicKey{},
		namespace:     namespace,
		authorize:     authorize,
	}, nil
}
