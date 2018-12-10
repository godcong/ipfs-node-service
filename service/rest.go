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
	//ctx, cancel := context.WithCancel(context.Background())
	s := &RestServer{
		Server: &http.Server{
			Addr:    addr,
			Handler: gin.Default(),
		},
		//Context:,
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
