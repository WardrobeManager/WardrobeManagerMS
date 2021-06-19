//
// definitions.go
//
// May 2021, Prashant Desai
//

package api

const Version = "1.0"

type NewWardrobeRequest struct {
	User        string `json:"user" binding:"required"`
	Id          string `json:"id" binding:"required"`
	Description string `json:"description" binding:"required"`
	MainImage   []byte `json:"main-image" binding:"required"`
	LabelImage  []byte `json:"label-image" binding:"required"`
}

type Wardrobe struct {
	Identifier  string `bson:"id"`
	MainFile    string `bson:"main-file"`
	LabelFile   string `bson:"label-file"`
	Description string `bson:"description"`
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
