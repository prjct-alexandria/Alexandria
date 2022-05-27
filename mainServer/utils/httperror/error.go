package httperror

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// Using http error handling as in:
// https://github.com/swaggo/swag/blob/master/example/celler/httputil/error.go

func NewError(ctx *gin.Context, status int, err error) {
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	fmt.Println(err)
	ctx.JSON(status, er)
}

type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}
