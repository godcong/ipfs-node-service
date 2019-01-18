package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// RestServer ...
type RestServer struct {
	*gin.Engine
	BackURL string
	Port    string
	server  *http.Server
}

// NewRestServer ...
func NewRestServer() *RestServer {
	s := &RestServer{
		Engine: gin.Default(),
		Port:   DefaultString(config.REST.Port, ":7780"),
	}
	return s
}

// Start ...
func (s *RestServer) Start() {
	if !config.REST.Enable {
		return
	}
	s.server = &http.Server{
		Addr:    s.Port,
		Handler: s.Engine,
	}
	go func() {
		log.Printf("Listening and serving HTTP on %s\n", s.Port)
		if err := s.server.ListenAndServe(); err != nil {
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()

}

// Stop ...
func (s *RestServer) Stop() {
	if err := s.server.Shutdown(nil); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
}

// Callback ...
func (s *RestServer) Callback(result *QueueResult) error {

	url := "/v0/ipfs/callback"

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
