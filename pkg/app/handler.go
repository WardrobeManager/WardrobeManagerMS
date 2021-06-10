//
// handler.go
//
// May 2021, Prashant Desai
//

package app

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"WardrobeManagerMS/pkg/api"
)

func (s *Server) addWardrobe(c *gin.Context) {

	username := c.Params.ByName("username")
	wardId := c.Params.ByName("id")

	glog.Infof("Add wardrobe for {user=%s}, {wardrobe-id=%s}", username, wardId)

	var newWd api.NewWardrobeRequest
	err := c.BindJSON(&newWd)
	if err != nil {
		glog.Errorf("Error decoding JSON : {err=%v} ", username, wardId, err)
		c.String(http.StatusUnprocessableEntity, fmt.Sprintf("error decoding JSON : %s", err))
		return
	}

	newWd.User = username
	newWd.Id = wardId
	err = s.ws.AddWardrobe(newWd)
	if err != nil {
		glog.Errorf("Error adding wardrobe, {err=%v} ", err)
		c.String(http.StatusUnprocessableEntity, fmt.Sprintf("error adding wardrobe: %s", err))
		return
	}

	c.String(http.StatusOK, "addUser")
}

func (s *Server) getAllWardrobe(c *gin.Context) {
	username := c.Params.ByName("username")

	glog.Infof("Get all wardrobe for {user=%s}", username)

	wards, err := s.ws.GetAllWardrobe(username)
	if err != nil {
		glog.Errorf("Error geting all wardrobe, {err=%v} ", err)
		c.String(http.StatusUnprocessableEntity, fmt.Sprintf("error: %s", err))
		return
	}

	c.JSON(http.StatusOK, &wards)
}

func (s *Server) getWardrobe(c *gin.Context) {
	username := c.Params.ByName("username")
	wardId := c.Params.ByName("id")

	glog.Infof("Get wardrobe for {user=%s}, {wardrobe-id=%s} ", username, wardId)

	wards, err := s.ws.GetWardrobe(username, wardId)
	if err != nil {
		glog.Errorf("Error get wardrobe,{err=%v}", err)
		c.String(http.StatusUnprocessableEntity, fmt.Sprintf("error: %s", err))
		return
	}

	c.JSON(http.StatusOK, &wards)
}

func (s *Server) deleteWardrobe(c *gin.Context) {
	username := c.Params.ByName("username")
	wardId := c.Params.ByName("id")

	glog.Infof("Delete wardrobe for {user=%s}, {wardrobe-id=%s} ", username, wardId)

	err := s.ws.DeleteWardrobe(username, wardId)
	if err != nil {
		glog.Errorf("Error deleting wardrobe, {err=%s}", err)
		c.String(http.StatusUnprocessableEntity, fmt.Sprintf("error: %s", err))
		return
	}

	c.String(http.StatusOK, "deleteWardrobe")
}
