package logic

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"shop/pkg/models"
)

func GetProductJSON(responseBody io.ReadCloser)(*[]byte, *models.Error) {
	defer responseBody.Close()
	body, bodyParseError := ioutil.ReadAll(responseBody)
	if bodyParseError != nil {
		return nil, &models.Error{
			ErrorString: bodyParseError.Error(),
			ErrorCode: http.StatusInternalServerError}
	}
	model := &models.Product{}
	unmarshalError := json.Unmarshal(body, model)
	if unmarshalError != nil {
		return nil, &models.Error{
			ErrorString: unmarshalError.Error(),
			ErrorCode: http.StatusInternalServerError}
	}
	jsonData, marshalError := json.Marshal(
		models.Product{
			Name: model.Name,
			Id: model.Id,
			Category: model.Category})
	if marshalError != nil {
		return nil, &models.Error{
			ErrorString: marshalError.Error(),
			ErrorCode: http.StatusInternalServerError}
	}
	return &jsonData, nil
}
