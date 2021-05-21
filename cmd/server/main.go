//
// main.go
//
// May 2021, Prashant Desai
//

package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"WardrobeManagerMS/pkg/api"
	repo "WardrobeManagerMS/pkg/repository"
)

var ws api.WardrobeService

func main() {

	r := gin.Default()

	mongoWardrobeRepo, err := repo.NewWardrobeRepository()
	if err != nil {
		fmt.Printf(" Initializing Mongo repository failed  : %v", err)
	}

	imageRepo, err1 := repo.NewFileImageRepository("/tmp")
	if err != nil {
		fmt.Printf(" Initializing file repository failed  : %v", err1)
	}

	ws, err = api.NewWardrobeService(mongoWardrobeRepo, imageRepo)
	if err != nil {
		fmt.Printf(" NewWardrobService failed : %v", err)
	}

	// add a wardrobe for a user
	r.POST("/wardrobe/:username/:id", addWardrobe)

	// get all wardrobe for a user
	r.GET("/wardrobe/:username", getAllWardrobe)

	// get a wardrobe for a user
	r.GET("/wardrobe/:username/:id", getWardrobe)

	// delete a wardrobe for a user
	r.DELETE("/wardrobe/:username/:id", deleteWardrobe)

	r.Run(":57400")
	fmt.Println("hello, welcome to user management MS")
}

func addWardrobe(c *gin.Context) {

	username := c.Params.ByName("username")
	wardId := c.Params.ByName("id")

	var newWd api.NewWardrobeRequest
	err := c.BindJSON(&newWd)
	if err != nil {
		c.String(http.StatusUnprocessableEntity, fmt.Sprintf("error decoding JSON : %s", err))
		return
	}

	newWd.User = username
	newWd.Id = wardId
	err = ws.AddWardrobe(newWd)
	if err != nil {
		c.String(http.StatusUnprocessableEntity, fmt.Sprintf("error adding wardrobe: %s", err))
		return
	}

	c.String(http.StatusOK, "addUser")
}

func getAllWardrobe(c *gin.Context) {
	username := c.Params.ByName("username")

	fmt.Printf("getAllWardrobe:%s", username)

	wards, err := ws.GetAllWardrobe(username)
	if err != nil {
		c.String(http.StatusUnprocessableEntity, fmt.Sprintf("error: %s", err))
		return
	}

	c.JSON(http.StatusOK, &wards)
}

func getWardrobe(c *gin.Context) {
	username := c.Params.ByName("username")
	wardId := c.Params.ByName("id")

	fmt.Printf("getWardrobe:%s:%s", username, wardId)
	c.JSON(http.StatusOK, "getWardrobe")
}

func deleteWardrobe(c *gin.Context) {
	username := c.Params.ByName("username")
	wardId := c.Params.ByName("id")

	err := ws.DeleteWardrobe(username, wardId)
	if err != nil {
		c.String(http.StatusUnprocessableEntity, fmt.Sprintf("error: %s", err))
		return
	}

	fmt.Printf("deleteWardrobe:%s:%s", username, wardId)
	c.String(http.StatusOK, "deleteWardrobe")
}
