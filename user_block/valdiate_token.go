package user_block

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/nuntiodev/block-proto/go_block"
)

// ValidateToken locally validates the JWT and returns a user with the corresponding user id
func (s *defaultSocialServiceClient) ValidateToken(ctx context.Context, jwtToken string, validateServerSide bool) (*go_block.User, error) {
	if validateServerSide {
		resp, err := s.userClient.ValidateToken(ctx, &go_block.UserRequest{
			TokenPointer: jwtToken,
		})
		if err != nil {
			return nil, err
		}
		return resp.User, nil
	} else {
		if jwtToken == "" {
			return nil, tokenIsEmptyErr
		}
		jwtPublicKey, err := s.getPublicKey()
		if err != nil {
			return nil, err
		}
		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(jwtPublicKey))
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
}
