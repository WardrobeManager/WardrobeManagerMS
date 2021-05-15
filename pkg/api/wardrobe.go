//
// definitions.go
//
// May 2021, Prashant Desai
//

package api

import (
	"crypto/md5"
	"encoding/hex"
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

func (w *wardrobeService) AddWardrobe(newWd NewWardrobeRequest) error {

	wc, err := w.db.Get(newWd.User)
	switch err := err.(type) {
	case nil:
	case *UserNotFound:
		wc = &WardrobeCloset{
			User:      newWd.User,
			Wardrobes: make([]Wardrobe, 1),
		}
	default:
		return fmt.Errorf("Database access failure : %w", err)
	}

	//Store image to file
	imageFile := genUniqImageFileName(newWd.User, newWd.Description)
	labelFile := genUniqLabelFileName(newWd.User, newWd.Description)

	for _, file := range []string{imageFile, labelFile} {
		_, err := w.imageDb.GetFile(file)
		switch err := err.(type) {
		case nil:
			// file with same name found
			return fmt.Errorf("Duplicate file name : %w", err)
		case *NoSuchFileOrDirectory:
			// this is good
		default:
			return fmt.Errorf("File system access error : %w", err)
		}
	}

	//Store files
	err = w.imageDb.AddFile(imageFile, newWd.MainImage)
	if err != nil {
		return fmt.Errorf("Error saving image to file system : %w", err)
	}

	err = w.imageDb.AddFile(labelFile, newWd.LabelImage)
	if err != nil {
		return fmt.Errorf("Error saving image to file system : %w", err)
	}

	//Update user
	wc.Wardrobes = append(wc.Wardrobes, Wardrobe{
		MainFile:    imageFile,
		LabelFile:   labelFile,
		Description: newWd.Description,
	})

	err = w.db.Update(newWd.User, wc)
	switch err := err.(type) {
	case nil:
	default:
		return fmt.Errorf("Database access failure : %w", err)
	}

	return nil
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

func (e *UserNotFound) Error() string {
	return fmt.Sprintf("User %s not found", e.User)
}

func (e *NoSuchFileOrDirectory) Error() string {
	return fmt.Sprintf("File %s not found", e.File)
}

func genUniqImageFileName(user string, filename string) string {
	stringToHash := []byte(user + "_image_" + filename)
	md5Bytes := md5.Sum(stringToHash)
	return hex.EncodeToString(md5Bytes[:])
}

func genUniqLabelFileName(user string, filename string) string {
	stringToHash := []byte(user + "_label_" + filename)
	md5Bytes := md5.Sum(stringToHash)
	return hex.EncodeToString(md5Bytes[:])

}
