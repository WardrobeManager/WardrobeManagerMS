//
// definitions.go
//
// May 2021, Prashant Desai
//

package api

const Version = "1.0"

type NewWardrobeRequest struct {
	User        string `json:"user"`
	Id          string `json:"id"`
	Description string `json:"description"`
	MainImage   []byte `json:"main-image"`
	LabelImage  []byte `json:"label-image"`
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
