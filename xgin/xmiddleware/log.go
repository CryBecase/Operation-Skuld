package xmiddleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"skuld/xgin"
	"skuld/xlogger"
)

type interceptWriter struct {
	buf *bytes.Buffer
	gin.ResponseWriter
}

func newInterceptWriter(w gin.ResponseWriter) *interceptWriter {
	return &interceptWriter{
		buf:            bytes.NewBufferString(""),
		ResponseWriter: w,
	}
}

func (w *interceptWriter) Write(data []byte) (int, error) {
	w.buf.Write(data)
	return w.ResponseWriter.Write(data)
}

func Log(logger xlogger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasSuffix(c.Request.RequestURI, "/ping") {
			return
		}

		requestId := rand.Int63()
		start := time.Now()

		bodydata, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error("Log io.ReadAll", "err", err)
			return
		}
		err = c.Request.Body.Close()
		if err != nil {
			logger.Error("Log c.Request.Body.Close", "err", err)
		}

		var bodylog interface{}
		if len(bodydata) > 10*1024 { // 10kb 不记录 body
			bodylog = fmt.Sprintf("body too large: %d byte", len(bodydata))
		} else {
			bodymap := make(map[string]interface{})
			err = json.Unmarshal(bodydata, &bodymap)
			if err != nil {
				bodylog = string(bodydata)
			} else {
				bodylog = bodymap
			}
		}

		logger.Info("request log", "request", map[string]interface{}{
			"request_id":     requestId,
			"request_host":   c.Request.Host,
			"request_uri":    c.Request.RequestURI,
			"request_header": c.Request.Header,
			"request_body":   bodylog,
		})

		c.Request.Body = io.NopCloser(bytes.NewReader(bodydata))
		w := newInterceptWriter(c.Writer)
		c.Writer = w

		defer func() {
			if r := recover(); r != nil {
				xgin.Failure(c, fmt.Errorf(recoverString(r)))
				var resplog interface{}
				respmap := make(map[string]interface{})
				err = json.Unmarshal(w.buf.Bytes(), &respmap)
				if err != nil {
					resplog = w.buf.String()
				} else {
					resplog = respmap
				}
				logger.Error("response panic log", "response", map[string]interface{}{
					"request_id":           requestId,
					"response_header":      c.Writer.Header(),
					"response_status_code": c.Writer.Status(),
					"response_body":        resplog,

					"panic": r,
					"stack": string(debug.Stack()),

					"request_host":   c.Request.Host,
					"request_uri":    c.Request.RequestURI,
					"request_header": c.Request.Header,
					"request_body":   bodylog,

					"duration": time.Since(start).Milliseconds(),
				})
			}
		}()

		c.Next()

		var resplog interface{}
		respmap := make(map[string]interface{})
		err = json.Unmarshal(w.buf.Bytes(), &respmap)
		if err != nil {
			resplog = w.buf.String()
		} else {
			resplog = respmap
		}
		logger.Info("response log", "response", map[string]interface{}{
			"request_id":           requestId,
			"response_header":      c.Writer.Header(),
			"response_status_code": c.Writer.Status(),
			"response_body":        resplog,

			"request_host":   c.Request.Host,
			"request_uri":    c.Request.RequestURI,
			"request_header": c.Request.Header,
			"request_body":   bodylog,

			"duration": time.Since(start).Milliseconds(),
		})
	}
}

func recoverString(r interface{}) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	case string:
		return v
	default:
		return "服务器异常"
	}
}
