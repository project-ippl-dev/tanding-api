package tools

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func StringEquals(str string) validation.RuleFunc {
	return func(value interface{}) error {
		val := value.(string)
		if val != str {
			return fmt.Errorf("those item didn't match, try again")
		}
		return nil
	}
}
