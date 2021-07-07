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

	/*
		//v2 : add a wardrobe for a user
		router.POST("/users/:username/wardrobes", s.addWardrobeV2)

		//v2 : get all wardrobe for a user
		router.GET("/users/:username/wardrobes", s.getAllWardrobeV2)

		//v2 : get a wardrobe for a user
		router.GET("/users/:username/wardrobes/:id", s.getWardrobeV2)

		//v2: delete a wardrobe for a user
		router.DELETE("/users/:username/wardrobs/:id", s.deleteWardrobeV2)

		//v2: api to get image
		router.GET("/images/:filename", nil)
	*/

	return router
}
