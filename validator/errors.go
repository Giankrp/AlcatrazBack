package validator

import "github.com/go-playground/validator/v10"

func ValidationErrors(err error) map[string]string {
	res := map[string]string{}
	for _, e := range err.(validator.ValidationErrors) {
		res[e.Field()] = e.ActualTag()
	}
	return res
}
