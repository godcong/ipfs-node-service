package service

import (
	"github.com/gin-gonic/gin"
	"github.com/godcong/node-service/config"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

// ContentTypeJSON ...
const ContentTypeJSON = "application/json"

// RestServer ...
type RestServer struct {
	*gin.Engine
	config *config.Configure
	server *http.Server
	Port   string
}

// NewRestServer ...
func NewRestServer(cfg *config.Configure) *RestServer {
	s := &RestServer{
		Engine: gin.Default(),
		config: cfg,
		Port:   config.DefaultString(cfg.REST.Port, ":7787"),
	}
	return s
}

// Start ...
func (s *RestServer) Start() {
	if !s.config.REST.Enable {
		return
	}

	Router(s.Engine)

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

type restBack struct {
	BackURL string
	Version string
}

// NewRestBack ...
func NewRestBack(cfg *config.Configure) StreamerCallback {
	return &restBack{
		BackURL: config.DefaultString(cfg.Callback.BackAddr, "localhost:7787"),
		Version: config.DefaultString(cfg.Callback.Version, "v0"),
	}
}

// Callback ...
func (s *restBack) Callback(result *QueueResult) error {
	back := filepath.Join(CheckPrefix(s.BackURL), s.Version, "ipfs/callback")
	log.Println(back)

	resp, err := http.Post(back, ContentTypeJSON, strings.NewReader(result.JSON()))
	if err != nil {
		log.Println(err)
		return err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	log.Println(string(bytes), err)
	return err
}

// CheckPrefix ...
func CheckPrefix(url string) string {
	if strings.Index(url, "http") != 0 {
		return "http://" + url
	}
	return url
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
