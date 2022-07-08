package xgin

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"skuld/ecode"
)

var emptyData = make(map[string]interface{})

type response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func (r *response) reset() {
	r.Code = 0
	r.Data = nil
	r.Msg = ""
}

var (
	resppool = sync.Pool{New: func() interface{} {
		return &response{
			Code: 0,
			Data: nil,
			Msg:  "",
		}
	}}
)

func Success(c *gin.Context, datas ...interface{}) {
	var data interface{}
	if len(datas) > 0 {
		data = datas[0]
	}
	result(c, http.StatusOK, "OK", data)
}

func Failure(c *gin.Context, err error) {
	c.Abort()
	errCode := c.GetInt("ErrorCode")
	code, msg := ecode.Info(err)
	if errCode > 0 {
		code = errCode
	}
	result(c, code, msg, emptyData)
}

func FailureData(c *gin.Context, err error, data interface{}) {
	c.Abort()
	errCode := c.GetInt("ErrorCode")
	code, msg := ecode.Info(err)
	if errCode > 0 {
		code = errCode
	}
	result(c, code, msg, data)
}

func result(c *gin.Context, code int, msg string, data interface{}) {
	httpCode := http.StatusOK
	if code > 100200 && code <= 100600 {
		httpCode = code - 100000
	}
	if data == nil {
		data = emptyData
	}
	resp := resppool.Get().(*response)
	resp.Code = code
	resp.Data = data
	resp.Msg = msg

	c.JSON(httpCode, resp)

	resp.reset()
	resppool.Put(resp)
}
