package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"shop/pkg/models"
)

type Repo struct {
	Connector *models.Connector
}

func (r *Repo) Verify(token string) *models.Error {
	accessToken := models.AccessToken{TokenStr: token}
	jsonToken, jsonTokenError := json.Marshal(accessToken)
	if jsonTokenError != nil {
		return &models.Error{
			ErrorString: jsonTokenError.Error(),
			ErrorCode: http.StatusBadRequest,
		}
	}

	_, authError := r.Connector.Post("validate", &jsonToken)
	if authError != nil {
		fmt.Println(authError.ErrorString)
		return authError
	}
	return nil
}
