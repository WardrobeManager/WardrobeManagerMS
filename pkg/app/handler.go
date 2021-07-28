//
// handler.go
//
// May 2021, Prashant Desai
//

package app

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"WardrobeManagerMS/pkg/api"
)

func (s *Server) addWardrobe(c *gin.Context) {

	username := c.Params.ByName("username")

	glog.Infof("add wardrobe for {user=%s}", username)

	var newWd api.NewWardrobeRequest
	err := c.Bind(&newWd)
	if err != nil {
		glog.Errorf("Error decoding Form : {err=%v} ", username, err)
		c.String(http.StatusUnprocessableEntity, fmt.Sprintf("error decoding JSON : %s", err))
		return
	}
	glog.Infof("done Bind for {user=%s}", username)

	newWd.User = username
	err = s.ws.AddWardrobe(newWd)
	if err != nil {
		glog.Errorf("Error adding wardrobe, {err=%v} ", err)
		c.String(http.StatusUnprocessableEntity, fmt.Sprintf("error adding wardrobe: %s", err))
		return
	}

	glog.Infof("done adding wardrobe for {user=%s}", username)

	c.String(http.StatusOK, "addUser")
}

func (s *Server) getAllWardrobes(c *gin.Context) {
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

func (s *Server) addOutfit(c *gin.Context) {

	username := c.Params.ByName("username")

	glog.Infof("add outfit for {user=%s}", username)

	var newOt api.NewOutfitRequest
	err := c.BindJSON(&newOt)
	if err != nil {
		glog.Errorf("Error decoding Form {users=%s}: {err=%v} ", username, err)
		c.String(http.StatusUnprocessableEntity, fmt.Sprintf("error decoding JSON : %s", err))
		return
	}
	glog.Infof("done Bind for {user=%s}", username)

	newOt.User = username
	err = s.ws.AddOutfit(newOt)
	if err != nil {
		glog.Errorf("Error adding outfit, {err=%v} ", err)
		c.String(http.StatusUnprocessableEntity, fmt.Sprintf("error adding outfit: %s", err))
		return
	}

	glog.Infof("done adding wardrobe for {user=%s}", username)

	c.String(http.StatusOK, "addUser")
}

func (s *Server) getAllOutfits(c *gin.Context) {
	username := c.Params.ByName("username")

	glog.Infof("Get all outfits for {user=%s}", username)

	outfits, err := s.ws.GetAllOutfits(username)
	if err != nil {
		glog.Errorf("Error geting all outfits, {err=%v} ", err)
		c.String(http.StatusUnprocessableEntity, fmt.Sprintf("error: %s", err))
		return
	}

	c.JSON(http.StatusOK, &outfits)
}

func (s *Server) getOutfit(c *gin.Context) {
	username := c.Params.ByName("username")
	otId := c.Params.ByName("id")

	glog.Infof("Get outfit for {user=%s}, {outfit-id=%s} ", username, otId)

	outfit, err := s.ws.GetOutfit(username, otId)
	if err != nil {
		glog.Errorf("Error get outfit,{err=%v}", err)
		c.String(http.StatusUnprocessableEntity, fmt.Sprintf("error: %s", err))
		return
	}

	c.JSON(http.StatusOK, &outfit)
}

func (s *Server) deleteOutfit(c *gin.Context) {
	username := c.Params.ByName("username")
	otId := c.Params.ByName("id")

	glog.Infof("Delete wardrobe for {user=%s}, {outfit-id=%s} ", username, otId)

	err := s.ws.DeleteWardrobe(username, otId)
	if err != nil {
		glog.Errorf("Error deleting outfit, {err=%s}", err)
		c.String(http.StatusUnprocessableEntity, fmt.Sprintf("error: %s", err))
		return
	}

	c.String(http.StatusOK, "deleteOutfit")
}

//utility
func printRequest(c *gin.Context) {

	// Print Header
	fmt.Println(c.Request.Host, c.Request.RemoteAddr, c.Request.RequestURI)

	// Save a copy of this request for debugging.
	requestDump, err := httputil.DumpRequest(c.Request, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))

	// Hex dump of request body
	body, _ := ioutil.ReadAll(c.Request.Body)
	println(hex.Dump(body))

	c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
}
