package user_block

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/go-blocks/nuntio_authorize"
	"github.com/nuntiodev/go-blocks/nuntio_options"
	"google.golang.org/grpc"
)

var (
	internalServerError   = errors.New("internal server error")
	invalidFindOptionsErr = errors.New("at least one find parameter is required")
	tokenIsEmptyErr       = errors.New("token is empty")
)

type PublicKey struct {
	publicKey string
	fetchedAt time.Time
	sync.Mutex
}

type UserBlock interface {
	Create() *CreateUserRequest
	UpdatePassword(findOptions *nuntio_options.FindOptions, password string) *UpdatePasswordUserRequest
	UpdateMetadata(findOptions *nuntio_options.FindOptions) *UpdateMetadataUserRequest
	UpdateEmail(findOptions *nuntio_options.FindOptions, email string) *UpdateEmailUserRequest
	UpdateOptionalId(findOptions *nuntio_options.FindOptions, optionalId string) *UpdateOptionalIdUserRequest
	UpdateImage(findOptions *nuntio_options.FindOptions, image string) *UpdateImageUserRequest
	UpdateSecurity(findOptions *nuntio_options.FindOptions) *UpdateSecurityUserRequest
	Get(findOptions *nuntio_options.FindOptions) *GetUserRequest
	GetAll() *GetAllUserRequest
	ValidateCredentials(findOptions *nuntio_options.FindOptions, password string) *ValidateCredentialsUserRequest
	Login(findOptions *nuntio_options.FindOptions) *LoginUserRequest
	PublicKeys() *PublicKeysUserRequest
	RefreshToken(refreshToken string) *RefreshTokenUserRequest
	ValidateToken(ctx context.Context, jwtToken string, forceValidateServerSide bool) (*go_block.User, error)
	BlockToken(token string) *BlockTokenUserRequest
	Delete(findOptions *nuntio_options.FindOptions) *DeleteUserRequest
	DeleteAll() *DeleteAllUserRequest
}

type defaultSocialServiceClient struct {
	userClient    go_block.UserServiceClient
	authorize     nuntio_authorize.Authorize
	publicKey     *PublicKey
	namespace     string
	encryptionKey string
}

func (c *defaultSocialServiceClient) getPublicKey() (string, error) {
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
		return "", err
	}
	publicKeyResp, err := c.userClient.PublicKeys(context.Background(), &go_block.UserRequest{
		CloudToken: accessToken,
	})
	if err != nil {
		return "", err
	}
	publicKey, ok := publicKeyResp.PublicKeys["public-jwt-key"]
	if !ok || len(publicKey) <= 10 {
		return "", errors.New("could not fetch public jwt key")
	}
	c.publicKey.publicKey = publicKey
	c.publicKey.fetchedAt = time.Now()
	return publicKey, nil
}

func New(apiUrl string, authorize nuntio_authorize.Authorize, encryptionKey, namespace string, dialOptions grpc.DialOption) (UserBlock, error) {
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
