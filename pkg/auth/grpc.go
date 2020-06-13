package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"log"
	"net"
	"shop/pkg/constants"
)

type ValidateServer struct { }

func (s ValidateServer) ValidateToken(_ context.Context, token *Token) (*Verification, error) {
	res := &Verification{Message: "token is invalid"}

	// check if access token exists
	accessTokenStr := token.Token
	if len(accessTokenStr) == 0 {
		log.Println("Authorization token is empty string")
		return res, errors.New("authorization token is empty string")
	}

	// parse access token into struct
	claim := jwt.MapClaims{}
	accessToken, accessTokenParseError := jwt.ParseWithClaims(
		accessTokenStr,
		claim,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(constants.SigningToken), nil
		},
	)
	if accessTokenParseError != nil {
		log.Printf("Failed to parse access token: %s\n", accessTokenParseError.Error())
		if accessTokenParseError.Error() != "Token is expired" {
			return res, errors.New(constants.InternalError)
		} else {
			return res, errors.New(constants.TokenIsExpired)
		}
	}

	// create reply with message and user email and role
	res.Message = constants.ValidAccessToken
	res.Role = claim["role"].(string)
	res.Email = claim["email"].(string)

	// check if access token has not been expired
	if accessToken.Valid {
		log.Printf("Token %s is valid \n", accessTokenStr)
	} else {
		log.Printf("Token %s is invalid \n", accessTokenStr)
		res.Message = constants.Unauthorized
	}
	return res, nil
}

func StartGrpcServer(address string) {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			func(
				ctx context.Context,
				req interface{},
				info *grpc.UnaryServerInfo,
				handler grpc.UnaryHandler,
			) (interface{}, error) {
				return handler(ctx, req)
			},
		),
	)
	RegisterAuthServer(server, ValidateServer{})

	lis, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("Error while net.Listen address: %s\n", address)
	}

	go func() {
		err := server.Serve(lis)
		if err != nil {
			fmt.Println(err)
		}
		log.Println("1111111111111111111111")
		//server.GracefulStop()
	}()
}
