package validations

import (
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"
	"mime/multipart"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

/*
 * order Validation - Error Message
 * -------------------------------
 * Kumpulan custom error message yang ditampilkan
 * ke response berdasarkan struct field dan validate tagnya
 */
var orderErrorMessages = map[string]string{
	"DestinationStartProvince|required"	: "Origin province is required",
	"DestinationStartCity|required"		: "Origin city is required",
	"DestinationStartDistrict|required"	: "Origin district is required",
	"DestinationStartAddress|required"	: "Origin address is required",
	"DestinationStartPostal|required"	: "Origin postal code is required",
	"DestinationStartLat|required"		: "Origin lattitude is required",
	"DestinationStartLong|required"		: "Origin longitude is required",
	"DestinationEndProvince|required"	: "Destination province is required",
	"DestinationEndCity|required"		: "Destination city is required",
	"DestinationEndDistrict|required"	: "Destination district is required",
	"DestinationEndAddress|required"	: "Destination address is required",
	"DestinationEndPostal|required"		: "Destination postal is required",
	"DestinationEndLat|required"		: "Destination lattitude is required",
	"DestinationEndLong|required"		: "Destination longitude is required",
	"TruckTypeID|required"				: "truck type id is required",
	"TotalVolume|required"				: "total volume is required",
	"TotalWeight|required"				: "total weight is required",
	"FixedPrice|required"				: "fixed price is required",
}

/*
 * Filesize Validation Rules
 * -------------------------------
 * Aturan input file order berdasarkan size
 * [field]: [size]
 */
var orderFileSizeRules = map[string]int{
	"avatar":        1024 * 1024, // 1MB
	"order_picture": 1024 * 1024,
}

var finishFileSizeRules = map[string]int{
	"arrived_picture": 1024 * 1024,
}

/*
 * Filesize Validation Rules
 * -------------------------------
 * Aturan format ekstensi file yang diperbolehkan
 * [field]: ext1|ext2|ext3...
 */
var orderFileExtRules = map[string]string{
	"avatar":        "jpg|jpeg|png|webp|bmp",
	"order_picture": "jpg|jpeg|png|webp|bmp",
}
var finishFileExtRules = map[string]string{
	"arrived_picture": "jpg|jpeg|png|webp|bmp",
}

/*
 * Order Validation - Validate Customer Create order request
 * -------------------------------
 * Validasi order saat registrasi berdasarkan validate tag
 * yang ada pada order request dan file rules diatas
 */
func ValidateCustomerCreateOrderRequest(validate *validator.Validate, customerCreateOrderReq entities.CustomerCreateOrderRequest, files map[string]*multipart.FileHeader) error {

	errors := []web.ValidationErrorItem{}

	validateOrderStruct(validate, customerCreateOrderReq, orderErrorMessages, &errors)
	validateOrderFiles(files, orderFileSizeRules, orderFileExtRules, &errors)

	if len(errors) > 0 {
		return web.ValidationError{
			Code:               400,
			ProductionMessage:  "Validation error",
			DevelopmentMessage: "Validation error",
			Errors:             errors,
		}
	}
	return nil
}

func ValidateAdminSetPriceOrderRequest(validate *validator.Validate, adminSetPriceOrderReq entities.AdminSetPriceOrderRequest) error {
	errors := []web.ValidationErrorItem{}
	validateOrderStruct(validate, adminSetPriceOrderReq, orderErrorMessages, &errors)
	if len(errors) > 0 {
		return web.ValidationError{
			Code:               400,
			ProductionMessage:  "Validation error",
			DevelopmentMessage: "Validation error",
			Errors:             errors,
		}
	}
	return nil
}

func ValidateUpdateOrderRequest(orderFiles map[string]*multipart.FileHeader) error {

	errors := []web.ValidationErrorItem{}

	validateOrderFiles(orderFiles, finishFileSizeRules, finishFileExtRules, &errors)
	if len(errors) > 0 {
		return web.ValidationError{
			Code:               400,
			ProductionMessage:  "Bad Request",
			DevelopmentMessage: "Validation error",
			Errors:             errors,
		}
	}
	return nil
}

func validateOrderStruct(validate *validator.Validate, request interface{}, errorMessages map[string]string, errors *[]web.ValidationErrorItem) {
	err := validate.Struct(request)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field, _ := reflect.TypeOf(request).FieldByName(err.Field())
			*errors = append(*errors, web.ValidationErrorItem{
				Field: field.Tag.Get("form"),
				Error: errorMessages[err.Field()+"|"+err.ActualTag()],
			})
		}
	}
}

func validateOrderFiles(files map[string]*multipart.FileHeader, fileSizeRules map[string]int, fileExtRules map[string]string, errors *[]web.ValidationErrorItem) {
	// File validation
	for field, file := range files {

		// required validation
		if file == nil {
			*errors = append(*errors, web.ValidationErrorItem{
				Field: field,
				Error: field + " must be filled",
			})
			continue
		}

		// Size validations
		if file.Size > int64(fileSizeRules[field]) {
			*errors = append(*errors, web.ValidationErrorItem{
				Field: field,
				Error: field + " size cannot more than " + strconv.Itoa(fileSizeRules[field]/1024) + " KB",
			})
		}

		// Extension validations
		fileExt := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
		allowedExt := strings.Split(fileExtRules[field], "|")
		fileExtAllowed := false
		for _, ext := range allowedExt {
			if fileExt == ext {
				fileExtAllowed = true
				break
			}
		}
		if !fileExtAllowed {
			*errors = append(*errors, web.ValidationErrorItem{
				Field: field,
				Error: field + " field must be type of " + strings.Join(allowedExt, ", "),
			})
		}
	}
}
