package logic

import (
	"../../common/constants"
	"../../common/logic"
	"../../common/models"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var AuthConnector = commonModels.Connector{
	Router: commonModels.Router{Host: commonLogic.GetUrl(constants.Protocol, constants.AuthHost, constants.AuthPort)},
	Mutex:  sync.Mutex{},
}

func IsAuthorised(token string) *commonModels.Error {
	accessToken := commonModels.AccessToken{TokenStr: token}
	jsonToken, jsonTokenError := json.Marshal(accessToken)
	if jsonTokenError != nil {
		return &commonModels.Error{
			ErrorString: jsonTokenError.Error(),
			ErrorCode: http.StatusBadRequest,
		}
	}

	_, authError := AuthConnector.Post("validate", &jsonToken)
	if authError != nil {
		fmt.Println(authError.ErrorString)
		return authError
	}
	return nil
}
