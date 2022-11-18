package genericcontrollers_gorm_gin

import (
	"fmt"
	genericcrud_repositories_gorm "github.com/danielcomboni/generic-crud/genericcrud_repositories"
	"github.com/danielcomboni/generic-crud/logging"
	"github.com/danielcomboni/generic-crud/models"
	"github.com/danielcomboni/generic-crud/responses"
	"github.com/danielcomboni/generic-crud/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

const BadRequest = http.StatusBadRequest
const InternalServerError = http.StatusInternalServerError
const Created = http.StatusCreated
const OK = http.StatusOK
const NotFound = http.StatusNotFound
const UnAuthorized = http.StatusUnauthorized

func Create[T any](model *T, c *gin.Context, fnServiceCreate func(t T) (T, responses.GenericResponse, error)) {

	logging.LogIncoming(model)

	//Validate the request body
	if err := c.BindJSON(&model); err != nil {
		logging.LogError(fmt.Sprintf("failed to bind incoming object: %v", err))
		c.JSON(BadRequest, responses.SetResponse(BadRequest, "error", err.Error()))
		return
	}

	//use the validator library to Validate required fields
	if validationErr := Validate.Struct(model); validationErr != nil {
		logging.LogError(fmt.Sprintf("failed to Validate incoming object: %v", validationErr))
		c.JSON(BadRequest, responses.SetResponse(BadRequest, "error", validationErr.Error()))
		return
	}

	// save (insert) to database
	created, res, err := fnServiceCreate(*model)
	if err != nil {
		logging.LogError(fmt.Sprintf("failed to save record: %v", err))
		c.JSON(InternalServerError, responses.SetResponse(InternalServerError, "error", err.Error()))
		return
	}

	if !utils.IsNullOrEmpty(res.Message) {
		logging.LogInfo(fmt.Sprintf("%v", res.Message))
		log.Println()
		c.JSON(res.Status, res)
		return
	}

	c.JSON(Created, responses.SetResponse(Created, "successful", created))

}

func CreateBatch[T any](model []T, c *gin.Context, fnServiceCreate func(t []T) ([]T, responses.GenericResponse, error)) {

	//Validate the request body
	if err := c.BindJSON(&model); err != nil {
		logging.LogError(fmt.Sprintf("failed to bind incoming object: %v", err))
		c.JSON(BadRequest, responses.SetResponse(BadRequest, "error", err.Error()))
		return
	}

	for _, t := range model {
		//use the validator library to Validate required fields
		if validationErr := Validate.Struct(t); validationErr != nil {
			logging.LogError(fmt.Sprintf("failed to Validate incoming object: %v", validationErr))
			c.JSON(BadRequest, responses.SetResponse(BadRequest, "error", validationErr.Error()))
			return
		}
	}

	// save (insert) to database
	created, res, err := fnServiceCreate(model)
	if err != nil {
		logging.LogError(fmt.Sprintf("failed to save record: %v", err))
		c.JSON(InternalServerError, responses.SetResponse(InternalServerError, "error", err.Error()))
		return
	}

	if !utils.IsNullOrEmpty(res.Message) {
		logging.LogIncoming(fmt.Sprintf("%v", res.Message))
		c.JSON(res.Status, res)
		return
	}

	c.JSON(Created, responses.SetResponse(Created, "successful", created))

}

func UpdateById[T any](model *T, c *gin.Context, fnServiceUpdate func(t T, id string) (T, error)) {
	id := c.Param("id")
	//Validate the request body
	if err := c.BindJSON(&model); err != nil {
		logging.LogError(fmt.Sprintf("failed to bind incoming object: %v", err))
		c.JSON(BadRequest, responses.SetResponse(BadRequest, "error", err.Error()))
		return
	}

	//use the validator library to Validate required fields
	if validationErr := Validate.Struct(model); validationErr != nil {
		logging.LogError(fmt.Sprintf("failed to Validate incoming object: %v", validationErr))
		c.JSON(BadRequest, responses.SetResponse(BadRequest, "error", validationErr.Error()))
		return
	}

	// save (insert) to database
	created, err := fnServiceUpdate(*model, id)
	if err != nil {
		logging.LogError(fmt.Sprintf("failed to update record: %v", err))
		c.JSON(InternalServerError, responses.SetResponse(InternalServerError, "error", err.Error()))
		return
	}

	c.JSON(Created, responses.SetResponse(Created, "successful", created))

}

func PatchById[T any](model *models.PatchByIdModel, c *gin.Context, fnServicePatch func(object models.PatchByIdModel) (T, error)) {
	//Validate the request body
	if err := c.BindJSON(&model); err != nil {
		logging.LogError(fmt.Sprintf("failed to bind incoming object: %v", err))
		c.JSON(BadRequest, responses.SetResponse(BadRequest, "error", err.Error()))
		return
	}

	//use the validator library to Validate required fields
	if validationErr := Validate.Struct(model); validationErr != nil {
		logging.LogError(fmt.Sprintf("failed to Validate incoming object: %v", validationErr))
		c.JSON(BadRequest, responses.SetResponse(BadRequest, "error", validationErr.Error()))
		return
	}

	// save (insert) to database
	created, err := fnServicePatch(*model)
	if err != nil {
		logging.LogError(fmt.Sprintf("failed to patch record: %v", err))
		c.JSON(InternalServerError, responses.SetResponse(InternalServerError, "error", err.Error()))
		return
	}

	c.JSON(Created, responses.SetResponse(Created, "successful", created))

}

func GetAll[T any](c *gin.Context, fnServiceGetAll func() ([]T, error)) {
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	sort := c.Request.URL.Query().Get("sort")
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))

	genericcrud_repositories_gorm.SetPagination(limit, page, sort)

	rows, err := fnServiceGetAll()
	if err != nil {

		c.JSON(InternalServerError, responses.SetResponse(InternalServerError, "error", err.Error()))
		return
	}

	c.JSON(OK, responses.SetResponse(OK, "successful", rows))
}

