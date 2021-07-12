//
// definitions.go
//
// May 2021, Prashant Desai
//

package api

import (
	"mime/multipart"
)

const Version = "1.0"

type NewWardrobeRequest struct {
	User           string
	Id             string
	Description    string `form:"description" binding:"required"`
	MainImage      []byte
	LabelImage     []byte
	MainImageMime  *multipart.FileHeader `form:"main-image" binding:"required"`
	LabelImageMime *multipart.FileHeader `form:"label-image" binding:"required"`
}

type Wardrobe struct {
	Identifier  string `bson:"id"`
	MainFile    string `bson:"main-file"`
	LabelFile   string `bson:"label-file"`
	Description string `bson:"description"`
	LabelText   string `bson:"label-text"`
}

type WardrobeCloset struct {
	User      string `bson:"user"`
	Wardrobes []Wardrobe
}

type LabelToTextRequest struct {
	User     string `json:"user" binding="required"`
	Id       string `json:"id" binding="required"`
	RawImage string `json:"raw-image" binding="required"`
}

type LabelToTextResponse struct {
	User string `json:"user" binding="required"`
	Id   string `json:"id" binding="required"`
	Text string `json:"text" binding="required"`
}

type GetWardrobeResponse struct {
	Id          string `json:"id" binding="required"`
	Description string `json:"description" binding:"required"`
	MainImage   string `json:"main-image-uri" binding:"required"`
	LabelImage  string `json:"label-image-uri" binding:"required"`
}

// Functions
type HandleFile func(filename string) error

// Error
type UserNotFound struct {
	User string
}

type NoSuchFileOrDirectory struct {
	File string
}

type ResourceUnavailable struct {
	Server string
}

type DuplicateFile struct {
	File string
}
