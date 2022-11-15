package genericcrud_repositories_gorm

import (
	"errors"
	"fmt"
	"github.com/danielcomboni/generic-crud/utils"
	"log"
	"reflect"

	"github.com/gobeam/stringy"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type Pagination struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Sort  string `json:"sort"`
}

var tablePagination *Pagination

func SetPagination(limit, page int, sort string) {
	tablePagination = &Pagination{
		Limit: limit,
		Page:  page,
		Sort:  sort,
	}
}

func paginationParams() (offset int, limit int) {
	page := tablePagination.Page
	if page == 0 {
		page = 1
	}

	pageSize := tablePagination.Limit

	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	offset = (page - 1) * pageSize
	limit = pageSize
	return offset, pageSize
}

func Create[T any](model *T, databaseInstance *gorm.DB) (T, error) {
	log.Print(fmt.Sprintf("\n\ncreating a new record: %v", reflect.TypeOf(*new(T)).Name()))
	var t T
	result := databaseInstance.Create(&model).Scan(&t)
	_, err := result.DB()
	if err != nil {
		log.Panic("failed to create")
		return *model, err
	}

	if utils.IsNullOrEmpty(utils.SafeGetFromInterface(&model, "$.id")) {
		log.Println(fmt.Sprintf("not saved:"))
		return *model, err
	}

	if result.RowsAffected > 0 {
		log.Println(fmt.Sprintf("saved to database: id: %v", utils.SafeGetFromInterface(t, "$.id")))
	}
	//t := utils.SafeGetFromInterfaceGenericAndDeserialize[T](&model, "$")
	return t, nil
}

func CreateBatch[T any](models []T, databaseInstance *gorm.DB) ([]T, error) {
	log.Println(fmt.Sprintf("\n\ncreating a new record: %v", reflect.TypeOf(*new(T)).Name()))
	var t []T
	result := databaseInstance.Create(&models).Scan(&t)
	_, err := result.DB()
	if err != nil {
		log.Println("failed to create in batch")
		return t, err
	}

	if len(t) == 0 {
		log.Println(fmt.Sprintf("not saved:"))
		return t, err
	}

	if result.RowsAffected > 0 {
		log.Println(fmt.Sprintf("saved to database: id: %v", utils.SafeGetFromInterface(t, "$.id")))
	}
	//t := utils.SafeGetFromInterfaceGenericAndDeserialize[T](&model, "$")
	return t, nil
}

func GetAll[T any](databaseInstance *gorm.DB) ([]T, error) {
	log.Println(fmt.Sprintf("\n\nretreiving collection: %v", reflect.TypeOf(*new(T)).Name()))
	var all []T
	offset, limit := paginationParams()
	log.Println(fmt.Sprintf("offset: %v, limit: %v", offset, limit))
	result := databaseInstance.Offset(offset).Limit(limit).Find(&all)
	_, err := result.DB()
	if err != nil {
		log.Println(fmt.Sprintf("failed to retrieve: %v %v", reflect.TypeOf(*new(T)).Name(), err))
		return all, err
	}
	return all, nil
}

func preloadsHandler(databaseInstance *gorm.DB, preloads ...string) *gorm.DB {

	instance := databaseInstance

	for _, s := range preloads {
		instance = instance.Preload(s)
	}

	return instance
}

// todo ... fields to be omitted
func omitsHandler(databaseInstance *gorm.DB, omits ...string) *gorm.DB {

	instance := databaseInstance

	for _, s := range omits {
		instance = instance.Omit(s)
	}

	return instance
}

func GetAllByFields[T any](databaseInstance *gorm.DB, queryMap map[string]interface{}, preloads ...string) ([]T, error) {
	log.Println(fmt.Sprintf("retreiving collection: %v\n", reflect.TypeOf(*new(T)).Name()))
	var all []T
	offset, limit := paginationParams()
	log.Println(fmt.Sprintf("offset: %v, limit: %v", offset, limit))

	var instance *gorm.DB

	if len(preloads) == 0 {
		instance = databaseInstance.Offset(offset).Limit(limit).Where(queryMap).Find(&all)
	} else {
		instance = preloadsHandler(databaseInstance, preloads...).Offset(offset).Limit(limit).Where(queryMap).Find(&all)
	}

	result := instance
	_, err := result.DB()
	if err != nil {
		log.Println(fmt.Sprintf("failed to retrieve: %v %v", reflect.TypeOf(*new(T)).Name(), err))
		return all, err
	}
	return all, nil
}

