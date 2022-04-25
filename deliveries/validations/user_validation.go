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
 * User Validation - Error Message
 * -------------------------------
 * Kumpulan custom error message yang ditampilkan
 * ke response berdasarkan struct field dan validate tagnya
 */
var userErrorMessages = map[string]string{
	"Name|required":        "Name field must be filled",
	"Email|required":       "Email field must be filled",
	"Email|email":          "Email field is not an email",
	"Password|required":    "Password field must be filled",
	"Gender|required":      "Gender field must be filled",
	"DOB|required":         "Date of birth field must be filled",
	"Address|required":     "Address must be filled",
	"PhoneNumber|required": "Phone number must be filled",
}

/*
 * Filesize Validation Rules
 * -------------------------------
 * Aturan input file user berdasarkan size
 * [field]: [size]
 */
var userFileSizeRules = map[string]int{
	"avatar": 1024 * 1024, // 1MB
}

/*
 * Filesize Validation Rules
 * -------------------------------
 * Aturan format ekstensi file yang diperbolehkan
 * [field]: ext1|ext2|ext3...
 */
var userFileExtRules = map[string]string{
	"avatar": "jpg|jpeg|png|webp|bmp",
}

/*
 * User Validation - Validate Create User Request
 * -------------------------------
 * Validasi user saat registrasi berdasarkan validate tag
 * yang ada pada user request dan file rules diatas
 */
func ValidateCreateCustomerRequest(validate *validator.Validate, userReq entities.CreateCustomerRequest, userFiles map[string]*multipart.FileHeader) error {

	errors := []web.ValidationErrorItem{}

	validateCustomerStruct(validate, userReq, &errors)
	validateCustomerFiles(userFiles, &errors)

	if len(errors) > 0 {
		return web.ValidationError{
			Code:    400,
			Message: "Validation error",
			Errors:  errors,
		}
	}
	return nil
}

/*
 * User Validation - Validate Update User Request
 * -------------------------------
 * Validasi user saat edit profile berdasarkan
 * file rules diatas
 */
func ValidateUpdateCustomerRequest(userFiles map[string]*multipart.FileHeader) error {

	errors := []web.ValidationErrorItem{}

	validateCustomerFiles(userFiles, &errors)
	if len(errors) > 0 {
		return web.ValidationError{
			Code:    400,
			Message: "Validation error",
			Errors:  errors,
		}
	}
	return nil
}

func validateCustomerStruct(validate *validator.Validate, userReq entities.CreateCustomerRequest, errors *[]web.ValidationErrorItem) {
	err := validate.Struct(userReq)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field, _ := reflect.TypeOf(userReq).FieldByName(err.Field())
			*errors = append(*errors, web.ValidationErrorItem{
				Field: field.Tag.Get("form"),
				Error: userErrorMessages[err.Field()+"|"+err.ActualTag()],
			})
		}
	}
}

func validateCustomerFiles(userFiles map[string]*multipart.FileHeader, errors *[]web.ValidationErrorItem) {
	// File validation
	for field, file := range userFiles {

		// Size validations
		if file.Size > int64(userFileSizeRules[field]) {
			*errors = append(*errors, web.ValidationErrorItem{
				Field: field,
				Error: field + " size cannot more than " + strconv.Itoa(userFileSizeRules[field]/1024) + " KB",
			})
		}

		// Extension validations
		fileExt := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
		allowedExt := strings.Split(userFileExtRules[field], "|")
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
