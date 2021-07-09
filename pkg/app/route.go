//
// route.go
//
// May 2021, Prashant Desai
//

package app

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) Routes() *gin.Engine {
	router := s.router

	//add a wardrobe for a user
	router.POST("/users/:username/wardrobes", s.addWardrobe)

	//get all wardrobe for a user
	router.GET("/users/:username/wardrobes", s.getAllWardrobe)

	//get a wardrobe for a user
	router.GET("/users/:username/wardrobes/:id", s.getWardrobe)

	//delete a wardrobe for a user
	router.DELETE("/users/:username/wardrobs/:id", s.deleteWardrobe)

	//api to get image
	router.GET("/images/:filename", s.getFile)

	return router
}
