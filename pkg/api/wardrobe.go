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
	DeleteAll(user string) error
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

	var addUser bool = false

	wc, err := w.db.Get(newWd.User)
	switch err := err.(type) {
	case nil:
	case *UserNotFound:
		addUser = true
		wc = &WardrobeCloset{
			User:      newWd.User,
			Wardrobes: make([]Wardrobe, 0),
		}
	case *ResourceUnavailable:
		return fmt.Errorf("Wardrobe db is unavailable : %w", err)
	default:
		return fmt.Errorf("Unknown error : %w", err)
	}

	//Store image to file
	imageFile := genUniqImageFileName(newWd.User, newWd.Description)
	labelFile := genUniqLabelFileName(newWd.User, newWd.Description)

	for _, file := range []string{imageFile, labelFile} {
		_, err := w.imageDb.GetFile(file)
		switch err := err.(type) {
		case NoSuchFileOrDirectory:
			// this is good
			break
		case nil:
			// file with same name found
			return &DuplicateFile{
				File: file,
			}
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
		Identifier:  newWd.Id,
		MainFile:    imageFile,
		LabelFile:   labelFile,
		Description: newWd.Description,
	})
	if addUser == true {
		err = w.db.Add(newWd.User, wc)
	} else {
		err = w.db.Update(newWd.User, wc)
	}
	switch err := err.(type) {
	case nil:
	default:
		return fmt.Errorf("Database access failure : %w", err)
	}

	return nil
}

func (w *wardrobeService) DeleteWardrobe(user string, id string) error {

	wc, err := w.db.Get(user)
	switch err := err.(type) {
	case nil:
		break
	case *UserNotFound:
		return fmt.Errorf("User not found %s : %w", user, err)
	case *ResourceUnavailable:
		return fmt.Errorf("Wardrobe db is unavailable : %w", err)
	default:
		return fmt.Errorf("Unknown error : %w", err)
	}

	if len(wc.Wardrobes) == 0 {
		return fmt.Errorf("Empty closet")
	}

	tmp := wc.Wardrobes[:0]
	for _, ward := range wc.Wardrobes {
		if ward.Identifier == id {
			err = w.imageDb.DeleteFile(ward.MainFile)
			if err != nil {
				fmt.Printf("Error deleting image file : %w", err)
			}

			err = w.imageDb.DeleteFile(ward.LabelFile)
			if err != nil {
				fmt.Printf("Error deleting label file : %w", err)
			}
		} else {
			tmp = append(tmp, ward)
		}
	}
	wc.Wardrobes = tmp

	err = w.db.Update(user, wc)
	switch err := err.(type) {
	case nil:
	default:
		return fmt.Errorf("Database access failure : %w", err)
	}

	return nil
}

func (w *wardrobeService) GetWardrobe(user string, id string) (*NewWardrobeRequest, error) {

	wc, err := w.db.Get(user)
	switch err := err.(type) {
	case nil:
		break
	case *UserNotFound:
		return nil, fmt.Errorf("User not found %s : %w", user, err)
	case *ResourceUnavailable:
		return nil, fmt.Errorf("Wardrobe db is unavailable : %w", err)
	default:
		return nil, fmt.Errorf("Unknown error : %w", err)
	}

	if len(wc.Wardrobes) == 0 {
		return nil, fmt.Errorf("Empty closet")
	}

	for _, ward := range wc.Wardrobes {
		if ward.Identifier == id {
			wardReq := &NewWardrobeRequest{
				User:        user,
				Id:          ward.Identifier,
				Description: ward.Description,
			}

			wardReq.MainImage, err = w.imageDb.GetFile(ward.MainFile)
			if err != nil {
				return nil, fmt.Errorf("Error accessing image : %w", err)
			}

			wardReq.LabelImage, err = w.imageDb.GetFile(ward.LabelFile)
			if err != nil {
				return nil, fmt.Errorf("Error accessing image : %w", err)
			}

			return wardReq, nil
		}
	}

	return nil, UserNotFound{User: user}
}

func (w *wardrobeService) GetAllWardrobe(user string) ([]NewWardrobeRequest, error) {

	wc, err := w.db.Get(user)
	switch err := err.(type) {
	case nil:
		break
	case *UserNotFound:
		return nil, fmt.Errorf("User not found %s : %w", user, err)
	case *ResourceUnavailable:
		return nil, fmt.Errorf("Wardrobe db is unavailable : %w", err)
	default:
		return nil, fmt.Errorf("Unknown error : %w", err)
	}

	if len(wc.Wardrobes) == 0 {
		return nil, fmt.Errorf("Empty closet")
	}

	wardReqs := make([]NewWardrobeRequest, 0)
	for _, ward := range wc.Wardrobes {
		wardReq := NewWardrobeRequest{
			User:        user,
			Id:          ward.Identifier,
			Description: ward.Description,
		}

		wardReq.MainImage, err = w.imageDb.GetFile(ward.MainFile)
		if err != nil {
			return nil, fmt.Errorf("Error accessing image : %w", err)
		}

		wardReq.LabelImage, err = w.imageDb.GetFile(ward.LabelFile)
		if err != nil {
			return nil, fmt.Errorf("Error accessing image : %w", err)
		}

		wardReqs = append(wardReqs, wardReq)
	}

	return wardReqs, nil
}

// Error codes
func (e UserNotFound) Error() string {
	return fmt.Sprintf("User %s not found", e.User)
}

func (e NoSuchFileOrDirectory) Error() string {
	return fmt.Sprintf("File %s not found", e.File)
}

func (e ResourceUnavailable) Error() string {
	return fmt.Sprintf("Service %s is down", e.Server)
}

func (e DuplicateFile) Error() string {
	return fmt.Sprintf("Duplicate file name %s", e.File)
}

/*
func (e DuplicateFile) Is(target error) bool {
	switch target.(type) {
	default:
		return false
	case *DuplicateFile:
		return true
	}

}
*/

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
