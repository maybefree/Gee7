package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"time"
)

type BodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w BodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w BodyLogWriter) WriteString (s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// 打印日志
func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestBody, _ := ctx.GetRawData()
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))

		bodyLogWriter := &BodyLogWriter{
			body: bytes.NewBufferString(""),
			ResponseWriter: ctx.Writer,
		}
		ctx.Writer = bodyLogWriter

		start := time.Now()
		// handler
		ctx.Next()
		// log
		end := time.Now()
		responseBody := bodyLogWriter.body.String()
		logField := map[string]interface{}{
			"uri": 				ctx.Request.URL.Path,
			"raw_query": 		ctx.Request.URL.RawQuery,
			"start_timestamp": 	start.Format("2006-01-02 15:04:05"),
			"end_timestamp":    end.Format("2006-01-02 15:04:05"),
			"server_name": 		ctx.Request.Host,
			"remote_addr":		ctx.ClientIP(),
			"proto":			ctx.Request.Proto,
			"referer":			ctx.Request.Referer(),
			"request_method":	ctx.Request.Method,
			"response_time":	end.Sub(start).Microseconds(),
			"content_type": 	ctx.Request.Header.Get("Content-Type"),
			"status":			ctx.Writer.Status(),
			"user_agent":		ctx.Request.UserAgent(),
			"request_body":		string(requestBody),
			"reponse_body":		responseBody,
			"response_err":		ctx.Errors.Last(),
		}

		bf := bytes.NewBuffer([]byte{})
		jsonEncoder := json.NewEncoder(bf)
		jsonEncoder.SetEscapeHTML(false)
		jsonEncoder.SetIndent("", "\t")
		jsonEncoder.Encode(logField)
		fmt.Println(bf.String())
	}
}




