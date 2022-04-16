package nuntio

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/nuntiodev/go-blocks/nuntio_authorize"
	"github.com/nuntiodev/go-blocks/nuntio_credentials"
	"github.com/nuntiodev/go-blocks/user_client"
	"google.golang.org/grpc"
)

var (
	// ENCRYPTION_KEY  is used to encrypt clients data under the given key
	ENCRYPTION_KEY = ""
	// API_KEY  is used to connect your application to Nuntio Cloud
	API_KEY = ""
	// API_URL  is the URL the SDK will try to connect to
	API_URL = "api.nuntio.io:443"
	// AUTHORIZE is used to override the default nuntio_authorize interface which is used to validate tokens
	// if you don't want any authorization, set it to nuntio_authorize.AUTHORIZE = nuntio_authorize.NoAuthorization.
	AUTHORIZE nuntio_authorize.Authorize
	// CREDENTIALS defines what security is passed to nuntio_credentials.Dial and (can be overwritten)
	// you can provide your own, or use nuntio_credentials.TRANSPORT_CREDENTIALS = nuntio_credentials.insecureTransportCredentials
	// if you want no transport credentials (do not use this in production as nothing will get encrypted).
	CREDENTIALS nuntio_credentials.TransportCredentials
	// NAMESPACE defines what namespace you want to use with Nuntio Blocks (only edit this if you know what you are doing)
	NAMESPACE = ""
)

var (
	// NoAuthorization disables the authentication interface.
	NoAuthorization = &nuntio_authorize.NoAuthorization{}
	// Insecure sets transport gRPC credentials to insecure.NewCredentials()
	Insecure = &nuntio_credentials.InsecureTransportCredentials{}
	// STORAGE_PROVIDER is used to override the default nuntio_storage provider
)

var (
	EmptyApiKeyErr = errors.New("api key is empty")
)

type Client struct {
	UserClient user_client.UserClient
}

func NewClient(ctx context.Context) (*Client, error) {
	// check if encryption and api key is valid hex
	if ENCRYPTION_KEY != "" {
		if _, err := hex.DecodeString(ENCRYPTION_KEY); err != nil {
			return nil, err
		}
	} else {
		fmt.Println("Your encryption key is empty. If you want to secure your users data, please provide a 256 byte AES encryption key in Hex.")
	}
	if API_KEY == "" {
		fmt.Println(EmptyApiKeyErr.Error())
	}
	// get dial security nuntio_options
	credentialsGenerator, err := nuntio_credentials.New(CREDENTIALS, API_URL)
	if err != nil {
		return nil, err
	}
	credentials, err := credentialsGenerator.GetTransportCredentials()
	if err != nil {
		return nil, err
	}
	dialOptions := grpc.WithTransportCredentials(credentials)
	// create authorization client
	auth, err := nuntio_authorize.New(ctx, API_URL, API_KEY, AUTHORIZE, dialOptions)
	if err != nil {
		return nil, err
	}
	// create user service client
	userClient, err := user_client.New(API_URL, auth, ENCRYPTION_KEY, NAMESPACE, dialOptions)
	if err != nil {
		return nil, err
	}
	return &Client{
		UserClient: userClient,
	}, nil
}
