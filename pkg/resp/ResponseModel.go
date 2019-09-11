package resp

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simple_project/pkg/e"
)

type ResponseModel struct {
	Context *gin.Context `json:"-"`
	Code    int          `json:"code"`
	Msg     string       `json:"msg"`
	Data    interface{}  `json:"data"`
}

func CreateResponse(ctx *gin.Context) *ResponseModel {
	return &ResponseModel{Context: ctx, Code: e.ERROR, Msg: e.GetMsg(e.ERROR), Data: ""}
}

func (res *ResponseModel) Ok() {
	res.Msg = e.GetMsg(e.SUCCESS)
	res.Code = e.SUCCESS
	res.Data = ""
	res.Context.JSON(http.StatusOK, res)
}

func (res *ResponseModel) Error() {
	res.ResponseJSON()
}

//ResponseJSON 通用返回数据格式
func (res *ResponseModel) ResponseJSON() {
	res.Context.JSON(http.StatusOK, res)
}
