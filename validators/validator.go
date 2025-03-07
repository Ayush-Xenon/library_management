package validators

import (
	"library_management/models"
	"regexp"
	"strings"
)

var response models.ValidateOutput

func ValidateEmail(email string) models.ValidateOutput {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	response.Message = "Invalid Email"
	response.Result = re.MatchString(email)
	return response
}

func ValidatePassword(password string) models.ValidateOutput {
	var hasMinLen = regexp.MustCompile(`.{8,}`)
	var hasNumber = regexp.MustCompile(`[0-9]`)
	var hasUpper = regexp.MustCompile(`[A-Z]`)
	var hasLower = regexp.MustCompile(`[a-z]`)
	var hasSpecial = regexp.MustCompile(`[!@#~$%^&*()_+|<>,.?/:;{}]`)

	response.Result = hasMinLen.MatchString(password) &&
		hasNumber.MatchString(password) &&
		hasUpper.MatchString(password) &&
		hasLower.MatchString(password) &&
		hasSpecial.MatchString(password)
	response.Message = "Password must contain at least 8 characters, a number, an uppercase letter, a lowercase letter and a special character"
	return response
}

func ValidatePhone(phone string) models.ValidateOutput {
	response.Message = "Invalid Contact No."
	if phone[0] == '+' {
		phone = phone[1:]
		response.Result = regexp.MustCompile(`^[0-9]{12}$`).MatchString(phone)
		return response
	}
	response.Result = regexp.MustCompile(`^[0-9]{10}$`).MatchString(phone)
	return response
}

func ValidateISBN(isbn string) models.ValidateOutput {
	response.Message = "Invalid ISBN Format (10 or 13 numeric digits)"
	var ten_dig = regexp.MustCompile(`^[0-9]{10}$`).MatchString(isbn)
	var thirteen_dig = regexp.MustCompile(`^[0-9]{13}$`).MatchString(isbn)
	response.Result = ten_dig || thirteen_dig
	return response
}
func ValidateName(name string) models.ValidateOutput {
	response.Message = "Invalid Name Format (3-50 characters, starts and ends with a letter and contains only letters and spaces and first word should be at least 3 characters)"
	if len(name) < 3 || len(name) > 50 {
		response.Result = false
		return response
	}
	response.Result = regexp.MustCompile(`^[a-zA-Z][a-zA-Z\s]*[a-zA-Z]$`).MatchString(name)

	words := strings.Fields(name)
	if len(words) < 1 || len(words[0]) < 3 {
		response.Result = false
		return response
	}
	return response
}
