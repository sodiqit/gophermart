package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidateJSONBody(ctx context.Context, body io.ReadCloser, dest interface{}) error {
	if err := json.NewDecoder(body).Decode(&dest); err != nil {
		return errors.New("invalid json body")
	}

	validate := validator.New()

	err := validate.StructCtx(ctx, dest)
	if err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			var errMessages []string
			for _, e := range errs {
				var errMsg string
				switch e.Tag() {
				case "min":
					errMsg = fmt.Sprintf("Field validation for '%s' failed on the '%s' tag, minimum length is %s", e.Field(), e.Tag(), e.Param())
				case "max":
					errMsg = fmt.Sprintf("Field validation for '%s' failed on the '%s' tag, maximum length is %s", e.Field(), e.Tag(), e.Param())
				default:
					errMsg = fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", e.Field(), e.Tag())
				}
				errMessages = append(errMessages, errMsg)
			}

			return errors.New(strings.Join(errMessages, "; "))
		}
		return err
	}

	return nil
}
