//
// server.go
//
// May 2021, Prashant Desai
//

package app

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"WardrobeManagerMS/pkg/api"
)

const port = "57401"

type Server struct {
	router *gin.Engine
	ws     api.WardrobeService
}

func NewWardrobeServer(router *gin.Engine, ws api.WardrobeService) *Server {
	return &Server{
		router: router,
		ws:     ws,
	}

}

func (s *Server) Run() error {
	// run function that initializes the routes
	r := s.Routes()

	// run the server through the router
	err := r.Run(":" + port)
	if err != nil {
		glog.Errorf("Server - there was an error calling Run on router: %v", err)
		return err
	}

	return nil
}
