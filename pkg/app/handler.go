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

	glog.Infof("add wardrobe for {user=%s}, {wardrobe-id=%s}", username, wardId)

	var newWd api.NewWardrobeRequest
	err := c.Bind(&newWd)
	if err != nil {
		glog.Errorf("Error decoding Form : {err=%v} ", username, wardId, err)
		c.String(http.StatusUnprocessableEntity, fmt.Sprintf("error decoding JSON : %s", err))
		return
	}
	glog.Infof("done Bind for {user=%s}, {wardrobe-id=%s}", username, wardId)

	newWd.User = username
	newWd.Id = wardId
	err = s.ws.AddWardrobe(newWd)
	if err != nil {
		glog.Errorf("Error adding wardrobe, {err=%v} ", err)
		c.String(http.StatusUnprocessableEntity, fmt.Sprintf("error adding wardrobe: %s", err))
		return
	}

	glog.Infof("done adding wardrobe for {user=%s}, {wardrobe-id=%s}", username, wardId)

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

func (s *Server) getFile(c *gin.Context) {
	filename := c.Params.ByName("filename")

	glog.Infof("Get file {filename=%s}", filename)

	fileHandler := func(filepath string) error {

		c.File(filepath)

		return nil
	}

	err := s.ws.GetFile(filename, fileHandler)
	if err != nil {
		glog.Errorf("Error retrieving file, {err=%s}", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		return
	}

}
