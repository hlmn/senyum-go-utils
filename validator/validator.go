package config

import (
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
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

	govalidator.AddCustomRule("base64Std", func(field string, rule string, message string, value interface{}) error {
		val, ok := value.(string)
		if !ok {
			return errors.New("The " + field + " is not base64 format")
		}

		_, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			return errors.New("The " + field + " is not base64 format")
		}

		return nil
	})

	govalidator.AddCustomRule("float_between", func(field string, rule string, message string, value interface{}) error {
		val, ok := value.(string)
		if !ok {
			return fmt.Errorf("The %s field is not a string", field)
		}

		if val == "" {
			return nil
		}

		data, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return fmt.Errorf("The %s field is not a float", field)
		}

		ruleData := strings.TrimPrefix(rule, "float_between:")

		if len(strings.Split(ruleData, ",")) != 2 {
			return fmt.Errorf("The %s field is invalid rule", field)
		}

		if strings.Split(ruleData, ",")[0] != "" {
			minValue := strings.Split(ruleData, ",")[0]

			min, err := strconv.ParseFloat(minValue, 64)
			if err != nil {
				return fmt.Errorf("The %s field is not a float for min value", field)
			}

			if data < min {
				return fmt.Errorf("The %s field must be minimum %f", field, min)
			}
		}

		if strings.Split(ruleData, ",")[1] != "" {
			maxValue := strings.Split(ruleData, ",")[1]

			max, err := strconv.ParseFloat(maxValue, 64)
			if err != nil {
				return fmt.Errorf("The %s field is not a float for max value", field)
			}

			if data > max {
				return fmt.Errorf("The %s field must be maximum %f", field, max)
			}
		}

		return nil
	})

	govalidator.AddCustomRule("reject_whitespace", func(field string, rule string, message string, value interface{}) error {
		val, ok := value.(string)
		if !ok {
			return errors.New("The " + field + "  is only filled by space")
		}
		if val == "" {
			return nil
		}
		newValue := strings.ReplaceAll(val, " ", "")

		if newValue == "" {
			return errors.New("The " + field + " is only filled by space")
		}

		return nil

	})
}
