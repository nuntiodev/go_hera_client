package softcorp

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/softcorp-io/go-blocks/softcorp_authorize"
	"github.com/softcorp-io/go-blocks/softcorp_credentials"
	"github.com/softcorp-io/go-blocks/user_client"
	"google.golang.org/grpc"
)

var (
	// ENCRYPTION_KEY  is used to encrypt clients data under the given key
	ENCRYPTION_KEY = ""
	// API_KEY  is used to connect your application to Softcorp Cloud
	API_KEY = ""
	// API_URL  is the URL the SDK will try to connect to
	API_URL = "api.softcorp.io:443"
	// AUTHORIZE is used to override the default softcorp_authorize interface which is used to validate tokens
	// if you don't want any authorization, set it to softcorp_authorize.AUTHORIZE = softcorp_authorize.NoAuthorization.
	AUTHORIZE softcorp_authorize.Authorize
	// CREDENTIALS defines what security is passed to softcorp_credentials.Dial and (can be overwritten)
	// you can provide your own, or use softcorp_credentials.TRANSPORT_CREDENTIALS = softcorp_credentials.insecureTransportCredentials
	// if you want no transport credentials (do not use this in production as nothing will get encrypted).
	CREDENTIALS softcorp_credentials.TransportCredentials
)

var (
	// NoAuthorization disables the authentication interface.
	NoAuthorization = &softcorp_authorize.NoAuthorization{}
	// Insecure sets transport gRPC credentials to insecure.NewCredentials()
	Insecure = &softcorp_credentials.InsecureTransportCredentials{}
	// STORAGE_PROVIDER is used to override the default softcorp_storage provider
)

var (
	EmptyApiKeyErr = errors.New("api key is empty")
)

type Client struct {
	UserClient user_client.UserClient
}

func NewClient(ctx context.Context) (*Client, error) {
	// check if encryption key is valid hex
	if ENCRYPTION_KEY != "" {
		if _, err := hex.DecodeString(ENCRYPTION_KEY); err != nil {
			return nil, err
		}
	}
	// get dial security softcorp_options
	credentialsGenerator, err := softcorp_credentials.New(CREDENTIALS)
	if err != nil {
		return nil, err
	}
	credentials, err := credentialsGenerator.GetTransportCredentials(API_URL)
	if err != nil {
		return nil, err
	}
	dialOptions := grpc.WithTransportCredentials(credentials)
	if API_KEY == "" {
		return nil, EmptyApiKeyErr
	}
	// create authorization client
	auth, err := softcorp_authorize.New(ctx, API_URL, API_KEY, AUTHORIZE, dialOptions)
	if err != nil {
		return nil, err
	}
	// create user service client
	userClient, err := user_client.New(API_URL, auth, ENCRYPTION_KEY, dialOptions)
	if err != nil {
		return nil, err
	}
	return &Client{
		UserClient: userClient,
	}, nil
}
