package responses

import "net/http"

type GenericResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

func SetResponse(status int, message string, data interface{}) GenericResponse {
	if data != nil {
		return GenericResponse{
			Status:  status,
			Message: message,
			Data:    map[string]interface{}{"result": data},
		}
	}
	return GenericResponse{
		Status:  status,
		Message: message,
		Data:    map[string]interface{}{"result": nil},
	}
}

const BadRequest = http.StatusBadRequest
const InternalServerError = http.StatusInternalServerError
const Created = http.StatusCreated
const OK = http.StatusOK
const NotFound = http.StatusNotFound
const UnAuthorized = http.StatusUnauthorized
const ConflictOrDuplicateOrAlreadyExists = http.StatusConflict
