package validation

import (
	"fmt"
	"log"
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// Регулярка для email
var nameRegex = regexp.MustCompile(`^[A-Za-zА-Яа-яЁё-]+$`)

func GetValidator() *validator.Validate {
	validate := validator.New()

	if err := validate.RegisterValidation("password", passwordValidator); err != nil {
		log.Println("cannot register password validator:", err)
	}

	if err := validate.RegisterValidation("name", nameValidator); err != nil {
		log.Println("cannot register name validator:", err)
	}

	return validate
}

func GetError(err error) string {
	errors := err.(validator.ValidationErrors)

	for _, e := range errors {
		switch e.Tag() {
		case "required":
			return "Все поля обязательны"
		case "email":
			return "Неверный формат email"
		case "password":
			return "Пароль должен включать хотя бы одну букву каждого регистра и цифру"
		case "name":
			return "Неверное имя/фамилия"
		case "min":
			return fmt.Sprintf("Поле %s должно содержать не менее %s символов", e.Field(), e.Param())
		case "max":
			return fmt.Sprintf("Поле %s должно содержать не более %s символов", e.Field(), e.Param())
		}
	}
	return "Неизвестная ошибка"
}

// nameValidator Валидация имя
func nameValidator(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	return nameRegex.MatchString(name)
}

// passwordValidator Валидация пароля
func passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var hasUpper, hasLower, hasDigit bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit {
		return false
	}

	return true
}
