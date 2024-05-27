package controller

import (
	"net/http"

	"example.com/m/v2/utils"
	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Code utils.ResCode `json:"code"`
	Msg  interface{}   `json:"msg"`
	Data interface{}   `json:"data,omitempty"`
}

func ResponseError(c *gin.Context, code utils.ResCode) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: utils.CodeSuccess,
		Msg:  utils.CodeSuccess.Msg(),
		Data: data,
	})
}

func ResponseErrorWithMsg(c *gin.Context, code utils.ResCode, msg interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
