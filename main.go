//
// main.go
//
// May 2021, Prashant Desai
//

package main

import (
	"fmt"
	"reflect"

	c "./controller"
)

func main() {

	controller := c.NewWardrobeController()

	fmt.Println(reflect.TypeOf(controller))
}
