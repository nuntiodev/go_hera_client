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
	// USER_API_URL is the URL of the API user service backend.
	USER_API_URL          = "api.softcorp.io:443"
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
	Create(ctx context.Context, password string, userOptions *softcorp_options.UserOptions, metadataOptions interface{}) (*go_block.User, error)
	UpdatePassword(ctx context.Context, findOptions *softcorp_options.FindOptions, password string) (*go_block.User, error)
	UpdateMetadata(ctx context.Context, findOptions *softcorp_options.FindOptions, metadataOptions interface{}) (*go_block.User, error)
	UpdateEmail(ctx context.Context, findOptions *softcorp_options.FindOptions, email string) (*go_block.User, error)
	UpdateOptionalId(ctx context.Context, findOptions *softcorp_options.FindOptions, optionalId string) (*go_block.User, error)
	UpdateImage(ctx context.Context, findOptions *softcorp_options.FindOptions, imageUrl string) (*go_block.User, error)
	UpdateSecurity(ctx context.Context, findOptions *softcorp_options.FindOptions, securityOptions *softcorp_options.SecurityOptions) (*go_block.User, error)
	Get(ctx context.Context, findOptions *softcorp_options.FindOptions) (*go_block.User, error)
	GetAll(ctx context.Context) ([]*go_block.User, error)
	ValidateCredentials(ctx context.Context, findOptions *softcorp_options.FindOptions, password string) (*go_block.User, error)
	Login(ctx context.Context, findOptions *softcorp_options.FindOptions, password string) (*go_block.Token, error)
	PublicKeys(ctx context.Context) (*go_block.Token, error)
	RefreshToken(ctx context.Context, refreshToken string) (*go_block.Token, error)
	ValidateToken(ctx context.Context, jwtToken string) (*go_block.User, error)
	BlockToken(ctx context.Context, token string) error
	Delete(ctx context.Context, findOptions *softcorp_options.FindOptions) error
	DeleteAll(ctx context.Context) error
}

type defaultSocialServiceClient struct {
	userClient    go_block.UserServiceClient
	authorize     softcorp_authorize.Authorize
	publicKey     *PublicKey
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

func New(authorize softcorp_authorize.Authorize, encryptionKey string, dialOptions grpc.DialOption) (UserClient, error) {
	// setup grpc connection to user service
	userClientConn, err := grpc.Dial(USER_API_URL, dialOptions)
	if err != nil {
		return nil, err
	}
	userClient := go_block.NewUserServiceClient(userClientConn)
	return &defaultSocialServiceClient{
		encryptionKey: encryptionKey,
		userClient:    userClient,
		publicKey:     &PublicKey{},
		authorize:     authorize,
	}, nil
}
