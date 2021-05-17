//
// api_test.go
//
// May 2021, Prashant Desai
//

package api_test

import (
	"WardrobeManagerMS/pkg/api"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"testing"
)

type mockWardRepo struct{}

func (m *mockWardRepo) Add(user string, wards *api.WardrobeCloset) error {
	if user == "WardrobeDbUnavailableUser" {
		return api.ResourceUnavailable{}
	}

	return nil
}

func (m *mockWardRepo) Get(user string) (*api.WardrobeCloset, error) {
	if user == "WardrobeDbUnavailableUser" {
		return nil, &api.ResourceUnavailable{
			Server: "someserver:57400",
		}
	}
	return &api.WardrobeCloset{}, nil
}

func (m *mockWardRepo) Update(user string, wards *api.WardrobeCloset) error {
	fmt.Printf("Updating user %s to repository\n", user)
	fmt.Println(wards)
	return nil
}

func (m *mockWardRepo) Delete(user string) error {
	return nil
}

type mockImageRepo struct{}

func (m *mockImageRepo) AddFile(name string, file []byte) error {
	fmt.Printf("Adding file to image folders %s\n", name)
	return nil
}

func (m *mockImageRepo) GetFile(name string) ([]byte, error) {
	fileName := tsGenUniqImageFileName("DuplicateImageFileUser", "DupLeggings")
	if name == fileName {
		return []byte{}, nil
	}

	return []byte{}, &api.NoSuchFileOrDirectory{File: name}
}

func (m *mockImageRepo) UpdateFile(name string, file []byte) error {
	return nil
}

func (m *mockImageRepo) DeleteFile(name string) error {
	return nil
}

func TestAddWardrobeService(t *testing.T) {

	mockWardrobe := &mockWardRepo{}
	mockImage := &mockImageRepo{}

	ws, err := api.NewWardrobeService(mockWardrobe, mockImage)
	if err != nil {
		t.Errorf(" NewWardrobService failed : %v", err)
	}

	cases := []struct {
		name     string
		newWd    api.NewWardrobeRequest
		expected error
	}{
		{
			name: "BasicAddNewWardrobeRequest",
			newWd: api.NewWardrobeRequest{
				User:        "foobar",
				Id:          "",
				Description: "Leggings",
				MainImage:   []byte{0xAA, 0xBB, 0xCC},
				LabelImage:  []byte{0xAA, 0xBB, 0xCC},
			},
			expected: nil,
		},
		{
			name: "WardrobeDBIsUnavailable",
			newWd: api.NewWardrobeRequest{
				User:        "WardrobeDbUnavailableUser",
				Id:          "",
				Description: "Leggings",
				MainImage:   []byte{0xAA, 0xBB, 0xCC},
				LabelImage:  []byte{0xAA, 0xBB, 0xCC},
			},
			expected: &api.ResourceUnavailable{
				Server: "someserver:57400",
			},
		},
		{
			name: "DuplicateImageFile",
			newWd: api.NewWardrobeRequest{
				User:        "DuplicateImageFileUser",
				Id:          "",
				Description: "DupLeggings",
				MainImage:   []byte{0xAA, 0xBB, 0xCC},
				LabelImage:  []byte{0xAA, 0xBB, 0xCC},
			},
			expected: &api.DuplicateFile{
				File: "",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := ws.AddWardrobe(c.newWd)

			if c.expected == nil {
				if err != nil {
					t.Errorf("Expected nil, got %v", err)
				}
			} else {
				if errors.As(err, &c.expected) == false {
					t.Errorf("Expected %v, got %v", c.expected, err)
				}
			}
		})
	}
}

func tsGenUniqImageFileName(user string, filename string) string {
	stringToHash := []byte(user + "_image_" + filename)
	md5Bytes := md5.Sum(stringToHash)
	return hex.EncodeToString(md5Bytes[:])
}
