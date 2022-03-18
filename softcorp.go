package softcorp

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/softcorp-io/cloud-sdk/authorize"
	"github.com/softcorp-io/cloud-sdk/credentials_generator"
	"github.com/softcorp-io/cloud-sdk/user_client"
	"google.golang.org/grpc"
)

var (
	// ENCRYPTION_KEY  is used to encrypt clients data under the given key
	ENCRYPTION_KEY = ""
	// API_KEY  is used to connect your application to Softcorp Cloud
	API_KEY = ""
	// STORAGE_PROVIDER is used to override the default storage provider
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
	// get dial security options
	credentialsGenerator, err := credentials_generator.New()
	if err != nil {
		return nil, err
	}
	credentials, err := credentialsGenerator.GetTransportCredentials()
	if err != nil {
		return nil, err
	}
	dialOptions := grpc.WithTransportCredentials(credentials)
	if API_KEY == "" {
		return nil, EmptyApiKeyErr
	}
	// create authorization client
	auth, err := authorize.New(ctx, API_KEY, dialOptions)
	if err != nil {
		return nil, err
	}
	// create user service client
	userClient, err := user_client.New(auth, ENCRYPTION_KEY, dialOptions)
	if err != nil {
		return nil, err
	}
	return &Client{
		UserClient: userClient,
	}, nil
}
