//
// main.go
//
// May 2021, Prashant Desai
//

package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"WardrobeManagerMS/pkg/api"
	"WardrobeManagerMS/pkg/app"
	repo "WardrobeManagerMS/pkg/repository"
)

const logFile = "/tmp/WM/gin.log"
const imageRepo = "/tmp/ImageDb"

const redisServer = ":6379"
const txChannel = "Label"
const rxChannel = "Text"

func init() {
	flag.Parse()
}

func main() {

	glog.Infof("Starting WM with {GIN-debug=%s}, {Image=%s}", logFile, imageRepo)
	r := gin.Default()

	f, err2 := os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err2 != nil {
		fmt.Printf("Failing to create log file : %v", err2)
	} else {
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	}

	mongoWardrobeRepo, err := repo.NewWardrobeRepository()
	if err != nil {
		glog.Errorf(" Initializing Mongo repository failed  : %v", err)
		return
	}

	imageRepo, err1 := repo.NewFileImageRepository("/tmp/ImageDb")
	if err1 != nil {
		glog.Errorf(" Initializing file repository failed  : %v", err1)
		return
	}

	ws, err2 := api.NewWardrobeService(mongoWardrobeRepo, imageRepo, redisServer, rxChannel, txChannel)
	if err2 != nil {
		glog.Errorf(" NewWardrobService failed : %v", err2)
		return
	}

	server := app.NewWardrobeServer(r, ws)

	// start the server
	err = server.Run()
	if err != nil {
		glog.Errorf(" Server run failed : %v", err)
		return
	}

}
