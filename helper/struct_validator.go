package helper

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
	"strings"
)

func ValidateInputs(dataSet interface{})(bool, map[string][]string){
	var validate *validator.Validate
	validate = validator.New()
	err := validate.Struct(dataSet)

	if err != nil{
		//Validation syntax is invalid
		if err, ok := err.(*validator.InvalidValidationError); ok {
			panic(err)
		}

		//Validation errors occurred
		errors := make(map[string][]string)

		reflected := reflect.ValueOf(dataSet)

		for _, err := range err.(validator.ValidationErrors) {
			field, _ := reflected.Type().FieldByName(err.StructField())
			var name string

			if name = field.Tag.Get("json"); name == "" {
				name = strings.ToLower(err.StructField())
			}

			switch err.Tag() {
			case "required":
				errors[name] = append(errors[name], "The "+name+" is required")
				break
			case "email":
				errors[name] = append(errors[name], "The "+name+" should be a valid email")
				break
			case "eqfield":
				errors[name] = append(errors[name], "The "+name+" should be equal to the "+err.Param())
				break
			default:
				errors[name] = append(errors[name], "The "+name+" is invalid")
				break
			}
		}
		return false,errors
	}
	return true, nil
}
