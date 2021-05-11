//
// definitions.go
//
// May 2021, Prashant Desai
//

package api

import (
	"fmt"
)

type WardrobeService interface {
	AddWardrobe(new NewWardrobeRequest) error
	DeleteWardrobe(user string, id string) error
	GetWardrobe(user string, id string) (*NewWardrobeRequest, error)
	GetAllWardrobe(user string) ([]NewWardrobeRequest, error)
}

type WardrobeRepository interface {
	Add(user string, wardrobes *WardrobeCloset) error
	Get(user string) (*WardrobeCloset, error)
	Update(user string, wardrobes *WardrobeCloset) error
	Delete(user string) error
}

type ImageRepository interface {
	AddFile(name string, file []byte) error
	GetFile(name string) ([]byte, error)
	UpdateFile(name string, file []byte) error
	DeleteFile(name string) error
}

type wardrobeService struct {
	db      WardrobeRepository
	imageDb ImageRepository
}

func NewWardrobeService(dbIn WardrobeRepository, imageDbIn ImageRepository) (WardrobeService, error) {
	service := &wardrobeService{
		db:      dbIn,
		imageDb: imageDbIn,
	}

	return service, nil
}

func (w *wardrobeService) AddWardrobe(new NewWardrobeRequest) error {
	//Get user

	//Store image to file

	//Update user

	return fmt.Errorf("Failed to AddWardrobe")
}

func (w *wardrobeService) DeleteWardrobe(user string, id string) error {
	return nil
}

func (w *wardrobeService) GetWardrobe(user string, id string) (*NewWardrobeRequest, error) {
	return &NewWardrobeRequest{}, nil
}

func (w *wardrobeService) GetAllWardrobe(user string) ([]NewWardrobeRequest, error) {
	return []NewWardrobeRequest{NewWardrobeRequest{}}, nil
}