func GetAllByClientId[T any](c *gin.Context, fnServiceGetAll func(id string) ([]T, error)) {
	id := c.Param("clientId")

	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	sort := c.Request.URL.Query().Get("sort")
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))

	genericcrud_repositories_gorm.SetPagination(limit, page, sort)

	rows, err := fnServiceGetAll(id)
	if err != nil {
		c.JSON(InternalServerError, responses.SetResponse(InternalServerError, "error", err.Error()))
		return
	}
	c.JSON(OK, responses.SetResponse(OK, "successful", rows))
}

func GetAllByOtherPathParamsId[T any](c *gin.Context, fnServiceGetAll func(pathParams ...genericcrud_repositories_gorm.PathParams) ([]T, error), pathParams ...string) {
	//id := c.Param("clientId")
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	sort := c.Request.URL.Query().Get("sort")
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	genericcrud_repositories_gorm.SetPagination(limit, page, sort)
	var params []genericcrud_repositories_gorm.PathParams
	for _, param := range c.Params {
		params = append(params, genericcrud_repositories_gorm.PathParams{
			Key:   param.Key,
			Value: param.Value,
		})
	}
	rows, err := fnServiceGetAll(params...)
	if err != nil {
		c.JSON(InternalServerError, responses.SetResponse(InternalServerError, "error", err.Error()))
		return
	}
	c.JSON(OK, responses.SetResponse(OK, "successful", rows))
}

func GetOneById[T any](c *gin.Context, fnServiceGetOneById func(id string) (T, error)) {
	id := c.Param("id")
	row, err := fnServiceGetOneById(id)
	if err != nil {
		c.JSON(InternalServerError, responses.SetResponse(InternalServerError, "error", err.Error()))
		return
	}
	if utils.IsNullOrEmpty(utils.SafeGetFromInterface(row, "$.id")) {
		c.JSON(OK, responses.SetResponse(NotFound, "not found", nil))
		return
	}
	c.JSON(OK, responses.SetResponse(OK, "successful", row))
}

func DeleteSoftlyById[T any](c *gin.Context, fnServiceDeleteSoftlyById func(id string) (int64, error)) {
	id := c.Param("id")
	rowsAffected, err := fnServiceDeleteSoftlyById(id)
	if err != nil {
		c.JSON(InternalServerError, responses.SetResponse(InternalServerError, "error", err.Error()))
		return
	}
	c.JSON(OK, responses.SetResponse(OK, "successful", rowsAffected))
}

func DeletePermanentlyById[T any](c *gin.Context, fnServiceDeletePermanentlyById func(id string) (int64, error)) {
	id := c.Param("id")
	rowsAffected, err := fnServiceDeletePermanentlyById(id)
	if err != nil {
		c.JSON(InternalServerError, responses.SetResponse(InternalServerError, "error", err.Error()))
		return
	}
	c.JSON(OK, responses.SetResponse(OK, "successful", rowsAffected))
}
