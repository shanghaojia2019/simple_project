package resp

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseModel struct {
	Context *gin.Context `json:"-"`
	Code    int          `json:"code"`
	Msg     string       `json:"msg"`
	Data    interface{}  `json:"data"`
}

//ResponseJSON 通用返回数据格式
func ResponseJSON(data *ResponseModel) {
	data.Context.JSON(http.StatusOK, data)
}