func GetOneById[T any](databaseInstance *gorm.DB, id string, preloads ...string) (T, error) {
	log.Println(fmt.Sprintf("\n\nretreiving single row of: %v by id: %v", reflect.TypeOf(*new(T)).Name(), id))
	var row T

	var instance *gorm.DB
	if len(preloads) == 0 {
		instance = databaseInstance.Where("id=?", id).Find(&row)
	} else {
		instance = preloadsHandler(databaseInstance, preloads...).Where("id=?", id).Find(&row)
	}

	result := instance

	_, err := result.DB()
	if err != nil {
		log.Println(fmt.Sprintf("failed to retrieve: %v %v", reflect.TypeOf(*new(T)).Name(), err))
		return row, err
	}
	return row, nil
}

func GetOneSoftDeletedById[T any](databaseInstance *gorm.DB, id string, preloads ...string) (T, error) {
	log.Println(fmt.Sprintf("\n\nretreiving single row of: %v by id: %v", reflect.TypeOf(*new(T)).Name(), id))
	var row T

	var instance *gorm.DB
	if len(preloads) == 0 {
		instance = databaseInstance.Unscoped().Where("id=?", id).Find(&row)
	} else {
		instance = preloadsHandler(databaseInstance, preloads...).Unscoped().Where("id=?", id).Find(&row)
	}

	result := instance

	_, err := result.DB()
	if err != nil {
		log.Println(fmt.Sprintf("failed to retrieve: %v %v", reflect.TypeOf(*new(T)).Name(), err))
		return row, err
	}
	return row, nil
}

func GetOneByModelPropertiesCheckIdPresence[T any](databaseInstance *gorm.DB, queryMap map[string]interface{}) (T, error) {
	log.Println(fmt.Sprintf("\n\nretreiving single row of: %v by values: %#v", reflect.TypeOf(*new(T)).Name(), queryMap))
	var row T
	result := databaseInstance.Where(queryMap).First(&row)
	_, err := result.DB()
	if err != nil {
		log.Println(fmt.Sprintf("failed to retrieve: %v %v", reflect.TypeOf(*new(T)).Name(), err))
		return row, err
	}

	r := reflect.ValueOf(row)
	f := reflect.Indirect(r).FieldByName("Id")
	if f.String() == "" {
		msg := "record not found"
		log.Println(msg)
		return row, errors.New(msg)
	}
	return row, nil
}

func PatchById[T any](databaseInstance *gorm.DB, id, columnName string, value interface{}) (T, error) {
	log.Println(fmt.Sprintf("\n\npatch column: %v row of: %v by id: %v", columnName, reflect.TypeOf(*new(T)).Name(), id))
	one, err := GetOneById[T](databaseInstance, id)
	var t2 T
	if err != nil {
		return t2, err
	}

	if err != nil {
		log.Println(fmt.Sprintf("failed to map structure: %v", err))
		return t2, err
	}
	//result := database.Instance.Where("id=?", id).Update(stringy.New(columnName).SnakeCase("?", "").ToLower(), value).Scan(&one)
	result := databaseInstance.Model(&one).Where("id=?", id).Update(stringy.New(columnName).SnakeCase("?", "").ToLower(), value).Scan(&one)
	rowsAffected := result.RowsAffected
	log.Println(fmt.Sprintf("rows affected: %v", utils.ConvertInt64ToStr(rowsAffected)))

	if result.Error != nil {
		log.Println(fmt.Sprintf("failed to patch env: %v", result.Error))
		return t2, result.Error
	}

	if result.RowsAffected == 0 {
		log.Println(fmt.Sprintf("not updated: affected rows: %v", result.RowsAffected))
		return t2, errors.New("not patched")
	}

	return one, nil
}

func UpdateById[T any](databaseInstance *gorm.DB, t T, id string) (T, error) {

	log.Println(fmt.Sprintf("\n\nupdating row of: %v by id: %v", reflect.TypeOf(*new(T)).Name(), id))
	one, err := GetOneById[T](databaseInstance, id)
	var t2 T
	if err != nil {
		return t2, err
	}

	err = mapstructure.Decode(t, &one)

	if err != nil {
		log.Println(fmt.Sprintf("failed to map structure: %v", err))
		return t2, err
	}

	// set the createdAt date and updatedAt

	result := databaseInstance.Where("id=?", id).Updates(&one).Scan(&one)
	rowsAffected := result.RowsAffected
	log.Println(fmt.Sprintf("rows affected: %v", utils.ConvertInt64ToStr(rowsAffected)))

	if result.Error != nil {
		log.Println(fmt.Sprintf("failed to update env: %v", result.Error))
		return t2, result.Error
	}

	if result.RowsAffected == 0 {
		log.Println(fmt.Sprintf("not updated: affected rows: %v", result.RowsAffected))
		return t2, errors.New("not updated")
	}

	return one, nil
}

