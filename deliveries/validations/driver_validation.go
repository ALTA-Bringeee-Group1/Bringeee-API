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
 * driver Validation - Error Message
 * -------------------------------
 * Kumpulan custom error message yang ditampilkan
 * ke response berdasarkan struct field dan validate tagnya
 */
var driverErrorMessages = map[string]string{
	"Name|required":              "Name field must be filled",
	"Email|required":             "Email field must be filled",
	"Email|email":                "Email field is not an email",
	"Password|required":          "Password field must be filled",
	"Gender|required":            "Gender field must be filled",
	"DOB|required":               "Date of birth field must be filled",
	"Address|required":           "Address must be filled",
	"Age|required":               "Age must be filled",
	"NIK|required":               "NIK must be filled",
	"PhoneNumber|required":       "Phone number must be filled",
	"TruckTypeID|required":       "Truck type must be filled",
	"VehicleIdentifier|required": "Vehicle identifier must be filled",
}

/*
 * Filesize Validation Rules
 * -------------------------------
 * Aturan input file driver berdasarkan size
 * [field]: [size]
 */
var driverFileSizeRules = map[string]int{
	"avatar":              1024 * 1024, // 1MB
	"ktp_file":            1024 * 1024, // 1MB
	"stnk_file":           1024 * 1024, // 1MB
	"driver_license_file": 1024 * 1024, // 1MB
	"vehicle_picture":     1024 * 1024, // 1MB
}

/*
 * Filesize Validation Rules
 * -------------------------------
 * Aturan format ekstensi file yang diperbolehkan
 * [field]: ext1|ext2|ext3...
 */
var driverFileExtRules = map[string]string{
	"avatar":              "jpg|jpeg|png|webp|bmp",
	"ktp_file":            "jpg|jpeg|png|webp|bmp",
	"stnk_file":           "jpg|jpeg|png|webp|bmp",
	"driver_license_file": "jpg|jpeg|png|webp|bmp",
	"vehicle_picture":     "jpg|jpeg|png|webp|bmp",
}

/*
 * driver Validation - Validate Create driver Request
 * -------------------------------
 * Validasi driver saat registrasi berdasarkan validate tag
 * yang ada pada driver request dan file rules diatas
 */
func ValidateCreateDriverRequest(validate *validator.Validate, driverReq entities.CreateDriverRequest, driverFiles map[string]*multipart.FileHeader) error {

	errors := []web.ValidationErrorItem{}

	validateDriverStruct(validate, driverReq, &errors)
	validateDriverFiles(driverFiles, &errors)

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

/*
 * driver Validation - Validate Update driver Request
 * -------------------------------
 * Validasi driver saat edit profile berdasarkan
 * file rules diatas
 */
func ValidateUpdateDriverRequest(driverFiles map[string]*multipart.FileHeader) error {

	errors := []web.ValidationErrorItem{}

	validateDriverFiles(driverFiles, &errors)
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

func validateDriverStruct(validate *validator.Validate, driverReq entities.CreateDriverRequest, errors *[]web.ValidationErrorItem) {
	err := validate.Struct(driverReq)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field, _ := reflect.TypeOf(driverReq).FieldByName(err.Field())
			*errors = append(*errors, web.ValidationErrorItem{
				Field: field.Tag.Get("form"),
				Error: driverErrorMessages[err.Field()+"|"+err.ActualTag()],
			})
		}
	}
}

func validateDriverFiles(driverFiles map[string]*multipart.FileHeader, errors *[]web.ValidationErrorItem) {
	// File validation
	for field, file := range driverFiles {
		if file != nil {

			// Size validations
			if file.Size > int64(driverFileSizeRules[field]) {
				*errors = append(*errors, web.ValidationErrorItem{
					Field: field,
					Error: field + " size cannot more than " + strconv.Itoa(driverFileSizeRules[field]/1024) + " KB",
				})
			}

			// Extension validations
			fileExt := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
			allowedExt := strings.Split(driverFileExtRules[field], "|")
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
		} else {
			*errors = append(*errors, web.ValidationErrorItem{
				Field: field,
				Error: field + " must be filled",
			})
		}
	}
}
