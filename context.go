package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	jsonContentType = "application/json"
	plainContentTyp = "text/plain"
	htmlContentType = "text/html"
)

// H 使得数据构建更加简洁
type H map[string]interface{}

// Context 上下文，每个 http 请求都会生成一个 Context
type Context struct {
	//ServeHTTP 所需要的两个参数
	ResponseWriter http.ResponseWriter
	Request        *http.Request

	Method string
	Path   string
}

//------------------------input-------------------------------------

func (c *Context) Postform(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}
func NewContext(ResponseWriter http.ResponseWriter, Request *http.Request) *Context {
	return &Context{
		Request:        Request,
		ResponseWriter: ResponseWriter,
		Method:         Request.Method,
		Path:           Request.URL.Path,
	}
}

//------------------------output-------------------------------------

// Status 写入请求头状态码
func (c *Context) Status(code int) {
	c.ResponseWriter.WriteHeader(code)
}

// SetHeader 重写请求头信息
func (c *Context) SetHeader(key, value string) {
	c.ResponseWriter.Header().Set(key, value)
}

// AddHeader 为该 Key 的请求头添加 value
func (c *Context) AddHeader(key, value string) {
	c.ResponseWriter.Header().Add(key, value)
}

// SetContentTpye 设置返回类型
func (c *Context) SetContentTpye(value string) {
	c.SetHeader("Content-Type", value)
}
func (c *Context) JSON(code int, obj any) {
	c.Status(code)
	c.SetContentTpye(jsonContentType)
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		http.Error(c.ResponseWriter, err.Error(), http.StatusInternalServerError)
	}
	_, err = c.ResponseWriter.Write(jsonBytes)
	if err != nil {
		http.Error(c.ResponseWriter, err.Error(), http.StatusInternalServerError)
	}
}
func (c *Context) HTML(code int, html string) {
	c.Status(code)
	c.SetContentTpye(htmlContentType)
	_, err := c.ResponseWriter.Write([]byte(html))
	if err != nil {
		http.Error(c.ResponseWriter, err.Error(), http.StatusInternalServerError)
	}
}
func (c *Context) String(code int, format string, values ...interface{}) {
	c.Status(code)
	c.SetContentTpye(plainContentTyp)
	_, err := c.ResponseWriter.Write([]byte(fmt.Sprintf(format, values...)))
	if err != nil {
		if err != nil {
			http.Error(c.ResponseWriter, err.Error(), http.StatusInternalServerError)
		}
	}
}
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	_, err := c.ResponseWriter.Write(data)
	if err != nil {
		http.Error(c.ResponseWriter, err.Error(), http.StatusInternalServerError)
	}
}