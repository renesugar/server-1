package error

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"net/http"
	"strings"
	"unicode"

	"github.com/gotify/server/model"
	"gopkg.in/go-playground/validator.v8"
)

// Handler creates a gin middleware for handling errors.
func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				switch e.Type {
				case gin.ErrorTypeBind:
					errs, ok := e.Err.(validator.ValidationErrors)

					if !ok {
						writeError(c, e.Error())
						return
					}

					var stringErrors []string
					for _, err := range errs {
						stringErrors = append(stringErrors, validationErrorToText(err))
					}
					writeError(c, strings.Join(stringErrors, "; "))
				default:
					writeError(c, e.Err.Error())
				}
			}
		}
	}
}

func validationErrorToText(e *validator.FieldError) string {
	runes := []rune(e.Field)
	runes[0] = unicode.ToLower(runes[0])
	fieldName := string(runes)
	switch e.Tag {
	case "required":
		return fmt.Sprintf("Field '%s' is required", fieldName)
	}
	return fmt.Sprintf("Field '%s' is not valid", fieldName)
}

func writeError(ctx *gin.Context, errString string) {
	status := http.StatusBadRequest
	if ctx.Writer.Status() != http.StatusOK {
		status = ctx.Writer.Status()
	}
	ctx.JSON(status, &model.Error{Error: http.StatusText(status), ErrorCode: status, ErrorDescription: errString})
}
