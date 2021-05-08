//
// controller.go
//
// May 2021, Prashant Desai
//

package controller

import "fmt"

type WardrobeController interface {
	AddOutfit(username string)
}

type wardrobeController struct {
}

func NewWardrobeController() WardrobeController {
	return &wardrobeController{}
}

func (m *wardrobeController) AddOutfit(username string) {
	fmt.Println(username)
}
