package errors

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hlmn/senyum-go-utils/structs"
	"github.com/spf13/cast"
)

type CustomHttpError struct {
	Code        int
	Description interface{}
	Internal    error
	Response    structs.Response
}

func (che *CustomHttpError) Error() string {
	if che.Internal == nil {
		return fmt.Sprintf("code=%d, description=%v", che.Code, che.Description)
	}
	return fmt.Sprintf("code=%d, description=%v, internal=%v", che.Code, che.Description, che.Internal)
}

func (che *CustomHttpError) SetInternal(err error) *CustomHttpError {
	che.Internal = err
	return che
}

func (che *CustomHttpError) SetResponseCode(code string) *CustomHttpError {
	che.Response.Code = code
	return che
}
func (che *CustomHttpError) OverrideResponseBody(response structs.Response) *CustomHttpError {
	che.Response = response
	return che
}

func NewHttpError(code int, description interface{}) *CustomHttpError {
	var response structs.Response
	if description == "" {
		description = http.StatusText(code)
	}
	he := &CustomHttpError{
		Code:        code,
		Description: description,
		Response:    response,
	}

	descStr, ok := description.(string)
	if ok {
		he.Description = descStr
	} else {
		err, ok := description.(error)
		if ok {
			he.Internal = err
		}
		he.Description = strings.ToUpper(strings.ReplaceAll(http.StatusText(code), " ", "_"))
	}

	if he.Response.Code == "" {
		he.Response.Code = fmt.Sprintf("%04d", he.Code)
	}
	he.Response.Description = cast.ToString(he.Description)

	return he
}
