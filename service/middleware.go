package service

import "github.com/gin-gonic/gin"

// Middleware ...
func Middleware() func(context *gin.Context) {
	return func(context *gin.Context) {
		//something
	}
}
