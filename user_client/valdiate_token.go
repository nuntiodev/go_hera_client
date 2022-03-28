package user_client

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/softcorp-io/block-proto/go_block"
)

func (s *defaultSocialServiceClient) ValidateToken(ctx context.Context, jwtToken string) (*go_block.User, error) {
	if jwtToken == "" {
		return nil, tokenIsEmptyErr
	}
	jwtPublicKey, err := s.getPublicKey()
	if err != nil {
		return nil, err
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(jwtPublicKey)
	if err != nil {
		return nil, err
	}
	token, err := jwt.ParseWithClaims(
		jwtToken,
		&go_block.CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return key, nil
		},
	)
	if err != nil {
		return nil, err
	}
	if token.Valid == false {
		return nil, errors.New("token is not valid")
	}
	claims, ok := token.Claims.(*go_block.CustomClaims)
	if !ok {
		return nil, errors.New("couldn't parse claims")
	}
	if err != nil {
		return nil, err
	}
	return &go_block.User{
		Id: claims.UserId,
	}, nil
}