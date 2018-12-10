package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// RestServer ...
type RestServer struct {
	*gin.Engine
	Server *http.Server
}

// NewRestServer ...
func NewRestServer(addr string) *RestServer {
	eng := gin.Default()
	s := &RestServer{
		Engine: eng,
		Server: &http.Server{
			Addr:    addr,
			Handler: eng,
		},
	}
	return s
}

// Start ...
func (s *RestServer) Start() {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: s.Engine,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()

}

// Stop ...
func (s *RestServer) Stop() {
	if err := s.Server.Shutdown(nil); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}

}

// JSON ...
func JSON(code int, msg string, detail ...gin.H) gin.H {
	if detail == nil {
		return gin.H{
			"code": code,
			"msg":  msg,
		}
	}
	return gin.H{
		"code":   code,
		"msg":    msg,
		"detail": detail[0],
	}
}