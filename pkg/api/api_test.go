//
// api_test.go
//
// May 2021, Prashant Desai
//

package api_test

import (
	"WardrobeManagerMS/pkg/api"
	"testing"
)

type mockWardRepo struct{}

func (m *mockWardRepo) Add(user string, wards *api.WardrobeCloset) error {
	return nil
}

func (m *mockWardRepo) Get(user string) (*api.WardrobeCloset, error) {
	return &api.WardrobeCloset{}, nil
}

func (m *mockWardRepo) Update(user string, wards *api.WardrobeCloset) error {
	return nil
}

func (m *mockWardRepo) Delete(user string) error {
	return nil
}

type mockImageRepo struct{}

func (m *mockImageRepo) AddFile(name string, file []byte) error {
	return nil
}

func (m *mockImageRepo) GetFile(name string) ([]byte, error) {
	return []byte{}, nil
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

	_, err := api.NewWardrobeService(mockWardrobe, mockImage)
	if err != nil {
		t.Errorf(" NewWardrobService failed : %v", err)
	}


}
