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

	// add a wardrobe for a user
	router.POST("/wardrobe/:username/:id", s.addWardrobe)

	// get all wardrobe for a user
	router.GET("/wardrobe/:username", s.getAllWardrobe)

	// get a wardrobe for a user
	router.GET("/wardrobe/:username/:id", s.getWardrobe)

	// delete a wardrobe for a user
	router.DELETE("/wardrobe/:username/:id", s.deleteWardrobe)

	return router
}
