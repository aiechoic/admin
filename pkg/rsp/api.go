package rsp

import (
	"github.com/aiechoic/admin/pkg/errs"
	"github.com/gin-gonic/gin"
)

func SendSuccess(c *gin.Context, data any) {
	c.JSON(200, Response{
		Success: true,
		Error:   "",
		Code:    0,
		Data:    data,
	})
}

func SendError(c *gin.Context, code errs.Code, err error) {
	if err == nil {
		err = code
	}
	c.JSON(200, Response{
		Success: false,
		Error:   err.Error(),
		Code:    code,
		Data:    nil,
	})
}
