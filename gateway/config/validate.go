package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func (c *Config) validate() error {
	v := validator.New()
	v.RegisterTagNameFunc(getTagName)
	err := v.Struct(c)

	failedKeys := make([]string, 0)
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, err := range errs {
			failedKeys = append(failedKeys, err.Field())
		}
	}

	if len(failedKeys) > 0 {
		return fmt.Errorf("config: failed_keys: %v", failedKeys)
	}
	return nil
}

func getTagName(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("env"), ",", 2)[0]

	if name == "" {
		name = fld.Tag.Get("yaml")
	}

	if name == "" {
		name = fld.Name
	}

	return name
}
