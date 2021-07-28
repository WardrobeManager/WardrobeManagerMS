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
	router.GET("/users/:username/wardrobes", s.getAllWardrobes)

	//get a wardrobe for a user
	router.GET("/users/:username/wardrobes/:id", s.getWardrobe)

	//delete a wardrobe for a user
	router.DELETE("/users/:username/wardrobs/:id", s.deleteWardrobe)

	//api to get image
	router.GET("/images/:filename", s.getFile)

	//add a outfit for a user
	router.POST("/users/:username/outfits", s.addOutfit)

	//get all outfits for a user
	router.GET("/users/:username/outfits", s.getAllOutfits)

	//get a wardrobe for a user
	router.GET("/users/:username/outfits/:id", s.getOutfit)

	//delete a wardrobe for a user
	router.DELETE("/users/:username/outfits/:id", s.deleteOutfit)

	return router
}
