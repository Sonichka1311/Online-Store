package auth

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"shop/pkg/constants"
	"shop/pkg/models"
	"strconv"
)

func Verify(token string) (*Verification, *models.Error) {
	grpcConn, connectError := grpc.Dial(
		constants.AuthHost + ":" + strconv.Itoa(constants.ValidatePort),
		grpc.WithInsecure(),
	)
	if connectError != nil {
		err, _ := models.NewError(errors.New(constants.InternalError), http.StatusInternalServerError)
		return nil, err
	}
	defer grpcConn.Close()

	client := NewAuthClient(grpcConn)
	verification, validateError := client.ValidateToken(context.Background(), &Token{Token: token})
	if validateError != nil {
		log.Printf("Error %s\n", validateError.Error())
		var err *models.Error
		if validateError.Error() == constants.InternalError {
			err, _ = models.NewError(errors.New(constants.InternalError), http.StatusInternalServerError)
		} else {
			err, _ = models.NewError(errors.New(constants.Unauthorized), http.StatusUnauthorized)
		}
		return nil, err
	} else if verification.Role != "admin" {
		err, _ := models.NewError(errors.New(constants.NoRight), http.StatusForbidden)
		return nil, err
	}

	return verification, nil
}
