package service

import "github.com/gin-gonic/gin"

// Router ...
func Router(engine *gin.Engine) error {
	group := engine.Group("/")

	group.Use()

	group.Handle()

	return nil
}