func DeleteHardById[T any](databaseInstance *gorm.DB, id string) (int64, error) {
	log.Println(fmt.Sprintf("\n\nhard deleting a row of: %v by id: %v", reflect.TypeOf(*new(T)).Name(), id))

	one, err := GetOneById[T](databaseInstance, id)
	var t2 T
	if err != nil {
		return 0, err
	}

	err = mapstructure.Decode(one, &t2)

	if err != nil {
		log.Println(fmt.Sprintf("failed to map structure: %v", err))
		return 0, err
	}

	r := databaseInstance.Delete(&one).Where("id=?", id)

	if r.Error != nil {
		log.Println(fmt.Sprintf("failed to delete row by id: %v", id))
		log.Println(fmt.Sprintf("err: %v", r.Error))
		return 0, r.Error
	}

	if r.RowsAffected <= 0 {
		log.Println(fmt.Sprintf("failed to delete row by id: %v", id))
		log.Println(fmt.Sprintf("number of rows deleted: %v", r.RowsAffected))
		return 0, r.Error
	}

	return r.RowsAffected, nil
}

func DeleteSoftById[T any](databaseInstance *gorm.DB, id string) (int64, error) {
	log.Println(fmt.Sprintf("\n\nsoft deleting a row of: %v by id: %v", reflect.TypeOf(*new(T)).Name(), id))
	one, err := GetOneById[T](databaseInstance, id)
	if err != nil {
		log.Println(fmt.Sprintf("failed to get record by id:% %v", id, err))
		return 0, err
	}

	if utils.IsNullOrEmpty(utils.SafeGetFromInterface(one, "$.id")) {
		msg := fmt.Sprintf("no record found with id: %v", id)
		log.Println(msg)
		return 0, errors.New(msg)
	}

	var t2 T

	err = mapstructure.Decode(one, &t2)

	if err != nil {
		log.Println(fmt.Sprintf("failed to map structure: %v", err))
		return 0, err
	}

	r := databaseInstance.Delete(&one).Where("id=?", id)

	if r.Error != nil {
		log.Println(fmt.Sprintf("failed to delete row by id: %v", id))
		log.Println(fmt.Sprintf("err: %v", r.Error))
		return 0, r.Error
	}

	if r.RowsAffected <= 0 {
		log.Println(fmt.Sprintf("failed to delete row by id: %v", id))
		log.Println(fmt.Sprintf("number of rows deleted: %v", r.RowsAffected))
		return 0, r.Error
	}

	return r.RowsAffected, nil
}

func DeletePermanentById[T any](databaseInstance *gorm.DB, id string) (int64, error) {
	log.Println(fmt.Sprintf("\n\nsoft deleting a row of: %v by id: %v", reflect.TypeOf(*new(T)).Name(), id))
	one, err := GetOneSoftDeletedById[T](databaseInstance, id)
	if err != nil {
		log.Println(fmt.Sprintf("failed to get record by id:% %v", id, err))
		return 0, err
	}

	if utils.IsNullOrEmpty(utils.SafeGetFromInterface(one, "$.id")) {
		msg := fmt.Sprintf("no record found with id: %v", id)
		log.Println(msg)
		return 0, errors.New(msg)
	}

	var t2 T

	err = mapstructure.Decode(one, &t2)

	if err != nil {
		log.Println(fmt.Sprintf("failed to map structure: %v", err))
		return 0, err
	}

	r := databaseInstance.Unscoped().Delete(&one).Where("id=?", id)

	if r.Error != nil {
		log.Println(fmt.Sprintf("failed to delete row by id: %v", id))
		log.Println(fmt.Sprintf("err: %v", r.Error))
		return 0, r.Error
	}

	if r.RowsAffected <= 0 {
		log.Println(fmt.Sprintf("failed to delete row by id: %v", id))
		log.Println(fmt.Sprintf("number of rows deleted: %v", r.RowsAffected))
		return 0, r.Error
	}

	return r.RowsAffected, nil
}
