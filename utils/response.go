package utils

import "github.com/gin-gonic/gin"

// Response 统一 API 响应结构体
// Code: 业务状态码，0 表示成功，非 0 表示错误（对应 HTTP 状态码）
// Msg:  提示信息
// Data: 响应数据，omitempty 表示无数据时不序列化该字段
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Success 返回成功响应（HTTP 200，业务码 0）
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{Code: 0, Msg: "ok", Data: data})
}

// Error 返回错误响应
// code 同时作为 HTTP 状态码和业务状态码
func Error(c *gin.Context, code int, msg string) {
	c.JSON(code, Response{Code: code, Msg: msg})
}
