package config

import (
	"errors"
	"regexp"
	"time"

	"github.com/thedevsaddam/govalidator"
)

func InitValidator() {
	govalidator.AddCustomRule("blacklist", func(field string, rule string, message string, value interface{}) error {
		val, ok := value.(string)
		if !ok {
			return errors.New("The " + field + " is not a string")
		}

		if val == "" {
			return nil
		}

		match, _ := regexp.MatchString(`^[^'"\[\]<>\{\}]+$`, val)
		if !match {
			return errors.New("The " + field + " contains unsafe characters")
		}

		return nil
	})

	govalidator.AddCustomRule("datetime", func(field string, rule string, message string, value interface{}) error {
		val, ok := value.(string)
		if !ok {
			return errors.New("The " + field + " is not datetime format")
		}

		_, err := time.Parse("2006-01-02 15:04:05", val)
		if err != nil {
			return errors.New("the " + field + " is not a datetime format")
		}

		return nil
	})
}
